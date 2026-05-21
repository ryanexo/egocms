package service

import (
    "context"
    "database/sql"
    
    contract2 `cms/internal/app/article/contract`
    `cms/internal/app/article/internal/assembler`
    `cms/internal/app/article/internal/domain`
    `cms/internal/app/article/internal/dto`
    `cms/internal/app/article/internal/errno`
    contract3 `cms/internal/app/articlemodel/contract`
    `cms/internal/infra/persistence/contract`
    "cms/internal/infra/persistence/datatype"
    "cms/internal/infra/persistence/model"
    "cms/internal/infra/persistence/query"
)

type ArticleService struct {
    txManager        contract.Transactor
    articleRepo      contract2.ArticleRepo
    articleModelRepo contract3.ArticleModelRepo
}

func NewArticleService(txManager contract.Transactor, articleRepo contract2.ArticleRepo, articleModelRepo contract3.ArticleModelRepo) *ArticleService {
    return &ArticleService{
        txManager:        txManager,
        articleRepo:      articleRepo,
        articleModelRepo: articleModelRepo,
    }
}

func (s ArticleService) Create(ctx context.Context, user *model.User, params dto.ArticleCreateParams) (datatype.SafeUint64, error) {
    artData := assembler.ToArticleCreateCommand(user, &params)
    
    content := domain.NewContent(artData.Description, artData.Content.Content)
    artData.Description = content.Description()
    artData.Content.Content = content.Content()
    
    err := s.txManager.Transaction(func(tx *query.Query) error {
        err := s.articleRepo.CloneWithQuery(tx).Create(ctx, artData)
        if err != nil {
            return err
        }
        
        artModelRepo := s.articleModelRepo.CloneWithQuery(tx)
        
        if artData.ModelID != nil {
            schema, err := artModelRepo.FindAllSchema(ctx, *artData.ModelID)
            if err != nil {
                return err
            }
            
            jsonData, modelData, err := s.ToArticleModelData(artData.ID, *artData.ModelID, schema, params.ModelData)
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

func (s ArticleService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Article, error) {
    artData, err := s.articleRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToArticleDTO(artData), nil
}

func (s ArticleService) FindByIDWithContent(ctx context.Context, id datatype.SafeUint64) (*dto.Article, error) {
    artData, err := s.articleRepo.FindByIDWithContent(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToArticleDTO(artData), nil
}

func (s ArticleService) Update(ctx context.Context, params dto.ArticleUpdateParams) error {
    artData, err := s.articleRepo.FindByID(ctx, params.ID)
    if err != nil {
        return err
    }
    
    content := domain.NewContent(params.Description, params.Content)
    
    artUpdateData := &model.Article{
        Base:        model.Base{ID: params.ID},
        Url:         params.Url,
        CategoryID:  params.CategoryID,
        Flag:        params.Flag,
        Title:       params.Title,
        Description: content.Description(),
        Target:      sql.NullString{String: params.Target, Valid: true},
    }
    
    return s.txManager.Transaction(func(tx *query.Query) error {
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
        
        if artData.ModelID != nil {
            artModelRepo := s.articleModelRepo.CloneWithQuery(tx)
            schema, txErr := artModelRepo.FindAllSchema(ctx, *artData.ModelID)
            if txErr != nil {
                return txErr
            }
            
            artJsonData, artTypedData, txErr := s.ToArticleModelData(artData.ID, *artData.ModelID, schema, params.ModelData)
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

func (s ArticleService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *query.Query) error {
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

func (s ArticleService) ChangeStatus(ctx context.Context, id datatype.SafeUint64, action func(status *domain.Status) error) error {
    data, err := s.articleRepo.FindByIDWithoutPreload(ctx, id)
    if err != nil {
        return err
    }
    status := domain.NewStatus(data.Status)
    if err = action(status); err != nil {
        return err
    }
    _, err = s.articleRepo.UpdateStatus(ctx, id, status.Value())
    return err
}

func (s ArticleService) ToArticleModelData(artID datatype.SafeUint64, modelID datatype.SafeUint64, allSchema []*model.ArticleModelSchema, data map[string]any) (*model.ArticleModelJsonData, []*model.ArticleModelData, error) {
    jsonResult := &model.ArticleModelJsonData{
        ArticleID: artID,
        ModelID:   modelID,
        Data:      make(map[string]any),
    }
    modelResult := make([]*model.ArticleModelData, 0, len(allSchema))
    
    for _, schema := range allSchema {
        value := data[schema.FieldKey]
        
        v, err := assembler.NewValue(schema.Type, value)
        if err != nil {
            return nil, nil, errno.ErrModelDataInvalidType.Format(schema.FieldName).ToError()
        }
        
        artModel := domain.NewModelValue(assembler.NewRules(schema), v)
        if err := artModel.IsValid(); err != nil {
            return nil, nil, err
        }
        
        modelData := assembler.NewModelData(&model.ArticleModelData{
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
