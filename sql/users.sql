create table users
(
	id int auto_increment primary key,
	name varchar(30) null,
	email varchar(30) null,
	password varchar(100) null,
	token varchar(40) null,
	gmt_create bigint null,
	bio varchar(256) null,
	avatar_url varchar(100) null
);

create unique index users_email_uindex
	on users (email);

create unique index users_id_uindex
	on users (id);


