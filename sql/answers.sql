create table answers
(
    id int primary key auto_increment,
    answer text null,
    view_count int default 0 null,
    upvote_count int default 0 null,
    downvote_count int default 0 null,
    comment_count int default 0 null,
    gmt_create bigint null,
    gmt_modified bigint null,
    creator_id int null,
    question_id int null,
    constraint answers_questions_id_fk
        foreign key (question_id) references questions (id),
    constraint answers_users_id_fk
        foreign key (creator_id) references users (id)
);

create unique index answers_id_uindex
    on answers (id);


