package dto

type ContentModelCreateParams struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ContentModelDefinitionUpdateParams struct {
    ContentModelId uint64 `json:"contentModelId"`
    
    Data []ContentModelDefinitionParams `json:"data"`
}

type ContentModelDefinitionParams struct {
    ID             *uint64        `json:"id"`
    ContentModelId uint64         `json:"contentModelId"`
    Field          string         `validate:"required,max=255" json:"field"`
    Name           string         `validate:"required,max=255" json:"name"`
    Description    string         `validate:"max=255" json:"description"`
    Sequence       int64          `json:"sequence"`
    Type           int16          `validate:"required" json:"type"`
    Required       int8           `validate:"oneof=0 1" json:"required"`
    ConfigData     map[string]any `copier:"-" json:"configData"`
}
