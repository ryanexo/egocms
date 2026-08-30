# authz

`authz` 是一个用于 Gin handler 的认证与权限中间件。

它只负责包装 handler，不持有 `*gin.Engine`，不注册路由，也不维护路由表。`Middleware` 是不可变值，派生公开或权限策略时不会修改原实例，运行期间不需要加锁。

## 快速开始

先注入凭证解析器和权限检查器，再为业务资源创建中间件：

```go
factory := authz.New(tokenParser, permissionChecker)
auth := factory.Resource("article")

// 需要 article:create 权限。
engine.POST("/create", auth.Permission("create").Wrap(createRouteHandler))

// 无需权限和登录凭证。
engine.POST("/view", auth.Public().Wrap(viewRouteHandler))

// 仅检查登录凭证。
engine.POST("/draft", auth.Wrap(draftRouteHandler))
```

`Resource` 由 `Factory` 创建，是因为凭证解析器和权限检查器需要通过项目依赖注入提供。这样可以避免包级可变状态，同时保持路由注册代码简洁。

## 策略

| 使用方式 | 验证凭证 | 检查权限 |
| --- | --- | --- |
| `auth.Wrap(handler)` | 是 | 否 |
| `auth.Public().Wrap(handler)` | 否 | 否 |
| `auth.Permission("create").Wrap(handler)` | 是 | 是 |

默认策略是验证凭证。权限策略会先验证凭证，再调用权限检查器；任一步骤失败都不会执行被包装的 handler。

`Public` 和 `Permission` 返回新的 `Middleware` 值，不会改变 `auth`：

```go
createAuth := auth.Permission("create")

// auth 仍然是仅验证凭证的默认策略。
engine.POST("/draft", auth.Wrap(draftRouteHandler))
engine.POST("/create", createAuth.Wrap(createRouteHandler))
```

## 依赖接口

调用方需要实现以下接口：

```go
type User interface {
    UserID() uint64
    Role() string
}

type TokenParser interface {
    Parse(ctx context.Context, token string) (User, error)
}

type PermissionChecker interface {
    Check(
        ctx context.Context,
        subject string,
        object string,
        action string,
    ) (bool, error)
}
```

权限检查参数的来源：

| 参数 | 来源 | 示例 |
| --- | --- | --- |
| `subject` | 当前用户的 `Role()` | `editor` |
| `object` | `Resource` 的参数 | `article` |
| `action` | `Permission` 的参数 | `create` |

## 凭证

需要认证的请求必须携带 Bearer 凭证：

```text
Authorization: Bearer <token>
```

请求头缺失、格式错误、解析器缺失或凭证解析失败时，中间件会中止请求并向 `gin.Context` 写入 `ErrAuthorized`。

## 权限

`Permission` 策略使用以下三元组检查权限：

```text
用户角色, 资源, 动作
```

例如：

```text
editor, article, create
```

权限检查器缺失、资源或动作为空、检查器返回错误或拒绝访问时，中间件会中止请求并写入 `ErrAccessDenied`。

## 当前用户

认证成功后，可以从 Gin 上下文或请求上下文读取当前用户：

```go
user, err := authz.GetCurrentUser(ctx)
user, err = authz.GetCurrentUser(ctx.Request.Context())

account, err := authz.CurrentUser[*Account](ctx)
```

公开路由不会设置当前用户。没有用户或用户类型不匹配时返回 `ErrAuthorized`。

## 错误响应

中间件只负责中止请求并通过 `gin.Context.Error` 传递错误，不直接生成 HTTP 响应。项目应使用统一错误处理中间件将以下错误转换为响应：

- `ErrAuthorized`：凭证验证失败。
- `ErrAccessDenied`：权限验证失败。
