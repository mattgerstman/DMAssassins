package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	//"text/tabwriter"
	"bufio"
	"os"
)


var err error
var logged_in_user User

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func addPlayerMenu() {
	fmt.Print("Email: ")
	email := readLine()
	fmt.Print("Password: ")
	plainPW := readLine()
	NewUser(email, plainPW, "muggle")
}

func loginMenu() {
	fmt.Print("Email: ")
	email := readLine()
	row := db.QueryRow(`SELECT user_id FROM dm_users WHERE email = $1`, email)

	var user_id string
	row.Scan(&user_id)

	logged_in_user.User_id = user_id

	fmt.Println("Logged in as:\t", user_id)
}

func killTargetMenu() {
	fmt.Println("Enter target secret")
	secret := readLine()
	logged_in_user.KillTarget(secret)

}

func assignTargets() {
	rows, _ := db.Query(`DELETE FROM dm_user_targets`)
	rows, _ = db.Query(`SELECT user_id FROM dm_users ORDER BY random()`)

	var user_id string
	var prev_user_id string
	var first_user_id string

	rows.Next()
	rows.Scan(&first_user_id)
	prev_user_id = first_user_id
	for rows.Next() {
		rows.Scan(&user_id)
		db.Query(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`, prev_user_id, user_id)
		prev_user_id = user_id
	}
	db.Query(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`, user_id, first_user_id)

	//fmt.Println(row.Columns())
}


func main() {

	StartServer()
	fmt.Println(err)

	//writer := tabwriter.NewWriter();

	// for {
	// 	fmt.Printf("Select a menu item:\nLogin\t\t(l)\nAdd Player\t(a)\nGet Player\t(g)\nAssign Targets\t(t)\n")
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Print("Enter Selection: ")
	// 	menu, _ := reader.ReadString('\n')
	// 	switch {
	// 	case menu[0] == 'l':
	// 		loginMenu()
	// 	case menu[0] == 'a':
	// 		addPlayerMenu()
	// 	case menu[0] == 't':
	// 		assignTargets()
	// 	case menu[0] == 'k':
	// 		killTargetMenu()
	// 	case menu[0] == 'q':
	// 		return
	// 	}

	// }

	//
	//
	//

	//addPlayer(name, string(password))
}
