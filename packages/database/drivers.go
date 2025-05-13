package database

import (
    "fmt"
    
    mysql2 "github.com/go-sql-driver/mysql"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var drivers = map[string]func(c *DBConfig) (gorm.Dialector, error){
    "mysql":  mysqlDriverBuilder,
    "sqlite": sqliteDriverBuilder,
}

func GetDriver(dbConfig *DBConfig) (gorm.Dialector, error) {
    builder, ok := drivers[dbConfig.Type]
    if !ok {
        return nil, fmt.Errorf("数据库驱动 %s 不存在", dbConfig.Type)
    }
    return builder(dbConfig)
}

func mysqlDriverBuilder(dbConfig *DBConfig) (gorm.Dialector, error) {
    return mysql.New(mysql.Config{
        DSNConfig: &mysql2.Config{
            User:      dbConfig.User,
            Passwd:    dbConfig.Pass,
            Addr:      dbConfig.Host,
            DBName:    dbConfig.Name,
            ParseTime: true,
            Params: map[string]string{
                "charset": "utf8mb4",
            },
        },
    }), nil
}

func sqliteDriverBuilder(dbConfig *DBConfig) (gorm.Dialector, error) {
    return sqlite.Open(fmt.Sprintf("%s?journal_mode=WAL", dbConfig.Host)), nil
}
