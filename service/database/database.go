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

// NUOVA ISTANZA APP DATABASE
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("è necessario il databse quando si crea un AppDatabase")
	}

	sqlStmt := `PRAGMA foreign_keys = ON`
	_, ferr := db.Exec(sqlStmt)
	if ferr != nil {
		return nil, fmt.Errorf("errore durante l'abilitazione delle chiavi esterne: %w", ferr)
	}

	// CONTROLLO SE LA TABELLA ESISTE ALTTRIENTI LA CREO
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Comments';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS "Comments" (
						"id_com"	INTEGER NOT NULL,
						"id_msg"	INTEGER NOT NULL,
						"login"	TEXT NOT NULL,
						"emoji"	TEXT NOT NULL,
						PRIMARY KEY("id_com" AUTOINCREMENT),
						CONSTRAINT "message" FOREIGN KEY("id_msg") REFERENCES "Messages"("id_msg") ON DELETE CASCADE,
						CONSTRAINT "user" FOREIGN KEY("login") REFERENCES "Users"("login") ON DELETE CASCADE,
						CONSTRAINT "comment" UNIQUE("id_msg", "login")
					);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Conversations';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS "Conversations" (
						"id_con"	INTEGER NOT NULL,
						"category"	TEXT NOT NULL,
						"name"	TEXT NOT NULL DEFAULT 'conversation',
						"photo"	BLOB,
						PRIMARY KEY("id_con" AUTOINCREMENT)
					);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database: %w", err)
		}

		sqlStmt = `INSERT INTO Conversations VALUES(0,'conversation', '', NULL)`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("ERRORE %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Messages';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS "Messages" (
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
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS "Users" (
						"login"	TEXT NOT NULL,
						"token"	TEXT NOT NULL,
						"username"	TEXT DEFAULT NULL,
						"userphoto"	BLOB,
						PRIMARY KEY("login")
					);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Forwarded';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS "Forwarded" (
						"id_msg_new"	INTEGER,
						"id_msg_old"	INTEGER,
						UNIQUE("id_msg_new","id_msg_old"),
						CONSTRAINT "messagenew" FOREIGN KEY("id_msg_new") REFERENCES "Messages"("id_msg") ON DELETE CASCADE
					);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Replyed';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS "Replyed" (
						"id_msg_new"	INTEGER,
						"id_msg_old"	INTEGER,
						UNIQUE("id_msg_new","id_msg_old"),
						CONSTRAINT "messagenew" FOREIGN KEY("id_msg_new") REFERENCES "Messages"("id_msg") ON DELETE CASCADE
					);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='UsersConversations';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS "UsersConversations" (
						"login"	TEXT,
						"id_con"	INTEGER,
						UNIQUE("login","id_con"),
						CONSTRAINT "conversation" FOREIGN KEY("id_con") REFERENCES "Conversations"("id_con") ON DELETE CASCADE,
						CONSTRAINT "user" FOREIGN KEY("login") REFERENCES "Users"("login") ON DELETE CASCADE
					);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='UsersMessages';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS "UsersMessages" (
						"login"	TEXT,
						"id_msg"	INTEGER,
						"received"	BOOLEAN,
						"read"	BOOLEAN,
						UNIQUE("login","id_msg"),
						CONSTRAINT "message" FOREIGN KEY("id_msg") REFERENCES "Messages"("id_msg") ON DELETE CASCADE,
						CONSTRAINT "user" FOREIGN KEY("login") REFERENCES "Users"("login") ON DELETE CASCADE
					);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}