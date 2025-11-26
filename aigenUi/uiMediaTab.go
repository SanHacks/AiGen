package aigenUi

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// UserMedia creates an improved media gallery tab
func UserMedia() *container.TabItem {
	imageDir := "dalleAssets"

	// Create a list to hold the image previews
	imageDataList := make([]ImageData, 0)

	// List files in the specified directory
	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageFile(path) {
			imageDataList = append(imageDataList, ImageData{Path: path, Info: info})
		}
		return nil
	})

	if err != nil {
		log.Printf("Error loading gallery: %v", err)
	}

	// Sort the imageDataList in descending order (newest first)
	sort.SliceStable(imageDataList, func(i, j int) bool {
		return imageDataList[i].Info.ModTime().After(imageDataList[j].Info.ModTime())
	})

	// Create grid of image cards
	imageCards := make([]fyne.CanvasObject, 0)
	for _, data := range imageDataList {
		card := createImageCard(data.Path)
		imageCards = append(imageCards, card)
	}

	// Empty state if no images
	if len(imageCards) == 0 {
		emptyLabel := widget.NewLabel("No images yet. Generate some AI images to see them here!")
		emptyLabel.Wrapping = fyne.TextWrapWord
		emptyLabel.Alignment = fyne.TextAlignCenter
		emptyContainer := container.NewCenter(emptyLabel)
		return container.NewTabItem("📷 Gallery", emptyContainer)
	}

	// Create a responsive grid (3 columns for desktop)
	grid := container.New(layout.NewGridLayout(3), imageCards...)

	// Create a scrollable container
	scrollable := container.NewVScroll(grid)

	// Header with info
	header := widget.NewCard("", "", 
		widget.NewLabel("🎨 Your AI-Generated Images"))

	content := container.NewBorder(header, nil, nil, nil, scrollable)

	return container.NewTabItem("📷 Gallery", content)
}

// createImageCard creates a card for each image with preview and actions
func createImageCard(imagePath string) fyne.CanvasObject {
	// Load image
	img := canvas.NewImageFromFile(imagePath)
	img.SetMinSize(fyne.NewSize(250, 250))
	img.FillMode = canvas.ImageFillContain

	// Get filename
	filename := filepath.Base(imagePath)

	// View button
	viewButton := widget.NewButtonWithIcon("View", theme.ZoomInIcon(), func() {
		showImageDialog(imagePath)
	})

	// Delete button
	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		confirmDelete(imagePath)
	})
	deleteButton.Importance = widget.DangerImportance

	// Button row
	buttons := container.NewHBox(
		layout.NewSpacer(),
		viewButton,
		deleteButton,
	)

	// Create card
	card := widget.NewCard("", filename, container.NewVBox(
		img,
		buttons,
	))

	return card
}

// showImageDialog shows the full-size image in a dialog
func showImageDialog(imagePath string) {
	img := canvas.NewImageFromFile(imagePath)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(800, 600))

	w := fyne.CurrentApp().Driver().AllWindows()[0]
	d := dialog.NewCustom(filepath.Base(imagePath), "Close", img, w)
	d.Resize(fyne.NewSize(900, 700))
	d.Show()
}

// confirmDelete asks for confirmation before deleting
func confirmDelete(imagePath string) {
	w := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm("Delete Image", 
		"Are you sure you want to delete this image?\n"+filepath.Base(imagePath), 
		func(confirm bool) {
			if confirm {
				err := os.Remove(imagePath)
				if err != nil {
					dialog.ShowError(err, w)
					log.Printf("Error deleting image: %v", err)
				} else {
					dialog.ShowInformation("Deleted", "Image deleted successfully", w)
					log.Printf("Deleted image: %s", imagePath)
				}
			}
		}, w)
}

// Check if a file has an image extension
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".bmp" || ext == ".webp"
}
