package cache

import (
    `fmt`
    `runtime`
    `sync`
    `time`
)

const forever int64 = 0

type value struct {
    object  any
    expires int64
}

func (s value) expired() bool {
    return s.expires != forever && time.Now().UnixNano() > s.expires
}

type segment struct {
    sync.RWMutex
    items map[string]value
}

type Cache struct {
    segments []*segment
    sgCount  int
    ttl      time.Duration
    interval time.Duration
    exit     chan bool
}

func (s *Cache) Flush() {
    for _, s := range s.segments {
        s.Lock()
        s.items = make(map[string]value)
        s.Unlock()
    }
}

func (s *Cache) segmentByKey(key string) *segment {
    var result int32
    for _, v := range key {
        result += v
    }
    return s.segments[int(result)%s.sgCount]
}

func (s *Cache) Delete(key string) {
    sg := s.segmentByKey(key)
    sg.Lock()
    delete(sg.items, key)
    sg.Unlock()
}

func (s *Cache) set(key string, data any, expires time.Duration) {
    sg := s.segmentByKey(key)
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
    s.set(key, data, time.Duration(forever))
}

func (s *Cache) Get(key string) (any, bool) {
    sg := s.segmentByKey(key)
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
    for _, v := range s.segments {
        v.Lock()
        for key, val := range v.items {
            if val.expired() {
                delete(v.items, key)
            }
        }
        v.Unlock()
    }
}

type Option interface {
    apply(*Cache)
}

type optionFunc func(*Cache)

func (f optionFunc) apply(cache *Cache) {
    f(cache)
}

func WithClearInterval(t time.Duration) Option {
    return optionFunc(func(c *Cache) {
        c.interval = t
    })
}

func WithDefaultTTL(t time.Duration) Option {
    return optionFunc(func(c *Cache) {
        c.ttl = t
    })
}

func WithSegmentCount(t int) Option {
    return optionFunc(func(c *Cache) {
        c.sgCount = t
    })
}

func makeSegments(count int) []*segment {
    segments := make([]*segment, count)
    for i := 0; i < count; i++ {
        segments[i] = &segment{items: make(map[string]value)}
    }
    return segments
}

func New(opts ...Option) *Cache {
    c := &Cache{
        ttl:      time.Hour,
        interval: time.Minute * 10,
        sgCount:  8,
        exit:     make(chan bool, 1),
    }
    for _, opt := range opts {
        opt.apply(c)
    }
    c.segments = makeSegments(c.sgCount)
    if c.interval > 0 {
        c.startMonitor()
        runtime.SetFinalizer(c, func(cache *Cache) {
            cache.stopMonitor()
        })
    }
    return c
}
