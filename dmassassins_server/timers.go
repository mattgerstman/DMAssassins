package main

import (
	"database/sql"
)

func (game *Game) SetKillTimer(tx *sql.Tx, time int) (appErr *ApplicationError) {
	return nil
}
