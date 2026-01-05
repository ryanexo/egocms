package service

import (
    `context`
    
    artmodelAssembler `dpcms/internal/app/assembler/artmodel`
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/persistence/repo`
)

type ArticleModel struct {
    persist *query.Query
}

func NewArticleModel(persist *query.Query) *ArticleModel {
    return &ArticleModel{persist}
}

func (s ArticleModel) CreateModel(ctx context.Context, params dto.ArticleModelCreateParams) (datatype.SafeUint64, error) {
    data := artmodelAssembler.BuildArticleModelCreateCommand(&params)
    err := repo.NewArticleModel(s.persist).Create(ctx, data)
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (s ArticleModel) UpdateModel(ctx context.Context, params dto.ArticleModelUpdateParams) error {
    data := artmodelAssembler.BuildArticleModelUpdateCommand(&params)
    return repo.NewArticleModel(s.persist).UpdateModel(ctx, data)
}

func (s ArticleModel) ReplaceSchema(ctx context.Context, params dto.ArticleModelSchemaUpdateParams) error {
    _, err := repo.NewArticleModel(s.persist).FindByID(ctx, params.ID)
    if err != nil {
        return err
    }
    
    schema := make([]*model.ArticleModelSchema, 0, len(params.Data))
    
    for _, item := range params.Data {
        tmpSchema := artmodelAssembler.BuildArticleModelSchemaModel(item)
        if item.ID != nil {
            tmpSchema.ModelId = params.ID
        }
        schema = append(schema, tmpSchema)
    }
    
    return s.persist.Transaction(func(tx *query.Query) error {
        return repo.NewArticleModel(tx).ReplaceSchema(ctx, params.ID, schema)
    })
}

func (s ArticleModel) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.ArticleModel, error) {
    data, err := repo.NewArticleModel(s.persist).FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return artmodelAssembler.BuildArticleModelDTO(data), nil
}

func (s ArticleModel) FindAllSchema(ctx context.Context, id datatype.SafeUint64) ([]*dto.ArticleModelSchemaParams, error) {
    allSchema, err := repo.NewArticleModel(s.persist).FindAllSchema(ctx, id)
    if err != nil {
        return nil, err
    }
    return artmodelAssembler.BuildArticleModelSchemaList(allSchema), nil
}

func (s ArticleModel) DeleteSchema(ctx context.Context, schemaId datatype.SafeUint64) error {
    _, err := repo.NewArticleModel(s.persist).DeleteSchema(ctx, schemaId)
    return err
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

func (s ArticleModel) List(ctx context.Context, params dto.ArticleModelListParams) (*dto.PaginatedResult[*dto.ArticleModel], error) {
    data, total, err := repo.NewArticleModel(s.persist).List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &dto.PaginatedResult[*dto.ArticleModel]{
        Total:    total,
        PageSize: params.PageSize,
        PageNo:   params.PageNo,
        List:     artmodelAssembler.BuildArticleModelListDTO(data),
    }, nil
}
