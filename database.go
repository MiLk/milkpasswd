package milkpasswd

import (
	"bytes"
	"errors"
	"os/user"
	"path"
	"time"

	"github.com/boltdb/bolt"
)

var dbPath string

const bucketName string = "milkpasswd"

func dbOpen() (*bolt.DB, error) {
	if dbPath == "" {
		currentUser, err := user.Current()
		if err != nil {
			return nil, err
		}
		dbPath = path.Join(currentUser.HomeDir, ".milkpasswd")
	}
	return bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
}

func setRecord(key string, value []byte) (err error) {
	db, err := dbOpen()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		err = b.Put([]byte(key), value)
		return err
	})
	if err != nil {
		return
	}
	return
}

func deleteRecord(key string) (err error) {
	db, err := dbOpen()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		err = b.Delete([]byte(key))
		return err
	})
	if err != nil {
		return
	}
	return
}

func listRecords() (data map[string][]byte, err error) {
	db, err := dbOpen()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		data = make(map[string][]byte)
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := string(k)
			data[key] = make([]byte, len(v))
			copy(data[key], v)
		}

		return nil
	})
	if err != nil {
		data = nil
		return
	}

	return
}

func searchRecords(prefix string) (data map[string][]byte, err error) {
	db, err := dbOpen()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		bPrefix := []byte(prefix)
		data = make(map[string][]byte)
		for k, v := c.Seek(bPrefix); bytes.HasPrefix(k, bPrefix); k, v = c.Next() {
			key := string(k)
			data[key] = make([]byte, len(v))
			copy(data[key], v)
		}

		return nil
	})
	if err != nil {
		data = nil
		return
	}

	return
}

func getRecord(key string) (data []byte, err error) {
	db, err := dbOpen()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(key))
		if v == nil {
			return errors.New("milkpasswd database: record not found in the database")
		}

		data = make([]byte, len(v))
		copy(data, v)

		return nil
	})
	if err != nil {
		data = nil
		return
	}

	return
}
