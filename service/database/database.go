package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase è l'interfaccia principale per tutte le operazioni sul DB
type AppDatabase interface {
	GetName(token string) (string, error)
	GetToken(name string) (string, error)
	SetToken(name string) (string, error)
	SetMyUsername(name string, token string) error
	GetUserInfo(token string) (string, []byte, error)
	GetUsers(token string, search string, groups bool) ([]User, error)
	GetMyConversations(token string) ([]ConversationShort, error)
	SetConversation(token string, category, receiver, convname string, photo []byte) (int, error)
	SendMessage(conid int, token, text, category string, photo []byte, replyID int) error
	IsAuthorized(token string, conid int) error
	IsMsgAuthorized(token string, mesid int) error
	IsAGroup(token string, conid int) error
	GetConversation(conid, start, limit int, token string) (Conversation, error)
	GetGroupInfo(conid int) (ConversationInfo, error)
	SetGroup(token string) (int, error)
	AddToGroup(logins []string, conid int) error
	LeaveGroup(token string, convid int) error
	SetGroupName(convid int, name string) error
	SetGroupPhoto(convid int, photo []byte) error
	ForwardMessage(token string, convid, mesid, forwardID int) error
	DeleteMessage(token string, mesid, convid int) (int, error)
	CommentMessage(mesid int, token, emoji string) (int, error)
	UncommentMessage(token string, mesid int) (int, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// NUOVA ISTANZA APPDATABASE
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("è necessario il databse quando si crea un AppDatabase")
	}

	_, err := db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return nil, fmt.Errorf("errore durante l'abilitazione delle chiavi esterne: %w", err)
	}

	// CONTROLLO SE LA TABELLA ESISTE, ALTRIMENTI LA CREO
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
			return nil, fmt.Errorf("errore durante la creazione della struttura del database %w", err)
		}

		sqlStmt = `INSERT INTO Conversations VALUES(0,'conversation', '', NULL)`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("errore durante la creazione della struttura del database %w", err)
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

