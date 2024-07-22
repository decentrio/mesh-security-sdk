package contract

import wasmvmtypes "github.com/CosmWasm/wasmvm/types"

type (
	SudoMsgProvider struct {
		VaultSudoMsg *StakeMsg `json:"vault,omitempty"`
	}

	StakeMsg struct {
		StakeRemote `json:"stake_remote"`
		StakeLocal  `json:"stake_local"`
	}

	StakeRemote struct {
		Contract string           `json:"contract"`
		Amount   wasmvmtypes.Coin `json:"amount"`
	}
	StakeLocal struct {
		Amount wasmvmtypes.Coin `json:"amount"`
	}
)
