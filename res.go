package main

import (
	"fmt"
	"fyne.io/fyne"
	"log"
)

var addIcon = loadResourceForSure("add.png")
var refreshIcon = loadResourceForSure("refresh.png")
var editSubscriptionIcon = loadResourceForSure("edit_subscription.png")
var deleteIcon = loadResourceForSure("delete.png")
var addProfileIcon = loadResourceForSure("add_profile.png")
var playIcon = loadResourceForSure("play.png")
var stopIcon = loadResourceForSure("stop.png")
var editProfileIcon = loadResourceForSure("edit_profile.png")
var settingsIcon = loadResourceForSure("settings.png")

func loadResourceForSure(filename string) (res fyne.Resource) {
	var err error
	res, err = fyne.LoadResourceFromPath(fmt.Sprintf("./res/%s", filename))
	if err != nil {
		log.Fatalln(err)
	}
	return
}
