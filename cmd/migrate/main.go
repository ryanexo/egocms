package main

import (
    `fmt`
    
    `cms/internal/config`
    `cms/internal/infra/db`
    `cms/internal/infra/persist/model`
    
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
    err = db.AutoMigrate(
        &model.Article{},
        &model.ArticleKeywords{},
        &model.ArticleContent{},
        &model.ArticleModel{},
        &model.ArticleModelJsonData{},
        &model.ArticleModelSchema{},
        &model.ArticleModelData{},
        &model.Category{},
        &model.CategorySeo{},
        &model.CategoryContext{},
        &model.Setting{},
        &model.Menu{},
        &model.MenuContext{},
        &model.User{},
        &model.Role{},
        &model.TokenBlacklist{},
        &model.File{},
        &model.Permission{},
    )
    if err != nil {
        panic(err)
    }
    fmt.Println("migrate done")
}

func initDB(c db.DBConfig) (*gorm.DB, error) {
    db, err := db.NewDB(c)
    if err != nil {
        return nil, err
    }
    return db, nil
}
