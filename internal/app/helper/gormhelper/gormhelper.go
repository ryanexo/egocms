package gormhelper

import (
    `errors`
    
    `dpcms/internal/app/erroz`
    
    `gorm.io/gorm`
)

func ReplaceNotFoundError(err error) error {
    if err == nil {
        return nil
    }
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return erroz.ErrDataNotFound.ToError()
    }
    return err
}
