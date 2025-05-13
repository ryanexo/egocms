package main

import (
    `dpcms/config`
    `dpcms/model`
    `dpcms/packages/database`
    "gorm.io/gen"
)

func main() {
    g := gen.NewGenerator(gen.Config{
        OutPath:       "./model/query",
        FieldNullable: true,
    })
    cfg := config.NewWithDefaultConfig()
    db, err := database.NewDB(cfg.DB)
    if err != nil {
        panic(err)
    }
    g.UseDB(db)
    
    g.ApplyBasic(
        model.User{},
        model.Category{},
        model.CategorySeo{},
        model.Menu{},
        model.Role{},
        model.RoleMenuRelation{},
        model.RoleUserRelation{},
        model.TokenBlacklist{},
    )
    
    g.ApplyInterface(
        func(table model.ClosureTable) {},
        model.CategoryContext{},
        model.MenuContext{},
    )
    
    g.Execute()
}
