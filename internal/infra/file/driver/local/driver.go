package local

import (
    "context"
    "errors"
    "io"
    "os"
    "path/filepath"
    "strings"
    
    `cms/internal/infra/file`
)

type localStorage struct {
    savePath string
}

var _ file.Driver = (*localStorage)(nil)

func (s localStorage) OpenReader(_ context.Context, path string) (io.ReadCloser, error) {
    fullPath, err := s.resolvePath(path)
    if err != nil {
        return nil, err
    }
    f, err := os.OpenFile(fullPath, os.O_RDONLY, 0o644)
    if err != nil {
        return nil, err
    }
    return f, nil
}

func (s localStorage) Create(_ context.Context, path string, r io.Reader, _ int64) error {
    fullPath, err := s.preparePath(path)
    if err != nil {
        return err
    }
    f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
    if err != nil {
        if errors.Is(err, os.ErrExist) {
            return nil
        }
        return err
    }
    _, copyErr := io.Copy(f, r)
    closeErr := f.Close()
    if copyErr != nil || closeErr != nil {
        _ = os.Remove(fullPath)
    }
    return errors.Join(copyErr, closeErr)
}

func (s localStorage) Put(_ context.Context, path string, r io.Reader, _ int64) error {
    fullPath, err := s.preparePath(path)
    if err != nil {
        return err
    }
    f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
    if err != nil {
        return err
    }
    _, copyErr := io.Copy(f, r)
    closeErr := f.Close()
    if copyErr != nil || closeErr != nil {
        _ = os.Remove(fullPath)
    }
    return errors.Join(copyErr, closeErr)
}

func (s localStorage) Delete(_ context.Context, path string) error {
    fullPath, err := s.resolvePath(path)
    if err != nil {
        return err
    }
    err = os.Remove(fullPath)
    if os.IsNotExist(err) {
        return nil
    }
    return err
}

func (s localStorage) Exists(_ context.Context, path string) (bool, error) {
    fullPath, err := s.resolvePath(path)
    if err != nil {
        return false, err
    }
    _, err = os.Stat(fullPath)
    if os.IsNotExist(err) {
        return false, nil
    }
    return err == nil, err
}

func (s localStorage) Stat(_ context.Context, path string) (file.FileInfo, error) {
    fullPath, err := s.resolvePath(path)
    if err != nil {
        return nil, err
    }
    data, err := os.Stat(fullPath)
    if err != nil {
        return nil, err
    }
    return newFileInfo(path, data), nil
}

func (s localStorage) resolvePath(targetPath string) (string, error) {
    if targetPath == "" {
        return "", errors.New("文件路径不能为空")
    }
    root, err := filepath.Abs(s.savePath)
    if err != nil {
        return "", err
    }
    fullPath, err := filepath.Abs(filepath.Join(root, filepath.ToSlash(targetPath)))
    if err != nil {
        return "", err
    }
    relativePath, err := filepath.Rel(root, fullPath)
    if err != nil {
        return "", err
    }
    if relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
        return "", errors.New("文件路径超出存储目录")
    }
    return fullPath, nil
}

func (s localStorage) preparePath(targetPath string) (string, error) {
    fullPath, err := s.resolvePath(targetPath)
    if err != nil {
        return "", err
    }
    dir := filepath.Dir(fullPath)
    err = os.MkdirAll(dir, 0o755)
    if err != nil {
        return "", err
    }
    err = os.Chmod(dir, 0o755)
    if err != nil {
        return "", err
    }
    return fullPath, nil
}
