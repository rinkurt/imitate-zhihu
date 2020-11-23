create table questions
(
	id int primary key auto_increment,
	title varchar(50) not null,
	content text not null,
	creator_id int not null,
	tag varchar(256) null,
	answer_count int default 0 null,
	comment_count int default 0 null,
	view_count int default 0 null,
	like_count int default 0 null,
	create_at bigint null,
	update_at bigint null
);

create unique index questions_id_uindex
	on questions (id);

alter table questions
    add constraint questions_users_id_fk
        foreign key (creator_id) references users (id);

