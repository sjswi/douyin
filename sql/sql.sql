create table comments
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(3)     null,
    updated_at  datetime(3)     null,
    deleted_at  datetime(3)     null,
    content     longtext        null,
    create_time datetime(3)     null,
    video_id    bigint unsigned null,
    user_id     bigint unsigned null
);

create index idx_comments_deleted_at
    on comments (deleted_at);

create table favorites
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    user_id    bigint unsigned null,
    video_id   bigint unsigned null,
    exist      tinyint(1)      null
);

create index idx_favorites_deleted_at
    on favorites (deleted_at);

create table messages
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(3)     null,
    updated_at  datetime(3)     null,
    deleted_at  datetime(3)     null,
    content     longtext        null,
    user_id     bigint unsigned null,
    target_id   bigint unsigned null,
    create_time datetime(3)     null
);

create index idx_messages_deleted_at
    on messages (deleted_at);

create table relations
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    user_id    bigint unsigned null,
    target_id  bigint unsigned null,
    type       bigint          null
);

create index idx_relations_deleted_at
    on relations (deleted_at);

create table users
(
    id             bigint unsigned auto_increment
        primary key,
    created_at     datetime(3) null,
    updated_at     datetime(3) null,
    deleted_at     datetime(3) null,
    name           longtext    null,
    follow_count   bigint      null,
    follower_count bigint      null,
    password       longtext    null,
    salt           longtext    null
);

create index idx_users_deleted_at
    on users (deleted_at);

create table videos
(
    id             bigint unsigned auto_increment
        primary key,
    created_at     datetime(3)     null,
    updated_at     datetime(3)     null,
    deleted_at     datetime(3)     null,
    author_id      bigint unsigned null,
    title          longtext        null,
    comment_count  bigint          null,
    favorite_count bigint          null,
    play_url       longtext        null,
    cover_url      longtext        null
);

create index idx_videos_deleted_at
    on videos (deleted_at);

