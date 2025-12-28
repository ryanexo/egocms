package service

import (
    `context`
    `database/sql`
    `strconv`
    
    artAssembler `dpcms/internal/app/assembler/article`
    `dpcms/internal/app/domain/article`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/repo`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
)

type Article struct {
    query *query.Query
}

func NewArticle(query *query.Query) *Article {
    return &Article{query}
}

func (s Article) Create(ctx context.Context, params dto.ArticleCreateParams) (*dto.ArticleDetail, error) {
    art := model.Article{
        Url:         params.Url,
        Title:       params.Title,
        Description: params.Description,
        Target:      sql.NullString{String: params.Target, Valid: true},
        Content:     model.ArticleContent{Content: params.Content},
    }
    
    if params.ModelId != nil {
        idStr := *(*string)(params.ModelId)
        idInt, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            return nil, err
        }
        art.ModelId = idInt
    }
    
    if params.Description != "" && params.Content != "" {
        contentStr := []rune(params.Content)
        art.Description = string(contentStr[:200])
    }
    
    err := s.query.Transaction(func(tx *query.Query) error {
        artRepo := repo.NewArticle(ctx, tx)
        txErr := artRepo.Create(&art)
        if txErr != nil {
            return txErr
        }
        
        jsonMap, modelData, err := s.transformModelDataWithVerify(ctx, art, params.ModelData)
        if err != nil {
            return err
        }
        
        return artRepo.SaveContentModelData(art, jsonMap, modelData)
    })
    
    if err != nil {
        return nil, err
    }
    
    return s.FindDetailById(art.ID)
}

func (s Article) FindBasicById(id uint64) (*dto.Article, error) {

}

func (s Article) FindDetailById(id uint64) (*dto.ArticleDetail, error) {}

func (s Article) transformModelDataWithVerify(ctx context.Context, art model.Article, data map[string]any) (map[string]any, []*model.ArticleModelData, error) {
    defDao := s.query.ArticleModelSchema
    defs, err := defDao.WithContext(ctx).Where(defDao.ModelId.Eq(art.ModelId)).Find()
    if err != nil {
        return nil, nil, err
    }
    
    jsonMapResult := make(map[string]any)
    modelResult := make([]*model.ArticleModelData, 0, len(defs))
    for _, def := range defs {
        value := data[def.FieldKey]
        
        v, err := artAssembler.NewValue(def.Type, value)
        if err != nil {
            return nil, nil, erroz.ArticleModelDataInvalidType.Format(def.FieldName).ToError()
        }
        
        if err := article.NewModelValue(artAssembler.NewRules(def), v).IsValid(); err != nil {
            return nil, nil, err
        }
        
        scannableModelData, persistModelData := artAssembler.NewModelData(model.ArticleModelData{
            ModelId:   art.ModelId,
            ArticleId: art.ID,
            FieldKey:  def.FieldKey,
            Type:      def.Type,
        })
        if err := v.Assign(scannableModelData); err != nil {
            return nil, nil, err
        }
        
        jsonMapResult[def.FieldKey] = value
        modelResult = append(modelResult, persistModelData)
    }
    
    return jsonMapResult, modelResult, nil
}
