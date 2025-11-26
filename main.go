package main

import (
	"aigen/aigenUi"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	_ "github.com/mattn/go-sqlite3"
)

//func updateTime(clock *widget.Label) {
//	formatted := time.Now().Format("Time: 03:04:05")
//	clock.SetText(formatted)
//}

func main() {
	//Run Setup Scripts
	setup()
	//Start Application
	mapungubwe := app.NewWithID("com.sage.aigen")

	// Set theme based on user preference (default to dark)
	mapungubwe.Settings().SetTheme(theme.DarkTheme())
	aigenUi.SwitchUp(mapungubwe)
	playWelcomeSound()

	mainContent, inputBoxContainer := mainApp(mapungubwe)
	window := mapungubwe.NewWindow(aigenUi.MainTitle)
	window.SetIcon(theme.MailAttachmentIcon())
	window.SetFixedSize(false)

	window.CenterOnScreen()
	window.Resize(aigenUi.WindowSize)
	window.SetPadded(true)
	
	// Border layout: main content with input at bottom
	window.SetContent(container.NewBorder(
		nil,
		inputBoxContainer, // Input at bottom
		nil,
		nil,
		mainContent, // Main content with sidebar
	))
	
	window.ShowAndRun()
	window.SetOnClosed(aigenUi.GoodBye(mapungubwe))
}
