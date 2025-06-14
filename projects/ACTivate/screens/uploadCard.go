package screens

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ethersphere/bee/v2/pkg/swarm"
)

type uploadedItem struct {
	Name      string
	Reference string
	Size      int64
	Timestamp time.Time
	Mimetype  string
}

func (i *index) showUploadCard() *widget.Card {
	upForm := i.uploadForm()
	listButton := i.listUploadsButton(fyne.NewSize(200, 100))
	return widget.NewCard("Upload", "upload content into swarm", container.NewVBox(upForm, listButton))
}

func (i *index) uploadForm() *widget.Form {
	filepath := ""
	mimetype := ""
	fileSize := int64(0)
	var pathBind = binding.BindString(&filepath)
	path := widget.NewEntry()
	path.Bind(pathBind)
	path.Disable()
	var file io.Reader
	openFileButton := widget.NewButton("File Open", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				i.showError(err)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()
			data, err := io.ReadAll(reader)
			if err != nil {
				i.showError(err)
				return
			}
			fileSize = int64(len(data))
			mimetype = reader.URI().MimeType()
			err = pathBind.Set(reader.URI().Name())
			if err != nil {
				i.showError(err)
				return
			}
			file = bytes.NewReader(data)
		}, i.Window)
		fd.Show()
	})

	upForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Add file", Widget: path, HintText: "Filepath"},
			{Text: "Choose File", Widget: openFileButton},
		},
	}
	upForm.OnSubmit = func() {
		go func() {
			defer func() {
				err := pathBind.Set("")
				if err != nil {
					i.logger.Log(fmt.Sprintf("failed to bind path: %s", err.Error()))
				}
				file = nil
			}()
			if file == nil {
				i.showError(fmt.Errorf("please select a file"))
				return
			}
			batchID := i.getPreferenceString(batchPrefKey)
			if batchID == "" {
				i.showError(fmt.Errorf("please select a batch of stamp"))
				return
			}
			filename := path.Text
			i.logger.Log(fmt.Sprintf("stamp selected: %s", batchID))
			i.showProgressWithMessage(fmt.Sprintf("Uploading %s", filename))
			ref, _, err := i.bl.AddFileBzz(context.Background(), batchID, filename, mimetype, false, swarm.ZeroAddress, false, 0, file)
			if err != nil {
				i.hideProgress()
				i.showError(err)
				return
			}
			i.logger.Log(fmt.Sprintf("reference of the uploaded file: %s", ref.String()))
			uploadedSrt := i.getPreferenceString(uploadsPrefKey)
			uploads := []uploadedItem{}
			if uploadedSrt != "" {
				err := json.Unmarshal([]byte(uploadedSrt), &uploads)
				if err != nil {
					i.showError(err)
				}
			}
			uploads = append(uploads, uploadedItem{
				Name:      filename,
				Reference: ref.String(),
				Timestamp: time.Now(),
				Size:      fileSize,
				Mimetype:  mimetype,
			})
			data, err := json.Marshal(uploads)
			if err != nil {
				i.hideProgress()
				i.showError(err)
				return
			}
			i.setPreference(uploadsPrefKey, string(data))
			d := dialog.NewCustomConfirm("Upload successful", "Ok", "Cancel", i.copyDialog(shortenHashOrAddress(ref.String()), ref.String()), func(b bool) {}, i.Window)
			i.hideProgress()
			d.Show()
		}()
	}

	return upForm
}

func (i *index) listUploadsButton(minSize fyne.Size) *widget.Button {
	button := widget.NewButton("All Uploads", func() {
		uploadedContent := container.NewVBox()
		uploadedContentWrapper := container.NewScroll(uploadedContent)
		uploadedSrt := i.getPreferenceString(uploadsPrefKey)
		uploads := []uploadedItem{}
		if uploadedSrt != "" {
			err := json.Unmarshal([]byte(uploadedSrt), &uploads)
			if err != nil {
				i.showError(err)
			}
			for _, v := range uploads {
				ref := v.Reference
				name := v.Name
				label := widget.NewLabel(fmt.Sprintf("%s\n%s", name, shortenHashOrAddress(ref)))
				label.Wrapping = fyne.TextWrapWord
				item := container.NewBorder(label, nil, nil, i.copyButton(ref))
				uploadedContent.Add(item)
			}
		}

		if len(uploads) == 0 {
			uploadedContent.Add(widget.NewLabel("Empty upload list"))
		}

		child := i.app.NewWindow("Uploaded content")
		size := child.Canvas().Content().Size()
		if size.Width < minSize.Width {
			size.Width = minSize.Width
		}
		if size.Height < minSize.Height {
			size.Height = minSize.Height
		}
		child.Resize(size)
		child.SetContent(uploadedContentWrapper)
		child.Show()
	})

	return button
}
