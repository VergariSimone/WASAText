package database

import "fmt"

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetMyUsername(name string, token string) error {
	var n int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM Users WHERE token!=? AND username=?`, token, name).Scan(&n)
	if err != nil {
		return err
	}

	if n > 0 {
		return fmt.Errorf("nome gia usato")
	}

	_, err = db.c.Exec("UPDATE Users set username=? WHERE token=?", name, token)
	return err

}

func (db *appdbimpl) SetMyUserphoto(photo []byte, token string) error {
	_, err := db.c.Exec("UPDATE Users set userphoto=? WHERE token=?", photo, token)
	return err
}

func (db *appdbimpl) SetConversation(token string, category string, receiver string, convname string, photo []byte) (int, error) {

	// Se esiste una conv fra i due ritorna quella
	var lid int
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) from UsersConversations uc1,  UsersConversations uc2, Conversations c
						  WHERE uc1.id_con=uc2.id_con AND c.id_con=uc1.id_con AND c.category='conversation'
						  AND uc1.login!=uc2.login
						  AND uc1.login=(SELECT login from Users WHERE token=?) 
						  AND uc2.login=?`, token, receiver).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		err := db.c.QueryRow(`SELECT uc1.id_con from UsersConversations uc1,  UsersConversations uc2, Conversations c
						  WHERE uc1.id_con=uc2.id_con AND c.id_con=uc1.id_con AND c.category='conversation'
						  AND uc1.login!=uc2.login
						  AND uc1.login=(SELECT login from Users WHERE token=?) 
						  AND uc2.login=?`, token, receiver).Scan(&lid)
		if err != nil {
			return 0, err
		}
		return lid, err
	}
	_, err = db.c.Exec(`INSERT INTO Conversations (id_con, category, name, photo) 
						values(NULL, ?, ?, ?)`, category, convname, photo)
	if err != nil {
		return 0, err
	}

	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&lid)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO UsersConversations (login, id_con) 
						values((SELECT login FROM Users WHERE token=?), ?)`, token, lid)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO UsersConversations (login, id_con) 
						values(?, ?)`, receiver, lid)
	if err != nil {
		return 0, err
	}

	return lid, err
}

func (db *appdbimpl) SendMessage(conid int, token string, text string, category string, photo []byte, msid int) error {

	// if category == "text" {
	_, err := db.c.Exec(`INSERT INTO Messages VALUES(NULL, ?, (SELECT login FROM Users WHERE token=?),
	CURRENT_TIMESTAMP, ?, ?, ?)`, conid, token, text, photo, category)
	if err != nil {
		return err
	}
	// } else {
	// 	_, err := db.c.Exec(`INSERT INTO Messages VALUES(NULL, ?, (SELECT login FROM Users WHERE token=?),
	// 	CURRENT_TIMESTAMP, ?, ?)`, conid, token, photo, category)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	var lid int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&lid)
	if err != nil {
		return err
	}
	rows, err := db.c.Query(`SELECT login FROM UsersConversations WHERE id_con=?`, conid)
	if err != nil {
		return err
	}
	var allm []string
	for rows.Next() {
		m := ""
		if err := rows.Scan(&m); err != nil {
			return err
		}
		allm = append(allm, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, a := range allm {
		_, cerr := db.c.Exec("INSERT INTO UsersMessages  VALUES(?, ?, 0, 0)", a, lid)
		if cerr != nil {
			return cerr
		}
	}
	if msid > 0 {
		_, err = db.c.Exec("INSERT INTO Replyed  VALUES(?, ?)", lid, msid)
		if err != nil {
			return err
		}
	}
	return err
}

func (db *appdbimpl) SetGroup(token string) (int, error) {
	_, err := db.c.Exec(`INSERT INTO Conversations (id_con, category, name, photo) 
						values(NULL, "group", concat("group created by ",(SELECT username FROM Users WHERE token=?)), NULL)`, token)
	if err != nil {
		return 0, err
	}
	var lid int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&lid)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO UsersConversations (login, id_con) 
						values((SELECT login FROM Users WHERE token=?), ?)`, token, lid)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO Messages (id_con, login, timestamp, text, photo, category)
						values(?,(SELECT login FROM Users WHERE token=?), CURRENT_TIMESTAMP, 'group created', NULL, 'text')`, lid, token)
	if err != nil {
		return 0, err
	}

	var lmsg int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&lmsg)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO UsersMessages 
						values((SELECT login FROM Users WHERE token=?), ?, 0, 0)`, token, lmsg)
	if err != nil {
		return 0, err
	}

	return lid, err
}

func (db *appdbimpl) AddToGroup(logins []string, conid int) error {
	var err error
	for _, c := range logins {
		var count int
		err = db.c.QueryRow(`SELECT COUNT(*) from UsersConversations where login=? and id_con=?`, c, conid).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		_, err = db.c.Exec(`INSERT INTO UsersConversations (login, id_con) 
							values(?, ?)`, c, conid)
		if err != nil {
			return err
		}
		var allmsg []int
		rows, err := db.c.Query(`SELECT id_msg FROM Messages WHERE id_con=?`, conid)
		if err != nil {
			return err
		}
		for rows.Next() {
			m := 0
			if err := rows.Scan(&m); err != nil {
				return err
			}
			allmsg = append(allmsg, m)
		}

		if err := rows.Err(); err != nil {
			return err
		}
		// Quando un nuovo utente viene aggiunto tutti i messaggi precedenti
		// devono risultare come letti e ricevuti
		for _, msg := range allmsg {
			_, cerr := db.c.Exec(`INSERT INTO UsersMessages VALUES(?, ?, 1, 1)`, c, msg)
			if cerr != nil {
				return cerr
			}
		}
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(token string, convid int) error {
	_, err := db.c.Exec(`DELETE FROM UsersConversations 
	WHERE login=(SELECT login FROM Users WHERE token=?) AND id_con=?`, token, convid)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`DELETE FROM UsersMessages 
	WHERE login=(SELECT login FROM Users WHERE token=?) AND 
	id_msg IN (SELECT id_msg FROM Messages WHERE id_con=?)`, token, convid)
	if err != nil {
		return err
	}

	var count int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM UsersConversations WHERE id_con=?`, convid).Scan(&count)
	if err != nil {
		return err
	}
	if count < 1 {
		_, err := db.c.Exec(`DELETE FROM Conversations WHERE id_con=?`, convid)
		if err != nil {
			return err
		}
	}

	return err
}

func (db *appdbimpl) SetGroupName(convid int, name string) error {
	_, err := db.c.Exec("UPDATE Conversations set name=? WHERE id_con=?", name, convid)
	return err
}

func (db *appdbimpl) SetGroupPhoto(convid int, photo []byte) error {
	_, err := db.c.Exec("UPDATE Conversations set photo=? WHERE id_con=?", photo, convid)
	return err
}

func (db *appdbimpl) ForwardMessage(token string, convid int, mesid int, forward int) error {
	_, err := db.c.Exec(`INSERT INTO Messages (id_con, login, timestamp, text, photo, category)
		SELECT ?, (SELECT login FROM Users WHERE token=?), CURRENT_TIMESTAMP, text, photo, category
		FROM Messages WHERE id_msg=?`, forward, token, mesid)
	if err != nil {
		return err
	}

	var lid int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&lid)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`INSERT INTO UsersMessages(login, id_msg, received, read)
		SELECT login, ?, 0, 0
		FROM UsersConversations WHERE id_con=?`, lid, forward)

	// rows, err := db.c.Query(`SELECT login FROM UsersConversations WHERE id_con=?`, forward)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`INSERT INTO Forwarded VALUES(?, ?)`, lid, forward)
	if err != nil {
		return err
	}

	/*
		for rows.Next() {
			m := ""
			if err := rows.Scan(&m); err != nil {
				return err
			}
			_, cerr := db.c.Exec(`INSERT INTO UsersMessages VALUES(?, ?, 0, 0)`, m, lid)
			if cerr != nil {
				return cerr
			}
		}
	*/
	return err

}

func (db *appdbimpl) DeleteMessage(token string, mesid int, convid int) (int, error) {
	var n int
	var s int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM Messages 
	WHERE login=(SELECT login FROM Users WHERE token=?) AND id_msg=?`, token, mesid).Scan(&n)
	if err != nil {
		return 0, err
	}
	err = db.c.QueryRow(`SELECT COUNT(*) FROM Messages 
	WHERE id_msg=?`, mesid).Scan(&s)
	if err != nil {
		return 0, err
	}
	if s > n {
		// n = numero di messaggi con mesid dell'utente, s = tutti mesid anche non dell'utente
		// se s > n il mesid non è dell'utente e non lo può cancellare
		return n - s, err
	}

	if n < 1 {
		return 0, fmt.Errorf("message not found")
	}

	err = db.c.QueryRow(`SELECT COUNT(*) FROM Replyed
	WHERE id_msg_old=?`, mesid).Scan(&s)
	if err != nil {
		return 0, err
	}

	if s > 0 {
		_, err = db.c.Exec(`UPDATE Messages set id_con=0 WHERE id_msg=?`, mesid)
		if err != nil {
			return 0, err
		}
		_, err = db.c.Exec(`DELETE FROM Messages WHERE id_con=0 AND 
							id_msg not in (SELECT id_msg_old FROM Replyed)`)
		if err != nil {
			return 0, err
		}
	} else {
		_, err = db.c.Exec(`DELETE FROM Messages WHERE id_msg=?`, mesid)
		if err != nil {
			return 0, err
		}
	}

	var m int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM Messages WHERE id_con=?`, convid).Scan(&m)
	if err != nil {
		return 0, err
	}
	if m < 1 {
		_, err = db.c.Exec(`DELETE FROM Conversations WHERE id_con=?`, convid)
		if err != nil {
			return 0, err
		}
	}

	return n, err
}

func (db *appdbimpl) CommentMessage(mesid int, token string, emoji string) (int, error) {

	_, err := db.c.Exec(`DELETE FROM Comments WHERE id_msg=? and login = (select login from Users where token=?)`, mesid, token)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO Comments VALUES(NULL, ?, (SELECT login FROM Users WHERE token=?),?)`,
		mesid, token, emoji)
	if err != nil {
		return 0, err
	}

	var lid int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&lid)
	if err != nil {
		return 0, err
	}

	return lid, err
}

func (db *appdbimpl) UncommentMessage(token string, mesid int) (int, error) {
	var n int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM Comments 
	WHERE login=(SELECT login FROM Users WHERE token=?) AND id_msg=?`, token, mesid).Scan(&n)
	if err != nil {
		return 0, err
	}
	// if n < 1 {
	// 	return 0, fmt.Errorf("comment not found")
	// }
	_, err = db.c.Exec(`DELETE FROM Comments WHERE login=(SELECT login FROM Users WHERE token=?) AND id_msg=?`, token, mesid)
	if err != nil {
		return 0, err
	}
	return n, err
}
