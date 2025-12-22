package erroz

var (
    // 成功
    OK = New("ok", "成功")
    
    // 服务端错误
    Unknown      = New("server.error", "服务器错误")
    DataNotFound = New("server.data.notfound", "数据不存在")
    
    // 客户端参数/验证错误
    ValidationFailed = New("client.param.error", "参数错误")
    
    // 认证/授权相关
    AuthorizationExpired = New("auth.expired", "授权已过期，请重新登陆")
    Unauthorized         = New("auth.unauthorized", "未授权")
    
    // 用户相关
    UserNotExists            = New("user.notfound.id", "用户不存在")
    UserUsernameNotExists    = New("user.notfound.username", "用户名不存在")
    UsernameExists           = New("user.exists.username", "用户名已存在")
    UserEmailExists          = New("user.exists.email", "邮箱已存在")
    UserWrongPassword        = New("user.password.wrong", "密码错误")
    UserWrongConfirmPassword = New("user.password.confirm_wrong", "两次密码输入不一致")
    UserNewPwdEqualsOldPwd   = New("user.password.equals_old", "新密码不能和旧密码相同")
    
    // 角色相关
    RoleCircularReference = New("role.inherit.circular", "角色 %s 已继承自当前角色，无法建立继承关系")
    
    // 菜单相关
    MenuCircularReferenceWhenMove = New("menu.move.circular", "目标菜单 %s 为当前菜单 %s 的子级")
    
    // 分类相关
    CategoryCircularReferenceWhenMove = New("category.move.circular", "目标分类 %s 为当前分类 %s 的子级")
    
    // 内容相关
    ArticlePublishStatusNotAllowed = New("article.publish.status_not_allowed", "内容非待审状态，无法发布")
    ArticleSubmitStatusNotAllowed  = New("article.submit.status_not_allowed", "内容非草稿状态，无法提交")
    ArticleOfflineStatusNotAllowed = New("article.offline.status_not_allowed", "内容非发布状态，无法下线")
    ArticleRejectStatusNotAllowed  = New("article.reject.status_not_allowed", "内容非待审状态，无法发布")
    ArticleAlreadyPublished        = New("article.publish.invalid", "文章已是发布状态")
    ArticleAlreadySubmitted        = New("article.submit.invalid", "文章已提交，请勿重复提交")
    ArticleAlreadyOffline          = New("article.offline.invalid", "文章已下线，请勿重复下线")
    ArticleAlreadyReject           = New("article.reject.invalid", "文章已拒审，请勿重复拒审")
    ArticleMissingSubmitPerm       = New("client.permission.article.missing.submit", "无提交权限")
    ArticleMissingPublishPerm      = New("client.permission.article.missing.publish", "无发布权限")
    ArticleMissingOfflinePerm      = New("client.permission.article.missing.offline", "无下线权限")
    ArticleMissingRejectPerm       = New("client.permission.article.missing.reject", "无审核权限")
    
    // 内容模型相关
    ContentModelJsonFormat = New("content_model.definition.invalid_json", "配置数据JSON不合法")
)
