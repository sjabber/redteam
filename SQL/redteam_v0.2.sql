create table if not exists target_info
(
	target_no int4 default nextval('target_info_target_no_seq'::regclass) not null
		constraint target_info_target_no_pk
			primary key,
	target_name text not null,
	target_email text not null,
	target_phone text,
	target_organize text,
	target_position text,
	target_classify text,
	created_time timestamp(6) default now(),
	modified_time timestamp(6) default now()
);

comment on column target_info.target_position is '//직급';

create table if not exists tag_info
(
	tag_no int4 default nextval('tag_info_tag_no_seq'::regclass) not null
		constraint tag_info_tag_no_pk
			primary key,
	tag_name text,
	target_no int4
		constraint tag_info_target_info_target_no_fk
			references target_info,
	created_t timestamp(6) default now(),
	modify_t timestamp(6) default now()
);

create unique index if not exists tag_info_tag_no_uindex
	on tag_info (tag_no);

create unique index if not exists target_info_target_no_uindex
	on target_info (target_no);

create table if not exists user_info
(
	user_no int4 default nextval('user_info_user_no_seq'::regclass) not null
		constraint user_info_user_no_pk
			primary key,
	user_email text not null
		constraint user_info_email_uindex
			unique,
	user_name text not null,
	user_pw text not null,
	created_time timestamp(6) default now(),
	modified_time timestamp(6) default now(),
	is_enabled int2 default 1,
	user_pw_hash varchar(255)
);

create table if not exists smtp_info
(
	smtp_no int4 default nextval('smtp_info_smtp_no_seq'::regclass) not null
		constraint smtp_info_smtp_no_pk
			primary key,
	user_no int4 not null
		constraint smtp_info_user_info_user_no_fk
			references user_info,
	smtp_host text default 'smtp.redteam.or.kr',
	smtp_port text default '587',
	protocol text default '1',
	tls text default '1',
	timeout text default '1000',
	smtp_id text,
	smtp_pw text,
	created_time timestamp(6) default now(),
	modify timestamp(6) default now()
);

create unique index if not exists smtp_info_smtp_no_uindex
	on smtp_info (smtp_no);

create table if not exists template_info
(
	tmp_no int4 default nextval('template_info_tmp_no_seq'::regclass) not null
		constraint template_info_pk
			primary key,
	user_no int4 not null
		constraint template_info_user_info_user_no_fk
			references user_info,
	tmp_division text,
	tmp_kind text,
	file_info text,
	tmp_name text,
	mail_title text,
	sender_name text,
	download_type text,
	created_time timestamp(6) default now(),
	modified_time timestamp(6) default now()
);

comment on column template_info.sender_name is 'email format';

create unique index if not exists template_info_tm_no_uindex
	on template_info (tmp_no);

create unique index if not exists user_info_user_no_uindex
	on user_info (user_no);

