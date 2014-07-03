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
