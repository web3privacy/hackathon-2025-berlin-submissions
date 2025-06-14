package screens

import (
	"fmt"
	"io"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ethereum/go-ethereum/common"
)

func (i *index) readAppData(filePath string) (string, error) {
	uri, err := storage.ParseURI("file://" + i.nodeConfig.path + filePath)
	if err != nil {
		i.logger.Log(fmt.Sprintf("failed to parse uri: %s", err.Error()))
		return "", err
	}

	reader, err := storage.Reader(uri)
	if err != nil {
		i.logger.Log(fmt.Sprintf("failed to read file: %s", err.Error()))
		return "", err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		i.logger.Log(fmt.Sprintf("failed to read data from file: %s", err.Error()))
		return "", err
	}
	return string(data), nil
}

func (i *index) printAppInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		i.logger.Log("No build info found")
		i.app.Metadata().Custom["beeVersion"] = "unknown"
		i.app.Metadata().Custom["beeliteVersion"] = "unknown"
	} else {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/ethersphere/bee/v2" {
				i.app.Metadata().Custom["beeVersion"] = dep.Version
			} else if dep.Path == "github.com/Solar-Punk-Ltd/bee-lite" {
				i.app.Metadata().Custom["beeliteVersion"] = dep.Version
			}
		}
	}
	meta := i.app.Metadata()

	i.logger.Log(fmt.Sprintf("\nApp ID: %s\n"+
		"App Name: %s\n"+
		"App Version: %s\n"+
		"Released: %t\n"+
		"Build number: %d\n"+
		"Commithash: %s\n"+
		"Bee version: %s\n"+
		"Bee-lite version: %s",
		meta.ID,
		meta.Name,
		meta.Version,
		meta.Release,
		meta.Build,
		meta.Custom["commithash"],
		meta.Custom["beeVersion"],
		meta.Custom["beeliteVersion"]))
}

func (i *index) showProgressWithMessage(message string) {
	i.progress = dialog.NewCustomWithoutButtons(message, widget.NewProgressBarInfinite(), i)
	i.progress.Show()
}

func (i *index) hideProgress() {
	i.progress.Hide()
}

func (i *index) showError(err error) {
	label := widget.NewLabel(err.Error())
	label.Wrapping = fyne.TextWrapWord
	d := dialog.NewCustom("Error", "       Close       ", label, i.Window)
	parentSize := i.Window.Canvas().Size()
	d.Resize(fyne.NewSize(parentSize.Width*90/100, 0))
	d.Show()
}

func (i *index) copyButton(s string) *widget.Button {
	return widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
		i.Window.Clipboard().SetContent(s)
	})
}

func (i *index) showErrorWithAddr(addr common.Address, err error) {
	header := container.NewHBox(widget.NewLabel(shortenHashOrAddress(addr.String())), i.copyButton(addr.String()))
	label := widget.NewLabel(err.Error())
	label.Wrapping = fyne.TextWrapWord
	content := container.NewBorder(header, label, nil, nil)
	d := dialog.NewCustom("Error", "       Close       ", content, i.Window)
	parentSize := i.Window.Canvas().Size()
	d.Resize(fyne.NewSize(parentSize.Width*90/100, 0))
	d.Show()
}

func shortenHashOrAddress(item string) string {
	return fmt.Sprintf("%s[...]%s", item[0:6], item[len(item)-6:])
}

func (i *index) copyDialog(info, data string) fyne.CanvasObject {
	return container.NewStack(container.NewBorder(nil, nil, nil, i.copyButton(data), widget.NewLabel(info)))
}

func (i *index) getPreferenceString(key string) string {
	if !i.nodeConfig.isKeyStoreMem {
		return i.app.Preferences().String(key)
	}
	return ""
}

func (i *index) getPreferenceBool(key string) bool {
	if !i.nodeConfig.isKeyStoreMem {
		return i.app.Preferences().Bool(key)
	}
	return false
}

func (i *index) setPreference(key string, value interface{}) {
	if !i.nodeConfig.isKeyStoreMem {
		switch valueType := value.(type) {
		case string:
			i.app.Preferences().SetString(key, valueType)
		case []string:
			i.app.Preferences().SetStringList(key, valueType)
		case bool:
			i.app.Preferences().SetBool(key, valueType)
		case []bool:
			i.app.Preferences().SetBoolList(key, valueType)
		case int:
			i.app.Preferences().SetInt(key, valueType)
		case []int:
			i.app.Preferences().SetIntList(key, valueType)
		case float64:
			i.app.Preferences().SetFloat(key, valueType)
		case []float64:
			i.app.Preferences().SetFloatList(key, valueType)
		case nil:
		default:
			i.logger.Log(fmt.Sprintf("Invalid type for preference: %T", value))
		}
	}
}
