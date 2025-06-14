package screens

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethersphere/bee/v2/pkg/api"
)

type nodeConfig struct {
	path           string
	password       string
	welcomeMessage string
	swapEnable     bool
	natAddress     string
	rpcEndpoint    string
	isKeyStoreMem  bool
}

func (i *index) showPasswordView() fyne.CanvasObject {
	i.intro.SetText("Initialise your swarm node with a strong password")
	content := container.NewStack()
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")
	nextButton := widget.NewButton("Next", func() {
		if passwordEntry.Text == "" {
			i.showError(fmt.Errorf("password cannot be blank"))
			return
		}

		i.nodeConfig.password = passwordEntry.Text
		content.Objects = []fyne.CanvasObject{i.showWelcomeMessageView()}
		content.Refresh()
	})
	nextButton.Importance = widget.HighImportance
	content.Objects = []fyne.CanvasObject{container.NewBorder(passwordEntry, nextButton, nil, nil)}
	i.content = content
	i.view = container.NewBorder(container.NewVBox(i.intro), nil, nil, nil, content)
	i.view.Refresh()
	return content
}

func setPlaceHolderText(config, defaultText string) string {
	placeHolder := config
	if placeHolder == "" {
		placeHolder = defaultText
	}
	return placeHolder
}

func (i *index) showWelcomeMessageView() fyne.CanvasObject {
	i.intro.SetText("Set the welcome message of your node (optional)")
	content := container.NewStack()
	welcomeMessageEntry := widget.NewEntry()
	welcomeMessageEntry.SetPlaceHolder(setPlaceHolderText(i.nodeConfig.welcomeMessage, defaultWelcomeMsg))

	nextButton := widget.NewButton("Next", func() {
		if welcomeMessageEntry.Text == "" {
			welcomeMessageEntry.SetText(defaultWelcomeMsg)
			i.logger.Log(fmt.Sprintf("Welcome message is blank, using default: %s", defaultWelcomeMsg))
		} else {
			i.logger.Log(fmt.Sprintf("Welcome message is: %s", welcomeMessageEntry.Text))
		}

		i.nodeConfig.welcomeMessage = welcomeMessageEntry.Text
		content.Objects = []fyne.CanvasObject{i.showNodeModeSelectionView()}
		content.Refresh()
	})

	backButton := widget.NewButton("Back", func() {
		content.Objects = []fyne.CanvasObject{i.showPasswordView()}
		content.Refresh()
	})
	backButton.Importance = widget.WarningImportance
	nextButton.Importance = widget.HighImportance
	content.Objects = []fyne.CanvasObject{container.NewBorder(welcomeMessageEntry, container.NewVBox(nextButton, backButton), nil, nil)}
	i.content = content
	i.view = container.NewBorder(container.NewVBox(i.intro), nil, nil, nil, content)
	i.view.Refresh()

	return content
}

func (i *index) getNodeModeRadio() *widget.RadioGroup {
	return widget.NewRadioGroup(
		[]string{api.LightMode.String(), api.UltraLightMode.String()},
		func(mode string) {
			if mode == api.LightMode.String() {
				i.nodeConfig.swapEnable = true
			} else {
				i.nodeConfig.swapEnable = false
			}
			i.logger.Log(fmt.Sprintf("Node mode selected: %s", mode))
		},
	)
}

func (i *index) showNodeModeSelectionView() fyne.CanvasObject {
	i.intro.SetText("Choose the type of your node")
	content := container.NewStack()
	nodeModeRadio := i.getNodeModeRadio()
	nextButton := widget.NewButton("Next", func() {
		if nodeModeRadio.Selected == "" {
			i.showError(fmt.Errorf("please select the node mode"))
			return
		}

		i.logger.Log(fmt.Sprintf("SWAP enable: %t, running in %s mode", i.nodeConfig.swapEnable, nodeModeRadio.Selected))
		content.Objects = []fyne.CanvasObject{i.showNATAddressView()}
		content.Refresh()
	})

	backButton := widget.NewButton("Back", func() {
		content.Objects = []fyne.CanvasObject{i.showWelcomeMessageView()}
		content.Refresh()
	})
	backButton.Importance = widget.WarningImportance
	nextButton.Importance = widget.HighImportance

	overlayAddr := i.getPreferenceString(overlayAddrPrefKey)
	infoLabel := widget.NewLabel("Overlay address is not saved, running in ultra-light mode")
	infoBox := container.NewVBox(infoLabel)
	if overlayAddr == "" {
		nodeModeRadio.SetSelected(api.UltraLightMode.String())
		nodeModeRadio.Disable()
	} else {
		addrCopyButton := i.copyButton(overlayAddr)
		infoLabel.SetText(shortenHashOrAddress(overlayAddr))
		infoBox.Add(addrCopyButton)
	}
	content.Objects = []fyne.CanvasObject{container.NewBorder(nodeModeRadio, container.NewVBox(nextButton, backButton), infoBox, nil)}
	i.content = content
	i.view = container.NewBorder(container.NewVBox(i.intro), nil, nil, nil, content)
	i.view.Refresh()

	return content
}

func (i *index) showNATAddressView() fyne.CanvasObject {
	i.intro.SetText("Set the NAT Address of your node (optional)")
	content := container.NewStack()
	natAdrrEntry := widget.NewEntry()
	natAdrrEntry.SetPlaceHolder(setPlaceHolderText(i.nodeConfig.natAddress, "123.123.123.123:1634"))
	nextButton := widget.NewButton("Next", func() {
		if natAdrrEntry.Text == "" {
			i.logger.Log("NAT address is blank")
		} else {
			i.logger.Log(fmt.Sprintf("Using NAT address: %s", natAdrrEntry.Text))
		}
		i.nodeConfig.natAddress = natAdrrEntry.Text

		// in ultra-light mode no RPC endpoint is necessary
		if i.nodeConfig.swapEnable {
			content.Objects = []fyne.CanvasObject{i.showRPCView()}
		} else {
			content.Objects = []fyne.CanvasObject{i.showStartView(true)}
		}

		content.Refresh()
	})

	backButton := widget.NewButton("Back", func() {
		content.Objects = []fyne.CanvasObject{i.showNodeModeSelectionView()}
		content.Refresh()
	})
	backButton.Importance = widget.WarningImportance
	nextButton.Importance = widget.HighImportance
	content.Objects = []fyne.CanvasObject{container.NewBorder(natAdrrEntry, container.NewVBox(nextButton, backButton), nil, nil)}
	i.content = content
	i.view = container.NewBorder(container.NewVBox(i.intro), nil, nil, nil, content)
	i.view.Refresh()

	return content
}

func (i *index) verifyRPCConnection(rpcEndpoint string) error {
	i.logger.Log(fmt.Sprintf("verifying RPC endpoint connection: %s", rpcEndpoint))
	// test endpoint is connectable
	eth, err := ethclient.Dial(rpcEndpoint)
	if err != nil {
		return fmt.Errorf("rpc endpoint is invalid or not reachable: %w", err)
	}
	// check connections
	_, err = eth.ChainID(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (i *index) showRPCView() fyne.CanvasObject {
	i.intro.SetText("Swarm mobile needs a RPC endpoint to start")
	content := container.NewStack()
	rpcEntry := widget.NewEntry()
	rpcEntry.SetPlaceHolder(setPlaceHolderText(i.nodeConfig.rpcEndpoint, defaultRPC))

	nextButton := widget.NewButton("Next", func() {
		if i.nodeConfig.swapEnable {
			if rpcEntry.Text == "" {
				rpcEntry.SetText(defaultRPC)
				i.logger.Log(fmt.Sprintf("RPC endpoint is blank, using default RPC: %s", defaultRPC))
			}
			i.nodeConfig.rpcEndpoint = rpcEntry.Text
		}
		content.Objects = []fyne.CanvasObject{i.showStartView(true)}
		content.Refresh()
	})

	backButton := widget.NewButton("Back", func() {
		content.Objects = []fyne.CanvasObject{i.showNATAddressView()}
		content.Refresh()
	})
	backButton.Importance = widget.WarningImportance
	nextButton.Importance = widget.HighImportance
	content.Objects = []fyne.CanvasObject{container.NewBorder(rpcEntry, container.NewVBox(nextButton, backButton), nil, nil)}
	i.content = content
	i.view = container.NewBorder(container.NewVBox(i.intro), nil, nil, nil, content)
	i.view.Refresh()

	return content
}

func (i *index) showStartView(firstStart bool) fyne.CanvasObject {
	i.intro.SetText("Start your Swarm node")
	i.intro.TextStyle.Bold = true
	content := container.NewStack()
	overlayAddr := i.getPreferenceString(overlayAddrPrefKey)

	startButton := widget.NewButton("Start", func() {
		if i.nodeConfig.path == "" && !i.nodeConfig.isKeyStoreMem {
			i.showError(fmt.Errorf("invalid app storage path"))
			return
		}

		if i.nodeConfig.password == "" {
			i.showError(fmt.Errorf("password is empty"))
			return
		}

		if i.nodeConfig.swapEnable {
			if i.nodeConfig.rpcEndpoint == "" {
				i.showError(fmt.Errorf("rpc endpoint is required in light mode"))
				return
			}
			err := i.verifyRPCConnection(i.nodeConfig.rpcEndpoint)
			if err != nil {
				i.logger.Log(fmt.Sprintf("rpc endpoint error: %s", err.Error()))
				i.showError(err)
				return
			}
			if overlayAddr == "" {
				i.showError(fmt.Errorf("Overlay address is not saved, need to start in ultra-light mode first"))
				return
			}
		} else {
			if i.nodeConfig.rpcEndpoint != "" {
				i.showError(fmt.Errorf("rpc endpoint must be empty in ultra-light mode"))
				return
			}
		}

		i.start(i.nodeConfig.path,
			i.nodeConfig.password,
			i.nodeConfig.welcomeMessage,
			i.nodeConfig.natAddress,
			i.nodeConfig.rpcEndpoint,
			i.nodeConfig.swapEnable)
		content.Refresh()
	})

	startButton.Importance = widget.HighImportance

	bottomBox := container.NewVBox()
	if overlayAddr != "" {
		infoLabel := widget.NewLabel(fmt.Sprintf("Warning: cannot continue in light-mode until there is\nat least min %s (for Gas) and at least min %s available on\n address: %s",
			NativeTokenSymbol, SwarmTokenSymbol, shortenHashOrAddress(overlayAddr)))
		bottomBox.Add(infoLabel)
		bottomBox.Add(i.copyButton(overlayAddr))
	}

	advancedView := i.showAdvancedSettings()
	content.Objects = []fyne.CanvasObject{container.NewBorder(startButton, bottomBox, advancedView, nil)}
	i.content = content
	i.view = container.NewBorder(container.NewVBox(i.intro), nil, nil, nil, content)
	i.view.Refresh()

	return content
}

func (i *index) showAdvancedSettings() fyne.CanvasObject {
	hyperlink := widget.NewHyperlink("How to fund your node", nil)
	err := hyperlink.SetURLFromString("https://docs.ethswarm.org/docs/installation/fund-your-node")
	if err != nil {
		i.logger.Log(fmt.Sprintf("failed to set hyperlink: %s", err.Error()))
	}
	modeDetail := container.NewVBox(i.getNodeModeRadio(), container.NewHBox(hyperlink))
	modeSwitchItem := &widget.AccordionItem{
		Title:  "Node mode",
		Detail: modeDetail,
		Open:   false,
	}
	if i.getPreferenceString(overlayAddrPrefKey) == "" {
		modeSwitchItem.Detail.Hide()
	}

	welcomeBind := binding.BindString(&i.nodeConfig.welcomeMessage)
	welcomeEntry := widget.NewEntryWithData(welcomeBind)
	welcomeEntry.OnChanged = func(s string) {
		err := welcomeBind.Set(s)
		if err != nil {
			i.logger.Log(fmt.Sprintf("failed to bind welcome message: %s", err.Error()))
		}
	}
	welcomeEntry.SetPlaceHolder(setPlaceHolderText(i.nodeConfig.welcomeMessage, defaultWelcomeMsg))
	welcomeMsgItem := &widget.AccordionItem{
		Title:  "Welcome message",
		Detail: welcomeEntry,
		Open:   false,
	}

	natAddrBind := binding.BindString(&i.nodeConfig.natAddress)
	natAdrrEntry := widget.NewEntryWithData(natAddrBind)
	natAdrrEntry.OnChanged = func(s string) {
		err := natAddrBind.Set(s)
		if err != nil {
			i.logger.Log(fmt.Sprintf("failed to bind NAT address: %s", err.Error()))
		}
	}
	natAdrrEntry.SetPlaceHolder(setPlaceHolderText(i.nodeConfig.natAddress, "123.123.123.123:1634"))
	natAddrItem := &widget.AccordionItem{
		Title:  "NAT Address",
		Detail: natAdrrEntry,
		Open:   false,
	}

	rpcBind := binding.BindString(&i.nodeConfig.rpcEndpoint)
	rpcEntry := widget.NewEntryWithData(rpcBind)
	rpcEntry.OnChanged = func(s string) {
		err := rpcBind.Set(s)
		if err != nil {
			i.logger.Log(fmt.Sprintf("failed to bind rpc endpoint: %s", err.Error()))
		}
	}
	rpcEntry.SetPlaceHolder(setPlaceHolderText(i.nodeConfig.rpcEndpoint, defaultRPC))
	rpcEndpointItem := &widget.AccordionItem{
		Title:  "RPC Endpoint",
		Detail: rpcEntry,
		Open:   false,
	}

	return container.NewBorder(container.NewVBox(
		widget.NewAccordion(modeSwitchItem, welcomeMsgItem, rpcEndpointItem, natAddrItem)),
		nil, nil, nil)
}
