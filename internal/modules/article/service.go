package article

import (
    "context"
    "database/sql"
    
    "cms/internal/infra/persistence"
    model2 `cms/internal/infra/persistence/gorm/model`
    `cms/internal/infra/persistence/gorm/gquery`
    domain2 `cms/internal/modules/article/domain`
    `cms/internal/modules/contenttype/domain/valueobject`
    "cms/internal/pkg/datatype"
    
    contract2 "cms/internal/modules/article/contract"
    "cms/internal/modules/article/internal/assembler"
    "cms/internal/modules/article/internal/dto"
    "cms/internal/modules/article/internal/errno"
    contract3 "cms/internal/modules/contenttype/contract"
)

type Service struct {
    transactor       persistence.Transactor
    articleRepo      contract2.ArticleRepo
    articleModelRepo contract3.ArticleModelRepo
}

func NewArticleService(transactor persistence.Transactor, articleRepo contract2.ArticleRepo, articleModelRepo contract3.ArticleModelRepo) *Service {
    return &Service{
        transactor,
        articleRepo,
        articleModelRepo,
    }
}

func (s Service) Create(ctx context.Context, user *model2.User, params dto.ArticleCreateParams) (uint64, error) {
    artData := assembler.ToArticleCreateCommand(user, &params)
    
    content := domain2.NewContent(artData.Summary, artData.Content.Content)
    artData.Summary = content.Description()
    artData.Content.Content = content.Content()
    
    err := s.transactor.Transaction(func(tx *gquery.Query) error {
        err := s.articleRepo.CloneWithQuery(tx).Create(ctx, artData)
        if err != nil {
            return err
        }
        
        artModelRepo := s.articleModelRepo.CloneWithQuery(tx)
        
        if artData.ContentTypeID != nil {
            schema, err := artModelRepo.FindAllSchema(ctx, *artData.ContentTypeID)
            if err != nil {
                return err
            }
            
            jsonData, modelData, err := s.ToArticleModelData(artData.ID, *artData.ContentTypeID, schema, params.ModelData)
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

func (s Service) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Article, error) {
    artData, err := s.articleRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToArticleDTO(artData), nil
}

func (s Service) FindByIDWithContent(ctx context.Context, id datatype.SafeUint64) (*dto.Article, error) {
    artData, err := s.articleRepo.FindByIDWithContent(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToArticleDTO(artData), nil
}

func (s Service) Update(ctx context.Context, params dto.ArticleUpdateParams) error {
    artData, err := s.articleRepo.FindByID(ctx, params.ID)
    if err != nil {
        return err
    }
    
    content := domain2.NewContent(params.Description, params.Content)
    
    artUpdateData := &model2.Article{
        Base:       model2.Base{ID: params.ID},
        Url:        params.Url,
        CategoryID: params.CategoryID,
        Flag:       params.Flag,
        Title:      params.Title,
        Summary:    content.Description(),
        Target:     sql.NullString{String: params.Target, Valid: true},
    }
    
    return s.transactor.Transaction(func(tx *gquery.Query) error {
        artRepo := s.articleRepo.CloneWithQuery(tx)
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
        
        if artData.ContentTypeID != nil {
            artModelRepo := s.articleModelRepo.CloneWithQuery(tx)
            schema, txErr := artModelRepo.FindAllSchema(ctx, *artData.ContentTypeID)
            if txErr != nil {
                return txErr
            }
            
            artJsonData, artTypedData, txErr := s.ToArticleModelData(artData.ID, *artData.ContentTypeID, schema, params.ModelData)
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

func (s Service) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return s.transactor.Transaction(func(tx *gquery.Query) error {
        artRepo := s.articleRepo.CloneWithQuery(tx)
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
        err = s.articleModelRepo.CloneWithQuery(tx).DeleteArticleData(ctx, id)
        if err != nil {
            return err
        }
        return nil
    })
}

func (s Service) ChangeStatus(ctx context.Context, id datatype.SafeUint64, action func(status *domain2.Status) error) error {
    data, err := s.articleRepo.FindByIDWithoutPreload(ctx, id)
    if err != nil {
        return err
    }
    status := domain2.NewStatus(data.Status)
    if err = action(status); err != nil {
        return err
    }
    _, err = s.articleRepo.UpdateStatus(ctx, id, status.Value())
    return err
}

func (s Service) ToArticleModelData(artID datatype.SafeUint64, modelID datatype.SafeUint64, allSchema []*model2.ContentTypeSchema, data map[string]any) (*model2.ContentTypeEntries, []*model2.ContentFieldValues, error) {
    jsonResult := &model2.ContentTypeEntries{
        ArticleID:     artID,
        ContentTypeID: modelID,
        Data:          make(map[string]any),
    }
    modelResult := make([]*model2.ContentFieldValues, 0, len(allSchema))
    
    for _, schema := range allSchema {
        value := data[schema.FieldKey]
        
        v, err := valueobject.NewValue(schema.Type, value)
        if err != nil {
            return nil, nil, errno.ErrModelDataInvalidType.Format(schema.FieldName).ToError()
        }
        
        artModel := domain2.NewContentFieldValue(domain2.NewRules(schema), v)
        if err := artModel.IsValid(); err != nil {
            return nil, nil, err
        }
        
        modelData := assembler.NewModelData(&model2.ContentFieldValues{
            ModelID:   modelID,
            ArticleID: artID,
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
