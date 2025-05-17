-- User role enum
create type user_role as enum ('user', 'admin');
-- User schema
create table if not exists "user" (
  id serial primary key,
  email varchar(255) unique,
  password text not null,
  role user_role default 'user',
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  deleted_at timestamp
);
-- RefreshToken schema
create table if not exists refresh_token (
  id serial primary key,
  user_id int not null references "user"(id),
  token text not null unique
);
-- Post schema
create table if not exists post (
  id serial primary key,
  user_id int not null references "user"(id),
  content text not null,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  deleted_at timestamp
);
-- Comment schema
create table if not exists comment (
  id serial primary key,
  user_id int not null references "user"(id),
  post_id int not null references post(id),
  content text not null,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  deleted_at timestamp
);
