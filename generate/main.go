package main

import (
    `dpcms/config`
    model2 `dpcms/model`
    `dpcms/packages/data`
    "gorm.io/gen"
)

func main() {
    g := gen.NewGenerator(gen.Config{
        OutPath:       "server/query",
        FieldNullable: true,
    })
    cfg := config.NewWithDefaultConfig()
    db, err := data.NewDB(cfg.DB)
    if err != nil {
        panic(err)
    }
    g.UseDB(db)
    
    g.ApplyBasic(
        model2.User{},
        model2.Category{},
        model2.CategorySeo{},
        model2.Menu{},
        model2.Role{},
        model2.RoleMenuRelation{},
        model2.RoleUserRelation{},
    )
    
    g.ApplyInterface(
        func(table model2.ClosureTable) {},
        model2.CategoryContext{},
        model2.MenuContext{},
    )
    
    g.Execute()
}
