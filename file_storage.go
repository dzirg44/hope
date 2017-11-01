package main

import (
	"io/ioutil"
	"os"
)
type FileStorage struct {
	Filename string
}

func (f  FileStorage) ReadStorage() ([]byte, error) {
	contents, err := ioutil.ReadFile(f.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return contents, err
		}
		return contents, err
	}

  return contents, err
}

func (f FileStorage) WriteStorage(contents []byte, perm os.FileMode) error {
	err := ioutil.WriteFile(f.Filename, contents, perm)
	return err
}

/*
func  JsonEncoder() {
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, nil
}
func JsonDecoder() {

}
*/