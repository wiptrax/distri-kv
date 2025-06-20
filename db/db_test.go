package db_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/wiptrax/dsitributed-kv-store/db"
)

func TestGetSet(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "kvdb")
	if err != nil {
		t.Fatalf("Could not create temp filr: %v", err)
	}

	name := f.Name()
	f.Close()
	defer os.Remove(name)

	db, close, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("Could not create a new dtabase: %v", err)
	}
	defer close()

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
