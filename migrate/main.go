package main

import (
    `dpcms/config`
    model2 `dpcms/model`
)

func main() {
    cfg := config.NewWithDefaultConfig()
    db, err := initDB(cfg.DB)
    if err != nil {
        panic(err)
    }
    if cfg.DB.Type == "mysql" {
        db.Set("gorm:table_options", "ENGINE=InnoDB")
    }
    err = db.AutoMigrate(
        &model2.Category{},
        &model2.CategorySeo{},
        &model2.CategoryContext{},
        &model2.Menu{},
        &model2.MenuContext{},
        &model2.User{},
        &model2.Role{},
        &model2.RoleMenuRelation{},
        &model2.RoleUserRelation{},
    )
    if err != nil {
        panic(err)
    }
}
