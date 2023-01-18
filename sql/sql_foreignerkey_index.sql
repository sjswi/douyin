

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
    create_time datetime(3)     null,
    constraint messages_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade,
    constraint messages_users_id_fk2
        foreign key (target_id) references users (id)
            on delete cascade
);

create index idx_messages_deleted_at
    on messages (deleted_at);

create index messages_target_id_index
    on messages (target_id);

create index messages_user_id_target_id_index
    on messages (user_id, target_id);

create table relations
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    user_id    bigint unsigned null,
    target_id  bigint unsigned null,
    type       bigint          null,
    exist      tinyint(1)      null,
    constraint relations_users_id_fk
        foreign key (target_id) references users (id)
            on delete cascade,
    constraint relations_users_id_fk2
        foreign key (user_id) references users (id)
            on delete cascade
);

create index idx_relations_deleted_at
    on relations (deleted_at);

create index relations_target_id_index
    on relations (target_id);

create index relations_user_id_target_id_index
    on relations (user_id, target_id);

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
    cover_url      longtext        null,
    constraint videos_users_id_fk
        foreign key (author_id) references users (id)
            on delete cascade
);

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
    user_id     bigint unsigned null,
    constraint comments_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade,
    constraint comments_videos_id_fk
        foreign key (video_id) references videos (id)
            on delete cascade
);

create index comments_user_id_index
    on comments (user_id)
    comment '用户id';

create index comments_video_id_index
    on comments (video_id)
    comment '视频id';

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
    exist      tinyint(1)      null,
    constraint favorites_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade,
    constraint favorites_videos_id_fk
        foreign key (video_id) references videos (id)
            on delete cascade
);

create index favorites_exist_index
    on favorites (exist);

create index favorites_user_id_video_id_index
    on favorites (user_id, video_id);

create index favorites_video_id_index
    on favorites (video_id);

create index idx_favorites_deleted_at
    on favorites (deleted_at);

create index idx_videos_deleted_at
    on videos (deleted_at);

create index videos_author_id_index
    on videos (author_id);

