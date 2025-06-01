package dbscopes

import (
    `gorm.io/gen`
)

type Pagination struct {
    PageNo   int `json:"pageNo"`
    PageSize int `json:"pageSize"`
}

func Paginate(pageNo, pageSize int) func(db gen.Dao) gen.Dao {
    return func(db gen.Dao) gen.Dao {
        if pageNo < 1 {
            pageNo = 1
        }
        if pageSize < 1 {
            pageSize = 20
        } else if pageSize > 1000 {
            pageSize = 1000
        }
        return db.Offset((pageNo - 1) * pageSize).Limit(pageSize)
    }
}
