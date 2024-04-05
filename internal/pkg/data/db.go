package data

import (
    `runtime`
    `time`

    `github.com/DATA-DOG/go-sqlmock`
    `gorm.io/driver/mysql`
    `gorm.io/gorm`
    `gorm.io/gorm/logger`
    `gorm.io/gorm/schema`
)

type DBOption func(*gorm.Config)

func DefaultConfig(prefix string) DBOption {
    return func(config *gorm.Config) {
        *config = gorm.Config{
            NamingStrategy: schema.NamingStrategy{
                TablePrefix:   prefix,
                SingularTable: true,
            },
            FullSaveAssociations:   false,
            Logger:                 logger.Default.LogMode(logger.Silent),
            PrepareStmt:            true,
            AllowGlobalUpdate:      false,
            QueryFields:            true,
            SkipDefaultTransaction: true,

            DisableForeignKeyConstraintWhenMigrating: true,
        }
    }
}

func resolveOptions(options ...DBOption) *gorm.Config {
    config := new(gorm.Config)
    for _, apply := range options {
        apply(config)
    }
    return config
}

func NewDB(dialector gorm.Dialector, options ...DBOption) (*gorm.DB, error) {
    config := resolveOptions(options...)
    db, err := gorm.Open(dialector, config)
    if err != nil {
        return nil, err
    }
    sqlDB, sqlDBErr := db.DB()
    if sqlDBErr != nil {
        return nil, err
    }

    sqlDB.SetConnMaxLifetime(time.Second * 30)
    sqlDB.SetConnMaxIdleTime(time.Minute)
    sqlDB.SetMaxOpenConns(runtime.NumCPU() * 2)

    return db, nil
}

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
    sqlDB, sqlm, err := sqlmock.New()
    if err != nil {
        panic(err)
    }
    sqlm.ExpectQuery(`SELECT VERSION\(\)`).WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.0.28-0ubuntu0.20.04.3"))
    db, err := NewDB(mysql.New(mysql.Config{Conn: sqlDB}), DefaultConfig(""))
    if err != nil {
        panic(err)
    }
    return db, sqlm
}
