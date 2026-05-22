package assembler

import (
	"cms/internal/pkg/datatype"
	"database/sql"

	"cms/internal/app/articlemodel/internal/dto"
	"cms/internal/infra/persistence/model"
	"cms/internal/util/types"

	"github.com/shopspring/decimal"
)

func ToArticleModelCreateCommand(data *dto.ArticleModelCreateParams) *model.ArticleModel {
	return &model.ArticleModel{
		Name:        data.Name,
		Description: data.Description,
	}
}

func ToArticleModelUpdateCommand(data *dto.ArticleModelUpdateParams) *model.ArticleModel {
	return &model.ArticleModel{
		Base: model.Base{
			ID: data.ID,
		},
		Name:        data.Name,
		Description: data.Description,
	}
}

func ToArticleModelDTO(data *model.ArticleModel) *dto.ArticleModel {
	return &dto.ArticleModel{
		Base: types.Base{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
		},
		Name:        data.Name,
		Description: data.Description,
	}
}

func ToArticleModelListDTO(data []*model.ArticleModel) []*dto.ArticleModel {
	result := make([]*dto.ArticleModel, 0, len(data))
	for _, item := range data {
		result = append(result, ToArticleModelDTO(item))
	}
	return result
}

func ToArticleModelSchemaDTO(data *model.ArticleModelSchema) *dto.ArticleModelSchemaParams {
	return &dto.ArticleModelSchemaParams{
		ID:          &data.ID,
		FieldKey:    data.FieldKey,
		FieldName:   data.FieldName,
		Description: data.Description,
		Sequence:    data.Sequence,
		Type:        data.Type,
		Required:    data.Required,
		MinLen:      data.MinLen,
		MaxLen:      data.MaxLen,
		MinValue:    datatype.SafeInt64(data.MinValue.Decimal.BigInt().Int64()),
		MaxValue:    datatype.SafeInt64(data.MaxValue.Decimal.BigInt().Int64()),
		MinTime:     data.MinTime.Time,
		MaxTime:     data.MaxTime.Time,
		Pattern:     data.Pattern,
		EnumOptions: data.EnumOptions,
		Hidden:      data.Hidden,
		Enable:      data.Enable,
	}
}

func ToArticleModelSchemaList(data []*model.ArticleModelSchema) []*dto.ArticleModelSchemaParams {
	result := make([]*dto.ArticleModelSchemaParams, 0, len(data))
	for _, item := range data {
		result = append(result, ToArticleModelSchemaDTO(item))
	}
	return result
}

func ToArticleModelSchemaModel(data *dto.ArticleModelSchemaParams) *model.ArticleModelSchema {
	return &model.ArticleModelSchema{
		FieldKey:  data.FieldKey,
		FieldName: data.FieldName,
		MinLen:    data.MinLen,
		MaxLen:    data.MaxLen,
		MinValue: decimal.NullDecimal{
			Decimal: decimal.NewFromInt(data.MinValue.Raw()),
			Valid:   true,
		},
		MaxValue: decimal.NullDecimal{
			Decimal: decimal.NewFromInt(data.MaxValue.Raw()),
			Valid:   true,
		},
		MinTime: sql.NullTime{
			Time:  data.MinTime,
			Valid: true,
		},
		MaxTime: sql.NullTime{
			Time:  data.MaxTime,
			Valid: true,
		},
		Pattern:     data.Pattern,
		Sequence:    data.Sequence,
		Type:        data.Type,
		EnumOptions: data.EnumOptions,
		Required:    data.Required,
		Hidden:      data.Hidden,
		Enable:      data.Enable,
	}
}
