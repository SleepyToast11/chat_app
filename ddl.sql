create sequence chat_rooms_id_seq
    as integer;

alter sequence chat_rooms_id_seq owner to chat;

create table users
(
    name text not null
        constraint users_pk_2
            unique,
    id   serial
        constraint users_pk
            primary key
);

alter table users
    owner to chat;

create table room_ids
(
    id serial
        constraint room_ids_pk
            primary key
);

alter table room_ids
    owner to chat;

create table messages
(
    id        integer default nextval('chat_rooms_id_seq'::regclass) not null
        constraint messages_pk
            primary key,
    timestamp timestamp                                              not null,
    room_id   integer                                                not null
        constraint messages_room_ids_id_fk
            references room_ids,
    message   text                                                   not null,
    user_id   integer                                                not null
        constraint messages___fk_user
            references users
);

alter table messages
    owner to chat;

alter sequence chat_rooms_id_seq owned by messages.id;

create table user_rooms
(
    id      serial
        constraint user_rooms_pk
            primary key,
    user_id integer not null
        constraint user_rooms_users_id_fk
            references users,
    room_id integer not null
        constraint user_rooms___fk_room_id
            references room_ids
);

alter table user_rooms
    owner to chat;

