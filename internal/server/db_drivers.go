package server

import (
    `fmt`
    `time`

    `GoBlog/internal/server/config`
    mysql2 `github.com/go-sql-driver/mysql`
    `gorm.io/driver/mysql`
    `gorm.io/driver/sqlite`
    `gorm.io/gorm`
)

type DialectConstructor func(config *config.Config) (gorm.Dialector, error)

var dials = map[string]DialectConstructor{
    "mysql":  buildMysqlDriver(),
    "sqlite": buildSqliteDriver(),
}

func buildMysqlDriver() DialectConstructor {
    return func(config *config.Config) (gorm.Dialector, error) {
        db := config.Database
        loc, err := time.LoadLocation(db.Timezone)
        if err != nil {
            return nil, err
        }

        return mysql.New(mysql.Config{
            DSNConfig: &mysql2.Config{
                User:      db.User,
                Passwd:    db.Pass,
                Addr:      db.Host,
                DBName:    db.Name,
                Loc:       loc,
                ParseTime: true,
                Params: map[string]string{
                    "charset": db.Charset,
                },
            },
        }), nil
    }
}

func buildSqliteDriver() DialectConstructor {
    return func(config *config.Config) (gorm.Dialector, error) {
        return sqlite.Open(fmt.Sprintf("%s?journal_mode=WAL", config.Database.Host)), nil
    }
}
