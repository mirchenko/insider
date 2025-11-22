create type messages_type as enum ('pending', 'sent', 'delivered', 'failed');

create table messages (
  id serial primary key,
  status messages_type default 'pending',
  content varchar(1000) not null,
  phone_number varchar(16) not null,
  reason text,
  external_message_id varchar(255),
  created_at timestamp default now(),
  updated_at timestamp default now(),
  deleted_at timestamp
);


create index idx_messages_status on messages(status);
create index idx_messages_created_at on messages(created_at);
