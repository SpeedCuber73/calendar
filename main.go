package main

import (
	"fmt"

	"github.com/SpeedCuber73/calendar/app"
	"github.com/SpeedCuber73/calendar/storage/simplestorage"
)

func main() {
	myStorage, _ := simplestorage.CreateSimpleStorage()
	myApp, _ := app.CreateApp(myStorage)
	events := myApp.ListEvents()

	fmt.Println(events)
}
