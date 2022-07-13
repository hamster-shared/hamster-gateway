package chain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/log"
	"reflect"
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
	PeerId      []types.U8
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
	Market_Money                                  []EventMarketMoney
}

func DecodeEventRecordsWithIgnoreError(e types.EventRecordsRaw, m *types.Metadata, t interface{}) error {
	fmt.Println(fmt.Sprintf("will decode event records from raw hex: %#x", e))

	// ensure t is a pointer
	ttyp := reflect.TypeOf(t)
	if ttyp.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer, but is " + fmt.Sprint(ttyp))
	}
	// ensure t is not a nil pointer
	tval := reflect.ValueOf(t)
	if tval.IsNil() {
		return errors.New("target is a nil pointer")
	}
	val := tval.Elem()
	typ := val.Type()
	// ensure val can be set
	if !val.CanSet() {
		return fmt.Errorf("unsettable value %v", typ)
	}
	// ensure val points to a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("target must point to a struct, but is " + fmt.Sprint(typ))
	}

	decoder := scale.NewDecoder(bytes.NewReader(e))

	// determine number of events
	n, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("found %v events", n))

	// iterate over events
	for i := uint64(0); i < n.Uint64(); i++ {
		fmt.Println(fmt.Sprintf("decoding event #%v", i))

		// decode Phase
		phase := types.Phase{}
		err := decoder.Decode(&phase)
		if err != nil {
			return fmt.Errorf("unable to decode Phase for event #%v: %v", i, err)
		}

		// decode EventID
		id := types.EventID{}
		err = decoder.Decode(&id)
		if err != nil {
			return fmt.Errorf("unable to decode EventID for event #%v: %v", i, err)
		}

		fmt.Println(fmt.Sprintf("event #%v has EventID %v", i, id))

		// ask metadata for method & event name for event
		moduleName, eventName, err := m.FindEventNamesForEventID(id)
		// moduleName, eventName, err := "System", "ExtrinsicSuccess", nil
		if err != nil {
			//return fmt.Errorf("unable to find event with EventID %v in metadata for event #%v: %s", id, i, err)
			log.Warn("unable to find event with EventID %v in metadata for event #%v: %s", id, i, err)
			continue
		}

		fmt.Println(fmt.Sprintf("event #%v is in module %v with event name %v", i, moduleName, eventName))

		// check whether name for eventID exists in t
		field := val.FieldByName(fmt.Sprintf("%v_%v", moduleName, eventName))
		if !field.IsValid() {
			fmt.Println(fmt.Sprintf("unable to find field %v_%v for event #%v with EventID %v ", moduleName, eventName, i, id))
			continue
		}

		// create a pointer to with the correct type that will hold the decoded event
		holder := reflect.New(field.Type().Elem())

		// ensure first field is for Phase, last field is for Topics
		numFields := holder.Elem().NumField()
		if numFields < 2 {
			return fmt.Errorf("expected event #%v with EventID %v, field %v_%v to have at least 2 fields "+
				"(for Phase and Topics), but has %v fields", i, id, moduleName, eventName, numFields)
		}
		phaseField := holder.Elem().FieldByIndex([]int{0})
		if phaseField.Type() != reflect.TypeOf(phase) {
			return fmt.Errorf("expected the first field of event #%v with EventID %v, field %v_%v to be of type "+
				"types.Phase, but got %v", i, id, moduleName, eventName, phaseField.Type())
		}
		topicsField := holder.Elem().FieldByIndex([]int{numFields - 1})
		if topicsField.Type() != reflect.TypeOf([]types.Hash{}) {
			return fmt.Errorf("expected the last field of event #%v with EventID %v, field %v_%v to be of type "+
				"[]types.Hash for Topics, but got %v", i, id, moduleName, eventName, topicsField.Type())
		}

		// set the phase we decoded earlier
		phaseField.Set(reflect.ValueOf(phase))

		// set the remaining fields
		for j := 1; j < numFields; j++ {
			err = decoder.Decode(holder.Elem().FieldByIndex([]int{j}).Addr().Interface())
			if err != nil {
				return fmt.Errorf("unable to decode field %v event #%v with EventID %v, field %v_%v: %v", j, i, id, moduleName,
					eventName, err)
			}
		}

		// add the decoded event to the slice
		field.Set(reflect.Append(field, holder.Elem()))

		fmt.Println(fmt.Sprintf("decoded event #%v", i))
	}
	return nil
}
