package screens

import (
	"context"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	beelite "github.com/Solar-Punk-Ltd/bee-lite"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethersphere/bee/v2/pkg/api"
	"github.com/ethersphere/bee/v2/pkg/transaction" // For transaction.Service, though might be nil
)

const (
	TestnetChainID         = 11155111
	TestnetNetworkID       = uint64(10)
	MainnetChainID         = 100
	MainnetNetworkID       = uint64(1)
	NativeTokenSymbol      = "xDAI"
	SwarmTokenSymbol       = "xBZZ"
	defaultRPC             = "wss://gnosis-mainnet.g.alchemy.com/v2/YtM4LIorMJrGNRWkvAOFWSKTDzhNsCMz"
	defaultTestRPC         = "https://eth-sepolia.g.alchemy.com/v2/atcICv4EFi9hXKew1D4LvnH36cm5-96S"
	defaultWelcomeMsg      = "Welcome from ACTivate!"
	defaultPassword        = "defaultpassword"
	defaultNatAddress      = ""
	defaultSwapEnable      = true
	dataContractAddressHex = "0x242A2174fa8d8586a784aBdB4fF03C3181E96bee"
	infoLogLevel           = "3"
	defaultDepth           = "21"
	defaultAmount          = "500000000"
	defaultImmutable       = true
	passwordPrefKey        = "password"
	welcomeMessagePrefKey  = "welcomeMessage"
	swapEnablePrefKey      = "swapEnable"
	natAddressPrefKey      = "natAddress"
	rpcEndpointPrefKey     = "rpcEndpoint"
	selectedStampPrefKey   = "selected_stamp"
	batchPrefKey           = "batch"
	uploadsPrefKey         = "uploads"
	overlayAddrPrefKey     = "overlayAddress"
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

	ethClient            *ethclient.Client
	contractSvc          DataContractInterface
	dataContractABI      abi.ABI // Store the parsed ABI here
	eventLogSubscription ethereum.Subscription
	eventMessageLabel    *widget.Label
}

func (i *index) initContract(txService transaction.Service) {
	rpcEndpoint := defaultRPC
	var err error
	i.ethClient, err = ethclient.DialContext(context.Background(), rpcEndpoint)
	if err != nil {
		i.logger.Log(fmt.Sprintf("Failed to connect to Ethereum client via %s: %v", rpcEndpoint, err))
	}

	i.dataContractABI, err = ParseContractABI()
	if err != nil {
		i.logger.Log(fmt.Sprintf("Failed to parse data contract ABI JSON: %v", err))
		i.showError(err)
	}

	dataContractAddr := common.HexToAddress(dataContractAddressHex)

	if i.dataContractABI.Events != nil { // Check if ABI was parsed and has events
		i.contractSvc = NewDataContract(
			i.bl.OverlayEthAddress(),
			dataContractAddr,
			i.dataContractABI,
			txService,
			true, // setGasLimit
		)
	} else {
		i.logger.Log("Data contract ABI not parsed or no events found, contractSvc not initialized.")
	}

	// Initialize UI elements for events
	i.eventMessageLabel = widget.NewLabel("Initializing event listener...")
	i.eventMessageLabel.Wrapping = fyne.TextWrapWord
	i.eventMessageLabel.Alignment = fyne.TextAlignCenter
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

	i.initContract(i.bl.TransactionService())

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

	// Use a VBox for menuContent to stack items vertically
	menuContent := container.NewVBox(infoCard)

	granteeList := i.showListCard()
	menuContent.Add(granteeList)

	// downloadCard := i.showDownloadCard()
	// menuContent.Add(downloadCard)

	// Add the event message label to the view
	if i.eventMessageLabel != nil {
		menuContent.Add(i.eventMessageLabel)
	} else {
		// Fallback, though it should be initialized in Make
		i.logger.Log("eventMessageLabel is nil in loadMenuView")
		menuContent.Add(widget.NewLabel("Event display not initialized."))
	}

	i.setupDataContractSubscription()

	i.content.Objects = []fyne.CanvasObject{container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewScroll(menuContent)),
	}
	i.content.Refresh()
}

func (i *index) setupDataContractSubscription() {
	if i.eventLogSubscription != nil {
		i.eventLogSubscription.Unsubscribe() // Unsubscribe from previous if any
	}

	logs := make(chan types.Log)
	var err error

	// Use a new context for the subscription goroutine, or manage it with the app's lifecycle
	subCtx, cancelSubCtx := context.WithCancel(context.Background())
	i.Window.SetOnClosed(func() { // Ensure cancellation when window closes
		cancelSubCtx()
	})

	i.eventLogSubscription, err = i.contractSvc.SubscribeDataSentToTarget(subCtx, i.ethClient, logs)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to subscribe to DataSentToTarget: %v", err)
		i.logger.Log(errMsg)
		if i.eventMessageLabel != nil {
			i.eventMessageLabel.SetText(errMsg)
		}
		cancelSubCtx() // Cancel context if subscription fails
		return
	}

	i.logger.Log("Successfully subscribed to DataSentToTarget events.")
	if i.eventMessageLabel != nil {
		i.eventMessageLabel.SetText("Subscribed. Waiting for 'DataSentToTarget' events...")
	}

	go func() {
		defer func() {
			if i.eventLogSubscription != nil {
				i.eventLogSubscription.Unsubscribe()
			}
			cancelSubCtx() // Ensure context is cancelled when goroutine exits
			i.logger.Log("Event listener goroutine stopped.")
		}()

		for {
			select {
			case <-subCtx.Done():
				i.logger.Log("Event listener context cancelled. Unsubscribing.")
				return
			case err := <-i.eventLogSubscription.Err():
				errMsg := fmt.Sprintf("Event subscription error: %v", err)
				i.logger.Log(errMsg)
				if i.eventMessageLabel != nil {
					i.eventMessageLabel.SetText("Subscription error. Check logs.")
				}
				// Depending on the error, you might want to attempt to resubscribe.
				// For now, we stop listening on error.
				return
			case vLog := <-logs:
				i.logger.Log(fmt.Sprintf("Received log: Block %d, TxHash %s, Topics %d, Data %d bytes", vLog.BlockNumber, vLog.TxHash.Hex(), len(vLog.Topics), len(vLog.Data)))

				eventName := "DataSentToTarget"
				eventAbi, ok := i.dataContractABI.Events[eventName]
				if !ok {
					i.logger.Log(fmt.Sprintf("Event %s not found in ABI. Cannot parse.", eventName))
					continue
				}

				// Check if this log is indeed for DataSentToTarget based on Topic[0]
				if len(vLog.Topics) == 0 || vLog.Topics[0] != eventAbi.ID {
					i.logger.Log(fmt.Sprintf("Received log does not match %s event signature. Skipping. Log Topic0: %s, Expected: %s", eventName, vLog.Topics[0].Hex(), eventAbi.ID.Hex()))
					continue
				}
				i.logger.Log(fmt.Sprintf("Processing '%s' event...", eventName))

				var targetAddr common.Address
				var ownerBytes []byte
				var actRefBytes []byte
				var topicString string

				// Unpack indexed fields from Topics
				topicIdx := 1 // Topics[0] is the event signature itself
				for _, input := range eventAbi.Inputs {
					if input.Indexed {
						if topicIdx < len(vLog.Topics) {
							if input.Name == "target" { // Assuming 'target' is the name in ABI
								targetAddr = common.BytesToAddress(vLog.Topics[topicIdx].Bytes())
							}
							// Add other indexed fields here if any, by checking input.Name or type
							topicIdx++
						} else {
							i.logger.Log(fmt.Sprintf("Warning: Mismatch count for indexed ABI inputs and log topics for event %s.", eventName))
							break
						}
					}
				}

				// Prepare to unpack non-indexed fields from Data
				var nonIndexedArgs abi.Arguments
				for _, input := range eventAbi.Inputs {
					if !input.Indexed {
						nonIndexedArgs = append(nonIndexedArgs, input)
					}
				}

				if len(nonIndexedArgs) > 0 {
					unpackedData, err := nonIndexedArgs.Unpack(vLog.Data)
					if err != nil {
						i.logger.Log(fmt.Sprintf("Failed to unpack non-indexed data for event %s: %v", eventName, err))
					} else {
						// Assign to variables based on the order of non-indexed args in ABI
						// This needs to be robust and match your ABI definition precisely.
						// Example: owner, actRef, topic
						currentUnpackedIdx := 0
						for _, arg := range nonIndexedArgs {
							if currentUnpackedIdx >= len(unpackedData) {
								break
							}
							switch arg.Name { // Or rely on order if names are not set/unique
							case "owner":
								if val, ok := unpackedData[currentUnpackedIdx].([]byte); ok {
									ownerBytes = val
								}
							case "actRef":
								if val, ok := unpackedData[currentUnpackedIdx].([]byte); ok {
									actRefBytes = val
								}
							case "topic":
								if val, ok := unpackedData[currentUnpackedIdx].(string); ok {
									topicString = val
								}
							}
							currentUnpackedIdx++
						}
					}
				}

				parsedMsg := fmt.Sprintf("'DataSentToTarget' Event! Block: %d.", vLog.BlockNumber)
				if targetAddr != (common.Address{}) {
					parsedMsg += fmt.Sprintf(" Target: %s.", targetAddr.Hex())
				}
				if len(ownerBytes) > 0 {
					parsedMsg += fmt.Sprintf(" Owner: 0x%x.", ownerBytes)
				}
				if len(actRefBytes) > 0 {
					parsedMsg += fmt.Sprintf(" ActRef: 0x%x.", actRefBytes)
				}
				if topicString != "" {
					parsedMsg += fmt.Sprintf(" Topic: '%s'.", topicString)
				}

				i.logger.Log("Formatted event message: " + parsedMsg)
				if i.eventMessageLabel != nil {
					i.eventMessageLabel.SetText(parsedMsg)
				}
			}
		}
	}()
}
