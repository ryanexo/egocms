CREATE TABLE `article`
(
    id                   BIGINT UNSIGNED AUTO_INCREMENT,
    created_at           DATETIME        NOT NULL,
    deleted_at           DATETIME                 DEFAULT NULL,
    updated_at           DATETIME        NOT NULL,

    current_version_id   BIGINT UNSIGNED NOT NULL,
    published_version_id BIGINT UNSIGNED NOT NULL,
    content_type_id      BIGINT UNSIGNED          DEFAULT NULL,
    author_id            BIGINT UNSIGNED NOT NULL,

    url                  VARCHAR(255)             DEFAULT NULL,
    slug                 VARCHAR(255)             DEFAULT NULL,
    title                VARCHAR(255)             DEFAULT NULL,
    summary              VARCHAR(500)             DEFAULT '',
    status               tinyint         NOT NULL DEFAULT 0,

    PRIMARY KEY (id),
    INDEX author_id (author_id),
    INDEX slug (slug),
    INDEX content_type_id (content_type_id),
    INDEX status (status)
);

CREATE TABLE `article_version`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    created_at DATETIME        NOT NULL,
    updated_at DATETIME        NOT NULL,
    deleted_at DATETIME                 DEFAULT NULL,

    article_id BIGINT UNSIGNED NOT NULL,
    version_no int,

    title      VARCHAR(255)    NOT NULL DEFAULT '',
    content    TEXT            NOT NULL,
    summary    VARCHAR(500)    NOT NULL DEFAULT '',
    change_log VARCHAR(500)             DEFAULT '',
    published  BOOLEAN                  DEFAULT false,

    PRIMARY KEY (id),
    UNIQUE KEY uk_article_version (version_no, article_id),
    INDEX article_id (article_id)
);

CREATE TABLE `article_publish`
(
    id           BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    created_at   DATETIME        NOT NULL,
    updated_at   DATETIME        NOT NULL,
    deleted_at   DATETIME                 DEFAULT NULL,

    publish_at   DATETIME        NOT NULL,

    article_id   BIGINT UNSIGNED NOT NULL,
    version_id   BIGINT UNSIGNED NOT NULL,
    publish_type TINYINT         NOT NULL DEFAULT 0,
    Remark       VARCHAR(255)             DEFAULT '',

    PRIMARY KEY (id),
    INDEX article_id (article_id),
    INDEX version_id (version_id)
);

CREATE TABLE `article_category_rel`
(
    article_id  BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (article_id, category_id),
    INDEX category_id (category_id)
);

CREATE TABLE `article_tag`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE `article_tag_rel`
(
    article_id BIGINT UNSIGNED NOT NULL,
    tag_id     BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (article_id, tag_id),
    INDEX tag_id (tag_id)
);

CREATE TABLE `article_comment`
(
    id         BIGINT UNSIGNED AUTO_INCREMENT,
    created_at DATETIME        NOT NULL,
    updated_at DATETIME        NOT NULL,
    deleted_at DATETIME                 DEFAULT NULL,

    parent_id  BIGINT UNSIGNED          DEFAULT NULL,
    article_id BIGINT UNSIGNED NOT NULL,
    author_id  BIGINT UNSIGNED NOT NULL,

    content    TEXT            NOT NULL DEFAULT '',

    PRIMARY KEY (id),
    INDEX author_id (author_id),
    INDEX article_id (article_id)
);