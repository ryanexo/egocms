package file

import (
    `io`
    `time`
)

type Driver interface {
    Read(path string) ([]byte, error)
    ReadStream(path string) (io.ReadCloser, error)
    Save(path string, data []byte) error
    SaveStream(path string, data io.Reader) error
    Delete(path string) error
    Exists(path string) (bool, error)
    URL() (string, error)
    List(dir string) ([]string, error)
    Stat(path string) (FileInfo, error)
}

type FileInfo struct {
    OriginName string
    Name       string
    Path       string
    Size       uint64
    CreatedAt  time.Time
    ModifiedAt time.Time
    IsDir      bool
}
