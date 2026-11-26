create table usr (
    id          bigserial primary key,
    created_at  timestamptz not null default now(),
    modified_at timestamptz not null default now(),
    email       text        not null unique,
    name        text        not null default '',
    campus      text        not null default '',
    role        text        not null default 'student' check (role in ('student', 'moderator', 'admin')),
    blocked     boolean     not null default false,
    telegram    text        not null default ''
);

create table listing (
    id             bigserial primary key,
    created_at     timestamptz not null default now(),
    modified_at    timestamptz not null default now(),
    owner_id       bigint      not null references usr (id),
    type           text        not null check (type in ('rent', 'exchange', 'sale')),
    status         text        not null default 'active' check (status in ('active', 'paused', 'archived')),
    title          text        not null,
    description    text        not null default '',
    category       text        not null,
    campus         text        not null default '',
    photos         text[]      not null default '{}',
    price          bigint      not null default 0 check (price >= 0),
    deposit        bigint      not null default 0 check (deposit >= 0),
    wanted         text        not null default '',
    available_from date,
    available_to   date,
    -- listing type rules from the domain requirements are enforced in the DB too
    check (type <> 'rent' or (price > 0 and available_from is not null and available_to >= available_from)),
    check (type <> 'sale' or price > 0),
    check (type <> 'exchange' or (wanted <> '' and price = 0))
);

create index listing_owner_id_idx on listing (owner_id);
create index listing_search_idx on listing (status, type, category, created_at desc);
