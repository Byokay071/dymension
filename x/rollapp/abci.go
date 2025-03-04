package rollapp

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension/x/rollapp/keeper"
	"github.com/dymensionxyz/dymension/x/rollapp/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called on every block.
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
}

// Called every block to finalize states that their dispute period over.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// check to see if there are pending  states to be finalized
	blockHeightToFinalizationQueue, found := k.GetBlockHeightToFinalizationQueue(ctx, uint64(ctx.BlockHeight()))

	if found {
		// finalize pending states
		for _, stateInfoIndex := range blockHeightToFinalizationQueue.FinalizationQueue {
			stateInfo, found := k.GetStateInfo(ctx, stateInfoIndex.RollappId, stateInfoIndex.Index)
			if !found {
				panic(fmt.Errorf("invariant broken: rollapp %s should have state to be finalizaed in height %d. update state index is %d \n",
					stateInfoIndex.RollappId, ctx.BlockHeight(), stateInfoIndex.Index))
			}
			stateInfo.Status = types.STATE_STATUS_FINALIZED
			// write the new status
			k.SetStateInfo(ctx, stateInfo)
		}
	}
}
