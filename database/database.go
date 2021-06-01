package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

const (
	defaultFile = "bolt.db"
)

type DB struct {
	boltDb *bolt.DB
}

func NewDb(file string) *DB {
	if strings.TrimSpace(file) == "" {
		file = defaultFile
	}
	var err error
	db, err := bolt.Open(file, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{boltDb: db}
}

func (db *DB) Shutdown() {
	db.boltDb.Close()
}

func (d *DB) CreateBucketIfNotExists(bucketType BucketType) error {
	err := d.boltDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketType))
		return err
	})

	if err != nil {
		return fmt.Errorf("Error creating bucket %s: %s", bucketType, err)
	}
	return nil
}

func (d *DB) Put(bucketType BucketType, key string, value interface{}) error {
	valueToWrite, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Error marshalling value '%v': %s", value, err)
	}
	err = d.boltDb.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketType))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), valueToWrite)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error writing key %s to bucket %s: %s", key, bucketType, err)
	}
	return nil
}

func (db *DB) Get(bucketType BucketType, key string) ([]byte, error) {
	var result []byte
	err := db.boltDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketType))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucketType)
		}

		result = bucket.Get([]byte(key))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error getting key %s from bucket %s: %s", key, bucketType, err)
	}
	return result, nil
}

func (db *DB) PutRandom(bucketType BucketType, value interface{}) error {
	return db.Put(bucketType, string(NewHashedKey(value)), value)
}

func (db *DB) GetRandom(bucketType BucketType) ([]byte, error) {
	var result []byte
	err := db.boltDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketType))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucketType)
		}

		cursor := bucket.Cursor()

		// Make sure at least one element exists
		k, _ := cursor.First()
		if k == nil {
			return nil
		}

		// Loop until a result is found. If the uuid in Seek is larger than the largest written
		// result nil is return and we should re-seek
		for {
			_, result = cursor.Seek([]byte(uuid.NewString()))
			if result != nil {
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error getting random item from bucket %s: %s", bucketType, err)
	}
	return result, nil
}
