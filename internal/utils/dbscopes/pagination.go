package dbscopes

import (
    `gorm.io/gen`
)

func Paginate(pageNo, pageSize int) func(db gen.Dao) gen.Dao {
    return func(db gen.Dao) gen.Dao {
        if pageNo < 1 {
            pageNo = 1
        }
        if pageSize < 1 {
            pageSize = 20
        }
        return db.Offset((pageNo - 1) * pageSize).Limit(pageSize)
    }
}
