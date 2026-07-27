DROP TABLE IF EXISTS `article`;
CREATE TABLE `article`
(
    id                   BIGINT UNSIGNED AUTO_INCREMENT,
    created_at           DATETIME        NOT NULL,
    deleted_at           DATETIME                 DEFAULT NULL,
    updated_at           DATETIME        NOT NULL,

    current_version_id   BIGINT UNSIGNED          DEFAULT NULL,
    published_version_id BIGINT UNSIGNED          DEFAULT NULL,
    content_type_id      BIGINT UNSIGNED          DEFAULT NULL,
    author_id            BIGINT UNSIGNED NOT NULL,

    url                  VARCHAR(255)             DEFAULT NULL,
    slug                 VARCHAR(255)             DEFAULT NULL,
    title                VARCHAR(255)             DEFAULT NULL,
    summary              VARCHAR(500)             DEFAULT '',
    status               tinyint         NOT NULL DEFAULT 0,

    PRIMARY KEY (id),
    INDEX author_id (author_id),
    INDEX current_version_id (current_version_id),
    INDEX published_version_id (published_version_id),
    INDEX slug (slug),
    INDEX content_type_id (content_type_id),
    INDEX status (status)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `article_version`;
CREATE TABLE `article_version`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    created_at DATETIME        NOT NULL,
    updated_at DATETIME        NOT NULL,
    deleted_at DATETIME                 DEFAULT NULL,

    article_id BIGINT UNSIGNED NOT NULL,
    version_no INT UNSIGNED    NOT NULL,

    title      VARCHAR(255)    NOT NULL DEFAULT '',
    content    TEXT            NOT NULL,
    summary    VARCHAR(500)    NOT NULL DEFAULT '',
    change_log VARCHAR(500)             DEFAULT '',

    PRIMARY KEY (id),
    UNIQUE KEY uk_article_version (version_no, article_id),
    INDEX article_id (article_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `article_publish`;
CREATE TABLE `article_publish`
(
    id           BIGINT UNSIGNED AUTO_INCREMENT,
    created_at   DATETIME        NOT NULL,
    updated_at   DATETIME        NOT NULL,
    deleted_at   DATETIME                 DEFAULT NULL,

    publish_at   DATETIME        NOT NULL,

    article_id   BIGINT UNSIGNED NOT NULL,
    version_id   BIGINT UNSIGNED NOT NULL,
    publish_type TINYINT         NOT NULL DEFAULT 0,
    remark       VARCHAR(255)             DEFAULT '',

    PRIMARY KEY (id),
    INDEX article_id (article_id),
    INDEX version_id (version_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `article_category_rel`;
CREATE TABLE `article_category_rel`
(
    article_id  BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (article_id, category_id),
    INDEX category_id (category_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `article_tag`;
CREATE TABLE `article_tag`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    name       VARCHAR(64) NOT NULL,
    created_at DATETIME    NOT NULL,
    updated_at DATETIME    NOT NULL,
    deleted_at DATETIME DEFAULT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY uk_article_tag_name (name)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `article_tag_rel`;
CREATE TABLE `article_tag_rel`
(
    article_id BIGINT UNSIGNED NOT NULL,
    tag_id     BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (article_id, tag_id),
    INDEX tag_id (tag_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `article_comment`;
CREATE TABLE `article_comment`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    created_at DATETIME        NOT NULL,
    updated_at DATETIME        NOT NULL,
    deleted_at DATETIME        DEFAULT NULL,

    parent_id  BIGINT UNSIGNED DEFAULT NULL,
    article_id BIGINT UNSIGNED NOT NULL,
    author_id  BIGINT UNSIGNED NOT NULL,

    content    TEXT            NOT NULL,

    PRIMARY KEY (id),
    INDEX author_id (author_id),
    INDEX article_id (article_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `content_type`;
CREATE TABLE `content_type`
(
    id          BIGINT UNSIGNED AUTO_INCREMENT,
    created_at  DATETIME     NOT NULL,
    updated_at  DATETIME     NOT NULL,
    deleted_at  DATETIME DEFAULT NULL,

    name        VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,

    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `content_type_entries`;
CREATE TABLE `content_type_entries`
(
    id              BIGINT UNSIGNED AUTO_INCREMENT,
    created_at      DATETIME        NOT NULL,
    updated_at      DATETIME        NOT NULL,
    deleted_at      DATETIME DEFAULT NULL,

    article_id      BIGINT UNSIGNED NOT NULL,
    content_type_id BIGINT UNSIGNED NOT NULL,
    data            JSON     DEFAULT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY uk_content_entries_article_type (article_id, content_type_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `content_type_schema`;
CREATE TABLE `content_type_schema`
(
    id              BIGINT UNSIGNED AUTO_INCREMENT,
    created_at      DATETIME        NOT NULL,
    updated_at      DATETIME        NOT NULL,
    deleted_at      DATETIME                 DEFAULT NULL,

    content_type_id BIGINT UNSIGNED NOT NULL,
    field_key       VARCHAR(255)    NOT NULL,
    field_name      VARCHAR(255)    NOT NULL,
    description     VARCHAR(255)    NOT NULL DEFAULT '',
    min_len         BIGINT UNSIGNED          DEFAULT 0,
    max_len         BIGINT UNSIGNED          DEFAULT 0,
    min_value       DECIMAL(10, 2)           DEFAULT NULL,
    max_value       DECIMAL(10, 2)           DEFAULT NULL,
    min_time        DATETIME                 DEFAULT NULL,
    max_time        DATETIME                 DEFAULT NULL,
    pattern         VARCHAR(255)             DEFAULT NULL,
    sequence        BIGINT          NOT NULL DEFAULT 0,
    type            SMALLINT        NOT NULL,
    enum_options    TEXT,
    required        TINYINT                  DEFAULT 0,
    visible         TINYINT                  DEFAULT 1,
    enable          TINYINT                  DEFAULT 1,

    PRIMARY KEY (id),
    UNIQUE KEY idx_field_key (content_type_id, field_key),
    INDEX sequence (sequence)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `content_field_values`;
CREATE TABLE `content_field_values`
(
    id           BIGINT UNSIGNED AUTO_INCREMENT,
    created_at   DATETIME        NOT NULL,
    updated_at   DATETIME        NOT NULL,
    deleted_at   DATETIME       DEFAULT NULL,

    model_id     BIGINT UNSIGNED NOT NULL,
    article_id   BIGINT UNSIGNED NOT NULL,
    field_key    VARCHAR(255)    NOT NULL,
    type         SMALLINT        NOT NULL,
    string_value VARCHAR(255)   DEFAULT NULL,
    bool_value   BOOLEAN        DEFAULT NULL,
    number_value DECIMAL(10, 2) DEFAULT NULL,
    time_value   DATETIME       DEFAULT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY uk_content_field_article_key (model_id, article_id, field_key)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    created_at DATETIME        NOT NULL,
    updated_at DATETIME        NOT NULL,
    deleted_at DATETIME DEFAULT NULL,

    parent_id  BIGINT UNSIGNED NOT NULL,
    sequence   BIGINT          NOT NULL,
    name       VARCHAR(255)    NOT NULL,
    path       VARCHAR(64)     NOT NULL,
    visible    TINYINT         NOT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY path (path),
    INDEX parent_id (parent_id),
    INDEX sequence (sequence)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `category_seo`;
CREATE TABLE `category_seo`
(
    id          BIGINT UNSIGNED AUTO_INCREMENT,
    created_at  DATETIME        NOT NULL,
    updated_at  DATETIME        NOT NULL,
    deleted_at  DATETIME DEFAULT NULL,

    category_id BIGINT UNSIGNED NOT NULL,
    title       VARCHAR(255)    NOT NULL,
    keywords    VARCHAR(255)    NOT NULL,
    description VARCHAR(255)    NOT NULL,

    PRIMARY KEY (id),
    INDEX category_id (category_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `category_context`;
CREATE TABLE `category_context`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    ancestor   BIGINT UNSIGNED NOT NULL,
    descendant BIGINT UNSIGNED NOT NULL,
    distance   BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (id),
    INDEX idx_category_context (ancestor, descendant)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `single_page`;
CREATE TABLE `single_page`
(
    id          BIGINT UNSIGNED AUTO_INCREMENT,
    created_at  DATETIME NOT NULL,
    updated_at  DATETIME NOT NULL,
    deleted_at  DATETIME DEFAULT NULL,

    name        VARCHAR(255),
    keywords    VARCHAR(255),
    description VARCHAR(500),
    content     TEXT,

    PRIMARY KEY (id),
    INDEX idx_name (name)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    created_at DATETIME        NOT NULL,
    updated_at DATETIME        NOT NULL,
    deleted_at DATETIME                 DEFAULT NULL,

    parent_id  BIGINT UNSIGNED NOT NULL,
    type       TINYINT         NOT NULL,
    name       VARCHAR(64)     NOT NULL,
    affix      TINYINT         NOT NULL DEFAULT 0,
    icon       VARCHAR(64)              DEFAULT NULL,
    url        VARCHAR(255)             DEFAULT NULL,
    sequence   BIGINT          NOT NULL DEFAULT 0,
    visible    TINYINT         NOT NULL DEFAULT 1,
    path       VARCHAR(255)    NOT NULL,
    template   VARCHAR(255)    NOT NULL,
    remark     VARCHAR(255)    NOT NULL DEFAULT '',

    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `menu_context`;
CREATE TABLE `menu_context`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    ancestor   BIGINT UNSIGNED NOT NULL,
    descendant BIGINT UNSIGNED NOT NULL,
    distance   BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (id),
    INDEX idx_menu_context (ancestor, descendant)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    id          BIGINT UNSIGNED AUTO_INCREMENT,
    created_at  DATETIME     NOT NULL,
    updated_at  DATETIME     NOT NULL,
    deleted_at  DATETIME              DEFAULT NULL,

    username    VARCHAR(255) NOT NULL,
    password    VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    verified_at DATETIME              DEFAULT NULL,
    ip          VARCHAR(255) NOT NULL DEFAULT '',
    status      TINYINT      NOT NULL DEFAULT 0,
    role_id     BIGINT UNSIGNED       DEFAULT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY uniq_user_username (username),
    UNIQUE KEY uniq_user_email (email),
    INDEX role_id (role_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `user_profile`;
CREATE TABLE `user_profile`
(
    id          BIGINT UNSIGNED AUTO_INCREMENT,
    created_at  DATETIME        NOT NULL,
    updated_at  DATETIME        NOT NULL,
    deleted_at  DATETIME     DEFAULT NULL,

    avatar      VARCHAR(255) DEFAULT NULL,
    user_id     BIGINT UNSIGNED NOT NULL,
    nickname    VARCHAR(255) DEFAULT NULL,
    gender      TINYINT      DEFAULT NULL COMMENT '0:male,1:female',
    description VARCHAR(255) DEFAULT NULL,
    country     VARCHAR(255) DEFAULT NULL,
    province    VARCHAR(255) DEFAULT NULL,
    city        VARCHAR(255) DEFAULT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY uniq_user_profile_user_id (user_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`
(
    id          BIGINT UNSIGNED AUTO_INCREMENT,
    created_at  DATETIME     NOT NULL,
    updated_at  DATETIME     NOT NULL,
    deleted_at  DATETIME              DEFAULT NULL,

    name        VARCHAR(64)  NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',

    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `file`;
CREATE TABLE `file`
(
    id            BIGINT UNSIGNED AUTO_INCREMENT,
    created_at    DATETIME        NOT NULL,
    updated_at    DATETIME        NOT NULL,
    deleted_at    DATETIME                 DEFAULT NULL,

    user_id       BIGINT UNSIGNED NOT NULL,
    original_name VARCHAR(255)    NOT NULL,
    ext           VARCHAR(32)     NOT NULL,
    path          VARCHAR(255)    NOT NULL,
    size          BIGINT          NOT NULL,
    driver        VARCHAR(32)              DEFAULT 'local',
    sha256        VARCHAR(64)     NOT NULL,
    is_image      TINYINT(1)      NOT NULL DEFAULT 0,

    PRIMARY KEY (id),
    INDEX idx_file_user_id (user_id),
    INDEX idx_file_path (path),
    INDEX idx_file_driver (driver),
    INDEX idx_file_uniq (sha256, size),
    INDEX idx_file_is_image (is_image)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission`
(
    id          BIGINT UNSIGNED AUTO_INCREMENT,
    created_at  DATETIME     NOT NULL,
    updated_at  DATETIME     NOT NULL,
    deleted_at  DATETIME        DEFAULT NULL,

    menu_id     BIGINT UNSIGNED DEFAULT NULL,
    name        VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    resource    VARCHAR(255) NOT NULL,
    action      VARCHAR(64)  NOT NULL,

    PRIMARY KEY (id),
    INDEX idx_perm_menu (menu_id),
    INDEX idx_perm_resource (resource, action)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `config`;
CREATE TABLE `config`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    created_at DATETIME     NOT NULL,
    updated_at DATETIME     NOT NULL,

    field      VARCHAR(64)  NOT NULL,
    value      VARCHAR(255) NOT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY uniq_field (field)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET utf8mb4;