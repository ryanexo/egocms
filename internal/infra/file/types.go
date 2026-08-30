package file

import (
    "context"
    `encoding/json`
    "io"
    "time"
)

type Driver interface {
    OpenReader(ctx context.Context, path string) (io.ReadCloser, error)
    
    // Create 目标文件必须不存在，否则返回 fs.ErrExist
    Create(ctx context.Context, path string, r io.Reader, size int64) error
    
    // Put 创建或覆盖目标文件
    Put(ctx context.Context, path string, r io.Reader, size int64) error
    
    // Delete 具有幂等语义：目标不存在时返回 nil
    Delete(ctx context.Context, path string) error
    
    Stat(ctx context.Context, path string) (FileInfo, error)
}

type FileInfo interface {
    Name() string
    Path() string
    Size() int64
    ModTime() time.Time
}

type DriverFactory interface {
    Name() string
    New(config json.RawMessage) (Driver, error)
}
