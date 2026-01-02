package main

import (
    `os`
    `path`
    
    `dpcms/internal/config`
    `dpcms/internal/infra/db`
    `dpcms/internal/infra/persistence/model`
    
    "gorm.io/gen"
)

func main() {
    outputPath := path.Clean("./internal/infra/persistence/query")
    err := os.RemoveAll(outputPath)
    if err != nil {
        panic(err)
    }
    
    g := gen.NewGenerator(gen.Config{
        OutPath:        outputPath,
        FieldNullable:  true,
        FieldCoverable: true,
    })
    cfg := config.NewWithBasicConfig()
    db, err := db.NewDB(cfg.DB)
    if err != nil {
        panic(err)
    }
    g.UseDB(db)
    
    g.ApplyBasic(
        model.User{},
        model.UserProfile{},
        model.Category{},
        model.CategorySeo{},
        model.Menu{},
        model.Role{},
        model.TokenBlacklist{},
        model.Article{},
        model.ArticleKeywords{},
        model.ArticleContent{},
        model.ArticleModelJsonData{},
        model.ArticleModel{},
        model.ArticleModelSchema{},
        model.ArticleModelData{},
    )
    
    g.ApplyInterface(
        func(table model.ClosureTable) {},
        model.CategoryContext{},
        model.MenuContext{},
    )
    
    g.Execute()
}
