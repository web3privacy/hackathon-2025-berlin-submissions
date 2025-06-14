package screens

import (
	"context"
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/ethersphere/bee/v2/pkg/swarm"
)

func (i *index) showDownloadCard() *widget.Card {
	dlForm := i.downloadForm()
	return widget.NewCard("Download", "download content from swarm", dlForm)
}

func (i *index) downloadForm() *widget.Form {
	hash := widget.NewEntry()
	hash.SetPlaceHolder("Swarm Hash")
	dlForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Swarm Hash", Widget: hash, HintText: "Swarm Hash"},
		},
		OnSubmit: func() {
			dlAddr, err := swarm.ParseHexAddress(hash.Text)
			if err != nil {
				i.showError(err)
				return
			}
			if hash.Text == "" {
				i.showError(fmt.Errorf("please enter a hash"))
				return
			}
			go func() {
				i.showProgressWithMessage(fmt.Sprintf("Downloading %s", shortenHashOrAddress(hash.Text)))
				ref, fileName, err := i.bl.GetBzz(context.Background(), dlAddr, nil, nil, nil)
				if err != nil {
					i.hideProgress()
					i.showError(err)
					return
				}
				hash.SetText("")
				data, err := io.ReadAll(ref)
				if err != nil {
					i.hideProgress()
					i.showError(err)
					return
				}
				i.hideProgress()
				saveFile := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
					if err != nil {
						i.showError(err)
						return
					}
					if writer == nil {
						return
					}
					_, err = writer.Write(data)
					if err != nil {
						i.showError(err)
						return
					}
					writer.Close()
				}, i.Window)
				saveFile.SetFileName(fileName)
				saveFile.Show()
			}()
		},
	}

	return dlForm
}
