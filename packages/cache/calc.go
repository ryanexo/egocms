package cache

import (
    `errors`
)

var (
    ErrType      = errors.New("被操作缓存数据类型错误")
    ErrNotExists = errors.New("缓存key不存在或已过期")
)

func Inc[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](cache *Cache, key string, delta T) (T, error) {
    sg := cache.segmentByKey(key)
    sg.Lock()
    defer sg.Unlock()
    val, ok := sg.items[key]
    if !ok {
        return 0, ErrNotExists
    }
    if val.expired() {
        delete(sg.items, key)
        return 0, ErrNotExists
    }
    oldVal, ok := val.object.(T)
    if !ok {
        return 0, ErrType
    }
    newVal := oldVal + delta
    val.object = newVal
    sg.items[key] = val
    return newVal, nil
}

func Dec[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](cache *Cache, key string, delta T) (T, error) {
    sg := cache.segmentByKey(key)
    sg.Lock()
    defer sg.Unlock()
    val, ok := sg.items[key]
    if !ok {
        return 0, ErrNotExists
    }
    if val.expired() {
        delete(sg.items, key)
        return 0, ErrNotExists
    }
    oldVal, ok := val.object.(T)
    if !ok {
        return 0, ErrType
    }
    newVal := oldVal - delta
    val.object = newVal
    sg.items[key] = val
    return newVal, nil
}
