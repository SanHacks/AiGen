package main

import (
	"aigen/aigenUi"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// mainApp is the main function that creates the tabs and the input box container
// It returns the tabs and the input box container
// The input box container is the container that contains the input box, send button, and voice chat button
// This is where All the fun happens :)
func mainApp(mapungubwe fyne.App) (fyne.CanvasObject, *fyne.Container) {
	// Create main chat tab
	chat, aiGen := ChatTab()
	
	// Create essential tabs only
	apiKeys := aigenUi.EasyConfigTab(mapungubwe) // API keys configuration
	gallery := aigenUi.UserMedia() // Media gallery
	browser := aigenUi.BrowserTab(mapungubwe) // Secure browser
	
	// Create notifications tab
	notifications := aigenUi.NotificationsTab()
	
	// Create sites/iframe tab
	sitesTab := createSitesTab(mapungubwe)
	
	// Modern tabs - only include valid tabs
	tabsItems := []*container.TabItem{}
	
	// Always add chat
	if aiGen != nil && aiGen.Content != nil {
		tabsItems = append(tabsItems, aiGen)
	}
	
	// Add browser if available
	if browser != nil && browser.Content != nil {
		tabsItems = append(tabsItems, browser)
	}
	
	// Add sites tab if available
	if sitesTab != nil && sitesTab.Content != nil {
		tabsItems = append(tabsItems, sitesTab)
	}
	
	// Add other tabs
	if apiKeys != nil && apiKeys.Content != nil {
		tabsItems = append(tabsItems, apiKeys)
	}
	if gallery != nil && gallery.Content != nil {
		tabsItems = append(tabsItems, gallery)
	}
	
	// Add notifications if available
	if notifications != nil && notifications.Content != nil {
		tabsItems = append(tabsItems, notifications)
	}
	
	tabs := container.NewAppTabs(tabsItems...)
	
	// Create modern input container
	inputContainer := createModernInputContainer(chat, tabs)
	
	// Add sidebar (ChatGPT/Gemini style) - Pass chat for new conversation button
	sidebar := aigenUi.ModernSidebarWithChat(mapungubwe, chat)
	
	// Main content with sidebar
	mainContent := container.NewHSplit(
		sidebar,      // Left sidebar
		tabs,         // Right content
	)
	mainContent.SetOffset(0.2) // 20% for sidebar
	
	// Show input by default on chat tab
	tabs.OnSelected = func(tab *container.TabItem) {
		if tab == aiGen {
			inputContainer.Show()
		} else {
			inputContainer.Hide()
		}
	}
	
	// Make sure input is visible initially
	inputContainer.Show()
	
	return mainContent, inputContainer
}
