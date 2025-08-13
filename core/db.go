package core

import (
	"fmt"
	"os"

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

// Open DB located in user's default config dir.
func InitDefaultDB() (*DB, error) {
	dbPath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return InitDB(dbPath + "/listly")
}

// Initialize the database. For testing purposes, a custom path can be provided,
// but it is recommended to use `os.UserConfigDir()`.
func InitDB(path string) (*DB, error) {
	// Create the directories for the DB
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return nil, err
	}

	// Open the DB
	dbPath := path + "/listly.db"
	db, err := bolt.Open(dbPath, 0700, nil)
	if err != nil {
		return nil, err
	}

	// Populate DB with the top level buckets if they don't yet exist
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

// A utility function that simplifies the usage of the default database.
func WithDefaultDB(fn func(db *DB) error) error {
	db, err := InitDefaultDB()
	if err != nil {
		return fmt.Errorf("could not initialize database due to the following error\n\t %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("could not close database due to the following error\n\t %v\n", err)
		}
	}()
	return fn(db)
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

		_, dataBucket, err := openList(allLists, name, false)
		if err != nil {
			return fmt.Errorf("failed to open list %s: %w", name, err)
		}

		data, err := getData(dataBucket)
		if err != nil {
			return fmt.Errorf("failed to get data for list %s: %w", name, err)
		}

		list.Tasks = data.Tasks
		list.TaskIds = data.TaskIds
		list.UsedIds = data.UsedIds

		// align list info with data
		list.Info.NumTasks = len(list.Tasks)
		list.Info.NumDone = 0
		list.Info.NumPending = 0
		for _, task := range list.Tasks {
			if task.Done {
				list.Info.NumDone++
			} else {
				list.Info.NumPending++
			}
		}

		return nil
	})
	return list, err
}

// get the name of the currently active list
func (db *DB) GetCurrentListName() (string, error) {
	var name string
	err := db.BoltDB.View(func(tx *bolt.Tx) error {
		name = getCurrListName(tx)
		return nil
	})
	return name, err
}

// set the name of the currently active list
func (db *DB) SetCurrentListName(name string) error {
	return db.BoltDB.Update(func(tx *bolt.Tx) error {
		if name == "" {
			return fmt.Errorf("cannot have empty name")
		}
		return setCurrListName(tx, name)
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

		// align list info with data
		list.Info.NumTasks = len(list.Tasks)
		list.Info.NumDone = 0
		list.Info.NumPending = 0
		for _, task := range list.Tasks {
			if task.Done {
				list.Info.NumDone++
			} else {
				list.Info.NumPending++
			}
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
		allLists := tx.Bucket([]byte("lists"))
		if allLists == nil {
			return fmt.Errorf("lists bucket not found")
		}

		// Get the old list bucket
		oldBucket := allLists.Bucket([]byte(oldName))
		if oldBucket == nil {
			return fmt.Errorf("old list %s not found", oldName)
		}

		// Create new bucket with newName
		newBucket, err := allLists.CreateBucket([]byte(newName))
		if err != nil {
			return fmt.Errorf("could not create new bucket %s due to the following error\n\t %w", newName, err)
		}

		// Recursively copy all keys/sub-buckets
		if err = copyBucket(oldBucket, newBucket); err != nil {
			return fmt.Errorf("could not copy bucket %s to %s due to the following error\n\t %w", oldName, newName, err)
		}

		// Delete old bucket
		err = allLists.DeleteBucket([]byte(oldName))
		if err != nil {
			return err
		}

		// update info in the new bucket
		listBucket := allLists.Bucket([]byte(newName))
		if listBucket == nil {
			return fmt.Errorf("list bucket %s not found", newName)
		}
		infoBucket := listBucket.Bucket([]byte("info"))
		if infoBucket == nil {
			return fmt.Errorf("info bucket not found for list %s", newName)
		}
		listInfo, err := getInfo(infoBucket)
		if err != nil {
			return err
		}
		listInfo.Name = newName
		err = saveInfo(infoBucket, listInfo)
		if err != nil {
			return err
		}

		// if the list being renamed is the current list, update the current list name
		currListName := getCurrListName(tx)
		if currListName == oldName {
			return setCurrListName(tx, newName)
		}
		return nil
	})
}

// remove the list with the given name
func (db *DB) DeleteLists(names []string) error {
	return db.BoltDB.Update(func(tx *bolt.Tx) error {
		allLists := tx.Bucket([]byte("lists"))
		if allLists == nil {
			return fmt.Errorf("lists bucket not found")
		}

		currListName := getCurrListName(tx)
		for _, name := range names {
			if name == currListName {
				setCurrListName(tx, "")
			}
			allLists.DeleteBucket([]byte(name))
		}
		return nil
	})
}

// remove all lists
func (db *DB) DeleteAllLists() error {
	return db.BoltDB.Update(func(tx *bolt.Tx) error {
		allLists := tx.Bucket([]byte("lists"))
		if allLists == nil {
			return fmt.Errorf("lists bucket not found")
		}

		setCurrListName(tx, "")
		return allLists.ForEach(func(k, v []byte) error {
			return allLists.DeleteBucket(k)
		})
	})
}

// Clean up completed tasks in the specified lists
func (db *DB) CleanLists(names []string) (int, error) {
	var totalRemoved int
	err := db.BoltDB.Update(func(tx *bolt.Tx) error {
		allBuckets := tx.Bucket([]byte("lists"))
		if allBuckets == nil {
			return fmt.Errorf("lists bucket not found")
		}

		for _, name := range names {
			listBucket := allBuckets.Bucket([]byte(name))
			if listBucket == nil {
				continue // skip if the list does not exist
			}

			numRemoved, err := cleanList(listBucket)
			if err != nil {
				return err
			}
			totalRemoved += numRemoved
		}
		return nil
	})
	return totalRemoved, err
}

// Clean up completed tasks in all lists
func (db *DB) CleanAllLists() (int, error) {
	var totalRemoved int
	err := db.BoltDB.Update(func(tx *bolt.Tx) error {
		allBuckets := tx.Bucket([]byte("lists"))
		if allBuckets == nil {
			return fmt.Errorf("lists bucket not found")
		}

		return allBuckets.ForEach(func(k, v []byte) error {
			listBucket := allBuckets.Bucket(k)
			if listBucket == nil {
				return nil // skip if the list does not exist
			}

			numRemoved, err := cleanList(listBucket)
			if err != nil {
				return err
			}
			totalRemoved += numRemoved
			return nil
		})
	})
	return totalRemoved, err
}

// Clean up completed tasks in the current list
func (db *DB) CleanCurrentList() (int, error) {
	var totalRemoved int
	err := db.BoltDB.Update(func(tx *bolt.Tx) error {
		name := getCurrListName(tx)

		allBuckets := tx.Bucket([]byte("lists"))
		if allBuckets == nil {
			return fmt.Errorf("lists bucket not found")
		}

		listBucket := allBuckets.Bucket([]byte(name))
		if listBucket == nil {
			return nil // skip if the list does not exist
		}
		numRemoved, err := cleanList(listBucket)
		if err != nil {
			return err
		}
		totalRemoved += numRemoved
		return nil
	})
	return totalRemoved, err
}

func (db *DB) ListExists(name string) (bool, error) {
	allInfo, err := db.GetInfo()
	if err != nil {
		return false, err
	}

	_, ok := allInfo[name]
	return ok, nil
}

func (db *DB) Close() error {
	return db.BoltDB.Close()
}
