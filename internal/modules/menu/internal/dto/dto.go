package dto

import (
    "cms/internal/infra/store/datatype"
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
)

type Menu struct {
    apitype.Base
    ParentID    jsontype.SafeUint64 `json:"parentId"`
    Type        int8                `json:"type"`
    Name        string              `json:"name"`
    Affix       datatype.BoolInt8   `json:"affix"`
    Icon        string              `json:"icon"`
    ExternalURL string              `json:"externalUrl"`
    Sequence    jsontype.SafeInt64  `json:"sequence"`
    Visible     datatype.BoolInt8   `json:"visible,omitempty" apitype:"boolean"`
    URI         string              `json:"uri"`
    Template    string              `json:"template"`
    Remark      string              `json:"remark"`
}

type MenuListQueryParams struct {
    apitype.Pagination
    Name     *string              `json:"name"`
    ParentID *jsontype.SafeUint64 `json:"parentId" apitype:"string"`
}

type MenuCreateParams struct {
    ParentID    jsontype.SafeUint64 `json:"parentId" apitype:"string"`
    Type        int8                `json:"type"`
    Name        string              `validate:"required" json:"name" label:"名称"`
    Affix       datatype.BoolInt8   `json:"affix"`
    Icon        string              `json:"icon"`
    ExternalURL string              `json:"externalUrl"`
    Sequence    jsontype.SafeInt64  `json:"sequence" apitype:"string"`
    URI         string              `validate:"required,alphanum" json:"uri"`
    Template    string              `json:"template"`
    Remark      string              `json:"remark"`
}

type MenuUpdateParams struct {
    ID          jsontype.SafeUint64 `validate:"required" json:"id" apitype:"string"`
    Type        int8                `json:"type"`
    Name        string              `json:"name"`
    Affix       datatype.BoolInt8   `json:"affix"`
    Icon        string              `json:"icon"`
    ExternalURL string              `json:"externalUrl"`
    Sequence    jsontype.SafeInt64  `json:"sequence" apitype:"string"`
    URI         string              `json:"uri"`
    Template    string              `json:"template"`
    Remark      string              `json:"remark"`
}

type MenuMoveParams struct {
    ID       jsontype.SafeUint64 `validate:"required" json:"id" apitype:"string"`
    TargetID jsontype.SafeUint64 `validate:"required" json:"targetId" apitype:"string"`
}
