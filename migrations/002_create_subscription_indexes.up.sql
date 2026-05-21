create index idx_subscriptions_user_service on subscriptions(user_id, service_name);
create index idx_subscriptions_start_date on subscriptions(start_date);