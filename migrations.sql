begin;
create table if not exists users(
    id serial primary key,
    fullname text,
    passwordhash text not null,
    email text unique not null
);

create table if not exists files{
    id serial primary key
    email text not null 
}

commit;
