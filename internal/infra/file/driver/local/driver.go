package local

import (
    `context`
    `io`
    `os`
    
    `dpcms/internal/infra/file`
)

type localDriver struct{}

var _ file.Driver = (*localDriver)(nil)

func (l localDriver) Read(_ context.Context, path string) (data []byte, err error) {
    err = operateFile(path, os.O_RDONLY, 0644, func(obj *os.File) (err error) {
        data, err = io.ReadAll(obj)
        return err
    })
    
    return
}

func (l localDriver) ReadSteam(_ context.Context, path string) (io.ReadCloser, error) {
    f, err := os.OpenFile(path, os.O_RDONLY, 0644)
    if err != nil {
        return nil, err
    }
    return f, nil
}

func (l localDriver) Write(_ context.Context, path string, data []byte) error {
    return operateFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644, func(obj *os.File) error {
        _, err := obj.Write(data)
        return err
    })
}

func (l localDriver) WriteStream(_ context.Context, path string) (io.WriteCloser, error) {
    f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return f, nil
}

func (l localDriver) Delete(_ context.Context, path string) error {
    return os.Remove(path)
}

func (l localDriver) Exists(_ context.Context, path string) (bool, error) {
    _, err := os.Stat(path)
    return err == nil, err
}

func (l localDriver) URL(_ context.Context, path string) (string, error) {
    return path, nil
}

func (l localDriver) Stat(_ context.Context, path string) (file.FileInfo, error) {
    data, err := os.Stat(path)
    if err != nil {
        return nil, err
    }
    return newFileInfo(data), nil
}

func operateFile(path string, flag int, perm os.FileMode, callback func(obj *os.File) error) (err error) {
    f, err := os.OpenFile(path, flag, perm)
    if err != nil {
        return err
    }
    defer f.Close()
    return callback(f)
}

func NewLocalDriver() file.Driver {
    return localDriver{}
}
