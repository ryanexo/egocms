package cache

import (
    `github.com/dgraph-io/ristretto/v2`
)

type ConfigCache struct {
    *ristretto.Cache[string, string]
}

func NewConfigCache() (*ConfigCache, error) {
    cache, err := ristretto.NewCache[string, string](&ristretto.Config[string, string]{
        BufferItems: 64,
        NumCounters: 1000,
        MaxCost:     5 << 20,
        Cost: func(value string) int64 {
            return int64(len(value))
        },
        TtlTickerDurationInSec: 60 * 10,
    })
    if err != nil {
        return nil, err
    }
    return &ConfigCache{cache}, nil
}
