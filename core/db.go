package core

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// schema:
// ------------------------------------------------------
// 		"currentList": "string",
// 		"lists": {
// 			"listName": {
// 				"info": {
// 					"name": "string",
// 					"numDone": int,
// 					"numPending": int,
// 					"numTasks": int
// 				},
// 				"data": {
// 					"taskIds": []int,
// 					"tasks": {
// 						"taskId": {
// 							"description": "string",
// 							"done": false
// 						},
// 						...
// 					}
// 				}
// 			},
// 			...
// 		}
// ------------------------------------------------------

type DB struct {
	BoltDB *bolt.DB
}

// ---------------------------- Setup Functions --------------------------------

// Initialize the database. For testing purposes, a custom path can be provided,
// but it is recommended to use `os.UserConfigDir()`.
func InitDB(path string) (*DB, error) {
	// load and/or create the database
	dbPath := path + "/listly.db"
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("currentList"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("lists"))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &DB{BoltDB: db}, nil
}

// ------------------------ business logic abstractions --------------------------------

// gets the metadata for every list stored in the database
func (db *DB) GetInfo() (map[string]ListInfo, error) {
	allInfo := make(map[string]ListInfo)

	err := db.BoltDB.View(func(tx *bolt.Tx) error {
		allLists := tx.Bucket([]byte("lists"))
		if allLists == nil {
			return fmt.Errorf("lists bucket not found - likely issue with database initialization")
		}

		// for every list in allLists, extract info and put it in allInfo
		return allLists.ForEach(func(k, v []byte) error {
			if v == nil { // k is a sub-bucket name
				list := allLists.Bucket(k)
				if list == nil {
					return fmt.Errorf("sub-bucket %s not found in lists bucket", k)
				}

				infoBucket := list.Bucket([]byte("info"))
				if infoBucket == nil {
					return fmt.Errorf("info bucket not found in list bucket %s", k)
				}

				info, err := getInfo(infoBucket)
				if err != nil {
					return err
				}
				allInfo[info.Name] = info
			}
			return nil
		})
	})

	return allInfo, err
}

// get a specific list by name
func (db *DB) GetList(name string) (List, error) {
	list := NewList(name)
	err := db.BoltDB.View(func(tx *bolt.Tx) error {
		allLists := tx.Bucket([]byte("lists"))
		if allLists == nil {
			return fmt.Errorf("lists bucket not found - likely issue with database initialization")
		}

		infoBucket, dataBucket, err := openList(allLists, name, false)
		if err != nil {
			return fmt.Errorf("failed to open list %s: %w", name, err)
		}

		info, err := getInfo(infoBucket)
		if err != nil {
			return fmt.Errorf("failed to get info for list %s: %w", name, err)
		}

		data, err := getData(dataBucket)
		if err != nil {
			return fmt.Errorf("failed to get data for list %s: %w", name, err)
		}

		list.Info = info
		list.Tasks = data.Tasks
		list.TaskIds = data.TaskIds
		list.UsedIds = data.UsedIds

		return nil
	})
	return list, err
}

// get the name of the currently active list
func (db *DB) GetCurrentListName() (string, error) {
	name := ""
	err := db.BoltDB.View(func(tx *bolt.Tx) error {
		currentList := tx.Bucket([]byte("currentList"))
		if currentList == nil {
			return fmt.Errorf("currentList bucket not found - likely issue with database initialization")
		}

		nameBytes := currentList.Get([]byte("name"))
		if nameBytes == nil {
			return fmt.Errorf("current list name not found")
		}

		name = string(nameBytes)
		return nil
	})
	return name, err
}

// set the name of the currently active list
func (db *DB) SetCurrentListName(name string) error {
	return db.BoltDB.Update(func(tx *bolt.Tx) error {
		currentList := tx.Bucket([]byte("currentList"))
		if currentList == nil {
			return fmt.Errorf("currentList bucket not found - likely issue with database initialization")
		}

		return currentList.Put([]byte("name"), []byte(name))
	})
}

// save the given list
func (db *DB) SaveList(list List) error {
	return db.BoltDB.Update(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket([]byte("lists"))
		if rootBucket == nil {
			return fmt.Errorf("lists bucket not found - likely issue with database initialization")
		}

		infoBucket, dataBucket, err := openList(rootBucket, list.Info.Name, true)
		if err != nil {
			return err
		}

		// save info into meta data bucket
		err = saveInfo(infoBucket, list.Info)
		if err != nil {
			return err
		}

		// save data into data bucket
		err = saveData(dataBucket, list)
		if err != nil {
			return err
		}

		return nil
	})
}

// Rename the list with the oldName to the newName. Since bbolt
// does not support renaming buckets, we need to completely
// recreate the bucket with the new name.
func (db *DB) RenameList(oldName, newName string) error {
	return db.BoltDB.Update(func(tx *bolt.Tx) error {
		listsBucket := tx.Bucket([]byte("lists"))
		if listsBucket == nil {
			return fmt.Errorf("lists bucket not found")
		}

		oldBucket := listsBucket.Bucket([]byte(oldName))
		if oldBucket == nil {
			return fmt.Errorf("old list %s not found", oldName)
		}

		// Create new bucket with newName
		newBucket, err := listsBucket.CreateBucket([]byte(newName))
		if err != nil {
			return err
		}

		// Recursively copy all keys/sub-buckets
		if err := copyBucket(oldBucket, newBucket); err != nil {
			return err
		}

		// Delete old bucket
		return listsBucket.DeleteBucket([]byte(oldName))
	})
}

// remove the list with the given name
func (db *DB) DeleteList(name string) error {
	return db.BoltDB.Update(func(tx *bolt.Tx) error {
		listsBucket := tx.Bucket([]byte("lists"))
		if listsBucket == nil {
			return fmt.Errorf("lists bucket not found")
		}

		return listsBucket.DeleteBucket([]byte(name))
	})
}
