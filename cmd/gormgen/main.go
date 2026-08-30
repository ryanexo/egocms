package main

import (
	"os"
	"path"

	categoryModel "cms/internal/app/category/model"
	fileModel "cms/internal/app/file/model"
	roleModel "cms/internal/app/role/model"
	settingModel "cms/internal/app/setting/model"
	userModel "cms/internal/app/user/model"
	"cms/internal/config"
	"cms/internal/infra/db"
	"cms/internal/infra/store/modeltype"
	model2 "cms/internal/public/model"

	"gorm.io/gen"
)

func main() {
	outputPath := path.Clean("./internal/infra/store/gorm/gquery")
	g := gen.NewGenerator(gen.Config{
		OutPath:        outputPath,
		FieldNullable:  true,
		FieldCoverable: true,
	})
	cfg, err := config.NewWithBasicConfig()
	if err != nil {
		panic(err)
	}
	gormDB, sqlDB, err := db.NewDB(cfg.DB)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	g.UseDB(gormDB)
	if err = os.RemoveAll(outputPath); err != nil {
		panic(err)
	}

	g.ApplyBasic(
		model2.Article{},
		model2.ArticleVersion{},
		model2.ArticlePublish{},
		model2.ArticleCategoryRel{},
		model2.ArticleTag{},
		model2.ArticleTagRel{},
		model2.ArticleComment{},
		modeltype.Base{},
		categoryModel.Category{},
		categoryModel.CategoryMeta{},
		categoryModel.CategoryContext{},
		settingModel.Setting{},
		model2.ContentType{},
		model2.ContentTypeEntries{},
		model2.ContentTypeSchema{},
		model2.ContentFieldValues{},
		fileModel.File{},
		fileModel.Attachment{},
		model2.Menu{},
		model2.MenuContext{},
		model2.Permission{},
		roleModel.Role{},
		model2.SinglePage{},
		userModel.User{},
		userModel.UserProfile{},
	)

	g.Execute()
}
