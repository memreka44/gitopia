package keeper

import (
	"encoding/binary"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gitopia/gitopia/x/gitopia/types"
)

// GetCommentCount get the total number of comment
func (k Keeper) GetCommentCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentCountKey))
	byteKey := types.KeyPrefix(types.CommentCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetCommentCount set the total number of comment
func (k Keeper) SetCommentCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentCountKey))
	byteKey := types.KeyPrefix(types.CommentCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendComment appends a comment in the store with a new id and update the count
func (k Keeper) AppendComment(
	ctx sdk.Context,
	creator string,
	parentId uint64,
	commentIid uint64,
	body string,
	attachments []string,
	diffHunk string,
	path string,
	system bool,
	authorId uint64,
	authorAssociation string,
	createdAt int64,
	updatedAt int64,
	commentType string,
	extensions string,
) uint64 {
	// Create the comment
	count := k.GetCommentCount(ctx)
	var comment = types.Comment{
		Creator:           creator,
		Id:                count,
		ParentId:          parentId,
		CommentIid:        commentIid,
		Body:              body,
		Attachments:       attachments,
		DiffHunk:          diffHunk,
		Path:              path,
		System:            system,
		AuthorId:          authorId,
		AuthorAssociation: authorAssociation,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		CommentType:       commentType,
		Extensions:        extensions,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	value := k.cdc.MustMarshalBinaryBare(&comment)
	store.Set(GetCommentIDBytes(comment.Id), value)

	// Update comment count
	k.SetCommentCount(ctx, count+1)

	return count
}

// SetComment set a specific comment in the store
func (k Keeper) SetComment(ctx sdk.Context, comment types.Comment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	b := k.cdc.MustMarshalBinaryBare(&comment)
	store.Set(GetCommentIDBytes(comment.Id), b)
}

// GetComment returns a comment from its id
func (k Keeper) GetComment(ctx sdk.Context, id uint64) types.Comment {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	var comment types.Comment
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetCommentIDBytes(id)), &comment)
	return comment
}

// HasComment checks if the comment exists in the store
func (k Keeper) HasComment(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	return store.Has(GetCommentIDBytes(id))
}

// GetCommentOwner returns the creator of the comment
func (k Keeper) GetCommentOwner(ctx sdk.Context, id uint64) string {
	return k.GetComment(ctx, id).Creator
}

// RemoveComment removes a comment from the store
func (k Keeper) RemoveComment(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	store.Delete(GetCommentIDBytes(id))
}

// GetAllComment returns all comment
func (k Keeper) GetAllComment(ctx sdk.Context) (list []types.Comment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Comment
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCommentIDBytes returns the byte representation of the ID
func GetCommentIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetCommentIDFromBytes returns ID in uint64 format from a byte array
func GetCommentIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
