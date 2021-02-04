create table profiles
(
    id          bigint auto_increment
        primary key,
    user_id     bigint       not null,
    name        varchar(30)  null,
    gender      int          null,
    description varchar(256) null,
    avatar_url  varchar(256) null,
    create_at   bigint       null,
    update_at   bigint       null,
    constraint profiles_user_id_uindex
        unique (user_id)
);

