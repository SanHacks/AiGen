package aigenUi

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// FinancialTab is the tab that contains the stock market
func FinancialTab(mapungubwe fyne.App) *container.TabItem {

	container.NewBorder(nil, nil, nil, nil, widget.NewLabel("FinTech Tab Content"))

	financeTab := container.NewTabItem("FinTech",
		&widget.Button{OnTapped: func() {
			//pop up window
			window := mapungubwe.NewWindow("FinTech")
			//Add content to the window and show it to the user image
			container.NewAdaptiveGrid(1, container.NewVBox(
				widget.NewCard("FinTech", "FinTech", widget.NewLabel("FinTech Tab Content")),
			))
			window.SetContent(widget.NewLabel("FinTech Tab Content"))
			window.Resize(fyne.NewSize(400, 300))
			window.Show()
		}})
	return financeTab
}
