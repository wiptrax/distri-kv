package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")

// Database is a open bolt database
type DataBase struct {
	db *bolt.DB
}

// NewDatabase returns an instance of a dtabse that we can work with
func NewDatabase(dbPath string) (db *DataBase, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	db = &DataBase{db: boltDb}
	closeFunc = boltDb.Close

	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating default bucket: %v", err)
	}

	return db, closeFunc, nil
}

func (d *DataBase) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

// SetKey sets the key to the requested value into the default database or return an error
func (d *DataBase) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

// GetKey gets the value of the requested fram a default databse
func (d *DataBase) GetKey(key string) ([]byte, error) {
	var result []byte

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return result, err
	}

	return nil, err
}
