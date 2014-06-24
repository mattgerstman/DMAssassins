package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"database/sql"
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	//"os"
	//"bufio"
)

var db *sql.DB
var err error

func addPlayer(name string, password string) {
	id := uuid.New()
	_, err := db.Query(`INSERT INTO dm_users (user_id, name, password) VALUES ($1,$2,$3)`, id, name, password)
	fmt.Println(err)
}

func clear(b []byte) {
    for i := 0; i < len(b); i++ {
        b[i] = 0;
    }
}

func Crypt(password []byte) ([]byte, error) {
    defer clear(password)
    return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func getPlayer(name string, plainPW string) {
	row := db.QueryRow(`SELECT user_id FROM dm_users where name = $1`, name)

	var user_id string;

	row.Scan(&user_id)

	fmt.Println(user_id)
	//fmt.Println(row.Columns())
	
}

func main() {

	db, err = sql.Open("postgres", "postgres://localhost?dbname=dmassassins&sslmode=disable")
	
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Name: ")
	name := "Matt"
	//name, _ := reader.ReadString('\n')
	//fmt.Print("Password: ")
	//plainPW, _ := reader.ReadString('\n')

	//password, _ := Crypt([]byte(plainPW))

	plainPW := ""
	getPlayer(name, plainPW )

	//addPlayer(name, string(password))
}