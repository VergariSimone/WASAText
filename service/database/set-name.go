package database

import "fmt"

// --- Aggiornamento dati utente

// aggiorna lo username dell'utente e controlla che il nome non sia già utilizzato da altri utenti
func (db *appdbimpl) SetMyUsername(name string, token string) error {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM Users WHERE token!=? AND username=?`, token, name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("nome già in uso")
	}
	_, err = db.c.Exec("UPDATE Users SET username=? WHERE token=?", name, token)
	return err
}

// aggiorna la foto dell'utente
func (db *appdbimpl) SetMyUserphoto(photo []byte, token string) error {
	_, err := db.c.Exec("UPDATE Users SET userphoto=? WHERE token=?", photo, token)
	return err
}

// --- Gestione conversazioni

// crea una nuova conversazione tra due utenti o ritorna quella esistente
func (db *appdbimpl) SetConversation(token, category, receiver, convname string, photo []byte) (int, error) {
	var convID, count int

	// controlla se esiste già una conversazione tra i due utenti
	err := db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM UsersConversations uc1, UsersConversations uc2, Conversations c
		WHERE uc1.id_con = uc2.id_con
		  AND c.id_con = uc1.id_con
		  AND c.category='conversation'
		  AND uc1.login != uc2.login
		  AND uc1.login=(SELECT login FROM Users WHERE token=?)
		  AND uc2.login=?`, token, receiver).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		// se esiste, restituisce id
		err := db.c.QueryRow(`
			SELECT uc1.id_con
			FROM UsersConversations uc1, UsersConversations uc2, Conversations c
			WHERE uc1.id_con = uc2.id_con
			  AND c.id_con = uc1.id_con
			  AND c.category='conversation'
			  AND uc1.login != uc2.login
			  AND uc1.login=(SELECT login FROM Users WHERE token=?)
			  AND uc2.login=?`, token, receiver).Scan(&convID)
		if err != nil {
			return 0, err
		}
		return convID, nil
	}

	// else nuova conversazione
	_, err = db.c.Exec(`INSERT INTO Conversations (id_con, category, name, photo) VALUES (NULL, ?, ?, ?)`, category, convname, photo)
	if err != nil {
		return 0, err
	}

	// prende id conv appena creata
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&convID)
	if err != nil {
		return 0, err
	}

	// inserisce i due utenti nella conv
	_, err = db.c.Exec(`INSERT INTO UsersConversations (login, id_con) VALUES ((SELECT login FROM Users WHERE token=?), ?)`, token, convID)
	if err != nil {
		return 0, err
	}
	_, err = db.c.Exec(`INSERT INTO UsersConversations (login, id_con) VALUES (?, ?)`, receiver, convID)
	if err != nil {
		return 0, err
	}

	return convID, nil
}

// invia un messaggio in una conversazione
func (db *appdbimpl) SendMessage(conID int, token, text, category string, photo []byte, replyToID int) error {

	_, err := db.c.Exec(`INSERT INTO Messages VALUES(NULL, ?, (SELECT login FROM Users WHERE token=?), CURRENT_TIMESTAMP, ?, ?, ?)`,
		conID, token, text, photo, category)
	if err != nil {
		return err
	}

	var msgID int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&msgID)
	if err != nil {
		return err
	}

	rows, err := db.c.Query(`SELECT login FROM UsersConversations WHERE id_con=?`, conID)
	if err != nil {
		return err
	}
	var allUsers []string
	for rows.Next() {
		var login string
		if err := rows.Scan(&login); err != nil {
			return err
		}
		allUsers = append(allUsers, login)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// inserisce il messaggio nella tabella UsersMessages per ciascun utente
	for _, login := range allUsers {
		_, err := db.c.Exec("INSERT INTO UsersMessages VALUES(?, ?, 0, 0)", login, msgID)
		if err != nil {
			return err
		}
	}

	// se  messaggio è risposta altro messaggio, registra la relazione
	if replyToID > 0 {
		_, err = db.c.Exec("INSERT INTO Replyed VALUES(?, ?)", msgID, replyToID)
		if err != nil {
			return err
		}
	}

	return nil
}

// --- Gestione gruppi

// SetGroup crea un nuovo gruppo con messaggio iniziale
func (db *appdbimpl) SetGroup(token string) (int, error) {
	_, err := db.c.Exec(`INSERT INTO Conversations (id_con, category, name, photo) 
		VALUES (NULL, "group", concat("group created by ", (SELECT username FROM Users WHERE token=?)), NULL)`, token)
	if err != nil {
		return 0, err
	}

	var groupID int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&groupID)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO UsersConversations (login, id_con) VALUES ((SELECT login FROM Users WHERE token=?), ?)`, token, groupID)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO Messages (id_con, login, timestamp, text, photo, category)
		VALUES (?, (SELECT login FROM Users WHERE token=?), CURRENT_TIMESTAMP, 'group created', NULL, 'text')`, groupID, token)
	if err != nil {
		return 0, err
	}

	var msgID int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&msgID)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`INSERT INTO UsersMessages VALUES((SELECT login FROM Users WHERE token=?), ?, 0, 0)`, token, msgID)
	if err != nil {
		return 0, err
	}

	return groupID, nil
}

// aggiunge utenti esistenti a un gruppo e marca come letti i messaggi precedenti
func (db *appdbimpl) AddToGroup(logins []string, conID int) error {
	for _, login := range logins {
		var count int
		err := db.c.QueryRow(`SELECT COUNT(*) FROM UsersConversations WHERE login=? AND id_con=?`, login, conID).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		_, err = db.c.Exec(`INSERT INTO UsersConversations (login, id_con) VALUES (?, ?)`, login, conID)
		if err != nil {
			return err
		}

		rows, err := db.c.Query(`SELECT id_msg FROM Messages WHERE id_con=?`, conID)
		if err != nil {
			return err
		}
		var msgIDs []int
		for rows.Next() {
			var msgID int
			if err := rows.Scan(&msgID); err != nil {
				return err
			}
			msgIDs = append(msgIDs, msgID)
		}
		if err := rows.Err(); err != nil {
			return err
		}

		for _, msgID := range msgIDs {
			_, err := db.c.Exec(`INSERT INTO UsersMessages VALUES(?, ?, 1, 1)`, login, msgID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// rimuove un utente dal gruppo e cancella la conversazione se non ci sono più membri
func (db *appdbimpl) LeaveGroup(token string, conID int) error {
	_, err := db.c.Exec(`DELETE FROM UsersConversations WHERE login=(SELECT login FROM Users WHERE token=?) AND id_con=?`, token, conID)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`DELETE FROM UsersMessages WHERE login=(SELECT login FROM Users WHERE token=?) AND id_msg IN (SELECT id_msg FROM Messages WHERE id_con=?)`, token, conID)
	if err != nil {
		return err
	}

	var count int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM UsersConversations WHERE id_con=?`, conID).Scan(&count)
	if err != nil {
		return err
	}

	if count < 1 {
		_, err := db.c.Exec(`DELETE FROM Conversations WHERE id_con=?`, conID)
		if err != nil {
			return err
		}
	}

	return nil
}

// cambia il nome del gruppo
func (db *appdbimpl) SetGroupName(conID int, name string) error {
	_, err := db.c.Exec("UPDATE Conversations SET name=? WHERE id_con=?", name, conID)
	return err
}

// cambia la foto del gruppo
func (db *appdbimpl) SetGroupPhoto(conID int, photo []byte) error {
	_, err := db.c.Exec("UPDATE Conversations SET photo=? WHERE id_con=?", photo, conID)
	return err
}

// --- Altre funzioni

// inoltra un messaggio in un'altra conversazione
func (db *appdbimpl) ForwardMessage(token string, conID int, mesID int, forward int) error {
	_, err := db.c.Exec(`INSERT INTO Messages (id_con, login, timestamp, text, photo, category)
		SELECT ?, (SELECT login FROM Users WHERE token=?), CURRENT_TIMESTAMP, text, photo, category
		FROM Messages WHERE id_msg=?`, forward, token, mesID)
	if err != nil {
		return err
	}

	var newMsgID int
	err = db.c.QueryRow(`SELECT last_insert_rowid()`).Scan(&newMsgID)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`INSERT INTO UsersMessages(login, id_msg, received, read)
		SELECT login, ?, 0, 0 FROM UsersConversations WHERE id_con=?`, newMsgID, forward)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`INSERT INTO Forwarded VALUES(?, ?)`, newMsgID, forward)
	if err != nil {
		return err
	}

	return nil
}

// cancella un messaggio di un utente
func (db *appdbimpl) DeleteMessage(token string, mesID int, conID int) (int, error) {

	var ownerLogin string
	var requesterLogin string

	// trova chi sta facendo la richiesta
	err := db.c.QueryRow(`SELECT login FROM Users WHERE token=?`, token).Scan(&requesterLogin)
	if err != nil {
		// token non valido
		return -1, nil 
	}

	// trova il proprietario del messaggio
	err = db.c.QueryRow(`SELECT login FROM Messages WHERE id_msg=? AND id_con=?`, mesID, conID).Scan(&ownerLogin)
	if err != nil {
		// Se err è sql.ErrNoRows, il messaggio non esiste (o ID conv errato)
		return 0, nil
	}

	// controllo Permessi
	if ownerLogin != requesterLogin {
		// messaggio esiste ma non è dell utente che vuole cancellarlo0
		return -1, nil
	}

	// gestione risposte 
	var replyCount int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM Replyed WHERE id_msg_old=?`, mesID).Scan(&replyCount)
	if err != nil {
		return 0, err
	}

	if replyCount > 0 {
		_, err = db.c.Exec(`UPDATE Messages SET id_con=0 WHERE id_msg=?`, mesID)
		if err != nil {
			return 0, err
		}
		_, err = db.c.Exec(`DELETE FROM Messages WHERE id_con=0 AND id_msg NOT IN (SELECT id_msg_old FROM Replyed)`)
		if err != nil {
			return 0, err
		}
	} else {
		_, err = db.c.Exec(`DELETE FROM Messages WHERE id_msg=?`, mesID)
		if err != nil {
			return 0, err
		}
	}

	// pulizia Conversazione (se vuota)
	var remaining int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM Messages WHERE id_con=?`, conID).Scan(&remaining)
	if err != nil {
		return 0, err
	}

	if remaining == 0 {
		// se non ci sono più messaggi, cancella l'intera conversazione
		_, err = db.c.Exec(`DELETE FROM Conversations WHERE id_con=?`, conID)
		if err != nil {
			return 0, err
		}
	}

	// success
	return 1, nil
}

// aggiunge un commento a un messaggio
func (db *appdbimpl) CommentMessage(mesID int, token string, emoji string) (int, error) {

	// recupera prima del login cosi piu sicuro
	var userLogin string
	err := db.c.QueryRow(`SELECT login FROM "Users" WHERE token = ?`, token).Scan(&userLogin)
	if err != nil {
		return 0, err // errore se utente non trovato
	}

	// rimozione commenti precednenti 
	_, err = db.c.Exec(`DELETE FROM "Comments" WHERE id_msg=? AND login=?`, mesID, userLogin)
	if err != nil {
		return 0, err
	}

	// inserimento nuovo commento
	res, err := db.c.Exec(`INSERT INTO "Comments" (id_msg, login, emoji) VALUES (?, ?, ?)`,
		mesID, userLogin, emoji)
	if err != nil {
		return 0, err
	}

	// recupero id creato
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

// rimuove un commento fatto dall'utente su un messaggio
// 1 = Cancellato con successo
// 0 = Nessun commento trovato da cancellare (o token non valido)
func (db *appdbimpl) UncommentMessage(token string, mesid int) (int, error) {

	// esegue la DELETE cercando la riga che corrisponde a quel messaggio
	// e a quell'utente (tramite subquery sul token)
	res, err := db.c.Exec(`
        DELETE FROM "Comments" 
        WHERE id_msg = ? 
        AND login = (SELECT login FROM "Users" WHERE token = ?)`,
		mesid, token)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	// se viene cancellata almeno una riga(successo)
	if rowsAffected > 0 {
		return 1, nil
	}

	// se non viene cancellato nulla (l'utente non aveva reagito o token errato)
	return 0, nil
}
