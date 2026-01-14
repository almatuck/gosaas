-- name: GetSubscriptionByUserID :one
SELECT id, user_id, stripe_customer_id, stripe_subscription_id, plan_id, status,
       current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
FROM subscriptions
WHERE user_id = sqlc.arg(user_id)
LIMIT 1;

-- name: GetSubscriptionByStripeSubID :one
SELECT id, user_id, stripe_customer_id, stripe_subscription_id, plan_id, status,
       current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
FROM subscriptions
WHERE stripe_subscription_id = sqlc.arg(stripe_subscription_id)
LIMIT 1;

-- name: GetSubscriptionByStripeCustomerID :one
SELECT id, user_id, stripe_customer_id, stripe_subscription_id, plan_id, status,
       current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
FROM subscriptions
WHERE stripe_customer_id = sqlc.arg(stripe_customer_id)
LIMIT 1;

-- name: CreateSubscription :one
INSERT INTO subscriptions (
    id, user_id, plan_id, status, created_at, updated_at
) VALUES (
    sqlc.arg(id), sqlc.arg(user_id), sqlc.arg(plan_id), sqlc.arg(status),
    strftime('%s', 'now'), strftime('%s', 'now')
)
RETURNING id, user_id, stripe_customer_id, stripe_subscription_id, plan_id, status,
          current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at;

-- name: UpdateSubscriptionStripeIDs :exec
UPDATE subscriptions
SET stripe_customer_id = sqlc.arg(stripe_customer_id),
    stripe_subscription_id = sqlc.arg(stripe_subscription_id),
    plan_id = sqlc.arg(plan_id),
    status = sqlc.arg(status),
    updated_at = strftime('%s', 'now')
WHERE user_id = sqlc.arg(user_id);

-- name: UpdateSubscriptionStatus :exec
UPDATE subscriptions
SET status = sqlc.arg(status),
    current_period_start = sqlc.arg(period_start),
    current_period_end = sqlc.arg(period_end),
    cancel_at_period_end = sqlc.arg(cancel_at_period_end),
    updated_at = strftime('%s', 'now')
WHERE stripe_subscription_id = sqlc.arg(stripe_subscription_id);

-- name: CancelSubscription :exec
UPDATE subscriptions
SET status = sqlc.arg(status),
    cancel_at_period_end = sqlc.arg(cancel_at_period_end),
    updated_at = strftime('%s', 'now')
WHERE user_id = sqlc.arg(user_id);

-- name: DowngradeToFree :exec
UPDATE subscriptions
SET plan_id = 'free',
    status = 'active',
    stripe_subscription_id = NULL,
    cancel_at_period_end = 0,
    updated_at = strftime('%s', 'now')
WHERE stripe_subscription_id = sqlc.arg(stripe_subscription_id);

-- name: DeleteSubscription :exec
DELETE FROM subscriptions WHERE user_id = sqlc.arg(user_id);
