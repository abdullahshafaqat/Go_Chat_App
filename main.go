package main

import (
	db "github.com/abdullahshafaqat/Go_Chat_App.git/db/postgres"
)

func main() {
	db := db.Database()
	defer db.Close()

}
