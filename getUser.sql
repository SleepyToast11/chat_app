WITH inserted AS (
    INSERT INTO users (name)
        VALUES (?)
        ON CONFLICT (name) DO NOTHING
        RETURNING id
)
SELECT id
FROM inserted
UNION ALL
SELECT id
FROM users
WHERE name = ?;

SELECT * FROM users;

SELECT u.name, ur.room_id from user_rooms ur JOIN users u ON ur.user_id = u.id where u.id = ?;


WITH UserRooms AS (
    SELECT room_id
    FROM user_rooms
    WHERE user_id = ?
)
SELECT u.name, ur.room_id
FROM user_rooms ur
         JOIN users u ON ur.user_id = u.id
WHERE ur.room_id IN (SELECT room_id FROM UserRooms) AND u.id != ?;

SELECT u.name, m.message, m.timestamp FROM messages m JOIN users u on m.user_id = u.id where room_id = ?;

INSERT INTO messages VALUES (DEFAULT, CURRENT_TIMESTAMP, ?, ?, ?);

