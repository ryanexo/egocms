package file

import (
    `context`
    `io`
    `time`
)

type Driver interface {
    Name() string
    
    Read(ctx context.Context, path string) ([]byte, error)
    OpenReader(ctx context.Context, path string) (io.ReadCloser, error)
    
    // Write
    // 文件已存在时返回fs.ErrExist
    Write(ctx context.Context, path string, data []byte) error
    // OpenWriter
    // 文件已存在时返回fs.ErrExist
    OpenWriter(ctx context.Context, path string) (io.WriteCloser, error)
    
    Delete(ctx context.Context, path string) error
    
    Exists(ctx context.Context, path string) (bool, error)
    URL(ctx context.Context, path string) (string, error)
    Stat(ctx context.Context, path string) (FileInfo, error)
}

type FileInfo interface {
    Name() string
    Path() string
    Size() int64
    ModTime() time.Time
    IsDir() bool
}

type DriverFactory interface {
    Setup(config map[string]any) (Driver, error)
}
