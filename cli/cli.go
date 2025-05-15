package main

import (
	"context"
	"fmt"
	"os"
	KDB "tigre/kdb"
)

func KDB_txt2list(txtName string, listName string) {

}

// run txt to go in terminal
func main() {
	db := KDB.ConnectDB()
	var ctx = context.Background()
	//fmt.Println(len(os.Args), os.Args)
	cmd := ""
	arg1 := ""
	arg2 := ""
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		cmd = os.Args[1]
	}
	if len(os.Args) > 2 {
		fmt.Println(os.Args[2])
		arg1 = os.Args[2]
	}
	if len(os.Args) > 3 {
		fmt.Println(os.Args[3])
		arg2 = os.Args[3]
	}

	switch cmd {
	case "T2L":
		if arg1 != "" && arg2 != "" {
			keyname := arg1
			targetfile := arg2
			KDB.Txt2List(db, keyname, targetfile)
		}
	default:
		fmt.Println("no valid command,", db.Ping(ctx))

	}
}
