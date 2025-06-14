package screens

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	beelite "github.com/Solar-Punk-Ltd/bee-lite"
	"github.com/ethersphere/bee/v2/pkg/api"
)

const (
	TestnetChainID        = 11155111
	TestnetNetworkID      = uint64(10)
	MainnetChainID        = 100
	MainnetNetworkID      = uint64(1)
	NativeTokenSymbol     = "xDAI"
	SwarmTokenSymbol      = "xBZZ"
	defaultRPC            = "https://gnosis.publicnode.com"
	defaultTestRPC        = "https://eth-sepolia.g.alchemy.com/v2/atcICv4EFi9hXKew1D4LvnH36cm5-96S"
	defaultWelcomeMsg     = "Welcome from ACTivate!"
	defaultPassword       = "defaultpassword"
	defaultNatAddress     = ""
	defaultSwapEnable     = true
	infoLogLevel          = "3"
	defaultDepth          = "21"
	defaultAmount         = "500000000"
	defaultImmutable      = true
	passwordPrefKey       = "password"
	welcomeMessagePrefKey = "welcomeMessage"
	swapEnablePrefKey     = "swapEnable"
	natAddressPrefKey     = "natAddress"
	rpcEndpointPrefKey    = "rpcEndpoint"
	selectedStampPrefKey  = "selected_stamp"
	batchPrefKey          = "batch"
	uploadsPrefKey        = "uploads"
	overlayAddrPrefKey    = "overlayAddress"
)

var (
	MainnetBootnodes = []string{
		"/dnsaddr/mainnet.ethswarm.org",
	}

	TestnetBootnodes = []string{
		"/dnsaddr/testnet.ethswarm.org",
	}
)

type logger struct{}

func (*logger) Write(p []byte) (int, error) {
	log.Println(string(p))
	return len(p), nil
}

func (*logger) Log(s string) {
	log.Println(s)
}

type index struct {
	fyne.Window
	app        fyne.App
	view       *fyne.Container
	content    *fyne.Container
	intro      *widget.Label
	progress   dialog.Dialog
	bl         *beelite.Beelite
	logger     *logger
	nodeConfig *nodeConfig
}

func Make(a fyne.App, w fyne.Window) fyne.CanvasObject {
	i := &index{
		Window:     w,
		app:        a,
		intro:      widget.NewLabel("ACTivate"),
		logger:     &logger{},
		nodeConfig: &nodeConfig{},
	}
	i.intro.Wrapping = fyne.TextWrapWord
	i.printAppInfo()

	i.nodeConfig.isKeyStoreMem = a.Driver().Device().IsBrowser()
	if i.nodeConfig.isKeyStoreMem {
		i.logger.Log("Running in browser, using in-memory keystore")
	} else {
		i.nodeConfig.path = a.Storage().RootURI().Path()
		i.logger.Log("App datadir path: " + i.nodeConfig.path)
	}

	i.nodeConfig.welcomeMessage = defaultWelcomeMsg
	i.nodeConfig.password = defaultPassword
	i.nodeConfig.natAddress = defaultNatAddress
	i.nodeConfig.rpcEndpoint = defaultRPC
	i.nodeConfig.swapEnable = defaultSwapEnable

	i.view = container.NewBorder(container.NewVBox(i.intro), nil, nil, nil, container.NewStack(i.showStartView(false)))
	i.view.Refresh()
	return i.view
}

func (i *index) start(path, password, welcomeMessage, natAddress, rpcEndpoint string, swapEnable bool) {
	if password == "" {
		i.showError(fmt.Errorf("password cannot be blank"))
		return
	}
	i.showProgressWithMessage("Starting Bee")

	err := i.initSwarm(path, welcomeMessage, password, natAddress, rpcEndpoint, swapEnable)
	i.hideProgress()
	if err != nil {
		if i.bl != nil {
			i.showErrorWithAddr(i.bl.OverlayEthAddress(), err)
		} else {
			i.showError(err)
		}
		return
	}

	if swapEnable {
		if i.bl.BeeNodeMode() != api.LightMode {
			i.showError(fmt.Errorf("swap is enabled but the current node mode is: %s", i.bl.BeeNodeMode()))
			return
		}
	} else if i.bl.BeeNodeMode() != api.UltraLightMode {
		i.showError(fmt.Errorf("swap disabled but the current node mode is: %s", i.bl.BeeNodeMode()))
		return
	}

	i.setPreference(welcomeMessagePrefKey, welcomeMessage)
	i.setPreference(swapEnablePrefKey, swapEnable)
	i.setPreference(natAddressPrefKey, natAddress)
	i.setPreference(rpcEndpointPrefKey, rpcEndpoint)
	i.loadMenuView()
	i.intro.SetText("")
	i.intro.Hide()
}

func (i *index) initSwarm(dataDir, welcomeMessage, password, natAddress, rpcEndpoint string, swapEnable bool) error {
	i.logger.Log(welcomeMessage)

	// isMainnet := rpcEndpoint == defaultRPC
	isMainnet := true
	networkID := MainnetNetworkID
	if !isMainnet {
		networkID = TestnetNetworkID
	}

	lo := &beelite.LiteOptions{
		FullNodeMode:             false,
		BootnodeMode:             false,
		Bootnodes:                MainnetBootnodes,
		DataDir:                  dataDir,
		WelcomeMessage:           welcomeMessage,
		BlockchainRpcEndpoint:    rpcEndpoint,
		SwapInitialDeposit:       "0",
		PaymentThreshold:         "100000000",
		SwapEnable:               swapEnable,
		ChequebookEnable:         true,
		UsePostageSnapshot:       false,
		Mainnet:                  isMainnet,
		NetworkID:                networkID,
		NATAddr:                  natAddress,
		CacheCapacity:            32 * 1024 * 1024,
		DBOpenFilesLimit:         50,
		DBWriteBufferSize:        32 * 1024 * 1024,
		DBBlockCacheCapacity:     32 * 1024 * 1024,
		DBDisableSeeksCompaction: false,
		RetrievalCaching:         true,
	}

	bl, err := beelite.Start(lo, password, infoLogLevel)
	if err != nil {
		return err
	}

	i.setPreference(passwordPrefKey, password)
	i.setPreference(overlayAddrPrefKey, bl.OverlayEthAddress().String())
	i.bl = bl
	return err
}

func (i *index) loadMenuView() {
	// only show certain views if the node mode is NOT ultra-light
	ultraLightMode := i.bl.BeeNodeMode() == api.UltraLightMode
	infoCard := i.showInfoCard(ultraLightMode)
	menuContent := container.NewGridWithColumns(1, infoCard)
	granteeList := i.showListCard()
	menuContent.Add(granteeList)

	// downloadCard := i.showDownloadCard()
	// menuContent.Add(downloadCard)

	i.content.Objects = []fyne.CanvasObject{container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewScroll(menuContent)),
	}
	i.content.Refresh()
}
