package database

// Riceve il login di un utente (name).
//Genera un token univoco usando MD5 del login.
//Controlla se l’utente esiste già nel database:
//		Se sì → aggiorna il token esistente.
//		Se no → crea un nuovo utente, cercando di usare eventuali valori precedenti di username e userphoto.
//	Restituisce il token generato, pronto per l’uso in autenticazione.

import (
	"crypto/md5"
	"encoding/hex"
)

func (db *appdbimpl) SetToken(name string) (string, error) {

	b := []byte(name)

	token := md5.Sum(b)

	stoken := hex.EncodeToString(token[:])

	var count int
	err := db.c.QueryRow(`SELECT count(*) FROM Users WHERE login=?`, name).Scan(&count)
	if err != nil {
		return stoken, err
	}

	if count > 0 {
		_, err = db.c.Exec(`UPDATE Users SET token=? WHERE login=?`, stoken, name)
	} else {
		_, err = db.c.Exec(`INSERT INTO Users (login, token, username, userphoto) 
			VALUES (?, ?, 
				coalesce((SELECT username FROM Users WHERE login=?), ?), 
				coalesce((SELECT userphoto FROM Users WHERE login=?), NULL)
			)`,
			name, stoken, name, name, name)
	}

	if err != nil {
		return stoken, err
	}

	return stoken, nil
}
