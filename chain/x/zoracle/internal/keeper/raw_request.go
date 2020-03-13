package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRawDataRequest is a function to save raw data request detail to the given request id and external id.
func (k Keeper) SetRawDataRequest(
	ctx sdk.Context, requestID types.RequestID, externalID types.ExternalID, rawDataRequest types.RawDataRequest,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.RawDataRequestStoreKey(requestID, externalID),
		k.cdc.MustMarshalBinaryBare(rawDataRequest),
	)
}

// GetRawDataRequest is a function to get raw data request detail by the given request id and external id.
func (k Keeper) GetRawDataRequest(
	ctx sdk.Context, requestID types.RequestID, externalID types.ExternalID,
) (types.RawDataRequest, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckRawDataRequestExists(ctx, requestID, externalID) {
		// TODO: fix error
		return types.RawDataRequest{}, types.ErrRawDataRequestNotFound(types.DefaultCodespace)
	}

	bz := store.Get(types.RawDataRequestStoreKey(requestID, externalID))
	var requestDetail types.RawDataRequest
	k.cdc.MustUnmarshalBinaryBare(bz, &requestDetail)
	return requestDetail, nil
}

// CheckRawDataRequestExists checks if the raw request data at this request id and external id
// presents in the store or not.
func (k Keeper) CheckRawDataRequestExists(ctx sdk.Context, requestID types.RequestID, externalID types.ExternalID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.RawDataRequestStoreKey(requestID, externalID))
}

// AddNewRawDataRequest checks all conditions before saving a new raw data request to the store.
func (k Keeper) AddNewRawDataRequest(
	ctx sdk.Context, requestID types.RequestID, externalID types.ExternalID, dataSourceID types.DataSourceID, calldata []byte,
) sdk.Error {
	if len(calldata) > int(k.MaxCalldataSize(ctx)) {
		return types.ErrBadDataValue(
			"AddNewRawDataRequest: Calldata size (%d) exceeds the maximum size (%d).",
			len(calldata),
			int(k.MaxCalldataSize(ctx)),
		)
	}

	request, err := k.GetRequest(ctx, requestID)
	if err != nil {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	if !k.CheckDataSourceExists(ctx, dataSourceID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	if k.CheckRawDataRequestExists(ctx, requestID, externalID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	ctx.GasMeter().ConsumeGas(
		k.GasPerRawDataRequestPerValidator(ctx)*uint64(len(request.RequestedValidators)),
		"RawDataRequest",
	)
	k.SetRawDataRequest(ctx, requestID, externalID, types.NewRawDataRequest(dataSourceID, calldata))

	return k.ValidateDataSourceCount(ctx, requestID)
}

// GetRawDataRequestIterator is a function to get iterator on all raw data request that belong to
// given request id
func (k Keeper) GetRawDataRequestIterator(ctx sdk.Context, requestID types.RequestID) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.RawDataRequestStoreKeyPrefix, requestID)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefix)
}

// GetRawDataRequestCount returns amount of raw data requests in given request.
func (k Keeper) GetRawDataRequestCount(ctx sdk.Context, requestID types.RequestID) int64 {
	iterator := k.GetRawDataRequestIterator(ctx, requestID)
	count := 0
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return int64(count)
}

// GetRawDataRequests returns a list of raw data requests in given request.
func (k Keeper) GetRawDataRequests(ctx sdk.Context, requestID types.RequestID) []types.RawDataRequest {
	iterator := k.GetRawDataRequestIterator(ctx, requestID)
	rawRequests := make([]types.RawDataRequest, 0)
	for ; iterator.Valid(); iterator.Next() {
		var rawRequest types.RawDataRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rawRequest)
		rawRequests = append(rawRequests, rawRequest)
	}
	return rawRequests
}

// GetRawDataRequestWithExternalIDs returns a list of raw data requests with external id in given request.
func (k Keeper) GetRawDataRequestWithExternalIDs(
	ctx sdk.Context, requestID types.RequestID,
) []types.RawDataRequestWithExternalID {
	iterator := k.GetRawDataRequestIterator(ctx, requestID)
	rawRequests := make([]types.RawDataRequestWithExternalID, 0)
	for ; iterator.Valid(); iterator.Next() {
		var rawRequest types.RawDataRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rawRequest)
		rawRequests = append(rawRequests,
			types.NewRawDataRequestWithExternalID(
				types.GetExternalIDFromRawDataRequestKey(iterator.Key()),
				rawRequest,
			))
	}
	return rawRequests
}
