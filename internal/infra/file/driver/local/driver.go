package local

import (
    `context`
    `io`
    `os`
    `path`
    
    `dpcms/internal/infra/file`
)

type localStorage struct {
    savePath string
}

var _ file.Driver = (*localStorage)(nil)

func (s localStorage) Name() string {
    return "local"
}

func (s localStorage) Read(_ context.Context, path string) (data []byte, err error) {
    err = s.operateFile(path, os.O_RDONLY, 0o644, func(obj *os.File) (err error) {
        data, err = io.ReadAll(obj)
        return err
    })
    
    return
}

func (s localStorage) OpenReader(_ context.Context, path string) (io.ReadCloser, error) {
    fullPath, err := s.getFullPath(path)
    if err != nil {
        return nil, err
    }
    f, err := os.OpenFile(fullPath, os.O_RDONLY, 0o644)
    if err != nil {
        return nil, err
    }
    return f, nil
}

func (s localStorage) Write(_ context.Context, path string, data []byte) error {
    return s.operateFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644, func(obj *os.File) error {
        _, err := obj.Write(data)
        return err
    })
}

func (s localStorage) OpenWriter(_ context.Context, path string) (io.WriteCloser, error) {
    fullPath, err := s.getFullPath(path)
    if err != nil {
        return nil, err
    }
    f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return f, nil
}

func (s localStorage) Delete(_ context.Context, path string) error {
    return os.Remove(path)
}

func (s localStorage) Exists(_ context.Context, path string) (bool, error) {
    _, err := os.Stat(path)
    return err == nil, err
}

func (s localStorage) URL(_ context.Context, path string) (string, error) {
    return path, nil
}

func (s localStorage) Stat(_ context.Context, path string) (file.FileInfo, error) {
    data, err := os.Stat(path)
    if err != nil {
        return nil, err
    }
    return newFileInfo(data), nil
}

func (s localStorage) getFullPath(targetPath string) (string, error) {
    fullPath := path.Join(s.savePath, "./", targetPath)
    dir := path.Dir(fullPath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return "", err
    }
    err = os.Chmod(dir, 0755)
    if err != nil {
        return "", err
    }
    return fullPath, nil
}

func (s localStorage) operateFile(path string, flag int, perm os.FileMode, callback func(obj *os.File) error) error {
    fullPath, err := s.getFullPath(path)
    if err != nil {
        return err
    }
    f, err := os.OpenFile(fullPath, flag, perm)
    if err != nil {
        return err
    }
    defer f.Close()
    return callback(f)
}
