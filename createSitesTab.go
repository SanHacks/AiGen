package main

import (
	"aigen/aigenRest"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fmt"
)

// createSitesTab creates the configurable iframe sites tab
func createSitesTab(app fyne.App) *container.TabItem {
	// Load sites
	sites, err := GetSites()
	if err != nil {
		sites = []Site{}
	}
	
	// Create site list
	list := widget.NewList(
		func() int { return len(sites) },
		func() fyne.CanvasObject {
			icon := widget.NewLabel("🛒")
			name := widget.NewLabel("Site Name")
			desc := widget.NewLabel("Description")
			vbox := container.NewVBox(name, desc)
			return container.NewHBox(icon, vbox)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(sites) {
				site := sites[id]
				if hbox, ok := obj.(*fyne.Container); ok && len(hbox.Objects) >= 2 {
					hbox.Objects[0].(*widget.Label).SetText(site.Icon)
					if vbox, ok := hbox.Objects[1].(*fyne.Container); ok && len(vbox.Objects) >= 2 {
						vbox.Objects[0].(*widget.Label).SetText(site.Name)
						vbox.Objects[1].(*widget.Label).SetText(site.Description)
					}
				}
			}
		},
	)
	
	// Create iframe for viewing sites
	iframe := widget.NewRichTextFromMarkdown("# Select a site to view\n\nChoose a site from the list to load it here.")
	iframe.Wrapping = fyne.TextWrapWord
	
	// Site selection handler
	list.OnSelected = func(id widget.ListItemID) {
		if id < len(sites) {
			site := sites[id]
			// Update iframe content with site info
			iframe.ParseMarkdown(fmt.Sprintf("# %s\n\n**URL:** %s\n\n**Description:** %s\n\n_Loading site: %s_", site.Name, site.URL, site.Description, site.URL))
			aigenRest.SendNotificationNow(fmt.Sprintf("Site selected: %s", site.Name))
		}
	}
	
	// Add site button
	addBtn := widget.NewButton("+ Add Site", func() {
		showAddSiteDialog(app, list)
	})
	
	// Layout: List on left, iframe on right
	content := container.NewHSplit(
		container.NewBorder(
			container.NewHBox(widget.NewLabel("📱 Sites"), addBtn),
			nil, nil, nil,
			list,
		),
		iframe,
	)
	content.SetOffset(0.3)
	
	tab := container.NewTabItem("🌍 Sites", content)
	return tab
}

func showAddSiteDialog(app fyne.App, list *widget.List) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Site name")
	
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://example.com")
	
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Description")
	
	iconEntry := widget.NewEntry()
	iconEntry.SetPlaceHolder("🛒")
	iconEntry.SetText("🌐")
	
	dialog := widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("URL", urlEntry),
		widget.NewFormItem("Description", descEntry),
		widget.NewFormItem("Icon", iconEntry),
	)
	
	// For now, just show the form
	dialog.OnSubmit = func() {
		AddSite(nameEntry.Text, urlEntry.Text, descEntry.Text, iconEntry.Text)
		list.Refresh()
		aigenRest.SendNotificationNow("Site added!")
	}
}

