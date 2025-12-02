
package database

import (
	b64 "encoding/base64"
	"fmt"
)

type Read struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Read     string `json:"read"`
}

type Received struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Received string `json:"received"`
}

type Cons struct {
	ConId        int    `json:"conid"`
	CategoryInfo string `json:"categoryinfo"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	Sender       string `json:"sender"`
	Toread       string `json:"toread"`
	Mtimestamp   string `json:"mtimestamp"`
	Mtext        string `json:"mtext"`
	Mphoto       string `json:"mphoto"`
	Mcategory    string `json:"mcategory"`
}

type Conversation struct {
	Conid        int    `json:"conid"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	Users        []User   `json:"user"`
	CategoryInfo string `json:"categoryinfo"`
	Messages     []Mes  `json:"message"`
}

type ConversationInfo struct {
	Conid        int    `json:"conid"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	Users        []User   `json:"user"`
	CategoryInfo string `json:"categoryinfo"`
}

type Mes struct {
	Mesid           int        `json:"mesid"`
	Conid           int        `json:"conid"`
	Timestamp       string     `json:"timestamp"`
	Text            string     `json:"text"`
	Photo           string     `json:"photo"`
	Username        string     `json:"username"`
	Login           string     `json:"login"`
	Category        string     `json:"category"`
	IsForwarded     int        `json:"isforwarded"`
	ReplyedText     string     `json:"replyedtext"`
	ReplyedPhoto    string     `json:"replyedphoto"`
	ReplyedUsername string     `json:"replyedusername"`
	ReceivedArray   []Received `json:"receivedarray"`
	ReadArray       []Read     `json:"readarray"`
	Comments        []Comm     `json:"comments"`
}

type Comm struct {
	Commid   int    `json:"commId"`
	Mesid    int    `json:"mesid"`
	Username string `json:"username"`
	Emoji    string `json:"emoji"`
}

type User struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type Usrs struct {
	Login        string `json:"login"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	CategoryInfo string `json:"categoryinfo"`
}

// UTILITÀ

func encodePhoto(b []byte) string {
	return b64.StdEncoding.EncodeToString(b)
}

// USERS

func (db *appdbimpl) GetName(token string) (string, error) {
	var name string
	err := db.c.QueryRow(
		`SELECT login FROM Users WHERE token=?`, token,
	).Scan(&name)

	return name, err
}

func (db *appdbimpl) GetUserInfo(token string) (string, []byte, error) {
	var username string
	var photo []byte

	err := db.c.QueryRow(
		`SELECT username, userphoto FROM Users WHERE token=?`, token,
	).Scan(&username, &photo)

	return username, photo, err
}

func (db *appdbimpl) GetUsers(token, search string, groups bool) ([]Usrs, error) {
	var users []Usrs
	var photo []byte

	groupQuery := ""
	if groups {
		groupQuery = `
			UNION
			SELECT DISTINCT c.id_con LL, c.name NN, coalesce(c.photo, '') PP, 'group' CC
			FROM Conversations c, UsersConversations uc 
			WHERE c.category = 'group'
			AND c.id_con = uc.id_con 
			AND uc.login = (SELECT login FROM Users WHERE token='` + token + `')`

		if search != "" {
			groupQuery += ` AND NN LIKE "%` + search + `%" `
		}
	}

	query := `
		SELECT DISTINCT login LL, username NN, coalesce(userphoto, '') PP, 'user' CC
		FROM Users`

	if search != "" {
		query += ` WHERE login LIKE "%` + search + `%" OR username LIKE "%` + search + `%"`
	}

	query += groupQuery + ` ORDER BY NN, LL`

	rows, err := db.c.Query(query)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user Usrs

		if err := rows.Scan(&user.Login, &user.Name, &photo, &user.CategoryInfo); err != nil {
			return users, err
		}

		user.Photo = encodePhoto(photo)
		users = append(users, user)
	}

	return users, rows.Err()
}

// GROUP

func (db *appdbimpl) GetGroupInfo(conid int) (ConInfo, error) {
	var conv ConInfo
	var photo []byte

	err := db.c.QueryRow(`
		SELECT id_con, category, name, coalesce(photo, '') 
		FROM Conversations WHERE id_con = ?`, conid,
	).Scan(&conv.Conid, &conv.CategoryInfo, &conv.Name, &photo)

	if err != nil {
		return conv, err
	}

	conv.Photo = encodePhoto(photo)

	rows, err := db.c.Query(`
		SELECT u.login, u.username, coalesce(u.userphoto, '') 
		FROM Users u, UsersConversations uc 
		WHERE uc.id_con=? AND uc.login=u.login`, conid)

	if err != nil {
		return conv, err
	}
	defer rows.Close()

	for rows.Next() {
		var user Us
		if err := rows.Scan(&user.Login, &user.Name, &photo); err != nil {
			return conv, err
		}

		user.Photo = encodePhoto(photo)
		conv.Users = append(conv.Users, user)
	}

	return conv, rows.Err()
}

// AUTH

func (db *appdbimpl) IsAuthorized(token string, conid int) error {
	var login string

	err := db.c.QueryRow(`
		SELECT login FROM UsersConversations 
		WHERE login=(SELECT login FROM Users WHERE token=?) AND id_con=?`,
		token, conid,
	).Scan(&login)

	if err != nil {
		return err
	}

	if login != "" {
		return nil
	}

	return fmt.Errorf("conversazione non autorizzata")
}

func (db *appdbimpl) IsMsgAuthorized(token string, msid int) error {
	var login string

	err := db.c.QueryRow(`
		SELECT login FROM UsersMessages
		WHERE login=(SELECT login FROM Users WHERE token=?) AND id_msg=?`,
		token, msid,
	).Scan(&login)

	if err != nil {
		return err
	}

	if login != "" {
		return nil
	}

	return fmt.Errorf("messaggio non autorizzato")
}

func (db *appdbimpl) IsAGroup(token string, conid int) error {
	var login string

	err := db.c.QueryRow(`
		SELECT uc.login 
		FROM UsersConversations uc, Conversations c
		WHERE uc.login=(SELECT login FROM Users WHERE token=?)
		AND uc.id_con=c.id_con
		AND c.id_con=?
		AND c.category='group'`, token, conid).Scan(&login)

	if login != "" {
		return err
	}

	return fmt.Errorf("non è un gruppo")
}

