package db

import (
	"bytes"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")
var replicaBucket = []byte("replication")

// Database is a open bolt database
type DataBase struct {
	db       *bolt.DB
	ReadOnly bool
}

// NewDatabase returns an instance of a dtabse that we can work with
func NewDatabase(dbPath string, readonly bool) (db *DataBase, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	db = &DataBase{db: boltDb, ReadOnly: readonly}
	closeFunc = boltDb.Close

	if err := db.createDefaultBuckets(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating default bucket: %v", err)
	}

	return db, closeFunc, nil
}

func (d *DataBase) createDefaultBuckets() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(defaultBucket); err != nil {
		return err
		}
		if _, err := tx.CreateBucketIfNotExists(replicaBucket); err != nil {
		return err
		}
		return nil
	})
}

// SetKeyOnReplica sets the key to the requested value into the default database and does not write
// to the replication queue.
// This method is intended to be used only on replicas.
func (d *DataBase) SetKeyOnReplica(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(defaultBucket).Put([]byte(key), value)
	})
}

func copyByteSlice(b []byte) []byte {
	if b == nil {
		return nil
	}
	res := make([]byte, len(b))
	copy(res, b)
	return res
}

// GetNextKeyForReplication returns the key and value for the keys that have
// changed and have not yet been applied to replicas.
// If there are no new keys, nil key and value will be returned.
func (d *DataBase) GetNextKeyForReplication() (key, value []byte, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(replicaBucket)
		k, v := b.Cursor().First()
		key = copyByteSlice(k)
		value = copyByteSlice(v)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return key, value, nil
}

// SetKey sets the key to the requested value into the default database or return an error
func (d *DataBase) SetKey(key string, value []byte) error {
	if d.ReadOnly {
		return errors.New("read-only-mode")
	}


	return d.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(defaultBucket).Put([]byte(key), value); err != nil {
			return err
		}

		return tx.Bucket(replicaBucket).Put([]byte(key), value)
	})
}

// DeleteReplicationKey deletes the key from the replication queue
// if the value matches the contents or if the key is already absent.
func (d *DataBase) DeleteReplicationKey(key, value []byte) (err error) {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(replicaBucket)

		v := b.Get(key)
		if v == nil {
			return errors.New("key does not exist")
		}

		if !bytes.Equal(v, value) {
			return errors.New("value does not match")
		}

		return b.Delete(key)
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

// DeleteExtraKeys deletes the keys that do not belongs to shard
func (d *DataBase) DeleteExtraKeys(isExtra func(string) bool) error {
	var keys []string

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.ForEach(func(k, v []byte) error {
			ks := string(k)
			if isExtra(ks) {
				keys = append(keys, string(ks))
			}
			return nil
		})
	})

	if err != nil {
		return err
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)

		for _, k := range keys {
			if err := b.Delete([]byte(k)); err != nil {
				return err
			}
		}
		return nil
	})
}
