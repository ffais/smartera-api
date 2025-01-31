package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileStore struct {
	storePath string
}

func NewFileStore(dirName string) *FileStore {
	err := os.Mkdir(dirName, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	return &FileStore{dirName}
}

func (f FileStore) Add(name string, user User) error {
	file, _ := json.MarshalIndent(user, "", " ")
	err := os.WriteFile(filepath.Join(f.storePath, (fmt.Sprintf("%s%s", name, ".json"))), file, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (f FileStore) Get(name string) (User, error) {
	var user User
	file, err := os.ReadFile(filepath.Join(f.storePath, (fmt.Sprintf("%s%s", name, ".json"))))
	if err != nil {
		return User{}, ErrNotFound
	}
	if err = json.Unmarshal(file, &user); err != nil {

		return User{}, ErrNotFound
	}
	return user, nil
}

func (f FileStore) Update(name string, user User) error {
	err := f.Add(name, user)
	if err != nil {
		return err
	}
	return nil
}
