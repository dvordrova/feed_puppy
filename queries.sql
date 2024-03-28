-- name: NewUser :exec
INSERT INTO users (telegram_id, name, language)
VALUES (?, ?, ?);

-- name: GetUser :one
SELECT * FROM users WHERE telegram_id = ?;

-- name: SetUserLanguage :exec
UPDATE users
SET language = ?
WHERE telegram_id = ?;

-- name: SetUserCurDog :exec
UPDATE users
SET current_dog = ?
WHERE id = ?;

-- name: SetUserState :exec
UPDATE users
SET state = ?
WHERE id = ?;

-- name: NewDog :one
INSERT INTO dogs (name, birth_date, sex, breed)
VALUES (?, ?, ?, ?)
RETURNING id;

-- name: GetDog :one
SELECT * FROM dogs WHERE id = ?;

-- name: NewAction :exec
INSERT INTO actions (user_id, dog_id, timestamp, action_id, additional_info)
VALUES (?, ?, ?, ?, ?);

-- name: NewSubscription :exec
INSERT INTO subscriptions (dog_id, user_id, type)
VALUES (?, ?, ?);

-- name: ClearSubscriptions :exec
DELETE FROM subscriptions
WHERE user_id = ? AND dog_id = ? and type = ?;

-- name: SelectDogsSubscribers :many
SELECT * FROM users
where id in (
    select user_id from subscriptions
    where dog_id = ?
);

-- name: SelectUsersDogs :many
Select * FROM dogs
where id in (
    select dog_id from subscriptions
    where type = 'owner' and user_id = ?
);

-- name: NewInvite :exec
INSERT INTO invite_subscriptions (hash, dog_id, type)
VALUES (?, ?, ?);

-- name: GetInvite :one
SELECT * FROM invite_subscriptions
WHERE used = false AND hash = ?;

-- name: MarkInviteUsed :exec
UPDATE invite_subscriptions
SET used = true
where id = ?;
