package service

import (
    `context`
    
    `dpcms/internal/app/articlemodel/internal/assembler`
    `dpcms/internal/app/articlemodel/internal/dto`
    `dpcms/internal/app/articlemodel/repo`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/types`
)

type ArticleModelService struct {
    Query *query.Query
}

func (s ArticleModelService) CreateModel(ctx context.Context, params dto.ArticleModelCreateParams) (datatype.SafeUint64, error) {
    data := assembler.BuildArticleModelCreateCommand(&params)
    err := repo.NewArticleModelRepo(s.Query).Create(ctx, data)
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (s ArticleModelService) UpdateModel(ctx context.Context, params dto.ArticleModelUpdateParams) error {
    data := assembler.BuildArticleModelUpdateCommand(&params)
    return repo.NewArticleModelRepo(s.Query).UpdateModel(ctx, data)
}

func (s ArticleModelService) ReplaceSchema(ctx context.Context, params dto.ArticleModelSchemaUpdateParams) error {
    _, err := repo.NewArticleModelRepo(s.Query).FindByID(ctx, params.ID)
    if err != nil {
        return err
    }
    
    schema := make([]*model.ArticleModelSchema, 0, len(params.Data))
    
    for _, item := range params.Data {
        tmpSchema := assembler.BuildArticleModelSchemaModel(item)
        if item.ID != nil {
            tmpSchema.ModelId = params.ID
        }
        schema = append(schema, tmpSchema)
    }
    
    return s.Query.Transaction(func(tx *query.Query) error {
        return repo.NewArticleModelRepo(tx).ReplaceSchema(ctx, params.ID, schema)
    })
}

func (s ArticleModelService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.ArticleModel, error) {
    data, err := repo.NewArticleModelRepo(s.Query).FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.BuildArticleModelDTO(data), nil
}

func (s ArticleModelService) FindAllSchema(ctx context.Context, id datatype.SafeUint64) ([]*dto.ArticleModelSchemaParams, error) {
    allSchema, err := repo.NewArticleModelRepo(s.Query).FindAllSchema(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.BuildArticleModelSchemaList(allSchema), nil
}

func (s ArticleModelService) DeleteSchema(ctx context.Context, schemaId datatype.SafeUint64) error {
    _, err := repo.NewArticleModelRepo(s.Query).DeleteSchema(ctx, schemaId)
    return err
}

func (s ArticleModelService) DeleteModel(ctx context.Context, modelId datatype.SafeUint64) error {
    return s.Query.Transaction(func(tx *query.Query) error {
        modelRepo := repo.NewArticleModelRepo(s.Query)
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

func (s ArticleModelService) List(ctx context.Context, params dto.ArticleModelListParams) (*types.PaginatedResult[*dto.ArticleModel], error) {
    data, total, err := repo.NewArticleModelRepo(s.Query).List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &types.PaginatedResult[*dto.ArticleModel]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.BuildArticleModelListDTO(data),
    }, nil
}
