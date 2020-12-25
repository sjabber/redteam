create table project_info
(
    p_no          serial not null
        constraint project_info_pk
            primary key,
    user_no        integer
        constraint project_info_template_info_user_no_fk
            references user_info,
    tml_no        integer
        constraint project_info_template_info_tmp_no_fk
            references template_info,
    tag1_no        integer
        constraint project_info_tag_info_tag1_no_fk
            references tag_info,
    tag2_no        integer
        constraint project_info_tag_info_tag2_no_fk
            references tag_info,
    tag3_no        integer
        constraint project_info_tag_info_tag3_no_fk
            references tag_info,
    p_name        text,
    p_description text,
    p_start_date  timestamp,
    p_end_date    timestamp,
    created_time  timestamp default now(),
    modified_time timestamp default now()
);

alter table public.project_info
    owner to postgres;

create unique index project_info_p_no_uindex
    on public.project_info (p_no);
