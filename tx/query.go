package tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/pkg/errors"
	hub "github.com/sentinel-official/hub/types"

	"github.com/sentinel-official/hub/x/vpn"
	"github.com/sentinel-official/hub/x/vpn/client/common"
)

func (t Tx) QueryAccount(_address string) (auth.Account, error) {
	address, err := sdk.AccAddressFromBech32(_address)
	if err != nil {
		return nil, err
	}

	account, err := t.Manager.CLI.GetAccount(address)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.Errorf("no account found")
	}

	return account, nil
}

func (t Tx) QueryNode(id string) (*vpn.Node, error) {
	return common.QueryNode(t.Manager.CLI.CLIContext, id)
}

func (t Tx) QueryNodesOfResolver(id string) ([]hub.NodeID, error) {
	return common.QueryNodesOfResolver(t.Manager.CLI.CLIContext, id)
}

func (t Tx) QuerySubscription(id string) (*vpn.Subscription, error) {
	return common.QuerySubscription(t.Manager.CLI.CLIContext, id)
}

func (t Tx) QuerySubscriptionByTxHash(hash string) (*vpn.Subscription, error) {
	res, err := t.Manager.QueryTx(hash)
	if err != nil {
		return nil, err
	}
	if !res.TxResult.IsOK() {
		return nil, errors.Errorf(res.TxResult.String())
	}

	var stdTx auth.StdTx
	if err := t.Manager.CLI.Codec.UnmarshalBinaryLengthPrefixed(res.Tx, &stdTx); err != nil {
		return nil, err
	}

	if len(stdTx.Msgs) != 1 || stdTx.Msgs[0].Type() != "start_subscription" {
		return nil, errors.Errorf("Invalid subscription transaction")
	}

	events := sdk.StringifyEvents(res.TxResult.Events)
	id := events[1].Attributes[0].Value
	return common.QuerySubscription(t.Manager.CLI.CLIContext, id)
}

func (t Tx) QuerySessionsCountOfSubscription(id string) (uint64, error) {
	return common.QuerySessionsCountOfSubscription(t.Manager.CLI.CLIContext, id)
}

func (t Tx) QuerySessionOfSubscription(id string, index uint64) (*vpn.Session, error) {
	return common.QuerySessionOfSubscription(t.Manager.CLI.CLIContext, id, index)
}
