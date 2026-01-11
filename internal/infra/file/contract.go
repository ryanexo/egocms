package file

import (
    `context`
    `io`
    `time`
)

type Driver interface {
    Read(ctx context.Context, path string) ([]byte, error)
    ReadSteam(ctx context.Context, path string) (io.ReadCloser, error)
    
    Write(ctx context.Context, path string, data []byte) error
    WriteStream(ctx context.Context, path string) (io.WriteCloser, error)
    
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
