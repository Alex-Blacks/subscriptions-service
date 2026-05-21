create table subscriptions(
    id bigint generated always as identity primary key,
    service_name varchar(100) not null,
    price integer not null check(price > 0),
    user_id uuid not null,
    start_date date not null default current_date,
    end_date date, 
    created_at timestamptz not null default now(),
    
    check (
        end_date is null
        or end_date >= start_date
    )
);