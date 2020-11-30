create table zhihu.users
(
    id        bigint auto_increment
        primary key,
    email     varchar(30)  null,
    password  varchar(100) null,
    create_at bigint       null,
    update_at bigint       null,
    constraint users_email_uindex
        unique (email)
);

