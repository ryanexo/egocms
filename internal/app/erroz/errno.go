package errno

import (
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/erroz/errno/internal/errmod`
    `dpcms/internal/app/erroz/errno/internal/errtype`
)

var (
    OK           = erroz.New(New(errmod.Server, errtype.OK, 0), "成功")
    Unknown      = erroz.New(New(errmod.Server, errtype.Unknown, 0), "系统异常")
    DataNotFound = erroz.New(New(errmod.Server, errtype.NotFound, 0), "数据不存在")
    
    ValidationFailed = erroz.New(New(errmod.Client, errtype.Parameter, 0), "参数错误")
    
    Unauthorized         = erroz.New(New(errmod.Server, errtype.NotFound, 0), "未授权")
    AuthorizationExpired = erroz.New(New(errmod.Server, errtype.InvalidState, 0), "授权已过期，请重新登陆")
    
    UserNotExists          = erroz.New(New(errmod.User, errtype.NotFound, 0), "用户不存在")
    UserNameNotExists      = erroz.New(New(errmod.User, errtype.NotFound, 1), "用户名不存在")
    UserNameExists         = erroz.New(New(errmod.User, errtype.Conflict, 0), "用户名已存在")
    UserEmailExists        = erroz.New(New(errmod.User, errtype.Conflict, 1), "邮箱已存在")
    UserWrongPasswd        = erroz.New(New(errmod.User, errtype.Parameter, 0), "密码错误")
    UserWrongConfirmPasswd = erroz.New(New(errmod.User, errtype.Parameter, 1), "两次密码输入不一致")
    UserEqualsOldPasswd    = erroz.New(New(errmod.User, errtype.Parameter, 2), "新密码不能和旧密码相同")
    
    RoleCircular = erroz.New(New(errmod.Role, errtype.Conflict, 0), "角色 %s 已继承自当前角色，无法建立继承关系")
    
    MenuCircular = erroz.New(New(errmod.Menu, errtype.Conflict, 0), "目标菜单 %s 为当前菜单 %s 的子级，无法移动")
    
    CategoryCircular = erroz.New(New(errmod.Category, errtype.Conflict, 0), "目标分类 %s 为当前分类 %s 的子级，无法移动")
    
    ArticlePublishStatusNotAllowed   = erroz.New(New(errmod.Article, errtype.InvalidState, 0), "内容非待审状态，无法发布")
    ArticleSubmitStatusNotAllowed    = erroz.New(New(errmod.Article, errtype.InvalidState, 1), "内容非草稿状态，无法提交")
    ArticleOfflineStatusNotAllowed   = erroz.New(New(errmod.Article, errtype.InvalidState, 2), "内容非发布状态，无法下线")
    ArticleRejectStatusNotAllowed    = erroz.New(New(errmod.Article, errtype.InvalidState, 3), "内容非待审状态，无法发布")
    ArticleRepublishStatusNotAllowed = erroz.New(New(errmod.Article, errtype.InvalidState, 4), "内容非发布状态，无法提交")
    ArticleAlreadyPublished          = erroz.New(New(errmod.Article, errtype.Conflict, 0), "文章已发布，请勿重复发布")
    ArticleAlreadyRepublish          = erroz.New(New(errmod.Article, errtype.Conflict, 1), "文章已提交，请勿重复提交")
    ArticleAlreadySubmitted          = erroz.New(New(errmod.Article, errtype.Conflict, 2), "文章已提交，请勿重复提交")
    ArticleAlreadyOffline            = erroz.New(New(errmod.Article, errtype.Conflict, 3), "文章已下线，请勿重复下线")
    ArticleAlreadyReject             = erroz.New(New(errmod.Article, errtype.Conflict, 4), "文章已拒审，请勿重复拒审")
    ArticleMissingSubmitPerm         = erroz.New(New(errmod.Article, errtype.Auth, 0), "无提交权限")
    ArticleMissingPublishPerm        = erroz.New(New(errmod.Article, errtype.Auth, 1), "无发布权限")
    ArticleMissingOfflinePerm        = erroz.New(New(errmod.Article, errtype.Auth, 2), "无下线权限")
    ArticleMissingRejectPerm         = erroz.New(New(errmod.Article, errtype.Auth, 3), "无审核权限")
    
    ArticleModelDataInvalidType  = erroz.New(New(errmod.ArticleModel, errtype.Parameter, 0), "%s数据类型不合法")
    ArticleModelDataMissingValue = erroz.New(New(errmod.ArticleModel, errtype.Parameter, 1), "%s必填")
)
