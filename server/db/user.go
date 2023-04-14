package db

import (
	"fmt"

	"go.etcd.io/bbolt"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"pass"`
	IsAdmin  bool   `json:"admin"`
	Banned   bool   `json:"ban"`
}

var ErrUserNotFound error = fmt.Errorf("user not found")
var userBucket []byte = []byte("user")

func GetUser(userID string) (User, error) {
	ret := User{}
	err := gDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(userBucket)
		if bucket == nil {
			return ErrUserNotFound
		}

		value := bucket.Get([]byte(userID))
		if value == nil {
			return ErrUserNotFound
		}

		if err := decodeValue(value, &ret); err != nil {
			return err
		}

		return nil
	})

	return ret, err
}

func PutUser(user User) error {
	return gDB.Update(func(tx *bbolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists(userBucket)
		data := encodeValue(user)
		return bucket.Put([]byte(user.ID), data)
	})
}

func BatchGetUser(users []string) (map[string]User, error) {
	ret := map[string]User{}
	err := gDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(userBucket)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			item := User{}
			if err := decodeValue(v, &item); err != nil {
				return err
			}
			ret[item.ID] = item
			return nil
		})
	})

	return ret, err
}

func GetAllUsers() ([]User, error) {
	ret := []User{}
	err := gDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(userBucket)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			item := User{}
			if err := decodeValue(v, &item); err != nil {
				return err
			}
			ret = append(ret, item)
			return nil
		})
	})
	return ret, err
}
