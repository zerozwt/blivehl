package db

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type Cookie struct {
	UserID   string `json:"id"`
	ExpireAt int64  `json:"expire"`
}

var ErrCookieNotFound error = fmt.Errorf("cookie not found")
var ErrCookieExpired error = fmt.Errorf("cookie expired")
var cookieBucket []byte = []byte("cookie")

func GetCookie(cookie string) (string, error) {
	ret := ""
	err := gDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(cookieBucket)
		if bucket == nil {
			return ErrCookieNotFound
		}

		value := bucket.Get([]byte(cookie))
		if value == nil {
			return ErrCookieNotFound
		}

		tmp := Cookie{}
		if err := decodeValue(value, &tmp); err != nil {
			return err
		}

		if tmp.ExpireAt < time.Now().Unix() {
			go DeleteCookie(cookie)
			return ErrCookieExpired
		}

		ret = tmp.UserID
		return nil
	})
	return ret, err
}

func DeleteCookie(cookie string) error {
	return gDB.Update(func(tx *bbolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists(cookieBucket)
		bucket.Delete([]byte(cookie))
		return nil
	})
}

func PutCookie(cookie, userID string, expire int64) error {
	return gDB.Update(func(tx *bbolt.Tx) error {
		data := encodeValue(Cookie{UserID: userID, ExpireAt: expire})
		bucket, _ := tx.CreateBucketIfNotExists(cookieBucket)
		bucket.Put([]byte(cookie), data)
		return nil
	})
}
