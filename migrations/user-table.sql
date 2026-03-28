create table if not exists users (
    username varchar(255) not null unique,
    password varchar(255) not null,
    id serial not null primary key
);
