package db

import (
    "runtime"
    "time"
    
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
)

type DBConfig struct {
    Type    string `json:"type" yaml:"type"`
    Host    string `json:"host" yaml:"host"`
    Name    string `json:"name" yaml:"name"`
    User    string `json:"user" yaml:"user"`
    Pass    string `json:"pass" yaml:"pass"`
    Charset string `json:"charset" yaml:"charset"`
}

func createBasicConfig() *gorm.Config {
    return &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
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

func buildDB(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
    db, err := gorm.Open(dialector, config)
    if err != nil {
        return nil, err
    }
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetConnMaxLifetime(time.Second * 30)
    sqlDB.SetConnMaxIdleTime(time.Minute)
    sqlDB.SetMaxOpenConns(runtime.NumCPU() * 2)
    
    return db, nil
}

func NewDB(config DBConfig) (*gorm.DB, error) {
    driver, err := GetDriver(config)
    if err != nil {
        return nil, err
    }
    return buildDB(driver, createBasicConfig())
}
