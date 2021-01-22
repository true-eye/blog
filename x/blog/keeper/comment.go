package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"

	"github.com/true-eye/blog/x/blog/types"
    "github.com/cosmos/cosmos-sdk/codec"
)

// GetCommentCount get the total number of comment
func (k Keeper) GetCommentCount(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	byteKey := []byte(types.CommentCountPrefix)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseInt(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to int64
		panic("cannot decode count")
	}

	return count
}

// SetCommentCount set the total number of comment
func (k Keeper) SetCommentCount(ctx sdk.Context, count int64)  {
	store := ctx.KVStore(k.storeKey)
	byteKey := []byte(types.CommentCountPrefix)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

// CreateComment creates a comment
func (k Keeper) CreateComment(ctx sdk.Context, msg types.MsgCreateComment) {
	// Create the comment
	count := k.GetCommentCount(ctx)
    var comment = types.Comment{
        Creator: msg.Creator,
        ID:      strconv.FormatInt(count, 10),
        Body: msg.Body,
        PostID: msg.PostID,
    }

	store := ctx.KVStore(k.storeKey)
	key := []byte(types.CommentPrefix + comment.ID)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(comment)
	store.Set(key, value)

	// Update comment count
    k.SetCommentCount(ctx, count+1)
}

// GetComment returns the comment information
func (k Keeper) GetComment(ctx sdk.Context, key string) (types.Comment, error) {
	store := ctx.KVStore(k.storeKey)
	var comment types.Comment
	byteKey := []byte(types.CommentPrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &comment)
	if err != nil {
		return comment, err
	}
	return comment, nil
}

// SetComment sets a comment
func (k Keeper) SetComment(ctx sdk.Context, comment types.Comment) {
	commentKey := comment.ID
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(comment)
	key := []byte(types.CommentPrefix + commentKey)
	store.Set(key, bz)
}

// DeleteComment deletes a comment
func (k Keeper) DeleteComment(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(types.CommentPrefix + key))
}

//
// Functions used by querier
//

func listComment(ctx sdk.Context, k Keeper) ([]byte, error) {
	var commentList []types.Comment
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CommentPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var comment types.Comment
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &comment)
		commentList = append(commentList, comment)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, commentList)
	return res, nil
}

func getComment(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	comment, err := k.GetComment(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, comment)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// Get creator of the item
func (k Keeper) GetCommentOwner(ctx sdk.Context, key string) sdk.AccAddress {
	comment, err := k.GetComment(ctx, key)
	if err != nil {
		return nil
	}
	return comment.Creator
}


// Check if the key exists in the store
func (k Keeper) CommentExists(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.CommentPrefix + key))
}
