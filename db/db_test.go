package db_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/wiptrax/dsitributed-kv-store/db"
)

func createTempDB(t *testing.T, readOnly bool) *db.DataBase{
	t.Helper()

	f, err := os.CreateTemp(os.TempDir(), "kvdb")
	if err != nil {
		t.Fatalf("Could not create temp filr: %v", err)
	}

	name := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(name) })

	db, close, err := db.NewDatabase(name, readOnly)
	if err != nil {
		t.Fatalf("Could not create a new dtabase: %v", err)
	}
	t.Cleanup(func() { close() })

	return db
}

func TestGetSet(t *testing.T) {
	db := createTempDB(t, false)

	if err := db.SetKey("testKey", []byte("testValue")); err != nil {
		t.Fatalf("Could not write key: %v", err)
	}

	value, err := db.GetKey("testKey")
	if err != nil {
		t.Fatalf(`Could not get key "testKey": %v`, err)
	}

	if !bytes.Equal(value, []byte("testValue")) {
		t.Errorf(`Unexpected value for key "tesKey": got %q, want %q`, value, "testValue")
	}
}

func TestDeleteReplicationKey(t *testing.T) {
	db := createTempDB(t, false)

	setKey(t, db, "testKey", "testValue")

	k, v, err := db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if !bytes.Equal(k, []byte("testKey")) || !bytes.Equal(v, []byte("testValue")) {
		t.Errorf(`GetNextKeyForReplication(): got %q, %q; want %q, %q`, k, v, "testKey", "testValue")
	}

	if err := db.DeleteReplicationKey([]byte("testKey"), []byte("Bad")); err == nil {
		t.Fatalf(`DeleteReplicationKey("testKey", "Bad"): got nil error, want non-nil error`)
	}

	if err := db.DeleteReplicationKey([]byte("testKey"), []byte("testValue")); err != nil {
		t.Fatalf(`DeleteReplicationKey("testKey", "testValue"): got %q, want nil error`, err)
	}

	k, v, err = db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if k != nil || v != nil {
		t.Errorf(`GetNextKeyForReplication(): got %v, %v; want nil, nil`, k, v)
	}
}

func TestSetReadOnly(t *testing.T) {
	db := createTempDB(t, true)

	if err := db.SetKey("testKey", []byte("Bad")); err == nil {
		t.Fatalf("SetKey(%q, %q): got nil error, want non-nil error", "testKey", []byte("Bad"))
	}
}

func setKey(t *testing.T, d *db.DataBase, key, value string) {
	t.Helper()

	if err := d.SetKey(key, []byte(value)); err != nil {
		t.Fatalf("SetKey(%q, %q) failed; %v", key, value, err)
	}
}

func getKey(t *testing.T, d *db.DataBase, key string) string {
	t.Helper()

	value, err := d.GetKey(key)
	if err != nil {
		t.Fatalf("GetKey(%q) failed; %v", key, err)
	}

	return string(value)
}

func TestDeleteExtraKeys(t *testing.T) {
	db :=  createTempDB(t, false)

	setKey(t, db, "testKey", "testValue")
	setKey(t, db, "testKey101", "testValue101")

	if err := db.DeleteExtraKeys(func(name string) bool { return name == "testKey101" }); err != nil {
		t.Fatalf("Could not delete extra keys: %v", err)
	}

	if value := getKey(t, db, "testKey"); value != "testValue" {
		t.Errorf(`Unexpected value for key "testKey": got %q, want %q`, value, "testValue")
	}

	if value := getKey(t, db, "testKey101"); value != "" {
		t.Errorf(`Unexpected value for key "testKey": got %q, want %q`, value, "")
	}
}
