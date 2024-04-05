-- 角色列表
DROP TABLE IF EXISTS bg_role;

CREATE TABLE bg_role
(
    id     INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name   VARCHAR(32)  NOT NULL COMMENT '角色名称',
    remark VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
    PRIMARY KEY (id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 角色权限
DROP TABLE IF EXISTS bg_role_permission;

CREATE TABLE bg_role_permission
(
    id       INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    role_id  INT UNSIGNED    NOT NULL COMMENT '角色ID',
    route_id INT UNSIGNED    NOT NULL COMMENT '权限ID',
    perm     BIGINT UNSIGNED NOT NULL COMMENT '权限值',
    PRIMARY KEY (id),
    UNIQUE relation (role_id, route_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 系统路由
DROP TABLE IF EXISTS bg_route;

CREATE TABLE bg_route
(
    id  INT UNSIGNED NOT NULL AUTO_INCREMENT,
    uri VARCHAR(255) NOT NULL COMMENT '路由标识',
    PRIMARY KEY (id),
    UNIQUE (uri)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 后台菜单表
DROP TABLE IF EXISTS bg_menu;

CREATE TABLE bg_menu
(
    id        INT UNSIGNED     NOT NULL AUTO_INCREMENT,
    root_id   INT UNSIGNED     NOT NULL COMMENT '根节点ID',
    parent_id INT UNSIGNED     NOT NULL COMMENT '父节点ID',
    name      VARCHAR(32)      NOT NULL COMMENT '菜单名称',
    route_id  VARCHAR(255)     NOT NULL COMMENT '对应路径',
    display   TINYINT UNSIGNED NOT NULL COMMENT '是否显示',
    sequence  INT UNSIGNED     NOT NULL COMMENT '排序',
    PRIMARY KEY (id),
    KEY (parent_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 文章数据
DROP TABLE IF EXISTS bg_article;

CREATE TABLE bg_article
(
    id           INT UNSIGNED     NOT NULL AUTO_INCREMENT,
    category_id  INT UNSIGNED     NOT NULL COMMENT '分类ID',
    user_id      INT UNSIGNED     NOT NULL COMMENT '用户ID',
    custom_url   VARCHAR(255)     NOT NULL DEFAULT '' COMMENT '自定义URL格式',
    flag         TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '属性,如置顶/幻灯',
    title        VARCHAR(255)     NOT NULL COMMENT '标题',
    description  VARCHAR(255)     NOT NULL DEFAULT '' COMMENT '描述',
    thumb        VARCHAR(255)     NOT NULL DEFAULT '' COMMENT '缩略图',
    seo_keywords VARCHAR(128)     NOT NULL DEFAULT '' COMMENT '关键词',
    click        INT UNSIGNED     NOT NULL DEFAULT 0 COMMENT '点击数',
    status       TINYINT UNSIGNED NOT NULL COMMENT '0:正常,1:审核,2:草稿箱',
    created_at   INT UNSIGNED     NOT NULL COMMENT '发布时间',
    updated_at   INT UNSIGNED     NOT NULL COMMENT '最后更新时间',
    deleted_at   INT UNSIGNED     NOT NULL DEFAULT 0 COMMENT '删除时间',
    PRIMARY KEY (id),
    UNIQUE (custom_url),
    KEY (category_id),
    KEY (user_id),
    KEY (created_at)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 文章内容
DROP TABLE IF EXISTS bg_article_data;

CREATE TABLE bg_article_data
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT,
    article_id INT UNSIGNED NOT NULL,
    title      VARCHAR(255) NOT NULL DEFAULT '' COMMENT '标题',
    data       LONGTEXT     NOT NULL COMMENT '内容',
    PRIMARY KEY (id),
    UNIQUE (article_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS bg_article_addition;

CREATE TABLE bg_article_addition
(
    id               INT UNSIGNED NOT NULL AUTO_INCREMENT,
    article_id       INT UNSIGNED NOT NULL,
    addition_id      INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '模型ID',
    addition_data_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '模型内容ID',
    PRIMARY KEY (id),
    UNIQUE (article_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 评论表
DROP TABLE IF EXISTS bg_reply;

CREATE TABLE bg_reply
(
    id         INT UNSIGNED  NOT NULL AUTO_INCREMENT,
    target_id  INT UNSIGNED  NOT NULL COMMENT '被回复数据ID,可能是文章/单页/他人回复',
    user_id    INT UNSIGNED  NOT NULL,
    data       VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '回复内容',
    created_at INT UNSIGNED  NOT NULL COMMENT '评论时间',
    updated_at INT UNSIGNED  NOT NULL COMMENT '更新时间',
    deleted_at INT UNSIGNED  NOT NULL COMMENT '删除时间',
    PRIMARY KEY (id),
    KEY (target_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 文章附加信息
DROP TABLE IF EXISTS bg_addition;

CREATE TABLE bg_addition
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name       VARCHAR(32)  NOT NULL COMMENT '模型名称',
    table_name VARCHAR(32)  NOT NULL COMMENT '数据表名',
    PRIMARY KEY (id),
    UNIQUE (table_name)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 附加信息字段
DROP TABLE IF EXISTS bg_addition_fields;

CREATE TABLE bg_addition_fields
(
    id          INT UNSIGNED     NOT NULL AUTO_INCREMENT,
    addition_id INT UNSIGNED     NOT NULL COMMENT '所属ID',
    name        VARCHAR(32)      NOT NULL COMMENT '名称',
    identifier  VARCHAR(16)      NOT NULL COMMENT '唯一标识符',
    type        TINYINT UNSIGNED NOT NULL COMMENT '字段类型,如单选,多选,输入框',
    message     VARCHAR(64)      NOT NULL DEFAULT '' COMMENT '提示信息',
    sequence    INT UNSIGNED     NOT NULL DEFAULT 0 COMMENT '展示顺序',
    required    TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否必填',
    filter      TINYINT UNSIGNED NOT NULL COMMENT '允许搜索',
    data        VARCHAR(255)     NOT NULL DEFAULT '' COMMENT '数据,如选项',
    PRIMARY KEY (id),
    UNIQUE (addition_id, identifier),
    KEY (filter),
    KEY (sequence)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 附件文件信息
DROP TABLE IF EXISTS bg_file;

CREATE TABLE bg_file
(
    id       INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    filename VARCHAR(255)    NOT NULL COMMENT '原文件名',
    path     VARCHAR(255)    NOT NULL COMMENT '文件路径',
    size     BIGINT UNSIGNED NOT NULL COMMENT '文件大小',
    ext      VARCHAR(16)     NOT NULL COMMENT '扩展名',
    is_image TINYINT         NOT NULL COMMENT '图片标识',
    disk     VARCHAR(32)     NOT NULL COMMENT '驱动程序',
    md5      CHAR(32)        NOT NULL COMMENT '文件MD5',
    sha1     CHAR(40)        NOT NULL COMMENT '文件SHA1',
    PRIMARY KEY (id),
    KEY (ext)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 附件上传信息
DROP TABLE IF EXISTS bg_file_record;

CREATE TABLE bg_file_record
(
    id             INT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id        INT UNSIGNED NOT NULL COMMENT '用户ID',
    file_id        INT UNSIGNED NOT NULL COMMENT '附件ID',
    remark         VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
    permission     VARCHAR(255) NOT NULL COMMENT '下载权限,逗号分隔多个用户组',
    download_count INT UNSIGNED NOT NULL COMMENT '下载次数',
    price          INT UNSIGNED NOT NULL COMMENT '附件价格',
    created_at     INT UNSIGNED NOT NULL COMMENT '上传时间',
    PRIMARY KEY (id),
    KEY record (user_id, file_id),
    KEY (created_at)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 分类表
DROP TABLE IF EXISTS bg_category;

CREATE TABLE bg_category
(
    id              INT UNSIGNED     NOT NULL AUTO_INCREMENT,
    sequence        INT UNSIGNED     NOT NULL COMMENT '排序权重',
    name            VARCHAR(32)      NOT NULL COMMENT '栏目名称',
    alias           VARCHAR(64)      NOT NULL COMMENT '栏目路径',
    root_id         INT UNSIGNED     NOT NULL COMMENT '根节点ID',
    parent_id       INT UNSIGNED     NOT NULL COMMENT '父节点ID',
    type            TINYINT UNSIGNED NOT NULL COMMENT '类别,0:栏目,1:单页,2:外链',
    display         TINYINT UNSIGNED NOT NULL COMMENT '是否显示,0:否,1:是',
    seo_title       VARCHAR(255)     NOT NULL DEFAULT '' COMMENT 'SEO标题',
    seo_keywords    VARCHAR(255)     NOT NULL DEFAULT '' COMMENT 'SEO关键词',
    seo_description VARCHAR(255)     NOT NULL DEFAULT '' COMMENT 'SEO简介',
    PRIMARY KEY (id),
    UNIQUE (alias),
    KEY relation (root_id, parent_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 单页内容表
DROP TABLE IF EXISTS bg_category_page;

CREATE TABLE bg_category_page
(
    id          INT UNSIGNED NOT NULL AUTO_INCREMENT,
    category_id INT UNSIGNED NOT NULL COMMENT '栏目ID',
    title       VARCHAR(255) NOT NULL COMMENT '单页标题',
    data        LONGTEXT     NOT NULL COMMENT '单页内容',
    PRIMARY KEY (id),
    UNIQUE (category_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 系统选项
DROP TABLE IF EXISTS bg_var;

CREATE TABLE bg_var
(
    id    INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name  VARCHAR(32)  NOT NULL COMMENT '唯一名称',
    value VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    UNIQUE (name)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 用户表
DROP TABLE IF EXISTS bg_user;

CREATE TABLE bg_user
(
    id          INT UNSIGNED NOT NULL AUTO_INCREMENT,
    username    VARCHAR(32)  NOT NULL COMMENT '用户名',
    password    VARCHAR(255) NOT NULL COMMENT '密码',
    email       VARCHAR(128) NOT NULL COMMENT '用户邮箱',
    status      TINYINT      NOT NULL COMMENT '账号状态,0:停用,-1:未验证,1:正常',
    group_id    INT          NOT NULL COMMENT '用户组',
    role_id     INT          NOT NULL COMMENT '用户角色',
    verified_at INT UNSIGNED NOT NULL COMMENT '邮箱认证时间',
    created_ip  VARCHAR(39)  NOT NULL COMMENT '注册IP',
    created_at  INT UNSIGNED NOT NULL COMMENT '注册时间',
    updated_at  INT UNSIGNED NOT NULL COMMENT '变更时间',
    PRIMARY KEY (id),
    UNIQUE (username),
    UNIQUE (email)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 登录日志
DROP TABLE IF EXISTS bg_log_login;

CREATE TABLE bg_log_login
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id    INT UNSIGNED NOT NULL,
    login_time INT UNSIGNED NOT NULL COMMENT '登录时间',
    login_ip   INT UNSIGNED NOT NULL COMMENT '登录IP',
    PRIMARY KEY (id),
    KEY (user_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci PARTITION BY HASH (id) PARTITIONS 10;

-- 积分变更日志
DROP TABLE IF EXISTS bg_log_points;

CREATE TABLE bg_log_points
(
    id       INT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id  INT UNSIGNED NOT NULL COMMENT '用户ID',
    admin_id INT UNSIGNED NOT NULL COMMENT '管理员ID,0为系统',
    type     TINYINT      NOT NULL COMMENT '操作类型, 0:减少, 1:增加',
    time     INT UNSIGNED NOT NULL COMMENT '变动时间',
    points   INT UNSIGNED NOT NULL COMMENT '积分数量',
    remark   VARCHAR(255) NOT NULL COMMENT '变动原因',
    PRIMARY KEY (id),
    KEY (user_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci PARTITION BY HASH (id) PARTITIONS 10;

-- 用户组
DROP TABLE IF EXISTS bg_group;

CREATE TABLE bg_group
(
    id   INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '用户组名称',
    perm BIGINT       NOT NULL DEFAULT 0 COMMENT '权限',
    PRIMARY KEY (id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 短消息
DROP TABLE IF EXISTS bg_message;

CREATE TABLE bg_message
(
    id           INT UNSIGNED     NOT NULL AUTO_INCREMENT,
    sender       INT UNSIGNED     NOT NULL COMMENT '发送人ID',
    receiver     INT UNSIGNED     NOT NULL COMMENT '接收人ID',
    created_at   INT UNSIGNED     NOT NULL COMMENT '发送时间',
    message      VARCHAR(255)     NOT NULL COMMENT '消息内容',
    status       TINYINT UNSIGNED NOT NULL COMMENT '状态,0:未读,1:已读',
    delete_state TINYINT UNSIGNED NOT NULL COMMENT '消息状态,0:正常,1:发送人已清理,2:接收人已清理,3:双方已清理',
    PRIMARY KEY (id),
    KEY (sender),
    KEY (receiver),
    KEY (delete_state)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 插件
DROP TABLE IF EXISTS bg_plugin;

CREATE TABLE bg_plugin
(
    id         INT UNSIGNED     NOT NULL AUTO_INCREMENT,
    active     TINYINT UNSIGNED NOT NULL COMMENT '激活状态',
    identifier VARCHAR(32)      NOT NULL DEFAULT '' COMMENT '标识符',
    weight     INT UNSIGNED     NOT NULL COMMENT '插件权重',
    PRIMARY KEY (id),
    KEY (identifier)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 离线任务
DROP TABLE IF EXISTS bg_progress;

CREATE TABLE bg_progress
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT,
    identifier VARCHAR(8)   NOT NULL COMMENT '唯一标识符',
    progress   TINYINT      NOT NULL COMMENT '进度百分比,-1为失败',
    PRIMARY KEY (id),
    UNIQUE (identifier)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 内容TAG
DROP TABLE IF EXISTS bg_tag;

CREATE TABLE bg_tag
(
    id     INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name   VARCHAR(32)  NOT NULL DEFAULT '' COMMENT 'TAG名称',
    remark VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'TAG备注',
    PRIMARY KEY (id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci PARTITION BY HASH (id) PARTITIONS 50;

-- 内容TAG关联信息
DROP TABLE IF EXISTS bg_tag_index;

CREATE TABLE bg_tag_index
(
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT,
    tag_id     INT UNSIGNED NOT NULL COMMENT '标签ID',
    content_id INT UNSIGNED NOT NULL COMMENT '内容ID',
    PRIMARY KEY (id),
    UNIQUE (tag_id, content_id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 自定义路由
DROP TABLE IF EXISTS bg_url_route;

CREATE TABLE bg_url_route
(
    id     INT UNSIGNED NOT NULL AUTO_INCREMENT,
    weight INT UNSIGNED NOT NULL COMMENT '权重',
    route  VARCHAR(255) NOT NULL DEFAULT '' COMMENT '路由',
    target VARCHAR(255) NOT NULL DEFAULT '' COMMENT '目标',
    PRIMARY KEY (id)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
