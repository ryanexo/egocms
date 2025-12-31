package erroz

import (
    `dpcms/internal/app/erroz/internal/errmod`
    `dpcms/internal/app/erroz/internal/errtype`
)

var (
    OK           = New(Code(errmod.Server, errtype.OK, 0), "成功")
    Unknown      = New(Code(errmod.Server, errtype.Unknown, 0), "系统异常")
    DataNotFound = New(Code(errmod.Server, errtype.NotFound, 0), "数据不存在")
    
    ValidationFailed = New(Code(errmod.Client, errtype.Parameter, 0), "参数错误")
    
    Unauthorized         = New(Code(errmod.Server, errtype.NotFound, 0), "未授权")
    AuthorizationExpired = New(Code(errmod.Server, errtype.InvalidState, 0), "授权已过期，请重新登陆")
    
    UserNotExists          = New(Code(errmod.User, errtype.NotFound, 0), "用户不存在")
    UserNameNotExists      = New(Code(errmod.User, errtype.NotFound, 1), "用户名不存在")
    UserNameExists         = New(Code(errmod.User, errtype.Conflict, 0), "用户名已存在")
    UserEmailExists        = New(Code(errmod.User, errtype.Conflict, 1), "邮箱已存在")
    UserWrongPasswd        = New(Code(errmod.User, errtype.Parameter, 0), "密码错误")
    UserWrongConfirmPasswd = New(Code(errmod.User, errtype.Parameter, 1), "两次密码输入不一致")
    UserEqualsOldPasswd    = New(Code(errmod.User, errtype.Parameter, 2), "新密码不能和旧密码相同")
    
    RoleCircular = New(Code(errmod.Role, errtype.Conflict, 0), "角色 %s 已继承自当前角色，无法建立继承关系")
    
    MenuCircular = New(Code(errmod.Menu, errtype.Conflict, 0), "目标菜单 %s 为当前菜单 %s 的子级，无法移动")
    
    CategoryCircular = New(Code(errmod.Category, errtype.Conflict, 0), "目标分类 %s 为当前分类 %s 的子级，无法移动")
    
    ArticlePublishStatusNotAllowed   = New(Code(errmod.Article, errtype.InvalidState, 0), "内容非待审状态，无法发布")
    ArticleSubmitStatusNotAllowed    = New(Code(errmod.Article, errtype.InvalidState, 1), "内容非草稿状态，无法提交")
    ArticleOfflineStatusNotAllowed   = New(Code(errmod.Article, errtype.InvalidState, 2), "内容非发布状态，无法下线")
    ArticleRejectStatusNotAllowed    = New(Code(errmod.Article, errtype.InvalidState, 3), "内容非待审状态，无法发布")
    ArticleRepublishStatusNotAllowed = New(Code(errmod.Article, errtype.InvalidState, 4), "内容非发布状态，无法提交")
    ArticleAlreadyPublished          = New(Code(errmod.Article, errtype.Conflict, 0), "文章已发布，请勿重复发布")
    ArticleAlreadyRepublish          = New(Code(errmod.Article, errtype.Conflict, 1), "文章已提交，请勿重复提交")
    ArticleAlreadySubmitted          = New(Code(errmod.Article, errtype.Conflict, 2), "文章已提交，请勿重复提交")
    ArticleAlreadyOffline            = New(Code(errmod.Article, errtype.Conflict, 3), "文章已下线，请勿重复下线")
    ArticleAlreadyReject             = New(Code(errmod.Article, errtype.Conflict, 4), "文章已拒审，请勿重复拒审")
    ArticleMissingSubmitPerm         = New(Code(errmod.Article, errtype.Auth, 0), "无提交权限")
    ArticleMissingPublishPerm        = New(Code(errmod.Article, errtype.Auth, 1), "无发布权限")
    ArticleMissingOfflinePerm        = New(Code(errmod.Article, errtype.Auth, 2), "无下线权限")
    ArticleMissingRejectPerm         = New(Code(errmod.Article, errtype.Auth, 3), "无审核权限")
    
    ArticleModelDataDisabled         = New(Code(errmod.ArticleModel, errtype.Parameter, 0), "%s参数不可用")
    ArticleModelDataInvalidType      = New(Code(errmod.ArticleModel, errtype.Parameter, 1), "%s数据类型不合法")
    ArticleModelDataMissingValue     = New(Code(errmod.ArticleModel, errtype.Parameter, 2), "%s必填")
    ArticleModelDataInvalidNumRange  = New(Code(errmod.ArticleModel, errtype.Parameter, 3), "%s必须在 %.2f 到 %.2f 之间")
    ArticleModelDataInvalidTimeRange = New(Code(errmod.ArticleModel, errtype.Parameter, 3), "%s必须在 %s 到 %s 之间")
    ArticleModelDataInvalidLen       = New(Code(errmod.ArticleModel, errtype.Parameter, 4), "%s长度必须在 %d 到 %d 之间")
    ArticleModelDataInvalidValue     = New(Code(errmod.ArticleModel, errtype.Parameter, 5), "%s值不合法")
)
