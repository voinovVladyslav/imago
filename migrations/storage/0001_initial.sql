create table if not exists file_registry (
    id text primary key,
    filepath text unique not null,
    created_at integer not null
);

