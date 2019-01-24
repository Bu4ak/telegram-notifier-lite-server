create table if not exists users
(
  id         bigserial primary key,
  chat_id    bigint unique      NOT NULL,
  token      varchar(64) unique NOT NULL,
  created_at timestamp default NULL
);