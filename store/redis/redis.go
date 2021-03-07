package redis

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/go-redis/cache/v8"
    "github.com/go-redis/redis/v8"

    "github.com/dbyington/manifest-server/store"
)

const defaultTTL = time.Hour * 3000 // These things don't tend to change regularly.

type Store struct {
	*cache.Cache
}

func New() *Store {
    urlOpt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
    if err != nil {
        panic(err)
    }
	r := redis.NewClient(urlOpt)

	log.Printf("connected to redis at: %s", os.Getenv("REDIS_URL"))
	return &Store{
		cache.New(&cache.Options{
			Redis:      r,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		}),
	}
}

func (s *Store) Put(ctx context.Context, key string, o *store.Object) error {
    if err := s.Set(&cache.Item{
        Ctx:   ctx,
        Key:   key,
        Value: o,
        TTL:   defaultTTL,
    }); err != nil {
        return err
    }

    return nil
}

func (s Store) Get(ctx context.Context, key string) (*store.Object, error) {
    var o *store.Object
    if err := s.Cache.Get(ctx, key, &o); err != nil {
        return nil, err
    }
    return o, nil
}
