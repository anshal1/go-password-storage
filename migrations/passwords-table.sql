create table if not exists passwords (
    id serial not null primary key,
    password varchar(255) not null,
    domain varchar(255) not null unique,
    userId int not null references users(id)
)
