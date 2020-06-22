CREATE EXTENSION pg_trgm;
CREATE EXTENSION pgcrypto;
CREATE EXTENSION btree_gin;

create table user_info
(
	user_no serial not null
		constraint user_info_user_no_pk
			primary key,
	user_id text not null,
	user_name text not null,
	user_pw text,
	created_time timestamp default now(),
	modify_time timestamp default now(),
	is_enabled smallint default 1
);

alter table user_info owner to redteamadmin;

create unique index user_info_email_uindex
	on user_info (user_id);

create unique index user_info_user_no_uindex
	on user_info (user_no);

create table smtp_info
(
	smtp_no serial not null
		constraint smtp_info_smtp_no_pk
			primary key,
	user_no integer not null
		constraint smtp_info_user_info_user_no_fk
			references user_info,
	smtp_host text default 'smtp.redteam.or.kr'::text,
	smtp_port text default '587'::text,
	protocol text default '1'::text,
	tls text default '1'::text,
	timeout text default '1000'::text,
	smtp_id text,
	smtp_pw text,
	created_time timestamp default now(),
	modify timestamp default now()
);

alter table smtp_info owner to redteamadmin;

create unique index smtp_info_smtp_no_uindex
	on smtp_info (smtp_no);

create table target_info
(
	target_no serial not null
		constraint target_info_target_no_pk
			primary key,
	user_no integer not null
		constraint target_info_user_info_user_no_fk
			references user_info,
	target_name text not null,
	target_email text not null,
	target_phone text,
	target_organize text,
	target_position text,
	target_classify text,
	created_time timestamp default now(),
	modified_time timestamp default now()
);

comment on column target_info.target_position is '//직급';

alter table target_info owner to redteamadmin;

create unique index target_info_target_no_uindex
	on target_info (target_no);

-- CREATE INDEX target_info_target_classify_idx_gin ON target_info USING gin (target_classify gin_trgm_ops);
CREATE INDEX target_info_tag_no_idx ON target_info USING btree (tag_no);
