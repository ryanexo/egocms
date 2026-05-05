package service

import (
    `context`
    
    `cms/internal/app/articlemodel/internal/assembler`
    `cms/internal/app/articlemodel/internal/dto`
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    `cms/internal/infra/persistence/query`
    `cms/internal/util/types`
)

type ArticleModelService struct {
    txManager contract.Transactor
    repo      contract.ArticleModelRepo
}

func NewArticleModelService(txManager contract.Transactor, repo contract.ArticleModelRepo) *ArticleModelService {
    return &ArticleModelService{
        txManager: txManager,
        repo:      repo,
    }
}

func (s ArticleModelService) CreateModel(ctx context.Context, params dto.ArticleModelCreateParams) (datatype.SafeUint64, error) {
    data := assembler.ToArticleModelCreateCommand(&params)
    err := s.repo.Create(ctx, data)
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (s ArticleModelService) UpdateModel(ctx context.Context, params dto.ArticleModelUpdateParams) error {
    data := assembler.ToArticleModelUpdateCommand(&params)
    return s.repo.UpdateModel(ctx, data)
}

func (s ArticleModelService) UpdateModelSchema(ctx context.Context, params dto.ArticleModelSchemaUpdateParams) error {
    _, err := s.repo.FindByID(ctx, params.ID)
    if err != nil {
        return err
    }
    
    schema := make([]*model.ArticleModelSchema, 0, len(params.Data))
    
    for _, item := range params.Data {
        tmpSchema := assembler.ToArticleModelSchemaModel(item)
        if item.ID != nil {
            tmpSchema.ModelID = params.ID
        }
        schema = append(schema, tmpSchema)
    }
    
    return s.txManager.Transaction(func(tx *query.Query) error {
        return s.repo.CloneWithQuery(tx).ReplaceSchema(ctx, params.ID, schema)
    })
}

func (s ArticleModelService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.ArticleModel, error) {
    data, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToArticleModelDTO(data), nil
}

func (s ArticleModelService) FindAllSchema(ctx context.Context, id datatype.SafeUint64) ([]*dto.ArticleModelSchemaParams, error) {
    allSchema, err := s.repo.FindAllSchema(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToArticleModelSchemaList(allSchema), nil
}

func (s ArticleModelService) DeleteSchema(ctx context.Context, schemaID datatype.SafeUint64) error {
    _, err := s.repo.DeleteSchema(ctx, schemaID)
    return err
}

func (s ArticleModelService) DeleteModel(ctx context.Context, modelID datatype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *query.Query) error {
        modelRepo := s.repo.CloneWithQuery(tx)
        txErr := modelRepo.DeleteModel(ctx, modelID)
        if txErr != nil {
            return txErr
        }
        txErr = modelRepo.DeleteAllSchema(ctx, modelID)
        if txErr != nil {
            return txErr
        }
        return nil
    })
}

func (s ArticleModelService) List(ctx context.Context, params dto.ArticleModelListParams) (*types.PaginatedResult[*dto.ArticleModel], error) {
    data, total, err := s.repo.List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &types.PaginatedResult[*dto.ArticleModel]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.ToArticleModelListDTO(data),
    }, nil
}
