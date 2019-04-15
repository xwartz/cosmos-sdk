// This document specifies the crisis module of the Cosmos SDK.
//
// State
//
// Due to the anticipated large gas cost requirement to verify an invariant (and
// potential to exceed the maximum allowable block gas limit) a constant fee is
// used instead of the standard gas consumption method. The constant fee is
// intended to be larger than the anticipated gas cost of running the invariant
// with the standard gas consumption method.
//
// The ConstantFee param is held in the global params store.
//
//  - Params: mint/params -> amino(sdk.Coin)
//
// Messages
//
// Blockchain invariants can be checked using the `MsgVerifyInvariant` message:
// 	type MsgVerifyInvariant struct {
// 		Sender         sdk.AccAddress
// 		InvariantRoute string
// 	}
//
// This message is expected to fail if the sender does not have enough coins for
// the constant fee or the invariant route is not registered.
//
// This message checks the invariant provided, and if the invariant is broken it
// panics, halting the blockchain. If the invariant is broken, the constant fee is
// never deducted as the transaction is never committed to a block (equivalent to
// being refunded). However, if the invariant is not broken, the constant fee will
// not be refunded.
package crisis
