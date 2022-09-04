package constants

import (
	"errors"
	"fmt"
)

type ErrType uint

const (
	DNE ErrType = iota
	KEY_EXIST
	EMPTY_KV
)

var ErrDne error = errors.New("key does not exist in data store")
var ErrKeyExist error = errors.New("value for key already exists")
var ErrEmptyKv error = errors.New("nothing to clear, datastore is empty")

func ErrDnef(key string) error {
	err := fmt.Errorf("key %s does not exist in data store", key)
	return err
}

func ErrKeyExistf(key string) error {
	err := fmt.Errorf("value for key %s already exists\nuse SimpleStore.Update to update values", key)
	return err
}

func ErrHandle(err ErrType) error {
	switch err {
	case DNE:
		return ErrDne
	case KEY_EXIST:
		return ErrKeyExist
	case EMPTY_KV:
		return ErrEmptyKv
	}
	return nil
}

func ErrHandlef(err ErrType, key string) error {
	switch err {
	case DNE:
		return ErrDnef(key)
	case KEY_EXIST:
		return ErrKeyExistf(key)
	}
	return nil
}
