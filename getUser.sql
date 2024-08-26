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

SELECT * FROM users where id != ? ORDER BY name; --

SELECT u.name, ur.room_id from user_rooms ur JOIN users u ON ur.user_id = u.id where u.id = ? ORDER BY room_id;


WITH UserRooms AS (
    SELECT room_id
    FROM user_rooms
    WHERE user_id = ?
)
SELECT u.name, u.id, ur.room_id
FROM user_rooms ur
         JOIN users u ON ur.user_id = u.id
WHERE ur.room_id IN (SELECT room_id FROM UserRooms) AND u.id != ? ORDER BY room_id; --

SELECT u.name, m.message, m.timestamp FROM messages m JOIN users u on m.user_id = u.id where room_id = ? ORDER BY timestamp;

INSERT INTO messages(id, timestamp, room_id, message, user_id) VALUES (DEFAULT, CURRENT_TIMESTAMP, ?, ?, ?);

