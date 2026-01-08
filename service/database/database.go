/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// INTERFACCIA PRINCIPALE PER TUTTE LE CRUD DEL DB
type AppDatabase interface {
	GetName(token string) (string, error)
	GetUsers(token string, search string, groups bool) ([]Usrs, error)
	SetMyUsername(name string, token string) error
	SetMyUserphoto(photo []byte, token string) error
	SetToken(name string) (string, error)
	GetToken(token string) (string, error)
	GetUserInfo(token string) (string, []byte, error)
	GetMyConversations(token string) ([]Cons, error)
	SetConversation(token string, category string, receiver string, convname string, photo []byte) (int, error)
	SendMessage(conid int, token string, text string, category string, photo []byte, msid int) error
	IsAuthorized(token string, conid int) error
	IsMsgAuthorized(token string, msid int) error
	IsAGroup(token string, conid int) error
	GetConversation(conid int, start int, limit int, token string) (Convers, error)
	GetGroupInfo(conid int) (ConInfo, error)
	SetGroup(token string) (int, error)
	AddToGroup(logins []string, conid int) error
	LeaveGroup(token string, convid int) error
	SetGroupName(convid int, name string) error
	SetGroupPhoto(convid int, photo []byte) error
	ForwardMessage(token string, convid int, mesid int, forward int) error
	DeleteMessage(token string, mesid int, convid int) (int, error)
	CommentMessage(mesid int, token string, emoji string) (int, error)
	UncommentMessage(token string, mesid int) (int, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New accetta un DB già aperto e crea le tabelle se mancano
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when creating AppDatabase")
	}

	// 1. Abilita le Foreign Keys
	// Questo è spesso il primo punto dove si nota se il file non è accessibile,
	// ma se main.go l'ha aperto, qui dovrebbe funzionare.
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("errore durante l'abilitazione delle chiavi esterne: %w", err)
	}

	// 2. Crea le tabelle se non esistono
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("errore durante la creazione delle tabelle: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// createTables contiene le query per creare la struttura del database
func createTables(db *sql.DB) error {
	
	// Tabella Users
	queryUsers := `CREATE TABLE IF NOT EXISTS "Users" (
		"login"	TEXT NOT NULL,
		"token"	TEXT NOT NULL,
		"username"	TEXT DEFAULT NULL,
		"userphoto"	BLOB,
		PRIMARY KEY("login")
	);`
	if _, err := db.Exec(queryUsers); err != nil {
		return err
	}

	// Tabella Conversations
	queryConversations := `CREATE TABLE IF NOT EXISTS "Conversations" (
		"id_con"	INTEGER NOT NULL,
		"category"	TEXT NOT NULL,
		"name"	TEXT NOT NULL DEFAULT 'conversation',
		"photo"	BLOB,
		PRIMARY KEY("id_con" AUTOINCREMENT)
	);`
	if _, err := db.Exec(queryConversations); err != nil {
		return err
	}

	// Conversazione Default (0)
	// Usa INSERT OR IGNORE per non dare errore se esiste già
	queryInsertConv := `INSERT OR IGNORE INTO Conversations (id_con, category, name, photo) VALUES(0,'conversation', '', NULL);`
	if _, err := db.Exec(queryInsertConv); err != nil {
		return fmt.Errorf("errore init conversazione default: %w", err)
	}

	// Tabella Messages
	queryMessages := `CREATE TABLE IF NOT EXISTS "Messages" (
		"id_msg"	INTEGER NOT NULL,
		"id_con"	INTEGER,
		"login"	TEXT,
		"timestamp"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		"text"	TEXT,
		"photo" BLOB,
		"category"  TEXT NOT NULL DEFAULT "text",
		PRIMARY KEY("id_msg" AUTOINCREMENT),
		CONSTRAINT "conversation" FOREIGN KEY("id_con") REFERENCES "Conversations"("id_con") ON DELETE CASCADE,
		CONSTRAINT "sender" FOREIGN KEY("login") REFERENCES "Users"("login") ON DELETE CASCADE
	);`
	if _, err := db.Exec(queryMessages); err != nil {
		return err
	}

	// Tabella Comments
	queryComments := `CREATE TABLE IF NOT EXISTS "Comments" (
		"id_com"	INTEGER NOT NULL,
		"id_msg"	INTEGER NOT NULL,
		"login"	TEXT NOT NULL,
		"emoji"	TEXT NOT NULL,
		PRIMARY KEY("id_com" AUTOINCREMENT),
		CONSTRAINT "message" FOREIGN KEY("id_msg") REFERENCES "Messages"("id_msg") ON DELETE CASCADE,
		CONSTRAINT "user" FOREIGN KEY("login") REFERENCES "Users"("login") ON DELETE CASCADE,
		CONSTRAINT "comment" UNIQUE("id_msg", "login")
	);`
	if _, err := db.Exec(queryComments); err != nil {
		return err
	}

	// Tabella Forwarded
	queryForwarded := `CREATE TABLE IF NOT EXISTS "Forwarded" (
		"id_msg_new"	INTEGER,
		"id_msg_old"	INTEGER,
		UNIQUE("id_msg_new","id_msg_old"),
		CONSTRAINT "messagenew" FOREIGN KEY("id_msg_new") REFERENCES "Messages"("id_msg") ON DELETE CASCADE
	);`
	if _, err := db.Exec(queryForwarded); err != nil {
		return err
	}

	// Tabella Replyed
	queryReplyed := `CREATE TABLE IF NOT EXISTS "Replyed" (
		"id_msg_new"	INTEGER,
		"id_msg_old"	INTEGER,
		UNIQUE("id_msg_new","id_msg_old"),
		CONSTRAINT "messagenew" FOREIGN KEY("id_msg_new") REFERENCES "Messages"("id_msg") ON DELETE CASCADE
	);`
	if _, err := db.Exec(queryReplyed); err != nil {
		return err
	}

	// Tabella UsersConversations
	queryUsersConv := `CREATE TABLE IF NOT EXISTS "UsersConversations" (
		"login"	TEXT,
		"id_con"	INTEGER,
		UNIQUE("login","id_con"),
		CONSTRAINT "conversation" FOREIGN KEY("id_con") REFERENCES "Conversations"("id_con") ON DELETE CASCADE,
		CONSTRAINT "user" FOREIGN KEY("login") REFERENCES "Users"("login") ON DELETE CASCADE
	);`
	if _, err := db.Exec(queryUsersConv); err != nil {
		return err
	}

	// Tabella UsersMessages
	queryUsersMsg := `CREATE TABLE IF NOT EXISTS "UsersMessages" (
		"login"	TEXT,
		"id_msg"	INTEGER,
		"received"	BOOLEAN,
		"read"	BOOLEAN,
		UNIQUE("login","id_msg"),
		CONSTRAINT "message" FOREIGN KEY("id_msg") REFERENCES "Messages"("id_msg") ON DELETE CASCADE,
		CONSTRAINT "user" FOREIGN KEY("login") REFERENCES "Users"("login") ON DELETE CASCADE
	);`
	if _, err := db.Exec(queryUsersMsg); err != nil {
		return err
	}

	return nil
}