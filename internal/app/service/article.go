package service

import (
    `context`
    `database/sql`
    
    artAssembler `dpcms/internal/app/assembler/article`
    artDomain `dpcms/internal/app/domain/article`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/persistence/repo`
)

type Article struct {
    query *query.Query
}

func NewArticle(query *query.Query) *Article {
    return &Article{query}
}

func (s Article) Create(ctx context.Context, params dto.ArticleCreateParams) error {
    artModel := &model.Article{
        Url:         params.Url,
        Title:       params.Title,
        Description: params.Description,
        Target:      sql.NullString{String: params.Target, Valid: true},
        Content:     model.ArticleContent{Content: params.Content},
    }
    
    idInt, err := params.ModelId.Uint64()
    if params.ModelId != nil {
        if err != nil {
            return err
        }
        artModel.ModelId = idInt
    }
    
    art := artDomain.NewArticle(artDomain.Draft{Description: params.Description, Content: params.Content})
    params.Description = art.GetDescription()
    
    return s.query.Transaction(func(tx *query.Query) error {
        err := repo.NewArticle(tx).Create(ctx, artModel)
        if err != nil {
            return err
        }
        
        artModelRepo := repo.NewArticleModel(tx)
        schema, err := artModelRepo.FindAllSchema(ctx, artModel.ModelId)
        if err != nil {
            return err
        }
        
        jsonData, modelData, err := s.buildArticleModelData(artModel, schema, params.ModelData)
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
        
        return nil
    })
}

func (s Article) FindById(ctx context.Context, id uint64) (*model.Article, error) {
    return repo.NewArticle(s.query).FindById(ctx, id)
}

func (s Article) FindByIdWithContent(ctx context.Context, id uint64) (*model.Article, error) {
    return repo.NewArticle(s.query).FindByIdWithContent(ctx, id)
}

func (s Article) Update(ctx context.Context, data dto.ArticleUpdateParams) error {
    artId, err := data.ID.Uint64()
    if err != nil {
        return err
    }
    
    catId, err := data.CategoryId.Uint64()
    if err != nil {
        return err
    }
    
    rawArt, err := repo.NewArticle(s.query).FindById(ctx, artId)
    if err != nil {
        return err
    }
    
    art := artDomain.NewArticle(artDomain.Draft{Description: data.Description, Content: data.Content})
    
    artUpdateData := &model.Article{
        Base:        model.Base{ID: artId},
        Url:         data.Url,
        CategoryID:  catId,
        Flag:        data.Flag,
        Title:       data.Title,
        Description: art.GetDescription(),
        Target:      sql.NullString{String: data.Target, Valid: true},
    }
    
    return s.query.Transaction(func(tx *query.Query) error {
        artRepo := repo.NewArticle(tx)
        txErr := artRepo.Update(ctx, artUpdateData)
        if txErr != nil {
            return txErr
        }
        
        txErr = artRepo.UpdateContent(ctx, artUpdateData.ID, data.Content)
        if txErr != nil {
            return txErr
        }
        
        txErr = artRepo.UpdateKeywords(ctx, artUpdateData.ID, data.Keywords)
        if txErr != nil {
            return txErr
        }
        
        artModelRepo := repo.NewArticleModel(tx)
        schema, txErr := artModelRepo.FindAllSchema(ctx, rawArt.ModelId)
        if txErr != nil {
            return txErr
        }
        
        artJsonData, artTypedData, txErr := s.buildArticleModelData(rawArt, schema, data.ModelData)
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
        
        return nil
    })
}

func (s Article) Delete(ctx context.Context, id uint64) error {
    return s.query.Transaction(func(tx *query.Query) error {
        artRepo := repo.NewArticle(s.query)
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
        err = repo.NewArticleModel(s.query).DeleteArticleData(ctx, id)
        if err != nil {
            return err
        }
        return nil
    })
}

func (s Article) buildArticleModelData(art *model.Article, schemas []*model.ArticleModelSchema, data map[string]any) (*model.ArticleModelJsonData, []*model.ArticleModelData, error) {
    jsonResult := &model.ArticleModelJsonData{
        ArticleId: art.ID,
        ModelId:   art.ModelId,
        Data:      make(map[string]any),
    }
    modelResult := make([]*model.ArticleModelData, 0, len(schemas))
    
    for _, schema := range schemas {
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
            ModelId:   art.ModelId,
            ArticleId: art.ID,
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
