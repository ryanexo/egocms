package model

import (
    "unicode/utf8"
    
    "cms/internal/app/file/errno"
    "cms/internal/infra/store/modeltype"
)

const (
    OriginalNameMaxLength = 255
    ExtensionMaxLength    = 32
    PathMaxLength         = 500
    DriverMaxLength       = 32
)

type File struct {
    modeltype.Base
    OriginalName string
    Ext          string
    Path         string
    Size         uint64
    Driver       string
    SHA256       [32]byte
}

func (file *File) SetOriginalName(name string) error {
    if name == "" {
        return errno.ErrOriginalNameRequired
    }
    if utf8.RuneCountInString(name) > OriginalNameMaxLength {
        return errno.ErrOriginalNameTooLong
    }
    file.OriginalName = name
    return nil
}

func (file *File) SetExt(ext string) error {
    if len(ext) > ExtensionMaxLength {
        return errno.ErrExtensionTooLong
    }
    file.Ext = ext
    return nil
}

func (file *File) SetStorage(driver, storagePath string) error {
    if driver == "" || utf8.RuneCountInString(driver) > DriverMaxLength {
        return errno.ErrInvalidDriver
    }
    if storagePath == "" || utf8.RuneCountInString(storagePath) > PathMaxLength {
        return errno.ErrInvalidPath
    }
    file.Driver = driver
    file.Path = storagePath
    return nil
}
