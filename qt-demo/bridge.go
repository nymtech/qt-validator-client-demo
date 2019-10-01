// bridge.go - bridging qml and go
// Copyright (C) 2019  Jedrzej Stuczynski.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	Curve "github.com/nymtech/amcl/version3/go/amcl/BLS381"
	"github.com/nymtech/nym-validator/client"
	"github.com/nymtech/nym-validator/client/config"
	coconut "github.com/nymtech/nym-validator/crypto/coconut/scheme"
	"github.com/nymtech/nym-validator/crypto/coconut/utils"
	"github.com/nymtech/nym-validator/nym/token"
	"github.com/therecipe/qt/core"
)

const (
	errNotificationTitle  = "Error"
	warnNotificationTitle = "Warning"
	infoNotificationTitle = "Notification"
)

type IssuedCredential struct {
	credential *coconut.Signature
	token      *token.Token
}

//go:generate qtmoc
type ConfigBridge struct {
	core.QObject

	_ string `property:"identifier"`
	_ string `property:"address"`
	_ string `property:"keyfile"`
	_ string `property:"ethereumNode"`
	_ string `property:"nymERC20"`
	_ string `property:"pipeAccount"`
}

//go:generate qtmoc
type QmlBridge struct {
	core.QObject
	cfg            *config.Config
	clientInstance *client.Client
	longtermSecret *Curve.BIG

	_ func()                                                                                        `constructor:"init"`
	_ func(file string)                                                                             `slot:"loadConfig,auto"`
	_ func()                                                                                        `slot:"confirmConfig,auto"`
	_ func(message, title string)                                                                   `signal:"displayNotification"`
	_ func(identifier, address string)                                                              `signal:"newNymValidator"`
	_ func(identifier, address string)                                                              `signal:"newTendermintValidator"`
	_ func(amount string)                                                                           `signal:"updateERC20NymBalance"`
	_ func(amount string)                                                                           `signal:"updateERC20NymBalancePending"`
	_ func() `signal:"ResetWaitingForEthereumLabel"`
	_ func(amount string)                                                                           `signal:"updateNymTokenBalance"`
	_ func(strigifiedSecret string)                                                                 `signal:"updateSecret"`
	_ func(values []string)                                                                         `signal:"populateValueComboBox"`
	_ func(sps []string)                                                                            `signal:"populateSPComboBox"`
	_ func(busyIndicator *core.QObject, mainLayoutObject *core.QObject)                             `slot:"forceUpdateBalances,auto"`
	_ func()                                                                                        `signal:"markSpentCredential"`
	_ func(amount string, busyIndicator *core.QObject, mainLayoutObject *core.QObject)              `slot:"sendToPipeAccount,auto"`
	_ func(amount string, busyIndicator *core.QObject, mainLayoutObject *core.QObject)              `slot:"redeemTokens,auto"`
	_ func(value string, busyIndicator *core.QObject, mainLayoutObject *core.QObject)               `slot:"getCredential,auto"`
	_ func(chosenSP, seqString string, busyIndicator *core.QObject, mainLayoutObject *core.QObject) `slot:"spendCredential,auto"`
	_ func(item CredentialListItem)                                                                 `signal:"addCredentialListItem"`
	_ func()                                                                                        `signal:"showNewKeyDialog"`
	_ func()                                                                                        `slot:"generateNewKey,auto"`
	_ func(accountExists bool)                                                                      `signal:"setAccountStatus"`
	_ func(busyIndicator *core.QObject, mainLayoutObject *core.QObject)                             `slot:"registerAccount,auto"`
	_ func(busyIndicator *core.QObject, mainLayoutObject *core.QObject)                             `slot:"getFaucetNym,auto"`
	_ func(seqString string) string                                                                 `slot:"randomizeCredential,auto"`
}

func enableAllObjects(objs []*core.QObject) {
	for _, obj := range objs {
		obj.SetProperty("enabled", core.NewQVariant1(true))
	}
}

func disableAllObjects(objs []*core.QObject) {
	for _, obj := range objs {
		obj.SetProperty("enabled", core.NewQVariant1(false))
	}
}

func toggleIndicatorAndObjects(indicator *core.QObject, objs []*core.QObject, run bool) {
	if run {
		if indicator != nil {
			indicator.SetProperty("running", core.NewQVariant1(true))
		}
		if len(objs) > 0 {
			disableAllObjects(objs)
		}
	} else {
		if indicator != nil {
			indicator.SetProperty("running", core.NewQVariant1(false))
		}
		if len(objs) > 0 {
			enableAllObjects(objs)
		}
	}
}

func (qb *QmlBridge) DisplayNotificationf(title string, fmtMessage string, a ...interface{}) {
	msg := fmtMessage
	if a != nil {
		msg = fmt.Sprintf(fmtMessage, a...)
	}
	qb.DisplayNotification(msg, title)
}

func (qb *QmlBridge) waitForERC20BalanceChange(ctx context.Context, expectedBalance uint64) {
	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return
	}
	retryTicker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-retryTicker.C:
			currentBalance, err := qb.clientInstance.GetCurrentERC20Balance()
			if err != nil {
				qb.DisplayNotificationf(errNotificationTitle, "failed to query for ERC20 Nym Balance: %v", err)
			} else {
				qb.UpdateERC20NymBalance(strconv.FormatUint(currentBalance, 10))
			}

			pendingBalance, err := qb.clientInstance.GetCurrentERC20PendingBalance()
			if err != nil {
				qb.DisplayNotificationf(errNotificationTitle, "failed to query for ERC20 Nym Balance (pending): %v", err)
			} else {
				qb.UpdateERC20NymBalancePending(strconv.FormatUint(pendingBalance, 10))
			}

			if currentBalance == expectedBalance {
				return
			}
		case <-ctx.Done():
			qb.DisplayNotificationf(errNotificationTitle, "failed to query for obtain current ERC20 balances: ctx timeout")
			return
		}
	}
}

func (qb *QmlBridge) updateBalances() {
	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return
	}
	erc20balance, err := qb.clientInstance.GetCurrentERC20Balance()
	if err != nil {
		qb.DisplayNotificationf(errNotificationTitle, "failed to query for ERC20 Nym Balance: %v", err)
	}
	pending, err := qb.clientInstance.GetCurrentERC20PendingBalance()
	if err != nil {
		qb.DisplayNotificationf(errNotificationTitle, "failed to query for ERC20 Nym Balance (pending): %v", err)
	}
	nymBalance, err := qb.clientInstance.GetCurrentNymBalance()
	if err != nil {
		qb.DisplayNotificationf(errNotificationTitle, "failed to query for Nym Token Balance: %v", err)
	}

	qb.UpdateERC20NymBalance(strconv.FormatUint(erc20balance, 10))
	qb.UpdateERC20NymBalancePending(strconv.FormatUint(pending, 10))
	qb.UpdateNymTokenBalance(strconv.FormatUint(nymBalance, 10))
}

func (qb *QmlBridge) loadConfig(file string) {
	// TODO: is that prefix always added?
	file = strings.TrimPrefix(file, "file://")

	cfg, err := config.LoadFile(file)
	if err != nil {
		qb.DisplayNotificationf(errNotificationTitle, "failed to load config file '%v': %v\n", file, err)
		return
	} else {
		fmt.Println("loaded config!")
	}

	configBridge.SetIdentifier(cfg.Client.Identifier)
	configBridge.SetKeyfile(cfg.Nym.AccountKeysFile)

	// TODO: later remove it, but for now it's temporary for demo sake
	privateKey, loadErr := ethcrypto.LoadECDSA(cfg.Nym.AccountKeysFile)
	if loadErr != nil {
		// qb.DisplayNotificationf(errNotificationTitle, "failed to load Nym keys: %v", loadErr)
		fmt.Printf("failed to load Nym keys: %v\n", loadErr)
		configBridge.SetAddress("could not load the key")
	} else {
		address := ethcrypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey)).Hex()
		configBridge.SetAddress(address)
	}

	// should have been detected during validation...
	if len(cfg.Nym.EthereumNodeAddresses) > 0 {
		configBridge.SetEthereumNode(cfg.Nym.EthereumNodeAddresses[0])
	} else {
		configBridge.SetEthereumNode("none specified")
	}
	configBridge.SetNymERC20(cfg.Nym.NymContract.Hex())
	configBridge.SetPipeAccount(cfg.Nym.PipeAccount.Hex())

	for i, addr := range cfg.Client.IAAddresses {
		qb.NewNymValidator(fmt.Sprintf("nymnode%v", i), addr)
	}

	for i, addr := range cfg.Nym.BlockchainNodeAddresses {
		qb.NewTendermintValidator(fmt.Sprintf("tendermintnode%v", i), addr)
	}

	qb.cfg = cfg

	if privateKey == nil || loadErr != nil {
		qb.ShowNewKeyDialog()
	}
}

func (qb *QmlBridge) confirmConfig() {
	if qb.clientInstance == nil {
		client, err := client.New(qb.cfg)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not use the config to create client instance: %v", err)
			return
		}
		qb.clientInstance = client
	}

	if qb.longtermSecret == nil {
		qb.longtermSecret = qb.clientInstance.RandomBIG()
		qb.UpdateSecret(utils.ToCoconutString(qb.longtermSecret))
	}
	valueList := make([]string, len(token.AllowedValues))
	for i, val := range token.AllowedValues {
		valueList[i] = strconv.FormatInt(val, 10) + "Nym"
	}
	qb.PopulateValueComboBox(valueList)

	// gui only cares about physical addresses (for now)
	spAddresses := make([]string, len(qb.cfg.Nym.ServiceProviders))
	i := 0
	for sp, _ := range qb.cfg.Nym.ServiceProviders {
		spAddresses[i] = sp
		i++
	}

	qb.PopulateSPComboBox(spAddresses)
	qb.SetAccountStatus(qb.checkIfAccountExists())
}

func (qb *QmlBridge) forceUpdateBalances(busyIndicator *core.QObject, mainLayoutObject *core.QObject) {
	go func() {
		toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, true)
		defer toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, false)

		qb.updateBalances()
	}()
}

func (qb *QmlBridge) sendToPipeAccount(amount string, busyIndicator *core.QObject, mainLayoutObject *core.QObject) {
	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return
	}

	go func() {
		toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, true)
		defer toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, false)

		amountInt64, err := strconv.ParseInt(amount, 10, 64)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not parse the value: %v", err)
			return
		}

		currentERC20Balance, err := qb.clientInstance.GetCurrentERC20Balance()
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "failed to query for ERC20 Nym Balance: %v", err)
			return
		}

		currentNymBalance, err := qb.clientInstance.GetCurrentNymBalance()
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "failed to query for Nym Token Balance: %v", err)
			return
		}

		// TODO:
		ctx := context.TODO()
		if err := qb.clientInstance.SendToPipeAccount(ctx, amountInt64); err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "failed to send %v to the pipe account: %v", amountInt64, err)
			return
		}

		// TODO: not the best option if multiple actions were taken concurrently, in future wait until block X is commited
		qb.waitForERC20BalanceChange(ctx, currentERC20Balance-uint64(amountInt64))

		if err := qb.clientInstance.WaitForBalanceChange(ctx, currentNymBalance+uint64(amountInt64)); err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "failed to query for Nym Token Balance: %v", err)
			return
		}

		qb.UpdateNymTokenBalance(strconv.FormatUint(currentNymBalance+uint64(amountInt64), 10))
		qb.ResetWaitingForEthereumLabel()
	}()
}

func (qb *QmlBridge) redeemTokens(amount string, busyIndicator *core.QObject, mainLayoutObject *core.QObject) {
	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return
	}

	go func() {
		toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, true)
		defer toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, false)

		amountInt64, err := strconv.ParseInt(amount, 10, 64)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not parse the value: %v", err)
			return
		}

		currentERC20Balance, err := qb.clientInstance.GetCurrentERC20Balance()
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "failed to query for ERC20 Nym Balance: %v", err)
			return
		}

		currentNymBalance, err := qb.clientInstance.GetCurrentNymBalance()
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "failed to query for Nym Token Balance: %v", err)
			return
		}

		// TODO:
		ctx := context.TODO()

		if err := qb.clientInstance.RedeemTokens(ctx, uint64(amountInt64)); err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "failed to redeem %v tokens: %v", amountInt64, err)
			return
		}

		if err := qb.clientInstance.WaitForBalanceChange(ctx, currentNymBalance-uint64(amountInt64)); err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "failed to query for Nym Token Balance: %v", err)
			return
		}

		qb.UpdateNymTokenBalance(strconv.FormatUint(currentNymBalance-uint64(amountInt64), 10))
		qb.waitForERC20BalanceChange(ctx, currentERC20Balance+uint64(amountInt64))
		qb.ResetWaitingForEthereumLabel()
	}()
}

func (qb *QmlBridge) getCredential(value string, busyIndicator *core.QObject, mainLayoutObject *core.QObject) {
	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return
	}

	go func() {
		toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, true)
		defer toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, false)

		value = strings.TrimSuffix(value, "Nym")
		valueInt64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not parse the value: %v", err)
			return
		}

		seq := qb.clientInstance.RandomBIG()

		token, err := token.New(seq, qb.longtermSecret, valueInt64)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not generate token for %v: %v", value, err)
			return
		}

		cred, err := qb.clientInstance.GetCredential(token)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not obtain credential for %v: %v", value, err)
			return
		}

		qb.updateBalances()

		fmt.Printf("obtained credential: %+v\n", cred)

		credBytes, err := cred.MarshalBinary()
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not marshal obtained credential: %v", err)
			return
		}

		seqString := utils.ToCoconutString(seq)

		item := CredentialListItem{
			credential: base64.StdEncoding.EncodeToString(credBytes),
			sequence:   seqString,
			value:      uint64(valueInt64),
		}

		issuedCredential := &IssuedCredential{
			credential: cred,
			token:      token, // encapsulates all attributes
		}

		// TODO: locking? - the map is being written in goroutine so in theory we might have concurrency issues
		// but then again, if button is pressed, the other parts of the gui are locked
		// in principle each credential has unique sequence number by which it can be identified
		credentialMap[seqString] = issuedCredential

		qb.AddCredentialListItem(item)
	}()
}

func (qb *QmlBridge) spendCredential(chosenSP, seqString string, busyIndicator *core.QObject, mainLayoutObject *core.QObject) {
	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return
	}

	go func() {
		toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, true)
		defer toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, false)

		spAddressRaw, ok := qb.cfg.Nym.ServiceProviders[chosenSP]
		spAddress := ethcommon.HexToAddress(spAddressRaw)
		if !ok {
			qb.DisplayNotificationf(errNotificationTitle, "No service provider with address %v exists", chosenSP)
			return
		}

		cred, ok := credentialMap[seqString]
		if !ok {
			qb.DisplayNotificationf(errNotificationTitle, "no credential exists for that sequence number (%v)", seqString)
			return
		}

		wasSuccessful, err := qb.clientInstance.SpendCredential(cred.token, cred.credential, chosenSP, spAddress, nil)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not spend the credential: %v", err)
			return
		}

		if wasSuccessful {
			qb.DisplayNotificationf(infoNotificationTitle, "We successfully managed to spend credential with value of %v Nyms at SP (%v) with address %v!", cred.token.Value(), chosenSP, spAddress.Hex())
		} else {
			qb.DisplayNotificationf(infoNotificationTitle, "We failed to spend credential with value of %v Nyms at SP (%v) with address %v: %v", cred.token.Value(), chosenSP, spAddress.Hex(), err)
		}

		// TODO: for demo sake, mark as spent (so you could see double-spent error), but in future just remove it
		qb.MarkSpentCredential()
	}()
}

func (qb *QmlBridge) generateNewKey() {
	fmt.Println("new keygen")
	pk, err := ethcrypto.GenerateKey()
	if err != nil {
		qb.DisplayNotificationf(errNotificationTitle, "could not generate a fresh key: %v", err)
		return
	}
	if err := ethcrypto.SaveECDSA(qb.cfg.Nym.AccountKeysFile, pk); err != nil {
		qb.DisplayNotificationf(errNotificationTitle, "could not save the new key: %v", err)
	}
}

func (qb *QmlBridge) checkIfAccountExists() bool {
	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return false
	}

	exists, err := qb.clientInstance.CheckAccountExistence()
	if err != nil {
		qb.DisplayNotificationf(errNotificationTitle, "could not check for account existence")
		return false
	}
	return exists
}

func (qb *QmlBridge) registerAccount(busyIndicator *core.QObject, mainLayoutObject *core.QObject) {
	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return
	}

	go func() {
		toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, true)
		defer toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, false)

		// fake non-existent credential
		accountCred := []byte("foo")
		if err := qb.clientInstance.RegisterAccount(accountCred); err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not register Nym account: %v", err)
			return
		}

		qb.SetAccountStatus(true)
	}()
}

func (qb *QmlBridge) getFaucetNym(busyIndicator *core.QObject, mainLayoutObject *core.QObject) {
	// for now just hardcode it
	var nyms int64 = 50

	if qb.clientInstance == nil {
		qb.DisplayNotificationf(errNotificationTitle, "nil client instance")
		return
	}

	go func(nyms int64) {
		toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, true)
		defer toggleIndicatorAndObjects(busyIndicator, []*core.QObject{mainLayoutObject}, false)

		ctx := context.TODO()
		erc20Hash, etherHash, err := qb.clientInstance.MakeFaucetRequest(ctx, nyms)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not send request to the faucet: %v", err)
			return
		}

		successERC20, err := qb.clientInstance.WaitForEthereumTxToResolve(ctx, erc20Hash)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not receive ERC20 Nym: %v", err)
			return
		}

		successEther, err := qb.clientInstance.WaitForEthereumTxToResolve(ctx, etherHash)
		if err != nil {
			qb.DisplayNotificationf(errNotificationTitle, "could not receive Ether: %v", err)
			return
		}

		if successERC20 && successEther {
			qb.updateBalances()
			qb.DisplayNotificationf(infoNotificationTitle, "Received %v Nym from the faucet (+ some Ether for transaction fees) from the faucet!", nyms)
		} else {
			qb.DisplayNotificationf(warnNotificationTitle, "unknown error when trying to receive funds from the faucet")
		}
		qb.ResetWaitingForEthereumLabel()
	}(nyms)
}

func (qb *QmlBridge) randomizeCredential(seqString string) string {
	cred, ok := credentialMap[seqString]
	if !ok {
		qb.DisplayNotificationf(errNotificationTitle, "no credential exists for that sequence number (%v)", seqString)
		return ""
	}

	rcred := qb.clientInstance.ForceReRandomizeCredential(cred.credential)
	if rcred != nil {
		// it should ALWAYS be not nil, it's just a sanity check
		credentialMap[seqString].credential = rcred
	}

	rCredBytes, err := rcred.MarshalBinary()
	if err != nil {
		qb.DisplayNotificationf(errNotificationTitle, "could not marshal randomized credential: %v", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(rCredBytes)
}

// this function will be automatically called, when you use the `NewQmlBridge` function
func (qb *QmlBridge) init() {
	// TODO: perhaps create client instance here?

}
