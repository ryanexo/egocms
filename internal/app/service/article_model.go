package service

import (
    `context`
    `strconv`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/util/copierutil`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/jinzhu/copier`
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

func (s ArticleModel) UpdateDefinition(ctx context.Context, params dto.ArticleModelSchemaUpdateParams) error {
    modelDao := s.persist.ArticleModel
    
    modelId, err := strconv.ParseUint(string(params.ModelId), 10, 64)
    if err != nil {
        return erroz.ConvertTypeFailed.Wrap(err).ToError()
    }
    
    _, err = s.persist.ArticleModel.WithContext(ctx).Where(modelDao.ID.Eq(modelId)).First()
    if err != nil {
        return err
    }
    
    defCount := len(params.Data)
    
    appendDefs := make([]*model.ArticleModelSchema, 0, defCount)
    updateDefs := make([]*model.ArticleModelSchema, 0, defCount)
    
    for _, defParams := range params.Data {
        tmpDef := &model.ArticleModelSchema{ModelId: modelId}
        err := copier.Copy(&tmpDef, defParams)
        if err != nil {
            return err
        }
        
        if defParams.ID != nil {
            updateDefs = append(updateDefs, tmpDef)
        } else {
            appendDefs = append(appendDefs, tmpDef)
        }
    }
    
    return s.persist.Transaction(func(tx *query.Query) error {
        defDao := tx.WithContext(ctx).ArticleModelSchema
        
        for _, def := range updateDefs {
            _, err := defDao.Updates(def)
            if err != nil {
                return err
            }
        }
        
        return defDao.CreateInBatches(appendDefs, 500)
    })
}

func (s ArticleModel) FindDefinition(ctx context.Context, id uint64) ([]*dto.ArticleModelSchemaParams, error) {
    defs, err := s.persist.WithContext(ctx).ArticleModelSchema.Where(s.persist.ArticleModelSchema.ModelId.Eq(id)).Find()
    if err != nil {
        return nil, err
    }
    
    result := make([]*dto.ArticleModelSchemaParams, 0, len(defs))
    err = copier.CopyWithOption(&result, &defs, copier.Option{
        Converters: []copier.TypeConverter{copierutil.WithUint64ToString(),
        },
    })
    if err != nil {
        return nil, err
    }
    return result, nil
}
