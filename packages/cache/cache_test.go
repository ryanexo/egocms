package cache

import (
    "runtime"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
)

var (
    instance = New(&Config{
        TTL:          time.Millisecond * 100,
        ScanInterval: time.Millisecond * 300,
        Partition:    4,
    })
)

func TestContainer_Set_Get(t *testing.T) {
    instance.Set("key1", "test1")
    val, ok := instance.Get("key1")
    assert.True(t, ok)
    assert.Equal(t, "test1", val)
    time.Sleep(time.Millisecond * 100)
    val, ok = instance.Get("key1")
    assert.False(t, ok)
    assert.Equal(t, nil, val)
    
    instance.SetWithTTL("key2", "test2", time.Millisecond*10)
    val, ok = instance.Get("key2")
    assert.True(t, ok)
    assert.Equal(t, "test2", val)
    time.Sleep(time.Millisecond * 300)
    val, ok = instance.Get("key2")
    assert.False(t, ok)
    assert.Equal(t, nil, val)
}

func TestContainer_Flush_Delete(t *testing.T) {
    instance.Set("k", "test")
    instance.Flush()
    n := 0
    for i := 0; i < instance.partition; i++ {
        n += len(instance.shards[i].items)
    }
    assert.Equal(t, 0, n)
    instance.Forever("k", "test")
    instance.Delete("k")
    _, ok := instance.Get("k")
    assert.False(t, false, ok)
}

func TestContainer_stopMonitor(t *testing.T) {
    baseGoroutineNum := runtime.NumGoroutine()
    key := "k"
    instance.SetWithTTL(key, "test", time.Millisecond*100)
    instance.stopMonitor()
    time.Sleep(time.Second)
    sg := instance.segmentByKey(key)
    v := sg.items[key]
    assert.Equal(t, baseGoroutineNum-1, runtime.NumGoroutine())
    assert.Equal(t, "test", v.object)
    
    instance.startMonitor()
    instance = nil
    sg = nil
    runtime.GC()
    runtime.Gosched()
    time.Sleep(time.Second)
    assert.Equal(t, baseGoroutineNum-1, runtime.NumGoroutine())
}

func Test_Inc_Dec(t *testing.T) {
    instance.SetWithTTL("k", uint8(15), time.Hour)
    newVal, err := Inc[uint8](instance, "k", 100)
    if err != nil {
        t.Error(err)
    }
    assert.Equal(t, newVal, uint8(115))
    newVal, err = Dec[uint8](instance, "k", 100)
    assert.Equal(t, newVal, uint8(15))
}
