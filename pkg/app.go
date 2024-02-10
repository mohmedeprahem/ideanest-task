package app

import (
	database "ideanest-task/pkg/database/mongodb"
)

func Run() {
	database.ConnectDB()
}
