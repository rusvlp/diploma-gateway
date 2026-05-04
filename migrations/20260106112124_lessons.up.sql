create table lessons
(
    id          uuid primary key,
    title       varchar(100)                 not null unique,
    description text,
    position    int                          not null,
    module_id   uuid references modules (id) not null,
    deleted_at  timestamp default null
);

create unique index uniq_lessons_module_position
    on lessons (module_id, position)
    where deleted_at is null;