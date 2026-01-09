package local

import (
    `os`
    `time`
    
    `dpcms/internal/infra/file`
)

type fileInfo struct {
    name    string
    path    string
    size    int64
    modTime time.Time
    isDir   bool
}

var _ file.FileInfo = (*fileInfo)(nil)

func (f fileInfo) Name() string {
    return f.name
}

func (f fileInfo) Path() string {
    return f.path
}

func (f fileInfo) Size() int64 {
    return f.size
}

func (f fileInfo) ModTime() time.Time {
    return f.modTime
}

func (f fileInfo) IsDir() bool {
    return f.isDir
}

func newFileInfo(data os.FileInfo) file.FileInfo {
    return fileInfo{
        name:    data.Name(),
        path:    data.Name(),
        size:    data.Size(),
        modTime: data.ModTime(),
        isDir:   data.IsDir(),
    }
}
