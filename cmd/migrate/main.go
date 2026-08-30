package main

import (
    `fmt`
    
    `cms/internal/config`
    `cms/internal/infra/db`
    
    `gorm.io/gorm`
)

func main() {
    cfg, err := config.NewWithBasicConfig()
    if err != nil {
        panic(err)
    }
    db, err := initDB(cfg.DB)
    if err != nil {
        panic(err)
    }
    if cfg.DB.Type == "mysql" {
        db.Set("gorm:table_options", "ENGINE=InnoDB")
    }
    err = db.AutoMigrate()
    if err != nil {
        panic(err)
    }
    fmt.Println("migrate done")
}

func initDB(c db.DBConfig) (*gorm.DB, error) {
    db, _, err := db.NewDB(c)
    if err != nil {
        return nil, err
    }
    return db, nil
}
