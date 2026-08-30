package setting

import (
	appsetting "cms/internal/app/setting"
	"cms/internal/infra/store/gorm/gquery"
	"cms/internal/modules/setting/internal/contract"
)

func NewSettingRepo(query *gquery.Query) contract.Repo {
	return appsetting.NewSettingRepo(query)
}
