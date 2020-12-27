create table if not exists project_info
(
	p_no serial not null
		constraint project_info_pk
			primary key,
	user_no integer
		constraint project_info_template_info_user_no_fk
			references user_info,
	tmp_no integer
		constraint project_info_template_info_tmp_no_fk
			references template_info,
	tag_1 integer
		constraint project_info_tag_info_tag1_no_fk
			references tag_info,
	tag_2 integer
		constraint project_info_tag_info_tag2_no_fk
			references tag_info,
	tag_3 integer
		constraint project_info_tag_info_tag3_no_fk
			references tag_info,
	p_name text,
	p_description text,
	start_date date,
	end_date date,
	created_time timestamp default now(),
	modified_time timestamp default now()
);

alter table project_info owner to postgres;

create unique index if not exists project_info_p_no_uindex
	on project_info (p_no);

