package main

import (
    `dpcms/internal/config`
    `dpcms/internal/database/model`
    `dpcms/internal/packages/database`
    
    "gorm.io/gen"
)

func main() {
    g := gen.NewGenerator(gen.Config{
        OutPath:        "./internal/database/query",
        FieldNullable:  true,
        FieldCoverable: true,
    })
    cfg := config.NewWithBasicConfig()
    db, err := database.NewDB(cfg.DB)
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
    )
    
    g.ApplyInterface(
        func(table model.ClosureTable) {},
        model.CategoryContext{},
        model.MenuContext{},
    )
    
    g.Execute()
}
