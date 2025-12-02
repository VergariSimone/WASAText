
package database

import (
	b64 "encoding/base64"
	"fmt"
)

// "WASAText/service/api"
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

type Convers struct {
	Conid        int    `json:"conid"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	Users        []Us   `json:"user"`
	CategoryInfo string `json:"categoryinfo"`
	Messages     []Mes  `json:"message"`
}

type ConInfo struct {
	Conid        int    `json:"conid"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	Users        []Us   `json:"user"`
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

type Us struct {
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

// FIXME aggiungere elementi messaggi in api
// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName(token string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT login FROM Users WHERE token=?", token).Scan(&name)
	return name, err
}

func (db *appdbimpl) GetUsers(token string, search string, groups bool) ([]Usrs, error) {
	var allusers []Usrs
	var p []byte
	var query string
	var groupQuery string

	groupQuery = ""
	if groups {
		groupQuery = ` UNION 
							SELECT DISTINCT c.id_con LL, c.name NN, coalesce(c.photo, '') PP, 'group' CC
							FROM Conversations c, UsersConversations uc 
							WHERE c.category='group' AND c.id_con=uc.id_con AND 
							uc.login=(SELECT login FROM Users WHERE token='` + token + `')`
		if search != "" {
			groupQuery = groupQuery + ` AND NN like "%` + search + `%" `
		}
	}

	if search != "" {
		query = `SELECT DISTINCT login LL, username NN, coalesce(userphoto, '') PP, 'user' CC
						  FROM Users WHERE login like "%` + search + `%" or username like "%` + search + `%"
						  ` + groupQuery + `
						  ORDER by NN, LL `
	} else {
		query = `SELECT DISTINCT login LL, username NN, coalesce(userphoto, '') PP, 'user' CC
						  FROM Users 
						  ` + groupQuery + ` 
						  ORDER by NN, LL `
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
		user.Photo = b64.StdEncoding.EncodeToString(p)

		allusers = append(allusers, user)
	}

	if err := rows.Err(); err != nil {
		return allusers, err
	}

	return allusers, err
}

func (db *appdbimpl) GetGroupInfo(conid int) (ConInfo, error) {
	var conv ConInfo
	var p []byte
	err := db.c.QueryRow(`
					SELECT id_con, category, name, coalesce(photo, '') 
					FROM Conversations WHERE id_con=?`, conid).Scan(
		&conv.Conid, &conv.CategoryInfo, &conv.Name, &p)
	if err != nil {
		return conv, err
	}
	conv.Photo = b64.StdEncoding.EncodeToString(p)

	rows, err := db.c.Query(`SELECT u.login, u.username, coalesce(u.userphoto, '') 
							 from Users u, UsersConversations uc where uc.id_con=? and uc.login=u.login`, conid)
	if err != nil {
		return conv, err
	}

	for rows.Next() {
		user := Us{}
		if err := rows.Scan(&user.Login, &user.Name, &p); err != nil {
			return conv, err
		}
		user.Photo = b64.StdEncoding.EncodeToString(p)

		conv.Users = append(conv.Users, user)
	}

	if err := rows.Err(); err != nil {
		return conv, err
	}

	return conv, err
}

func (db *appdbimpl) GetUserInfo(token string) (string, []byte, error) {
	var username string
	var photo []byte

	err := db.c.QueryRow("SELECT username, userphoto FROM Users WHERE token=?", token).Scan(&username, &photo)
	return username, photo, err
}

func (db *appdbimpl) GetMyConversations(token string) ([]Cons, error) {
	conShort := []Cons{}

	rows, err := db.c.Query(`SELECT	c.id_con, c.category, 
								CASE WHEN c.category = 'group'
	   							THEN c.name
	   							ELSE 
									(SELECT coalesce(username, login) FROM Users WHERE login=
									(SELECT login FROM UsersConversations WHERE id_con=c.id_con AND login!=
									(SELECT login FROM Users WHERE token=?)))
	   							END AS Name, 
								CASE WHEN c.category = 'group'
	   							THEN c.photo
	   							ELSE 
									(SELECT userphoto FROM Users WHERE login=
									(SELECT login FROM UsersConversations WHERE id_con=c.id_con AND login!=
									(SELECT login FROM Users WHERE token=?)))
	   							END AS Photo, m1.timestamp, coalesce(m1.text, ""), coalesce(m1.photo, ""), m1.category,
									(SELECT username FROM Users Where login=m1.login),
									(select count(*) from UsersMessages 
										where id_msg in (
											select id_msg 
											from Messages 
											where id_con=c.id_con
											) 
										and login=(SELECT login FROM Users WHERE token=?) 
										and read = 0)
								from Users u, UsersConversations uc, Conversations c, Messages m1
								where u.token = ? and uc.login = u.login and uc.id_con = c.id_con and 
									c.id_con = m1.id_con and m1.timestamp = (
										select max(timestamp) from Messages m2
										where id_con = m1.id_con
									) ORDER BY m1.timestamp DESC`,
		token, token, token, token)

	if err != nil {
		return nil, err
	}

	var allmsg []int
	var p []byte
	var q []byte
	// var t string
	for rows.Next() {
		conv := Cons{}
		if err := rows.Scan(&conv.ConId, &conv.CategoryInfo, &conv.Name, &p,
			&conv.Mtimestamp, &conv.Mtext, &q, &conv.Mcategory, &conv.Sender, &conv.Toread); err != nil {
			return nil, err
		}
		conv.Photo = b64.StdEncoding.EncodeToString(p)
		conv.Mphoto = b64.StdEncoding.EncodeToString(q)
		conv.Mtimestamp = conv.Mtimestamp[:10] + " " + conv.Mtimestamp[11:16]

		conShort = append(conShort, conv)

		allmsg = append(allmsg, conv.ConId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, idmsg := range allmsg {
		_, err = db.c.Exec(`UPDATE UsersMessages SET received=1 WHERE login=(SELECT login FROM Users WHERE token=?)
		AND received!=1 AND id_msg in (SELECT id_msg FROM Messages WHERE id_con=?)`, token, idmsg)
		if err != nil {
			return nil, err
		}
	}

	return conShort, err
}

func (db *appdbimpl) IsAuthorized(token string, conid int) error {
	var name string
	err := db.c.QueryRow(`SELECT login FROM UsersConversations 
	WHERE login=(SELECT login FROM Users WHERE token=?) AND id_con=?`, token, conid).Scan(&name)
	if err != nil {
		return err
	}
	if name != "" {
		return nil
	}
	return fmt.Errorf("unautharized conversation")
}

func (db *appdbimpl) IsMsgAuthorized(token string, msid int) error {
	var name string
	err := db.c.QueryRow(`SELECT login FROM UsersMessages 
	WHERE login=(SELECT login FROM Users WHERE token=?) AND id_msg=?`, token, msid).Scan(&name)
	if err != nil {
		return err
	}
	if name != "" {
		return nil
	}
	return fmt.Errorf("unautharized message")
}

func (db *appdbimpl) IsAGroup(token string, conid int) error {
	var name string
	err := db.c.QueryRow(`SELECT uc.login FROM UsersConversations uc, Conversations c
	WHERE uc.login=(SELECT login FROM Users WHERE token=?) AND uc.id_con=c.id_con AND c.id_con=? AND c.category="group"`, token, conid).Scan(&name)
	if name != "" {
		return err
	}
	return fmt.Errorf("not a group")
}

func (db *appdbimpl) GetConversation(conid int, start int, limit int, token string) (Convers, error) {
	con := Convers{}
	err := db.c.QueryRow(`SELECT category, 
								CASE WHEN category = 'group'
	   							THEN name
	   							ELSE 
									(SELECT coalesce(username, login) FROM Users WHERE login=
									(SELECT login FROM UsersConversations WHERE id_con=? AND login!=
									(SELECT login FROM Users WHERE token=?)))
	   							END AS Name , coalesce(photo, '') FROM Conversations WHERE id_con=?`,
		conid, token, conid).Scan(&con.CategoryInfo, &con.Name, &con.Photo)
	if err != nil {
		return con, err
	}
	// mes := []Mes{}
	var p []byte
	var q []byte
	var n int
	rows, err := db.c.Query(`SELECT m.id_msg, m.id_con, m.timestamp, coalesce(m.text, ""), coalesce(m.photo, ""), u.username, m.login, m.category, 
	(SELECT COUNT(*) FROM Forwarded WHERE id_msg_new=m.id_msg) AS forwarded,
	(SELECT COUNT(*) FROM Replyed WHERE id_msg_new=m.id_msg) AS replyed
	FROM Messages m, Users u WHERE m.id_con=? AND m.login=u.login ORDER BY m.timestamp DESC limit ?, ?`,
		conid, start, start+limit)
	if err != nil {
		return con, err
	}
	for rows.Next() {
		m := Mes{}
		if err := rows.Scan(&m.Mesid, &m.Conid, &m.Timestamp, &m.Text, &p, &m.Username, &m.Login, &m.Category, &m.IsForwarded, &n); err != nil {
			return con, err
		}
		m.Photo = b64.StdEncoding.EncodeToString(p)

		m.Timestamp = m.Timestamp[:10] + " " + m.Timestamp[11:16]
		crows, cerr := db.c.Query(`SELECT c.id_com, c.id_msg, u.username, c.emoji
		FROM comments c, Users u WHERE c.id_msg = ? AND u.login=c.login ORDER BY u.username`, m.Mesid)
		if cerr != nil {
			return con, cerr
		}
		for crows.Next() {
			c := Comm{}
			if coerr := crows.Scan(&c.Commid, &c.Mesid, &c.Username, &c.Emoji); coerr != nil {
				return con, coerr
			}
			m.Comments = append(m.Comments, c)
		}
		if coerr := crows.Err(); coerr != nil {
			return con, coerr
		}

		crows, cerr = db.c.Query(`SELECT distinct um.login, u.username, um.read, um.received 
								  FROM UsersMessages um, Users u WHERE u.login=um.login and um.id_msg = ?
								  ORDER by u.username, um.login`, m.Mesid)

		if cerr != nil {
			return con, cerr
		}
		for crows.Next() {
			rcv := Received{}
			rea := Read{}
			login := ""
			username := ""
			if coerr := crows.Scan(&login, &username, &rea.Read, &rcv.Received); coerr != nil {
				return con, coerr
			}
			rcv.Name = login
			rea.Name = login
			rcv.Username = username
			rea.Username = username

			m.ReadArray = append(m.ReadArray, rea)
			m.ReceivedArray = append(m.ReceivedArray, rcv)
		}
		if cerr := crows.Err(); cerr != nil {
			return con, cerr
		}
		if n == 0 {
			m.ReplyedPhoto = ""
			m.ReplyedText = ""
			m.ReplyedUsername = ""
		} else {
			rerr := db.c.QueryRow(`SELECT coalesce(m.text, ""), coalesce(m.photo, ""), coalesce(u.username, "") 
								FROM Messages m, Users u
								WHERE m.id_msg=(SELECT id_msg_old FROM Replyed WHERE id_msg_new=?) AND m.login=u.login`,
				m.Mesid).Scan(&m.ReplyedText, &q, &m.ReplyedUsername)
			if rerr != nil {
				return con, rerr
			}
			m.ReplyedPhoto = b64.StdEncoding.EncodeToString(q)
		}
		con.Messages = append(con.Messages, m)
	}
	con.Conid = conid

	if err := rows.Err(); err != nil {
		return con, err
	}

	_, err = db.c.Exec(`UPDATE UsersMessages SET received=1, read=1 WHERE login=(SELECT login FROM Users WHERE token=?)
		AND read!=1 AND id_msg in (SELECT id_msg FROM Messages WHERE id_con=?)`, token, conid)
	if err != nil {
		return con, err
	}

	return con, err
}