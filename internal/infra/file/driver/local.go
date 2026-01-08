package driver

import (
    `fmt`
    `io`
    `os`
    `path`
    
    `dpcms/internal/infra/file`
)

type Local struct {
    Output string
}

var _ file.Driver = (*Local)(nil)

func (s Local) Read(path string) ([]byte, error) {
    return os.ReadFile(path)
}

func (s Local) ReadStream(path string) (io.ReadCloser, error) {
    return os.Open(path)
}

func (s Local) Save(path string, data []byte) (err error) {
    f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    
    defer func(file *os.File) {
        closeErr := file.Close()
        if closeErr != nil {
            if err == nil {
                err = closeErr
            } else {
                err = fmt.Errorf("%v; file close error: %w", err, closeErr)
            }
        }
    }(f)
    
    _, err = f.Write(data)
    if err != nil {
        return err
    }
    
    return nil
}

func (s Local) SaveStream(path string, data io.Reader) error {
    // TODO implement me
    panic("implement me")
}

func (s Local) Delete(path string) error {
    // TODO implement me
    panic("implement me")
}

func (s Local) Exists(path string) (bool, error) {
    // TODO implement me
    panic("implement me")
}

func (s Local) URL() (string, error) {
    // TODO implement me
    panic("implement me")
}

func (s Local) List(dir string) ([]string, error) {
    // TODO implement me
    panic("implement me")
}

func (s Local) Stat(path string) (file.FileInfo, error) {
    // TODO implement me
    panic("implement me")
}

func NewLocal(output string) file.Driver {
    return Local{Output: path.Clean(output)}
}
