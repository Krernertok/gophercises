package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const DefaultBucket = "URLShort"

func AddData(dbName string) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(DefaultBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		pairs := [][]string{
			[]string{"/urlshort", "https://github.com/gophercises/urlshort"},
			[]string{"/red", "https://www.reddit.com"},
		}

		for _, pair := range pairs {
			err = b.Put([]byte(pair[0]), []byte(pair[1]))
			if err != nil {
				return fmt.Errorf("add key: %s", pair[0])
			}
		}

		return nil
	})
}

func GetPathMap(dbName string, bucket []byte) map[string]string {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil
	}
	defer db.Close()

	pathMap := make(map[string]string)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		b.ForEach(func(k, v []byte) error {
			pathMap[string(k)] = string(v)
			return nil
		})

		return nil
	})

	return pathMap
}
