package chain

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type EventProviderRegisterResourceSuccess struct {
	Phase            types.Phase
	AccountId        types.AccountID
	Index            types.U64
	PeerId           string
	Cpu              types.U64
	Memory           types.U64
	System           string
	CpuModel         string
	PriceHour        types.U128
	RentDurationHour types.U32
	Topics           []types.Hash
}

type EventResourceOrderCreateOrderSuccess struct {
	Phase         types.Phase
	AccountId     types.AccountID
	OrderIndex    types.U64
	ResourceIndex types.U64
	Duration      types.U32
	PublicKey     string
	Topics        []types.Hash
}

type EventResourceOrderOrderExecSuccess struct {
	Phase          types.Phase
	AccountId      types.AccountID
	OrderIndex     types.U64
	ResourceIndex  types.U64
	AgreementIndex types.U64
	Topics         []types.Hash
}

type EventResourceOrderReNewOrderSuccess struct {
	Phase          types.Phase
	AccountId      types.AccountID
	OrderIndex     types.U64
	ResourceIndex  types.U64
	AgreementIndex types.U64
	Topics         []types.Hash
}

type EventResourceOrderWithdrawLockedOrderPriceSuccess struct {
	Phase      types.Phase
	AccountId  types.AccountID
	OrderIndex types.U64
	OrderPrice types.U128
	Topics     []types.Hash
}

type EventResourceOrderFreeResourceProcessed struct {
	Phase      types.Phase
	OrderIndex types.U64
	PeerId     string
	Topics     []types.Hash
}

type EventResourceOrderFreeResourceApplied struct {
	Phase      types.Phase
	AccountId  types.AccountID
	OrderIndex types.U64
	Cpu        types.U64
	Memory     types.U64
	Duration   types.U32
	DeployType types.U32
	PublicKey  string
	Topics     []types.Hash
}

type EventMarketMoney struct {
	Phase  types.Phase
	Money  types.U128
	Topics []types.Hash
}

type EventElectionProviderMultiPhase_SignedPhaseStarted struct {
	Phase       types.Phase
	SignedPhase types.U32
	Topics      []types.Hash
}

type EventRegisterGatewayNodeSuccess struct {
	Phase       types.Phase
	AccountId   types.AccountID
	BlockNumber types.BlockNumber
	Topics      []types.Hash
}

type MyEventRecords struct {
	types.EventRecords
	Provider_RegisterResourceSuccess              []EventProviderRegisterResourceSuccess //nolint:stylecheck,golint
	ResourceOrder_CreateOrderSuccess              []EventResourceOrderCreateOrderSuccess //nolint:stylecheck,golint
	ResourceOrder_OrderExecSuccess                []EventResourceOrderOrderExecSuccess
	ResourceOrder_ReNewOrderSuccess               []EventResourceOrderReNewOrderSuccess
	ResourceOrder_WithdrawLockedOrderPriceSuccess []EventResourceOrderWithdrawLockedOrderPriceSuccess
	Gateway_RegisterGatewayNodeSuccess            []EventRegisterGatewayNodeSuccess
	ResourceOrder_FreeResourceProcessed           []EventResourceOrderFreeResourceProcessed
	ResourceOrder_FreeResourceApplied             []EventResourceOrderFreeResourceApplied
	ElectionProviderMultiPhase_SignedPhaseStarted []EventElectionProviderMultiPhase_SignedPhaseStarted
	Market_StakingSuccess                         []EventMarket_StakingSuccess
}

type EventMarket_StakingSuccess struct {
	Phase            types.Phase
	AccountId        types.AccountID
	MarketUserStatus types.U8
	BalanceOf        types.U128
	Topics           []types.Hash
}

type EventMarket_YES struct {
	Phase  types.Phase
	Yes    types.U8
	Topics []types.Hash
}

type EventMarket_Money struct {
	Phase     types.Phase
	BalanceOf types.U128
	Topics    []types.Hash
}

type MarketUserStatus struct {
	IsProvider bool
	IsGateway  bool
	IsClient   bool
}

func (m *MarketUserStatus) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	fmt.Println(b)

	if err != nil {
		return err
	}

	if b == 0 {
		m.IsProvider = true
	} else if b == 1 {
		m.IsGateway = true
	} else if b == 2 {
		m.IsClient = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *MarketUserStatus) Encode(encoder scale.Encoder) error {
	var err1 error
	if m.IsProvider {
		err1 = encoder.PushByte(0)
	} else if m.IsGateway {
		err1 = encoder.PushByte(1)
	} else if m.IsClient {
		err1 = encoder.PushByte(2)
	}
	if err1 != nil {
		return err1
	}
	return nil
}
