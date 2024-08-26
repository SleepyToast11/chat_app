package data_accessor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"sort"
)

/*
*
* For Better query format see GetUser.sql file!
*
 */

type DBVar struct {
	User         string `json:"userId"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Database     string `json:"database"`
	DatabaseType string `json:""`
}

type User struct {
	Name string `json:"userName"`
	Id   int    `json:"userId"`
}

type Users struct {
	Users []User `json:"users"`
}

type Message struct {
	Message string `json:"message"`
	Name    string `json:"userName"`
	Time    string `json:"time"`
}

type Room struct {
	Id    int    `json:"userId"`
	Users []User `json:"users"`
}

type Rooms struct {
	Rooms []Room `json:"rooms"`
}

const secretFilepath = "secret.json"

func getDBVar() (*DBVar, error) {
	file, err := os.Open("secret.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var dbVar DBVar
	err = decoder.Decode(&dbVar)
	if err != nil {
		return nil, err
	}
	return &dbVar, nil
}

func connect(dbVar DBVar) (*sql.DB, error) {
	connStr := fmt.Sprintf("userId=%s password=%s host=%s dbname=%s sslmode=disable", dbVar.User, dbVar.Password, dbVar.Host, dbVar.Database)
	db, err := sql.Open(dbVar.DatabaseType, connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetUserId(name string, db sql.DB) (int64, error) {
	//tries to insert a userId and return it or return the already created userId
	rows, err := db.Query("WITH inserted AS (INSERT INTO users (name)VALUES (?) ON CONFLICT (name) DO NOTHING RETURNING id) SELECT id FROM inserted UNION ALL SELECT id FROM users WHERE name = ?;", name, name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, fmt.Errorf("User %s not found", name)
	}
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getUserList(db sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

type iRoom struct {
	userId   int
	userName string
	roomId   int
}

/*
its possible to get the room id of only the user and then make
another func that gets that room and return the room filled with
user, but this requires multiple hit to the database as
1 for the user room and 1 for every room that this user is part of
,so might put undue stress on the db
*/
func getUserRooms(db sql.DB, userId int) ([]Room, error) {
	rows, err := db.Query("WITH UserRooms AS ( SELECT room_id FROM user_rooms WHERE user_id = ?  )  SELECT u.id, u.name, ur.room_id  FROM user_rooms ur JOIN users u ON ur.user_id = u.id  WHERE ur.room_id IN (SELECT room_id FROM UserRooms) AND u.id != ? ORDER BY room_id;", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var irooms []iRoom
	for rows.Next() {
		var temp iRoom
		err := rows.Scan(&temp)
		if err != nil {
			return nil, err
		}
		irooms = append(irooms, temp)
	}
	var roomMap map[int]Room
	for _, v := range irooms {
		_, ok := roomMap[v.roomId]
		if !ok {
			var user = User{v.userName, v.userId}
			var users = []User{user}
			roomMap[v.roomId] = Room{v.roomId, users}
		} else {
			tempUsers := roomMap[v.roomId].Users
			tempUsers = append(tempUsers, User{v.userName, v.userId})
			roomMap[v.roomId] = Room{v.roomId, tempUsers}
		}
	}
	var rooms []Room
	for _, val := range roomMap {
		rooms = append(rooms, val)
	}
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Id < rooms[j].Id
	})
	return rooms, nil
}

func getConv(db sql.DB, roomId int) ([]Message, error) {
	rows, err := db.Query("SELECT u.name, m.message, m.timestamp FROM messages m JOIN users u on m.user_id = u.id where room_id = ? ORDER BY timestamp;", roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Messages []Message
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.Message, &message.Name, &message.Time)
		if err != nil {
			return nil, err
		}
		Messages = append(Messages, message)
	}
	return Messages, nil
}

func insertMessage(db sql.DB, message Message, roomId int, userId int) error {
	err := db.QueryRow("INSERT INTO messages(id, timestamp, room_id, message, user_id) VALUES (DEFAULT, CURRENT_TIMESTAMP, ?, ?, ?);", roomId, message.Message, userId).Err()
	return err
}
