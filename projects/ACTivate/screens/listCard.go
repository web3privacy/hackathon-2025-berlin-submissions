package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (i *index) showListCard() *container.Scroll {
	sampleData := []string{
		"Entry 1",
		"Entry 2",
		"Entry 3",
	}

	stringList := widget.NewList(
		func() int {
			return len(sampleData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("grantee list")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(sampleData[id])
		},
	)

	scrollableStringList := container.NewScroll(stringList)
	scrollableStringList.SetMinSize(fyne.NewSize(200, 150))

	return scrollableStringList
}
