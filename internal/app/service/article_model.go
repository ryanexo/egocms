package service

import (
    `context`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/persistence/repo`
    
    `github.com/jinzhu/copier`
)

type ArticleModel struct {
    persist *query.Query
}

func NewArticleModel(persist *query.Query) *ArticleModel {
    return &ArticleModel{persist}
}

func (s ArticleModel) Create(ctx context.Context, params dto.ArticleModelCreateParams) (*model.ArticleModel, error) {
    m := model.ArticleModel{Name: params.Name, Description: params.Description}
    err := s.persist.ArticleModel.WithContext(ctx).Create(&m)
    return &m, err
}

func (s ArticleModel) Update(ctx context.Context, params dto.ArticleModelUpdateParams) error {
    m := &model.ArticleModel{}
    err := copier.Copy(m, &params)
    if err != nil {
        return err
    }
    return repo.NewArticleModel(s.persist).Update(ctx, m)
}

func (s ArticleModel) UpdateSchema(ctx context.Context, params dto.ArticleModelSchemaUpdateParams) error {
    modelDao := s.persist.ArticleModel
    
    _, err := s.persist.ArticleModel.WithContext(ctx).Where(modelDao.ID.Eq(params.ID.Raw())).First()
    if err != nil {
        return err
    }
    
    schemaCount := len(params.Data)
    appendSchema := make([]*model.ArticleModelSchema, 0, schemaCount)
    updateSchema := make([]*model.ArticleModelSchema, 0, schemaCount)
    
    for _, schemaParams := range params.Data {
        tmpSchema := &model.ArticleModelSchema{ModelId: params.ID}
        err := copier.Copy(&tmpSchema, &schemaParams)
        if err != nil {
            return err
        }
        
        if schemaParams.ID != nil {
            updateSchema = append(updateSchema, tmpSchema)
        } else {
            appendSchema = append(appendSchema, tmpSchema)
        }
    }
    
    return s.persist.Transaction(func(tx *query.Query) error {
        defDao := tx.WithContext(ctx).ArticleModelSchema
        
        for _, def := range updateSchema {
            _, err := defDao.Updates(def)
            if err != nil {
                return err
            }
        }
        
        return defDao.CreateInBatches(appendSchema, 500)
    })
}

func (s ArticleModel) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.ArticleModel, error) {
    return repo.NewArticleModel(s.persist).FindByID(ctx, id)
}

func (s ArticleModel) FindAllSchema(ctx context.Context, id datatype.SafeUint64) ([]*dto.ArticleModelSchemaParams, error) {
    allSchema, err := repo.NewArticleModel(s.persist).FindAllSchema(ctx, id)
    if err != nil {
        return nil, err
    }
    result := make([]*dto.ArticleModelSchemaParams, 0, len(allSchema))
    
    err = copier.Copy(&result, &allSchema)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}

func (s ArticleModel) DeleteSchema(ctx context.Context, schemaId datatype.SafeUint64) error {
    return repo.NewArticleModel(s.persist).DeleteSchema(ctx, schemaId)
}

func (s ArticleModel) DeleteModel(ctx context.Context, modelId datatype.SafeUint64) error {
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
