package main

import (
    "os"
    "path"
    
    "cms/internal/config"
    "cms/internal/infra/db"
    "cms/internal/infra/persist/model"
    
    "gorm.io/gen"
)

func main() {
    outputPath := path.Clean("./internal/infra/persist/query")
    err := os.RemoveAll(outputPath)
    if err != nil {
        panic(err)
    }
    
    g := gen.NewGenerator(gen.Config{
        OutPath:        outputPath,
        FieldNullable:  true,
        FieldCoverable: true,
    })
    cfg, err := config.NewWithBasicConfig()
    if err != nil {
        panic(err)
    }
    gormDB, err := db.NewDB(cfg.DB)
    if err != nil {
        panic(err)
    }
    g.UseDB(gormDB)
    
    g.ApplyBasic(
        model.User{},
        model.UserProfile{},
        model.Category{},
        model.CategorySeo{},
        model.Menu{},
        model.Permission{},
        model.Role{},
        model.TokenBlacklist{},
        model.Article{},
        model.ArticleKeywords{},
        model.ArticleContent{},
        model.ArticleModelJsonData{},
        model.ArticleModel{},
        model.ArticleModelSchema{},
        model.ArticleModelData{},
        model.File{},
        model.Setting{},
    )
    
    g.ApplyInterface(
        func(table model.ClosureTable) {},
        model.CategoryContext{},
        model.MenuContext{},
    )
    
    g.Execute()
}
