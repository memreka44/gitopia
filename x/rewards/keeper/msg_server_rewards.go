package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gitopia/gitopia/x/rewards/types"
)

func (k msgServer) CreateReward(goCtx context.Context, msg *types.MsgCreateReward) (*types.MsgCreateRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)

	if params.EvaluatorAddress != msg.Creator {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "user (%v) doesn't have permission to perform this operation", msg.Creator)
	}

	// Check if the value already exists
	_, isFound := k.GetReward(
		ctx,
		msg.Recipient,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var reward = types.Reward{
		Creator:   msg.Creator,
		Recipient: msg.Recipient,
		Amount:    msg.Amount,
	}

	k.SetReward(
		ctx,
		reward,
	)
	return &types.MsgCreateRewardResponse{}, nil
}
