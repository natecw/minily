create extension if not exists pgcrypto with schema public;

create table if not exists urls (
    short varchar(50) not null, 
    long varchar(8000) not null, 
    alias varchar(20) null, 
    expiration timestamp with time zone null
);

