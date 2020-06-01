package tokencacher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

type File struct {
	path string

	lock sync.Mutex
}

type cache struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewFile(path string) *File {
	if len(path) == 0 {
		panic("path required")
	}

	return &File{
		path: path,
	}
}

func (f *File) Get() (string, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	contents, err := ioutil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	if len(contents) == 0 {
		return "", ErrTokenNotExist
	}

	c := &cache{}
	err = json.Unmarshal(contents, c)
	if err != nil {
		return "", fmt.Errorf("Error unmarshaling token: %s", err)
	}

	if time.Now().After(c.ExpiresAt) {
		return "", ErrTokenExpired
	}

	return c.Token, nil
}

func (f *File) Set(token string, expiresAt time.Time) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	c := &cache{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.path, b, 0600)
	if err != nil {
		return err
	}

	return nil
}
