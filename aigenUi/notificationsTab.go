package aigenUi

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

// Notification represents a single notification
type Notification struct {
	ID        int
	Message   string
	Timestamp time.Time
	Read      bool
}

var notificationsList []Notification = []Notification{}
var notificationListWidget *widget.List

// NotificationsTab creates a tab for viewing all notifications
func NotificationsTab() *container.TabItem {
	notificationListWidget = widget.NewList(
		func() int {
			return len(notificationsList)
		},
		func() fyne.CanvasObject {
			message := widget.NewLabel("Notification message")
			message.Wrapping = fyne.TextWrapWord
			time := widget.NewLabel("00:00")
			time.Importance = widget.LowImportance
			
			return container.NewBorder(
				nil, nil,
				message,
				time,
				nil,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(notificationsList) {
				notif := notificationsList[id]
				if border, ok := obj.(*fyne.Container); ok {
					// Update message
					messageLabel := border.Objects[2].(*widget.Label)
					messageLabel.SetText(notif.Message)
					
					// Update time
					timeLabel := border.Objects[3].(*widget.Label)
					timeLabel.SetText(notif.Timestamp.Format("15:04"))
				}
			}
		},
	)
	
	// Clear all button
	clearBtn := widget.NewButton("Clear All", func() {
		notificationsList = []Notification{}
		notificationListWidget.Refresh()
	})
	
	// Main content
	content := container.NewBorder(
		container.NewHBox(
			widget.NewLabel("📢 Notifications"),
			container.NewHBox(nil, clearBtn),
		),
		nil, nil, nil,
		notificationListWidget,
	)
	
	tab := container.NewTabItem("Notifications", content)
	return tab
}

// AddNotification adds a new notification to the list
func AddNotification(message string) {
	notification := Notification{
		ID:        len(notificationsList) + 1,
		Message:   message,
		Timestamp: time.Now(),
		Read:      false,
	}
	
	notificationsList = append([]Notification{notification}, notificationsList...)
	
	if notificationListWidget != nil {
		notificationListWidget.Refresh()
	}
}

// GetNotifications returns all notifications
func GetNotifications() []Notification {
	return notificationsList
}

