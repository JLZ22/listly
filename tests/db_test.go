package core_test

import (
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
	list.Info.Name = "list1"

	err := db.SaveList(list)
	require.NoError(t, err)

	got, err := db.GetList("list1")
	require.NoError(t, err)
	require.Equal(t, "list1", got.Info.Name)
}

func TestRenameList(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	// Setup list bucket manually for rename test
	err := db.BoltDB.Update(func(tx *bbolt.Tx) error {
		lists := tx.Bucket([]byte("lists"))
		l, err := lists.CreateBucket([]byte("old"))
		if err != nil {
			return err
		}
		_, err = l.CreateBucket([]byte("info"))
		return err
	})
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

	err = db.DeleteList("todelete")
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

	_, err := db.GetCurrentListName()
	require.Error(t, err)
	require.Contains(t, err.Error(), "current list name not found")
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

func TestDeleteList_NotFound(t *testing.T) {
	db, cleanup := setupTempDB(t)
	defer cleanup()

	// Try deleting a bucket that doesn't exist
	err := db.DeleteList("nope")
	require.Error(t, err)
	require.Contains(t, err.Error(), "bucket not found")
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
