package service

import (
    `context`
    
    `dpcms/internal/app/articlemodel/internal/assembler`
    `dpcms/internal/app/articlemodel/internal/dto`
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/types`
)

type ArticleModelService struct {
    TxManager contract.TxManager
    Repo      ArticleModelRepo
}

func (s ArticleModelService) CreateModel(ctx context.Context, params dto.ArticleModelCreateParams) (datatype.SafeUint64, error) {
    data := assembler.BuildArticleModelCreateCommand(&params)
    err := s.Repo.Create(ctx, data)
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (s ArticleModelService) UpdateModel(ctx context.Context, params dto.ArticleModelUpdateParams) error {
    data := assembler.BuildArticleModelUpdateCommand(&params)
    return s.Repo.UpdateModel(ctx, data)
}

func (s ArticleModelService) ReplaceSchema(ctx context.Context, params dto.ArticleModelSchemaUpdateParams) error {
    _, err := s.Repo.FindByID(ctx, params.ID)
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
    
    return s.TxManager.Transaction(func(tx *query.Query) error {
        return s.Repo.CloneWithQuery(tx).ReplaceSchema(ctx, params.ID, schema)
    })
}

func (s ArticleModelService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.ArticleModel, error) {
    data, err := s.Repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.BuildArticleModelDTO(data), nil
}

func (s ArticleModelService) FindAllSchema(ctx context.Context, id datatype.SafeUint64) ([]*dto.ArticleModelSchemaParams, error) {
    allSchema, err := s.Repo.FindAllSchema(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.BuildArticleModelSchemaList(allSchema), nil
}

func (s ArticleModelService) DeleteSchema(ctx context.Context, schemaId datatype.SafeUint64) error {
    _, err := s.Repo.DeleteSchema(ctx, schemaId)
    return err
}

func (s ArticleModelService) DeleteModel(ctx context.Context, modelId datatype.SafeUint64) error {
    return s.TxManager.Transaction(func(tx *query.Query) error {
        modelRepo := s.Repo.CloneWithQuery(tx)
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
    data, total, err := s.Repo.List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &types.PaginatedResult[*dto.ArticleModel]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.BuildArticleModelListDTO(data),
    }, nil
}
