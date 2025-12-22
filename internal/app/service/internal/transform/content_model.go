package transform

import (
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/model`
    
    `github.com/bytedance/sonic`
    `github.com/jinzhu/copier`
)

func ContentModelDefinitionToDTO(obj ...*model.ContentModelDefinition) ([]*dto.ContentModelDefinitionParams, error) {
    result := make([]*dto.ContentModelDefinitionParams, 0, len(obj))
    for _, def := range obj {
        tmpDef := &dto.ContentModelDefinitionParams{}
        if err := copier.Copy(&tmpDef, def); err != nil {
            return nil, err
        }
        if def.ConfigJSON != "" {
            configMap := map[string]any{}
            if err := sonic.UnmarshalString(def.ConfigJSON, &configMap); err != nil {
                return nil, err
            }
            tmpDef.ConfigData = configMap
        }
        result = append(result, tmpDef)
    }
    return result, nil
}
