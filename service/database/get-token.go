package database

import (
	"fmt"
)

func (db *appdbimpl) GetToken(token string) (string, error) {

	// CONTROLLO RIGHE ESISTONO X TOKEN FORNITO
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE token=?", token).Scan(&count)
	if err != nil {
		
		// ERRORE DURANTE QUERY
		return token, err
	}

	
	if count < 1 {
		return token, fmt.Errorf("token non autorizzato: %s", token)
	}

	// TOKEN VALIDO
	return token, nil
}
