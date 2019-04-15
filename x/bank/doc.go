// This document specifies the bank module of the Cosmos SDK.
//
// The bank module is responsible for handling multi-asset coin transfers between
// accounts and tracking special-case pseudo-transfers which must work differently
// with particular kinds of accounts (notably delegating/undelegating for vesting
// accounts). It exposes several interfaces with varying capabilities for secure
// interaction with other modules which must alter user balances.
//
// This module is used in the Cosmos Hub.
//
// State
//
// Presently, the bank module has no inherent state — it simply reads and writes
// accounts using the AccountKeeper from the auth module.
//
// This implementation choice is intended to minimize necessary state reads/writes,
// since we expect most transactions to involve coin amounts (for fees), so storing
// coin data in the account saves reading it separately.
//
// Keepers
//
// The bank module provides three different exported keeper interfaces which can be
// passed to other modules which need to read or update account balances. Modules
// should use the least-permissive interface which provides the functionality they require.
//
// Note that you should always review the bank module code to ensure that permissions
// are limited in the way that you expect.
//
// Common Types
//
// Input represents an input of a multiparty transfer:
//
// 	type Input struct {
// 		Address AccAddress
// 		Coins   Coins
// 	}
//
// Output represents an output of a multiparty transfer:
//
// 	type Output struct {
// 		Address AccAddress
// 		Coins   Coins
// 	}
//
// BaseKeeper provides full-permission access: the ability to arbitrary modify any
// account's balance and mint or burn coins:
//
// 	type BaseKeeper interface {
// 		SetCoins(addr AccAddress, amt Coins)
// 		SubtractCoins(addr AccAddress, amt Coins)
// 		AddCoins(addr AccAddress, amt Coins)
// 		InputOutputCoins(inputs []Input, outputs []Output)
// }
//
// setCoins fetches an account by address, sets the coins on the account,
// and saves the account:
//
// 	setCoins(addr AccAddress, amt Coins)
// 		account = accountKeeper.getAccount(addr)
// 		if account == nil
// 			fail with "no account found"
// 		account.Coins = amt
// 		accountKeeper.setAccount(account)
//
// subtractCoins fetches the coins of an account, subtracts the provided amount,
// and saves the account. This decreases the total supply:
//
// 	subtractCoins(addr AccAddress, amt Coins)
// 		oldCoins = getCoins(addr)
// 		newCoins = oldCoins - amt
// 		if newCoins < 0
// 			fail with "cannot end up with negative coins"
// 		setCoins(addr, newCoins)
//
// addCoins fetches the coins of an account, adds the provided amount, and saves
// the account. This increases the total supply:
//
// 	addCoins(addr AccAddress, amt Coins)
// 		oldCoins = getCoins(addr)
// 		newCoins = oldCoins + amt
// 		setCoins(addr, newCoins)
//
// inputOutputCoins transfers coins from any number of input accounts to any number of
// output accounts:
//
// 	inputOutputCoins(inputs []Input, outputs []Output)
// 		for input in inputs
// 			subtractCoins(input.Address, input.Coins)
// 		for output in outputs
// 			addCoins(output.Address, output.Coins)
//
// SendKeeper provides access to account balances and the ability to transfer coins
// between accounts, but not to alter the total supply (mint or burn coins):
//
// 	type SendKeeper interface {
// 		SendCoins(from AccAddress, to AccAddress, amt Coins)
// 	}
//
// sendCoins transfers coins from one account to another.
//
// 	sendCoins(from AccAddress, to AccAddress, amt Coins)
// 		subtractCoins(from, amt)
// 		addCoins(to, amt)
//
// ViewKeeper provides read-only access to account balances but no balance alteration
// functionality. All balance lookups are O(1):
//
// 	type ViewKeeper interface {
// 		GetCoins(addr AccAddress) Coins
// 		HasCoins(addr AccAddress, amt Coins) bool
// 	}
//
// getCoins returns the coins associated with an account:
//
// 	getCoins(addr AccAddress)
// 		account = accountKeeper.getAccount(addr)
//		if account == nil
// 			return Coins{}
// 		return account.Coins
//
// hasCoins returns whether or not an account has at least the provided amount
// of coins:
//
// 	hasCoins(addr AccAddress, amt Coins)
// 		account = accountKeeper.getAccount(addr)
// 		coins = getCoins(addr)
// 		return coins >= amt
//
// Messages
//
// MsgSend:
//
// 	type MsgSend struct {
// 		Inputs  []Input
// 		Outputs []Output
// 	}
//
// handleMsgSend just runs inputOutputCoins:
//
// 	handleMsgSend(msg MsgSend)
// 		inputSum = 0
// 		for input in inputs
// 			inputSum += input.Amount
// 		outputSum = 0
// 		for output in outputs
// 			outputSum += output.Amount
// 		if inputSum != outputSum:
// 			fail with "input/output amount mismatch"
// 		return inputOutputCoins(msg.Inputs, msg.Outputs)
//
// Tags
//
// MsgSend tags:
//
//	| Key         | Value                     |
//	|-------------|---------------------------|
//	| action      | send                      |
//	| category    | bank                      |
//	| sender      | {senderAccountAddress}    |
//	| recipient   | {recipientAccountAddress} |
//
// MsgMultiSend tags:
//
//	| Key         | Value                     |
//	|-------------|---------------------------|
//	| action      | multisend                 |
//	| category    | bank                      |
//	| sender      | {senderAccountAddress}    |
//	| recipient   | {recipientAccountAddress} |
//
package bank
