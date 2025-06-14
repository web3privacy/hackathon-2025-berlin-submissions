package screens

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (i *index) showInfoCard(ultraLightMode bool) *widget.Card {
	addressContent := i.addressContent()
	walletDataButton := i.walletDataButton()
	infoContent := container.NewVBox(addressContent)
	if !ultraLightMode {
		// batchRadio := i.batchRadio()
		// stampsContent := i.stampsContent(batchRadio)
		// buyBatchButton := i.buyBatchButton(batchRadio)
		balanceContent := i.balanceContent()
		infoContent = container.NewVBox(addressContent, balanceContent, nil, nil)
	}
	infoContent.Add(walletDataButton)

	infoCard := widget.NewCard("Info",
		fmt.Sprintf("Connected with %d peers", i.bl.ConnectedPeerCount()), infoContent)

	// auto reload
	go func() {
		for {
			time.Sleep(time.Second * 5)
			if i.bl != nil {
				infoCard.SetSubTitle(fmt.Sprintf("Connected with %d peers", i.bl.ConnectedPeerCount()))
			}
		}
	}()

	return infoCard
}

func (i *index) walletDataButton() *widget.Button {
	button := widget.NewButton("Read wallet data ", func() {
		key, err := i.readAppData("/keys/swarm.key")
		if err != nil {
			i.showError(fmt.Errorf("failed to read swarm.key: %s", err.Error()))
			return
		}

		d := dialog.NewCustomConfirm("Wallet data", "Ok", "Cancel", i.copyDialog("Do not share this with anyone!\nCopy and save your key file.", key), func(b bool) {}, i.Window)
		d.Show()
	})
	button.Importance = widget.DangerImportance

	return button
}

func (i *index) addressContent() *fyne.Container {
	addrCopyButton := i.copyButton(i.bl.OverlayEthAddress().String())
	addrHeader := container.NewHBox(widget.NewLabel("Overlay address:"))
	addr := container.NewHBox(
		widget.NewLabel(i.bl.OverlayEthAddress().String()),
		addrCopyButton,
	)
	return container.NewVBox(addrHeader, addr)
}

func (i *index) balanceContent() *fyne.Container {
	chequebookBalance, err := i.bl.ChequebookBalance()
	if err != nil {
		i.logger.Log(fmt.Sprintf("Cannot get chequebook balance: %s", err.Error()))
		return container.NewHBox(widget.NewLabel("Cannot get chequebook balance"))
	}

	balanceContent := container.NewHBox(widget.NewLabel(
		fmt.Sprintf("Chequebook balance: %s %s", chequebookBalance.String(), SwarmTokenSymbol)))

	// auto reload
	go func() {
		for {
			time.Sleep(time.Second * 60)
			chequebookBalance, err := i.bl.ChequebookBalance()
			if err != nil {
				i.logger.Log(fmt.Sprintf("Cannot get chequebook balance: %s", err.Error()))
			} else {
				balanceContent = container.NewHBox(widget.NewLabel(
					fmt.Sprintf("Chequebook balance: %s %s", chequebookBalance.String(), SwarmTokenSymbol)))
			}
		}
	}()

	return balanceContent
}

func (i *index) stampsContent(batchRadio *widget.RadioGroup) *fyne.Container {
	stampsHeader := container.NewHBox(widget.NewLabel("Postage stamps:"))
	stamps := i.bl.GetUsableBatches()

	if len(stamps) != 0 {
		selectedStamp := i.getPreferenceString(selectedStampPrefKey)
		for _, v := range stamps {
			batchRadio.Append(shortenHashOrAddress(hex.EncodeToString(v.ID())))
		}

		batchRadio.SetSelected(selectedStamp)
	}
	return container.NewVBox(stampsHeader, batchRadio)
}

func (i *index) batchRadio() *widget.RadioGroup {
	return widget.NewRadioGroup([]string{}, func(s string) {
		if s == "" {
			i.setPreference(selectedStampPrefKey, "")
			i.setPreference(batchPrefKey, "")
			return
		}
		batches := i.bl.GetUsableBatches()
		for _, v := range batches {
			stamp := hex.EncodeToString(v.ID())
			if s[0:6] == stamp[0:6] {
				i.setPreference(selectedStampPrefKey, s)
				i.setPreference(batchPrefKey, stamp)
			}
		}
	})
}

func (i *index) buyBatchButton(batchRadio *widget.RadioGroup) *widget.Button {
	return widget.NewButton("Buy a postage batch", func() {
		child := i.app.NewWindow("Buying a postage batch")
		depthStr := defaultDepth
		amountStr := defaultAmount
		isImmutable := defaultImmutable
		label := ""
		content := container.NewStack()
		buyBatchContent := i.buyBatchForm(&depthStr, &amountStr, &label, &isImmutable)
		size := child.Canvas().Content().MinSize()
		if size.Width < 250 {
			size.Width = 250
		}
		if size.Height < 100 {
			size.Height = 100
		}
		child.Resize(size)

		buyButton := widget.NewButton("Buy", func() {
			amount, ok := big.NewInt(0).SetString(amountStr, 10)
			if !ok {
				i.showError(fmt.Errorf("invalid amountStr"))
				return
			}
			depth, err := strconv.ParseUint(depthStr, 10, 8)
			if err != nil {
				i.showError(fmt.Errorf("invalid depthStr %s", err.Error()))
				return
			}
			child.Close()
			i.showProgressWithMessage(fmt.Sprintf("Buying a postage batch\ndepth: %s, amount: %s, label: \"%s\", immutable: %t", depthStr, amountStr, label, isImmutable))
			hash, id, err := i.bl.BuyStamp(amount, depth, label, isImmutable)
			if err != nil {
				i.hideProgress()
				i.showError(err)
				return
			}
			i.logger.Log(fmt.Sprintf("Batch created: %s", hash.String()))
			i.hideProgress()
			batchRadio.Append(shortenHashOrAddress(hex.EncodeToString(id)))
		})
		buyButton.Importance = widget.HighImportance
		content.Objects = []fyne.CanvasObject{container.NewBorder(buyBatchContent, container.NewVBox(buyButton), nil, nil)}
		child.SetContent(content)
		child.Show()
	})
}

func (i *index) buyBatchForm(depthStr, amountStr, label *string, isImmutable *bool) fyne.CanvasObject {
	depthBind := binding.BindString(depthStr)
	amountBind := binding.BindString(amountStr)
	labelBind := binding.BindString(label)
	immutableBind := binding.BindBool(isImmutable)

	amountEntry := widget.NewEntryWithData(amountBind)
	amountEntry.OnChanged = func(s string) {
		err := amountBind.Set(s)
		if err != nil {
			i.logger.Log(fmt.Sprintf("failed to bind amount: %s", err.Error()))
		}
	}
	depthEntry := widget.NewEntryWithData(depthBind)
	depthEntry.OnChanged = func(s string) {
		err := depthBind.Set(s)
		if err != nil {
			i.logger.Log(fmt.Sprintf("failed to bind depth: %s", err.Error()))
		}
	}
	labelEntry := widget.NewEntryWithData(labelBind)
	labelEntry.OnChanged = func(s string) {
		err := labelBind.Set(s)
		if err != nil {
			i.logger.Log(fmt.Sprintf("failed to bind label: %s", err.Error()))
		}
	}
	labelEntry.SetPlaceHolder("My first batch")
	immutableCheck := widget.NewCheck("Immutable", func(b bool) {
		err := immutableBind.Set(b)
		if err != nil {
			i.logger.Log(fmt.Sprintf("failed to bind immutable: %s", err.Error()))
		}
	})
	immutableCheck.Checked = true

	optionsForm := widget.NewForm()
	optionsForm.Append(
		"Depth",
		container.NewStack(depthEntry),
	)
	optionsForm.Append(
		"Amount",
		container.NewStack(amountEntry),
	)
	optionsForm.Append(
		"Label",
		container.NewStack(labelEntry),
	)
	optionsForm.Append(
		"",
		immutableCheck,
	)

	return container.NewStack(optionsForm)
}
