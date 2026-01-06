package service

import (
	"context"
	"database/sql"

	artAssembler "dpcms/internal/app/assembler/article"
	artDomain "dpcms/internal/app/domain/article"
	"dpcms/internal/app/dto"
	"dpcms/internal/app/erroz"
	"dpcms/internal/infra/persistence/datatype"
	"dpcms/internal/infra/persistence/model"
	"dpcms/internal/infra/persistence/query"
	"dpcms/internal/infra/persistence/repo"
)

type Article struct {
	persist *query.Query
}

func NewArticle(query *query.Query) *Article {
	return &Article{query}
}

func (s Article) Create(ctx context.Context, user *model.User, params dto.ArticleCreateParams) (datatype.SafeUint64, error) {
	artData := artAssembler.BuildArticleCreateCommand(user, &params)

	content := artDomain.NewContent(artData.Description, artData.Content.Content)
	artData.Description = content.Description()
	artData.Content.Content = content.Content()

	err := s.persist.Transaction(func(tx *query.Query) error {
		err := repo.NewArticle(tx).Create(ctx, artData)
		if err != nil {
			return err
		}

		artModelRepo := repo.NewArticleModel(tx)

		if artData.ModelID != nil {
			schema, err := artModelRepo.FindAllSchema(ctx, *artData.ModelID)
			if err != nil {
				return err
			}

			jsonData, modelData, err := s.buildArticleModelData(artData.ID, *artData.ModelID, schema, params.ModelData)
			if err != nil {
				return err
			}

			err = artModelRepo.CreateModelTypedData(ctx, modelData)
			if err != nil {
				return err
			}

			err = artModelRepo.CreateModelJsonData(ctx, jsonData)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return artData.ID, nil
}

func (s Article) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Article, error) {
	artData, err := repo.NewArticle(s.persist).FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return artAssembler.BuildArticleDTO(artData), nil
}

func (s Article) FindByIDWithContent(ctx context.Context, id datatype.SafeUint64) (*dto.Article, error) {
	artData, err := repo.NewArticle(s.persist).FindByIDWithContent(ctx, id)
	if err != nil {
		return nil, err
	}
	return artAssembler.BuildArticleDTO(artData), nil
}

func (s Article) Update(ctx context.Context, params dto.ArticleUpdateParams) error {
	artData, err := repo.NewArticle(s.persist).FindByID(ctx, params.ID)
	if err != nil {
		return err
	}

	content := artDomain.NewContent(params.Description, params.Content)

	artUpdateData := &model.Article{
		Base:        model.Base{ID: params.ID},
		Url:         params.Url,
		CategoryID:  params.CategoryId,
		Flag:        params.Flag,
		Title:       params.Title,
		Description: content.Description(),
		Target:      sql.NullString{String: params.Target, Valid: true},
	}

	return s.persist.Transaction(func(tx *query.Query) error {
		artRepo := repo.NewArticle(tx)
		txErr := artRepo.Update(ctx, artUpdateData)
		if txErr != nil {
			return txErr
		}

		txErr = artRepo.UpdateContent(ctx, artUpdateData.ID, content.Content())
		if txErr != nil {
			return txErr
		}

		txErr = artRepo.ReplaceKeywords(ctx, artUpdateData.ID, params.Keywords)
		if txErr != nil {
			return txErr
		}

		if artData.ModelID != nil {
			artModelRepo := repo.NewArticleModel(tx)
			schema, txErr := artModelRepo.FindAllSchema(ctx, *artData.ModelID)
			if txErr != nil {
				return txErr
			}

			artJsonData, artTypedData, txErr := s.buildArticleModelData(artData.ID, *artData.ModelID, schema, params.ModelData)
			if txErr != nil {
				return txErr
			}

			txErr = artModelRepo.UpdateModelTypedData(ctx, artTypedData)
			if txErr != nil {
				return txErr
			}

			txErr = artModelRepo.UpdateModelJsonData(ctx, artJsonData)
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})
}

func (s Article) Delete(ctx context.Context, id datatype.SafeUint64) error {
	return s.persist.Transaction(func(tx *query.Query) error {
		artRepo := repo.NewArticle(s.persist)
		err := artRepo.DeleteArticle(ctx, id)
		if err != nil {
			return err
		}
		err = artRepo.DeleteContent(ctx, id)
		if err != nil {
			return err
		}
		err = artRepo.DeleteKeywords(ctx, id)
		if err != nil {
			return err
		}
		err = repo.NewArticleModel(s.persist).DeleteArticleData(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s Article) buildArticleModelData(artId datatype.SafeUint64, modelId datatype.SafeUint64, allSchema []*model.ArticleModelSchema, data map[string]any) (*model.ArticleModelJsonData, []*model.ArticleModelData, error) {
	jsonResult := &model.ArticleModelJsonData{
		ArticleId: artId,
		ModelId:   modelId,
		Data:      make(map[string]any),
	}
	modelResult := make([]*model.ArticleModelData, 0, len(allSchema))

	for _, schema := range allSchema {
		value := data[schema.FieldKey]

		v, err := artAssembler.NewValue(schema.Type, value)
		if err != nil {
			return nil, nil, erroz.ArticleModelDataInvalidType.Format(schema.FieldName).ToError()
		}

		artModel := artDomain.NewModelValue(artAssembler.NewRules(schema), v)
		if err := artModel.IsValid(); err != nil {
			return nil, nil, err
		}

		modelData := artAssembler.NewModelData(&model.ArticleModelData{
			ModelId:   modelId,
			ArticleId: artId,
			FieldKey:  schema.FieldKey,
			FieldName: schema.FieldName,
			Type:      schema.Type,
		})
		if err := artModel.Assign(modelData); err != nil {
			return nil, nil, err
		}

		jsonResult.Data[schema.FieldKey] = value
		modelResult = append(modelResult, modelData.Model())
	}

	return jsonResult, modelResult, nil
}
