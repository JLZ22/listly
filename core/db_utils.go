package core

import (
	"encoding/binary"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// ------------------------------ Datatype Conversion Helper Functions ------------------------------

func boolToBytes(b bool) []byte {
	var v uint8
	if b {
		v = 1
	} else {
		v = 0
	}
	return []byte{v}
}

func bytesToBool(b []byte) bool {
	return len(b) > 0 && b[0] == 1
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func intsToBytes(ints []int) []byte {
	b := make([]byte, 8*len(ints)) // 8 bytes per int64
	for i, v := range ints {
		binary.BigEndian.PutUint64(b[i*8:(i+1)*8], uint64(v))
	}
	return b
}

func bytesToInts(b []byte) []int {
	n := len(b) / 8
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = int(binary.BigEndian.Uint64(b[i*8 : (i+1)*8]))
	}
	return ints
}

// ------------------------------------- Transaction Helper Functions ---------------------------------

// Save fields from the given Task struct into the bucket
func saveTask(bucket *bolt.Bucket, task Task) error {
	taskBucket, err := bucket.CreateBucketIfNotExists(itob(task.Id))
	if err != nil {
		return err
	}
	taskBucket.Put([]byte("id"), itob(task.Id))
	taskBucket.Put([]byte("description"), []byte(task.Description))
	taskBucket.Put([]byte("done"), boolToBytes(task.Done))
	return nil
}

// Populate fields of Task struct by reading from the given bucket.
func getTask(bucket *bolt.Bucket) (Task, error) {
	task := &Task{}

	description := bucket.Get([]byte("description"))
	if description == nil {
		return Task{}, fmt.Errorf("description not found")
	}
	done := bucket.Get([]byte("done"))
	if done == nil {
		return Task{}, fmt.Errorf("done not found")
	}
	id := bucket.Get([]byte("id"))
	if id == nil {
		return Task{}, fmt.Errorf("id not found")
	}

	task.Id = btoi(id)
	task.Description = string(description)
	task.Done = bytesToBool(done)

	return *task, nil
}

// Save fields from ListInfo struct into the bucket
func saveInfo(bucket *bolt.Bucket, info ListInfo) error {
	if info.Name == "" {
		return fmt.Errorf("name is required")
	}
	err := bucket.Put([]byte("name"), []byte(info.Name))
	if err != nil {
		return err
	}
	err = bucket.Put([]byte("numDone"), itob(info.NumDone))
	if err != nil {
		return err
	}
	err = bucket.Put([]byte("numPending"), itob(info.NumPending))
	if err != nil {
		return err
	}
	return bucket.Put([]byte("numTasks"), itob(info.NumTasks))
}

// Populate and return a ListInfo struct by reading fields from the given bucket.
func getInfo(bucket *bolt.Bucket) (ListInfo, error) {
	info := ListInfo{
		Name:       "",
		NumDone:    -1,
		NumPending: -1,
		NumTasks:   -1,
	}

	name := bucket.Get([]byte("name"))
	if len(name) == 0 {
		return info, fmt.Errorf("name not found")
	}
	numDone := bucket.Get([]byte("numDone"))
	if numDone == nil {
		return info, fmt.Errorf("numDone not found")
	}
	numPending := bucket.Get([]byte("numPending"))
	if numPending == nil {
		return info, fmt.Errorf("numPending not found")
	}
	numTasks := bucket.Get([]byte("numTasks"))
	if numTasks == nil {
		return info, fmt.Errorf("numTasks not found")
	}

	info.Name = string(name)
	info.NumDone = btoi(numDone)
	info.NumPending = btoi(numPending)
	info.NumTasks = btoi(numTasks)
	return info, nil
}

// Save data part of List struct into given bucket.
func saveData(bucket *bolt.Bucket, list List) error {
	// save ordered task ids
	err := bucket.Put([]byte("taskIds"), intsToBytes(list.TaskIds))
	if err != nil {
		return err
	}

	// save tasks
	taskBucket, err := bucket.CreateBucketIfNotExists([]byte("tasks"))
	if err != nil {
		return err
	}
	for _, task := range list.Tasks {
		err = saveTask(taskBucket, *task)
		if err != nil {
			return err
		}
	}

	// skip saving usedIds because it is cheaper to reconstruct it when loading
	return nil
}

// Populate non-info fields of List struct by reading from the given bucket.
func getData(bucket *bolt.Bucket) (List, error) {
	list := NewList("")

	taskIdsBytes := bucket.Get([]byte("taskIds"))
	if taskIdsBytes == nil {
		return list, fmt.Errorf("taskIds not found")
	}

	taskBucket := bucket.Bucket([]byte("tasks"))
	if taskBucket == nil {
		return list, fmt.Errorf("tasks bucket not found")
	}

	taskIds := bytesToInts(taskIdsBytes)
	tasksMap := make(map[int]*Task)
	usedIds := make(map[int]struct{})
	for _, id := range taskIds {
		currTaskBucket := taskBucket.Bucket(itob(id))
		if currTaskBucket == nil {
			return list, fmt.Errorf("task bucket %d not found", id)
		}
		task, err := getTask(currTaskBucket)
		if err != nil {
			return list, err
		}
		usedIds[id] = struct{}{}
		tasksMap[id] = &task
	}

	list.TaskIds = taskIds
	list.Tasks = tasksMap
	list.UsedIds = usedIds
	return list, nil
}

// Retrieve the info and data sub-buckets associated with the bucket with
// the given list name, and return an error for missing buckets depending
// on the existOkay flag.
func openList(b *bolt.Bucket, listName string, notExistOk bool) (infoBucket, dataBucket *bolt.Bucket, err error) {
	creationFn := func(b *bolt.Bucket, fieldName string) (*bolt.Bucket, error) {
		if notExistOk {
			return b.CreateBucketIfNotExists([]byte(fieldName))
		}
		out := b.Bucket([]byte(fieldName))
		if out == nil {
			return nil, fmt.Errorf("bucket %s not found", fieldName)
		}
		return out, nil
	}

	// get the bucket that stores the list we are interested in
	listBucket, err := creationFn(b, listName)
	if err != nil {
		return nil, nil, err
	}

	// get meta data bucket
	infoBucket, err = creationFn(listBucket, "info")
	if err != nil {
		return nil, nil, err
	}

	// get actual data bucket
	dataBucket, err = creationFn(listBucket, "data")
	if err != nil {
		return nil, nil, err
	}

	return infoBucket, dataBucket, nil
}

func copyBucket(src, dst *bolt.Bucket) error {
	return src.ForEach(func(k, v []byte) error {
		if v == nil {
			subSrc := src.Bucket(k)
			if subSrc == nil {
				return fmt.Errorf("sub-bucket %s not found", k)
			}
			subDst, err := dst.CreateBucket(k)
			if err != nil {
				return err
			}
			return copyBucket(subSrc, subDst)
		}
		return dst.Put(k, v)
	})
}

func getCurrListName(tx *bolt.Tx) string {
	currentList := tx.Bucket([]byte("currentList"))
	if currentList == nil {
		return ""
	}

	nameBytes := currentList.Get([]byte("name"))
	if nameBytes == nil {
		return ""
	}

	return string(nameBytes)
}

func setCurrListName(tx *bolt.Tx, name string) error {
	currentList := tx.Bucket([]byte("currentList"))
	if currentList == nil {
		return fmt.Errorf("currentList bucket not found - likely issue with database initialization")
	}

	return currentList.Put([]byte("name"), []byte(name))
}

// delete all tasks that are marked as done
func cleanList(b *bolt.Bucket) error {
	dataBucket := b.Bucket([]byte("data"))
	if dataBucket == nil {
		return fmt.Errorf("data bucket not found")
	}

	// remove tasks that are done
	var numRemoved int
	err := dataBucket.ForEach(func(k, v []byte) error {
		if v == nil {
			taskBucket := dataBucket.Bucket(k)
			if taskBucket == nil {
				return fmt.Errorf("task bucket %s not found", string(k))
			}

			task, err := getTask(taskBucket)
			if err != nil {
				return nil // skip if task bucket is empty or doesn't exist
			}

			if task.Done {
				numRemoved++
				return dataBucket.DeleteBucket(k)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// update info with number of total and done tasks
	infoBucket := b.Bucket([]byte("info"))
	if infoBucket == nil {
		return fmt.Errorf("info bucket not found")
	}

	info, err := getInfo(infoBucket)
	if err != nil {
		return err
	}

	info.NumTasks -= numRemoved
	info.NumDone = 0
	return saveInfo(infoBucket, info)
}
