package cache

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

const noExpiry int64 = 0

type value struct {
    object  any
    expires int64
}

func (s value) expired() bool {
    return s.expires != noExpiry && time.Now().UnixNano() > s.expires
}

type Config struct {
    TTL          time.Duration `json:"ttl,omitempty" yaml:"ttl,omitempty"`
    ScanInterval time.Duration `json:"scanInterval,omitempty" yaml:"scanInterval,omitempty"`
    Partition    int           `json:"partition,omitempty" yaml:"partition,omitempty"`
}

type shard struct {
    sync.RWMutex
    items map[string]value
}

type Cache struct {
    shards    []*shard
    partition int
    ttl       time.Duration
    interval  time.Duration
    exit      chan bool
}

func (s *Cache) Flush() {
    for _, s := range s.shards {
        s.Lock()
        s.items = make(map[string]value)
        s.Unlock()
    }
}

func (s *Cache) getSegmentByKey(key string) *shard {
    var result int32
    for _, v := range key {
        result += v
    }
    return s.shards[int(result)%s.partition]
}

func (s *Cache) Delete(key string) {
    sg := s.getSegmentByKey(key)
    sg.Lock()
    delete(sg.items, key)
    sg.Unlock()
}

func (s *Cache) set(key string, data any, expires time.Duration) {
    sg := s.getSegmentByKey(key)
    sg.Lock()
    defer sg.Unlock()
    val := value{object: data}
    if expires == 0 {
        val.expires = 0
    } else {
        val.expires = time.Now().Add(expires).UnixNano()
    }
    sg.items[key] = val
}

func (s *Cache) Set(key string, data any) {
    s.set(key, data, s.ttl)
}

func (s *Cache) SetWithTTL(key string, data any, ttl time.Duration) {
    s.set(key, data, ttl)
}

func (s *Cache) Forever(key string, data any) {
    s.set(key, data, time.Duration(noExpiry))
}

func (s *Cache) Get(key string) (any, bool) {
    sg := s.getSegmentByKey(key)
    sg.RLock()
    canUnlock := true
    defer func() {
        if canUnlock {
            defer sg.RUnlock()
        }
    }()
    result, ok := sg.items[key]
    if !ok {
        return nil, false
    }
    if result.expired() {
        canUnlock = false
        sg.RUnlock()
        s.Delete(key)
        return nil, false
    }
    return result.object, true
}

func (s *Cache) startMonitor() {
    go func() {
        tick := time.NewTicker(s.interval)
        select {
        case <-s.exit:
            tick.Stop()
            fmt.Println("stop")
            return
        
        case <-tick.C:
            s.ClearExpires()
        }
    }()
}

func (s *Cache) stopMonitor() {
    s.exit <- true
}

func (s *Cache) ClearExpires() {
    for _, v := range s.shards {
        v.Lock()
        for key, val := range v.items {
            if val.expired() {
                delete(v.items, key)
            }
        }
        v.Unlock()
    }
}

func makeSegments(count int) []*shard {
    segments := make([]*shard, count)
    for i := 0; i < count; i++ {
        segments[i] = &shard{items: make(map[string]value)}
    }
    return segments
}

func New(cfg *Config) *Cache {
    c := &Cache{
        ttl:       cfg.TTL,
        interval:  cfg.ScanInterval,
        partition: cfg.Partition,
        exit:      make(chan bool, 1),
    }
    
    c.shards = makeSegments(c.partition)
    if c.interval > 0 {
        c.startMonitor()
        runtime.SetFinalizer(c, func(cache *Cache) {
            cache.stopMonitor()
        })
    }
    return c
}
