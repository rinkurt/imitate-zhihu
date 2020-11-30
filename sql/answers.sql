create table zhihu.answers
(
    id             bigint auto_increment
        primary key,
    content        text          null,
    view_count     int default 0 null,
    upvote_count   int default 0 null,
    downvote_count int default 0 null,
    comment_count  int default 0 null,
    create_at      bigint        null,
    update_at      bigint        null,
    creator_id     int           null,
    question_id    int           null
);

create index answers_creator_id_index
    on zhihu.answers (creator_id);

create index answers_question_id_index
    on zhihu.answers (question_id);

