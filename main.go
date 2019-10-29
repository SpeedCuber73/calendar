package main

import (
	"fmt"
	"time"

	"github.com/SpeedCuber73/calendar/app"
	simplestorage "github.com/SpeedCuber73/calendar/simple-storage"
)

func main() {
	myStorage, _ := simplestorage.CreateSimpleStorage()
	myApp, _ := app.CreateApp(myStorage)
	myEvent := &app.Event{
		ID:        0,
		Name:      "test event",
		StartAt:   time.Now(),
		EndAt:     time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	}

	// add one event
	myApp.AddNewEvent(myEvent)
	events := myApp.ListAllEvents()
	fmt.Printf("%s\n", events)

	// update it
	fmt.Println("Updating...")
	myEvent.Name = "updated"
	myApp.ChangeEvent(myEvent.ID, myEvent)
	events = myApp.ListAllEvents()
	fmt.Printf("%s\n", events)

	// and remove it
	fmt.Println("Removing...")
	myApp.RemoveEvent(0)
	events = myApp.ListAllEvents()
	fmt.Printf("%s\n", events)
}
