create table answer_votes
(
    id        bigint auto_increment
        primary key,
    is_upvote tinyint(1) null,
    answer_id bigint     null,
    user_id   bigint     null,
    update_at bigint     null
);