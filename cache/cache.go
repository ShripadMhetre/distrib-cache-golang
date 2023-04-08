package cache

import "time"

type Cache interface {
	Set([]byte, []byte, time.Duration) error
	Get([]byte) (string, error)
	Delete([]byte) error
	Exists([]byte) bool
}
