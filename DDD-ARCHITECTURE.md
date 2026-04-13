# 轻量 DDD 目录布局建议

## 现有架构分析

当前项目已经有一些 DDD 的雏形：
- `domain` 目录按业务模块划分（article、user、category 等）
- `internal/domain/{module}/internal/domain` 有值对象（Value、Status、Content）
- `contract` 定义仓储接口，`provider` 实现仓储
- `service` 作为应用层

但存在以下问题：
1. **Service 层直接依赖基础设施层**（`persist/model`、`persist/query`）
2. **领域对象贫血**：`model.Article` 是纯数据结构，业务逻辑在 Service 中
3. **领域层依赖泄露**：`domain/types.go` 引用了 `infra/persist/datatype`

---

## 建议目录结构

```
internal/
├── domain/                          # 领域层（核心）
│   ├── article/
│   │   ├── article.go               # 聚合根
│   │   ├── content.go               # 值对象
│   │   ├── status.go                # 值对象（状态机）
│   │   ├── errors.go                # 领域错误
│   │   ├── events.go                # 领域事件（可选）
│   │   └── repository.go            # 仓储接口定义
│   ├── user/
│   │   ├── user.go
│   │   └── repository.go
│   └── shared/                      # 共享值对象/基类
│       ├── base.go
│       └── types.go
│
├── application/                     # 应用层
│   ├── article/
│   │   ├── service.go               # 应用服务
│   │   ├── dto.go                   # 输入输出DTO
│   │   └── assembler.go             # DTO-领域对象转换
│   └── user/
│
├── infrastructure/                  # 基础设施层
│   ├── persistence/                 # 持久化
│   │   ├── model/                   # ORM 模型（PO）
│   │   ├── query/                   # 查询构建器
│   │   └── repository/              # 仓储实现
│   │       ├── article_repo.go
│   │       └── user_repo.go
│   ├── cache/
│   ├── logger/
│   └── ...
│
├── interfaces/                      # 接口层（可选，或保持现有的 httpserver）
│   └── http/
│       ├── controller/
│       │   └── article.go
│       └── middleware/
│
├── bootstrap/                       # 依赖注入/启动（保持）
├── config/                          # 配置（保持）
└── middleware/                      # 中间件（保持）
```

---

## 核心改动点

### 1. 领域对象富化

```go
// internal/app/article/article.go
package article

import "time"

// Article 聚合根
type Article struct {
    id          ArticleID
    title       Title
    content     Content
    status      Status
    categoryID  CategoryID
    createdAt   time.Time
    updatedAt   time.Time
}

// 工厂方法
func NewArticle(title Title, content Content, categoryID CategoryID) *Article {
    return &Article{
        title:      title,
        content:    content,
        status:     StatusDraft,
        categoryID: categoryID,
        createdAt:  time.Now(),
        updatedAt:  time.Now(),
    }
}

// 业务行为
func (a *Article) Publish() error {
    if err := a.status.CanTransitionTo(StatusPublished); err != nil {
        return err
    }
    a.status = StatusPublished
    a.updatedAt = time.Now()
    return nil
}

func (a *Article) UpdateContent(content Content) error {
    if a.status == StatusPublished {
        return ErrCannotUpdatePublishedArticle
    }
    a.content = content
    a.updatedAt = time.Now()
    return nil
}

// Getter（只读）
func (a *Article) ID() ArticleID       { return a.id }
func (a *Article) Title() Title        { return a.title }
func (a *Article) Status() Status      { return a.status }
```

### 2. 仓储接口属于领域层

```go
// internal/app/article/repository.go
package article

type ArticleRepository interface {
    Save(ctx context.Context, article *Article) error
    FindByID(ctx context.Context, id ArticleID) (*Article, error)
    Delete(ctx context.Context, id ArticleID) error
}
```

### 3. 仓储实现属于基础设施层

```go
// internal/infrastructure/persistence/repository/article_repo.go
package repository

type articleRepo struct {
    db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) article.ArticleRepository {
    return &articleRepo{db: db}
}

func (r *articleRepo) Save(ctx context.Context, agg *article.Article) error {
    // 领域对象 -> PO 转换
    po := toPO(agg)
    return r.db.WithContext(ctx).Create(po).Error
}

func (r *articleRepo) FindByID(ctx context.Context, id article.ArticleID) (*article.Article, error) {
    var po ArticlePO
    if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
        return nil, err
    }
    // PO -> 领域对象 转换
    return toDomain(&po), nil
}
```

### 4. 应用服务编排

```go
// internal/application/article/service.go
package article

type ArticleService struct {
    repo      article.ArticleRepository
    txManager contract.TxManager
}

func (s *ArticleService) Create(ctx context.Context, cmd CreateArticleCommand) (ArticleID, error) {
    // 1. 创建领域对象
    title, _ := article.NewTitle(cmd.Title)
    content, _ := article.NewContent(cmd.Description, cmd.Body)
    
    agg := article.NewArticle(title, content, cmd.CategoryID)
    
    // 2. 持久化
    if err := s.repo.Save(ctx, agg); err != nil {
        return article.ArticleID{}, err
    }
    
    return agg.ID(), nil
}
```

---

## 轻量化权衡

| 完整 DDD | 轻量 DDD（建议） |
|---------|----------------|
| 聚合、实体、值对象严格区分 | 聚合根 + 值对象即可 |
| 领域事件、事件溯源 | 暂不引入 |
| CQRS 分离 | 保持单一模型 |
| 应用层命令/查询分离 | 可用 DTO 替代 |
| 六边形架构/洋葱架构 | 分层架构即可 |

---

## 迁移建议

1. **先改造一个模块**：选择 `article` 作为试点
2. **建立领域模型**：把 `internal/domain/article/internal/domain` 中的逻辑富化
3. **剥离基础设施依赖**：`domain` 包不依赖 `infra`
4. **逐步迁移**：其他模块参照模式迁移

---

## 依赖方向

```
interfaces（接口层）
    ↓
application（应用层）
    ↓
domain（领域层）
    ↑
infrastructure（基础设施层）
```

关键原则：
- `domain` 层不依赖任何外层
- `infrastructure` 实现由 `domain` 定义的接口
- `application` 编排领域对象，不包含业务规则
- `interfaces` 处理 HTTP/gRPC 等外部请求
