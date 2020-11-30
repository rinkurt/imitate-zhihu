create table zhihu.questions
(
    id            bigint auto_increment
        primary key,
    title         varchar(50)   not null,
    content       text          not null,
    creator_id    int           not null,
    tag           varchar(256)  null,
    answer_count  int default 0 null,
    comment_count int default 0 null,
    view_count    int default 0 null,
    like_count    int default 0 null,
    create_at     bigint        null,
    update_at     bigint        null
);

create index questions_creator_id_index
    on zhihu.questions (creator_id);

