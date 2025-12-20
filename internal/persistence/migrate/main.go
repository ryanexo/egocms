package main

import (
    `fmt`
    
    `dpcms/internal/config`
    `dpcms/internal/database/model`
    `dpcms/internal/infra/database`
    
    `gorm.io/gorm`
)

func main() {
    cfg := config.NewWithBasicConfig()
    db, err := initDB(cfg.DB)
    if err != nil {
        panic(err)
    }
    if cfg.DB.Type == "mysql" {
        db.Set("gorm:table_options", "ENGINE=InnoDB")
    }
    err = db.AutoMigrate(
        &model.Category{},
        &model.CategorySeo{},
        &model.CategoryContext{},
        &model.Menu{},
        &model.MenuContext{},
        &model.User{},
        &model.Role{},
        &model.TokenBlacklist{},
    )
    if err != nil {
        panic(err)
    }
    fmt.Println("migrate done")
}

func initDB(c database.DBConfig) (*gorm.DB, error) {
    db, err := database.NewDB(c)
    if err != nil {
        return nil, err
    }
    return db, nil
}
