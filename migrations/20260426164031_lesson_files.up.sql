create table lesson_files
(
    id           uuid primary key,
    file_name    text                         not null,
    object_name  text                         not null,
    content_type text,
    size         bigint,
    lesson_id    uuid references lessons (id) not null
)