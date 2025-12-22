package service

import (
    `context`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/service/internal/transform`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/bytedance/sonic`
    `github.com/jinzhu/copier`
)

type ContentModel struct {
    persist *query.Query
}

func NewContentModel(i *infra.Infra) *ContentModel {
    return &ContentModel{persist: i.Query}
}

func (s ContentModel) Create(ctx context.Context, params dto.ContentModelCreateParams) (*model.ContentModel, error) {
    m := model.ContentModel{Name: params.Name, Description: params.Description}
    err := s.persist.ContentModel.WithContext(ctx).Create(&m)
    return &m, err
}

func (s ContentModel) UpdateDefinition(ctx context.Context, params dto.ContentModelDefinitionUpdateParams) error {
    modelDao := s.persist.ContentModel
    
    _, err := s.persist.ContentModel.WithContext(ctx).Where(modelDao.ID.Eq(params.ContentModelId)).First()
    if err != nil {
        return err
    }
    
    defCount := len(params.Data)
    
    appendDefs := make([]*model.ContentModelDefinition, 0, defCount)
    updateDefs := make([]*model.ContentModelDefinition, 0, defCount)
    
    for i, defParams := range params.Data {
        tmpDef := &model.ContentModelDefinition{}
        err := copier.Copy(&tmpDef, defParams)
        if err != nil {
            return err
        }
        
        if defParams.ConfigData != nil {
            json, err := sonic.MarshalString(params.Data[i])
            if err != nil {
                return err
            }
            tmpDef.ConfigJSON = json
        }
        
        if defParams.ID != nil {
            updateDefs = append(updateDefs, tmpDef)
        } else {
            appendDefs = append(appendDefs, tmpDef)
        }
    }
    
    return s.persist.Transaction(func(tx *query.Query) error {
        defDao := tx.WithContext(ctx).ContentModelDefinition
        
        for _, def := range updateDefs {
            _, err := defDao.Updates(def)
            if err != nil {
                return err
            }
        }
        
        return defDao.CreateInBatches(appendDefs, 500)
    })
}

func (s ContentModel) FindDefinition(ctx context.Context, id uint64) ([]*dto.ContentModelDefinitionParams, error) {
    defs, err := s.persist.WithContext(ctx).ContentModelDefinition.Where(s.persist.ContentModelDefinition.ContentModelId.Eq(id)).Find()
    if err != nil {
        return nil, err
    }
    return transform.ContentModelDefinitionToDTO(defs...)
}
