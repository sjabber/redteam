create table public.project_info
(
    p_no          serial not null
        constraint project_info_pk
            primary key,
    tml_no        integer
        constraint project_info_template_info_tmp_no_fk
            references public.template_info,
    tag_no        integer
        constraint project_info_tag_info_tag_no_fk
            references public.tag_info,
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
