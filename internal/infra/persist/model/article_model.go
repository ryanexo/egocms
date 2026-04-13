package model

import (
    `database/sql`
    
    `cms/internal/infra/persist/datatype`
    
    `github.com/shopspring/decimal`
)

type ArticleModel struct {
    Base
    Name        string                `gorm:"type:varchar(255);not null;comment:'模型名称'"`
    Description string                `gorm:"type:varchar(255);not null;comment:'模型说明'"`
    Schema      []*ArticleModelSchema `gorm:"foreignKey:ModelID;referenceKey:ID"`
}

type ArticleModelJsonData struct {
    Base
    ArticleID datatype.SafeUint64 `gorm:"index:idx_artid_modelid,priority:1"`
    ModelID   datatype.SafeUint64 `gorm:"index:idx_artid_modelid,priority:2"`
    Data      datatype.JSONMap    `gorm:"type:text;comment:'文章的模型json数据'"`
}

type ArticleModelSchema struct {
    Base
    ModelID     datatype.SafeUint64 `gorm:"uniqueIndex:idx_field_key,priority:1;"`
    FieldKey    string              `gorm:"type:varchar(255);not null;uniqueIndex:idx_field_key,priority:2;comment:'字段唯一标识'"`
    FieldName   string              `gorm:"type:varchar(255);not null;comment:'字段名称'"`
    Description string              `gorm:"type:varchar(255);not null;default:'';comment:'字段说明'"`
    MinLen      datatype.SafeUint64 `gorm:"default:0;comment:'最小长度'"`
    MaxLen      datatype.SafeUint64 `gorm:"default:0;comment:'最大长度'"`
    MinValue    decimal.NullDecimal `gorm:"decimal(10,2);comment:'最小值'"`
    MaxValue    decimal.NullDecimal `gorm:"decimal(10,2);comment:'最大值'"`
    MinTime     sql.NullTime        `gorm:"comment:'最小时间'"`
    MaxTime     sql.NullTime        `gorm:"comment:'最大时间'"`
    Pattern     string              `gorm:"type:varchar(255);comment:'正则'"`
    Sequence    datatype.SafeInt64  `gorm:"index;default:0;comment:'字段排序'"`
    Type        int16               `gorm:"type:smallint;not null;comment:'字段值类型'"`
    EnumOptions datatype.EnumValues `gorm:"type:text;comment:'字段枚举json'"`
    Required    datatype.BoolInt8   `gorm:"type:tinyint;default:0;comment:'是否必填'"`
    Hidden      datatype.BoolInt8   `gorm:"type:tinyint;default:0;comment:'是否隐藏'"`
    Enable      datatype.BoolInt8   `gorm:"type:tinyint;default:1;comment:'是否启用'"`
}

type ArticleModelData struct {
    Base
    ModelID     datatype.SafeUint64 `gorm:"index:idx_artid_modelid_fieldkey,priority:1"`
    ArticleID   datatype.SafeUint64 `gorm:"index:idx_artid_modelid_fieldkey,priority:2"`
    FieldName   string              `gorm:"type:varchar(255);comment:'冗余字段名称'"`
    FieldKey    string              `gorm:"type:varchar(255);index:idx_artid_modelid_fieldkey,priority:3;comment:'字段唯一标识'"`
    Type        int16               `gorm:"type:smallint;not null;comment:'字段值类型'"`
    ValueString sql.NullString      `gorm:"type:varchar(255);comment:'字段字符串内容'"`
    ValueBool   sql.NullBool        `gorm:"type:bool;comment:'字段布尔内容'"`
    ValueNumber decimal.NullDecimal `gorm:"type:decimal(10,2);comment:'字段数字内容'"`
    ValueTime   sql.NullTime        `gorm:"type:datetime;comment:'字段时间内容'"`
}
