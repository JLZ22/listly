package core_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jlz22/listly/core"
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
)

// -- Minimal stub types and helpers to satisfy dependencies --

type List = core.List
type ListInfo = core.ListInfo

// -- Tests --

func setupTempDB(t *testing.T) (*core.DB, func()) {
	tmpDir := t.TempDir()
	db, err := core.InitDB(tmpDir)
	require.NoError(t, err)
	return db, func() {
		db.BoltDB.Close()
		os.RemoveAll(tmpDir)
	}
}

func TestInitDB(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	require.NotNil(t, db.BoltDB)

	// Check buckets exist
	err := db.BoltDB.View(func(tx *bbolt.Tx) error {
		require.NotNil(t, tx.Bucket([]byte("lists")))
		require.NotNil(t, tx.Bucket([]byte("currentList")))
		return nil
	})
	require.NoError(t, err)
}

func TestSetAndGetCurrentListName(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	err := db.SetCurrentListName("mylist")
	require.NoError(t, err)

	name, err := db.GetCurrentListName()
	require.NoError(t, err)
	require.Equal(t, "mylist", name)
}

func TestSaveListAndGetList(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	list := core.NewList("list1")

	for i := 0; i < 5; i++ {
		list.AddNewTask(fmt.Sprintf("task%d", i+1), i%2 == 0)
		require.Equal(t, fmt.Sprintf("task%d", i+1), list.Tasks[list.TaskIds[i]].Description)
		require.Equal(t, list.Tasks[list.TaskIds[i]].Done, i%2 == 0)
		require.Equal(t, list.TaskIds[i], list.Tasks[list.TaskIds[i]].Id)
	}

	err := db.SaveList(list)
	require.NoError(t, err)

	got, err := db.GetList("list1")
	require.NoError(t, err)
	for i := 0; i < 5; i++ {
		require.Equal(t, fmt.Sprintf("task%d", i+1), got.Tasks[got.TaskIds[i]].Description)
		require.Equal(t, i%2 == 0, got.Tasks[got.TaskIds[i]].Done)
		require.Equal(t, got.TaskIds[i], list.TaskIds[i])
		require.Equal(t, got.TaskIds[i], got.Tasks[got.TaskIds[i]].Id)
	}
}

func TestRenameList(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	// Setup list bucket manually for rename test
	dummyList := core.NewList("old")
	err := db.SaveList(dummyList)
	require.NoError(t, err)

	err = db.RenameList("old", "new")
	require.NoError(t, err)

	err = db.BoltDB.View(func(tx *bbolt.Tx) error {
		lists := tx.Bucket([]byte("lists"))
		require.Nil(t, lists.Bucket([]byte("old")))
		require.NotNil(t, lists.Bucket([]byte("new")))
		return nil
	})
	require.NoError(t, err)
}

func TestDeleteList(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	// Setup list bucket manually for delete test
	err := db.BoltDB.Update(func(tx *bbolt.Tx) error {
		lists := tx.Bucket([]byte("lists"))
		_, err := lists.CreateBucket([]byte("todelete"))
		return err
	})
	require.NoError(t, err)

	err = db.DeleteLists([]string{"todelete"})
	require.NoError(t, err)

	err = db.BoltDB.View(func(tx *bbolt.Tx) error {
		lists := tx.Bucket([]byte("lists"))
		require.Nil(t, lists.Bucket([]byte("todelete")))
		return nil
	})
	require.NoError(t, err)
}

func TestGetInfoEmpty(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	infos, err := db.GetInfo()
	require.NoError(t, err)
	require.Empty(t, infos)
}

func TestGetCurrentListName_NoNameSet(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	curr, err := db.GetCurrentListName()
	require.NoError(t, err)
	require.Equal(t, "", curr) // No current list should return empty string
}

func TestGetList_NotFound(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	_, err := db.GetList("doesnotexist")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to open list")
}

func TestRenameList_OldNotFound(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	err := db.RenameList("missing", "new")
	require.Error(t, err)
	require.Contains(t, err.Error(), "old list missing not found")
}

func TestGetInfo_MissingInfoBucket(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	// Create a list bucket without "info"
	err := db.BoltDB.Update(func(tx *bbolt.Tx) error {
		lists := tx.Bucket([]byte("lists"))
		_, err := lists.CreateBucket([]byte("badlist"))
		return err
	})
	require.NoError(t, err)

	_, err = db.GetInfo()
	require.Error(t, err)
	require.Contains(t, err.Error(), "info bucket not found")
}
