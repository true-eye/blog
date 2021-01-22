package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/true-eye/blog/x/blog/types"
	"github.com/true-eye/blog/x/blog/keeper"
)

func handleMsgCreateComment(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateComment) (*sdk.Result, error) {
	k.CreateComment(ctx, msg)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
