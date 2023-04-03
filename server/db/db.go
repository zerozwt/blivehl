package db

import "go.etcd.io/bbolt"

var gDB *bbolt.DB

func InitDB(file string) error {
	var err error
	gDB, err = bbolt.Open(file, 0644, nil)
	return err
}
