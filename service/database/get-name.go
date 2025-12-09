package database

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// STRUTTURE PER MESS, CONV, UTENTI

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
	Photo        string `json:"photo"` // binhex
	Sender       string `json:"sender"`
	Toread       string `json:"toread"`
	Mtimestamp   string `json:"mtimestamp"`
	Mtext        string `json:"mtext"`
	Mphoto       string `json:"mphoto"` // binhex
	Mcategory    string `json:"mcategory"`
}

type Convers struct {
	Conid        int    `json:"conid"`
	Name         string `json:"name"`
	Photo        string `json:"photo"` // binhex
	Users        []Us   `json:"user"`
	CategoryInfo string `json:"categoryinfo"`
	Messages     []Mes  `json:"message"`
}

type ConInfo struct {
	Conid        int    `json:"conid"`
	Name         string `json:"name"`
	Photo        string `json:"photo"` // binhex
	Users        []Us   `json:"user"`
	CategoryInfo string `json:"categoryinfo"`
}

type Mes struct {
	Mesid           int        `json:"mesid"`
	Conid           int        `json:"conid"`
	Timestamp       string     `json:"timestamp"`
	Text            string     `json:"text"`
	Photo           string     `json:"photo"` // binhex
	Username        string     `json:"username"`
	Login           string     `json:"login"`
	Category        string     `json:"category"`
	IsForwarded     int        `json:"isforwarded"`
	ReplyedText     string     `json:"replyedtext"`
	ReplyedPhoto    string     `json:"replyedphoto"` // binhex
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

type Us struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Photo string `json:"photo"` // binhex
}

type Usrs struct {
	Login        string `json:"login"`
	Name         string `json:"name"`
	Photo        string `json:"photo"` // binhex
	CategoryInfo string `json:"categoryinfo"`
}

// FUNZIONI X DB

// restituisce il login dell'utente dal token
func (db *appdbimpl) GetName(token string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT login FROM Users WHERE token=?", token).Scan(&name)
	return name, err
}

// restituisce tutti gli utenti (opzionalmente anche i gruppi) filtrando per search
func (db *appdbimpl) GetUsers(token string, search string, groups bool) ([]Usrs, error) {
	var allusers []Usrs
	var p []byte
	var query, groupQuery string

	// Se vogliamo includere anche i gruppi
	if groups {
		groupQuery = ` UNION 
			SELECT DISTINCT c.id_con LL, c.name NN, coalesce(c.photo, '') PP, 'group' CC
			FROM Conversations c, UsersConversations uc 
			WHERE c.category='group' AND c.id_con=uc.id_con 
			AND uc.login=(SELECT login FROM Users WHERE token='` + token + `')`
		if search != "" {
			groupQuery += ` AND NN like "%` + search + `%" `
		}
	}

	if search != "" {
		query = `SELECT DISTINCT login LL, username NN, coalesce(userphoto, '') PP, 'user' CC
				FROM Users WHERE login like "%` + search + `%" or username like "%` + search + `%"
				` + groupQuery + ` ORDER by NN, LL`
	} else {
		query = `SELECT DISTINCT login LL, username NN, coalesce(userphoto, '') PP, 'user' CC
				FROM Users ` + groupQuery + ` ORDER by NN, LL`
	}

	rows, err := db.c.Query(query)
	if err != nil {
		return allusers, err
	}

	for rows.Next() {
		var user Usrs
		if err := rows.Scan(&user.Login, &user.Name, &p, &user.CategoryInfo); err != nil {
			return allusers, err
		}
		// CONVERSIONE FOTO IN BINHEX
		user.Photo = hex.EncodeToString(p)
		allusers = append(allusers, user)
	}

	if err := rows.Err(); err != nil {
		return allusers, err
	}

	return allusers, nil
}

// restituisce le info di un gruppo e i suoi utenti
func (db *appdbimpl) GetGroupInfo(conid int) (ConInfo, error) {
	var conv ConInfo
	var p []byte

	err := db.c.QueryRow(`SELECT id_con, category, name, coalesce(photo, '') 
						  FROM Conversations WHERE id_con=?`, conid).Scan(&conv.Conid, &conv.CategoryInfo, &conv.Name, &p)
	if err != nil {
		return conv, err
	}

	// conversione foto binhex
	conv.Photo = hex.EncodeToString(p)

	rows, err := db.c.Query(`SELECT u.login, u.username, coalesce(u.userphoto, '') 
							 FROM Users u, UsersConversations uc 
							 WHERE uc.id_con=? AND uc.login=u.login`, conid)
	if err != nil {
		return conv, err
	}

	for rows.Next() {
		user := Us{}
		if err := rows.Scan(&user.Login, &user.Name, &p); err != nil {
			return conv, err
		}
		user.Photo = hex.EncodeToString(p)
		conv.Users = append(conv.Users, user)
	}

	if err := rows.Err(); err != nil {
		return conv, err
	}

	return conv, nil
}

// restituisce username e foto di un utente dal token
func (db *appdbimpl) GetUserInfo(token string) (string, []byte, error) {
	var username string
	var photo []byte

	err := db.c.QueryRow("SELECT username, userphoto FROM Users WHERE token=?", token).Scan(&username, &photo)
	return username, photo, err
}
