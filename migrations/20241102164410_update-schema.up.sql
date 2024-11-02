drop table urls;

create table if not exists urls (
    short_code varchar(50) primary key, 
    long_url varchar(8000) not null, 
    alias varchar(20) null, 
    expiration timestamp with time zone null,
    created_by varchar(30) null
);
