package store

import (
	"context"

	"github.com/go-redis/cache/v8"
)

var ErrCacheMiss = cache.ErrCacheMiss

type Object struct {
	ContentLength int64  `json:"-"`
	Etag          string `json:"-"`
	HashType      string `json:"hashType"`
	Id            string `json:"id"`
	Json          []byte `json:"asEncodedJson"`
	Plist         string `json:"asEncodedPlist"`
	Title         string `json:"title"`
}

type Store interface {
	Put(context.Context, string, *Object) error
	Get(context.Context, string) (*Object, error)
}
