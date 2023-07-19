create table if not exists orders
(
    uuid       uuid not null
        constraint order_pk
            primary key,
    user_uuid  uuid,
    created_at timestamp default current_timestamp
);

create table if not exists order_item
(
    order_uuid   uuid
        constraint order_item_order_uuid_fk
            references orders
            on update cascade on delete cascade,
    count        integer,
    product_uuid uuid
);