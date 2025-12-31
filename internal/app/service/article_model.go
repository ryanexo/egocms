package service

import (
    `context`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/util/copierutil`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/persistence/repo`
)

type ArticleModel struct {
    persist *query.Query
}

func NewContentModel(persist *query.Query) *ArticleModel {
    return &ArticleModel{persist}
}

func (s ArticleModel) Create(ctx context.Context, params dto.ArticleModelCreateParams) (*model.ArticleModel, error) {
    m := model.ArticleModel{Name: params.Name, Description: params.Description}
    err := s.persist.ArticleModel.WithContext(ctx).Create(&m)
    return &m, err
}

func (s ArticleModel) UpdateSchemas(ctx context.Context, params dto.ArticleModelSchemaUpdateParams) error {
    modelDao := s.persist.ArticleModel
    
    modelId, err := params.ModelId.Uint64()
    if err != nil {
        return err
    }
    
    _, err = s.persist.ArticleModel.WithContext(ctx).Where(modelDao.ID.Eq(modelId)).First()
    if err != nil {
        return err
    }
    
    schemaCount := len(params.Data)
    appendSchemas := make([]*model.ArticleModelSchema, 0, schemaCount)
    updateSchemas := make([]*model.ArticleModelSchema, 0, schemaCount)
    
    for _, schemaParams := range params.Data {
        tmpSchema := &model.ArticleModelSchema{ModelId: modelId}
        err := copierutil.CopyWithIdConverter(&tmpSchema, &schemaParams)
        if err != nil {
            return err
        }
        
        if schemaParams.ID != nil {
            updateSchemas = append(updateSchemas, tmpSchema)
        } else {
            appendSchemas = append(appendSchemas, tmpSchema)
        }
    }
    
    return s.persist.Transaction(func(tx *query.Query) error {
        defDao := tx.WithContext(ctx).ArticleModelSchema
        
        for _, def := range updateSchemas {
            _, err := defDao.Updates(def)
            if err != nil {
                return err
            }
        }
        
        return defDao.CreateInBatches(appendSchemas, 500)
    })
}

func (s ArticleModel) FindAllSchema(ctx context.Context, id uint64) ([]*dto.ArticleModelSchemaParams, error) {
    schemas, err := repo.NewArticleModel(s.persist).FindAllSchema(ctx, id)
    if err != nil {
        return nil, err
    }
    result := make([]*dto.ArticleModelSchemaParams, 0, len(schemas))
    
    err = copierutil.CopyWithIdConverter(&result, &schemas)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}

func (s ArticleModel) DeleteSchema(ctx context.Context, schemaId uint64) error {
    return repo.NewArticleModel(s.persist).DeleteSchema(ctx, schemaId)
}

func (s ArticleModel) DeleteModel(ctx context.Context, modelId uint64) error {
    return s.persist.Transaction(func(tx *query.Query) error {
        modelRepo := repo.NewArticleModel(s.persist)
        txErr := modelRepo.DeleteModel(ctx, modelId)
        if txErr != nil {
            return txErr
        }
        txErr = modelRepo.DeleteAllSchema(ctx, modelId)
        if txErr != nil {
            return txErr
        }
        return nil
    })
}
