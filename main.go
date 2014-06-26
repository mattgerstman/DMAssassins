package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"database/sql"
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	//"text/tabwriter"
	 "os"
	 "bufio"
	 "strings"
	 "strconv"
)

var db *sql.DB
var err error
var logged_in_user string

func addPlayer(name string) { //, password string) {
	id := uuid.New()
	_, err := db.Query(`INSERT INTO dm_users (user_id, email) VALUES ($1,$2)`, id, name)
	fmt.Println(err)
}

func addPlayerMenu() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	addPlayer(strings.TrimSpace(name))//, string(password))
}

func loginMenu() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	name, _ := reader.ReadString('\n')
	row := db.QueryRow(`SELECT user_id FROM dm_users WHERE name = $1`, strings.TrimSpace(name));
	row.Scan(&logged_in_user)
	fmt.Println("Logged in as:\t", logged_in_user);
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func killTargetMenu() {
	fmt.Print("Enter Target Pin: ");
	pin64, _ := strconv.ParseInt(readLine(), 10, 0);
	pin := int (pin64)

	row := db.QueryRow(`SELECT pin FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user)
	var target_pin int;
	row.Scan(&target_pin)
	if pin == target_pin {
		row = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user)
		var new_target_id string
		row.Scan(&new_target_id)
		row = db.QueryRow(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1)`, logged_in_user)
		row = db.QueryRow(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2`, new_target_id, logged_in_user)
	} else {
		fmt.Println("Invalid pin\n", target_pin, "\n")
	}
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

func assignTargets() {
	rows, _ := db.Query(`DELETE FROM dm_user_targets`)
	rows, _ = db.Query(`SELECT user_id FROM dm_users ORDER BY random()`)

	var user_id string
	var prev_user_id string
	var first_user_id string

	rows.Next();
	rows.Scan(&first_user_id)
	prev_user_id = first_user_id
	for rows.Next() {
		rows.Scan(&user_id);
		db.Query(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`, prev_user_id, user_id)
		prev_user_id = user_id;
	}
	db.Query(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`, user_id, first_user_id)

	//fmt.Println(row.Columns())
}

func connect() {
	db, err = sql.Open("postgres", "postgres://localhost?dbname=dmassassins&sslmode=disable")
}

func main() {

	connect()
	fmt.Println(err)

	//writer := tabwriter.NewWriter();

	for {
		fmt.Printf("Select a menu item:\nLogin\t\t(l)\nAdd Player\t(a)\nGet Player\t(g)\nAssign Targets\t(t)\n")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Selection: ")
		menu, _ := reader.ReadString('\n')
		switch {
			case menu[0] == 'l':
				loginMenu()
			case menu[0] == 'a':
				addPlayerMenu()
			case menu[0] == 't':
				assignTargets()
			case menu[0] == 'k':
				killTargetMenu()
			case menu[0] == 'q':
				return
		}	

	}
		
				// 
		// 
		// //plainPW, _ := reader.ReadString('\n')
		// //password, _ := Crypt([]byte(strings.TrimSpace(plainPW)))
		// 

	//addPlayer(name, string(password))
}