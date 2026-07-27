package main

import (
    "os"
    "path"
    
    "cms/internal/config"
    "cms/internal/infra/db"
    "cms/internal/infra/persistence/gorm/model"
    
    "gorm.io/gen"
)

func main() {
    outputPath := path.Clean("./internal/infra/persistence/gorm/gquery")
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
        model.Article{},
        model.ArticleVersion{},
        model.ArticlePublish{},
        model.ArticleCategoryRel{},
        model.ArticleTag{},
        model.ArticleTagRel{},
        model.ArticleComment{},
        model.Base{},
        model.Category{},
        model.CategorySeo{},
        model.CategoryContext{},
        model.Config{},
        model.ContentType{},
        model.ContentTypeEntries{},
        model.ContentTypeSchema{},
        model.ContentFieldValues{},
        model.File{},
        model.Menu{},
        model.MenuContext{},
        model.Permission{},
        model.Role{},
        model.SinglePage{},
        model.User{},
        model.UserProfile{},
    )
    
    g.Execute()
}
