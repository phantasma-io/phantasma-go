package event

import (
	"encoding/hex"
	"math/big"

	crypto "github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/io"
)

type EventKind uint
type TypeAuction uint

const (
	Unknown            EventKind = 0
	ChainCreate        EventKind = 1
	TokenCreate        EventKind = 2
	TokenSend          EventKind = 3
	TokenReceive       EventKind = 4
	TokenMint          EventKind = 5
	TokenBurn          EventKind = 6
	TokenStake         EventKind = 7
	TokenClaim         EventKind = 8
	AddressRegister    EventKind = 9
	AddressLink        EventKind = 10
	AddressUnlink      EventKind = 11
	OrganizationCreate EventKind = 12
	OrganizationAdd    EventKind = 13
	OrganizationRemove EventKind = 14
	GasEscrow          EventKind = 15
	GasPayment         EventKind = 16
	AddressUnregister  EventKind = 17
	OrderCreated       EventKind = 18
	OrderCancelled     EventKind = 19
	OrderFilled        EventKind = 20
	OrderClosed        EventKind = 21
	FeedCreate         EventKind = 22
	FeedUpdate         EventKind = 23
	FileCreate         EventKind = 24
	FileDelete         EventKind = 25
	ValidatorPropose   EventKind = 26
	ValidatorElect     EventKind = 27
	ValidatorRemove    EventKind = 28
	ValidatorSwitch    EventKind = 29
	PackedNFT          EventKind = 30
	ValueCreate        EventKind = 31
	ValueUpdate        EventKind = 32
	PollCreated        EventKind = 33
	PollClosed         EventKind = 34
	PollVote           EventKind = 35
	ChannelCreate      EventKind = 36
	ChannelRefill      EventKind = 37
	ChannelSettle      EventKind = 38
	LeaderboardCreate  EventKind = 39
	LeaderboardInsert  EventKind = 40
	LeaderboardReset   EventKind = 41
	PlatformCreate     EventKind = 42
	ChainSwap          EventKind = 43
	ContractRegister   EventKind = 44
	ContractDeploy     EventKind = 45
	AddressMigration   EventKind = 46
	ContractUpgrade    EventKind = 47
	Log                EventKind = 48
	Inflation          EventKind = 49
	OwnerAdded         EventKind = 50
	OwnerRemoved       EventKind = 51
	DomainCreate       EventKind = 52
	DomainDelete       EventKind = 53
	TaskStart          EventKind = 54
	TaskStop           EventKind = 55
	CrownRewards       EventKind = 56
	Infusion           EventKind = 57
	Crowdsale          EventKind = 58
	OrderBid           EventKind = 59
	Custom             EventKind = 64
)

var eventLookup = map[EventKind]string{
	Unknown:            `Unknown`,
	ChainCreate:        `ChainCreate`,
	TokenCreate:        `TokenCreate`,
	TokenSend:          `TokenSend`,
	TokenReceive:       `TokenReceive`,
	TokenMint:          `TokenMint`,
	TokenBurn:          `TokenBurn`,
	TokenStake:         `TokenStake`,
	TokenClaim:         `TokenClaim`,
	AddressRegister:    `AddressRegister`,
	AddressLink:        `AddressLink`,
	AddressUnlink:      `AddressUnlink`,
	OrganizationCreate: `OrganizationCreate`,
	OrganizationAdd:    `OrganizationAdd`,
	OrganizationRemove: `OrganizationRemove`,
	GasEscrow:          `GasEscrow`,
	GasPayment:         `GasPayment`,
	AddressUnregister:  `AddressUnregister`,
	OrderCreated:       `OrderCreated`,
	OrderCancelled:     `OrderCancelled`,
	OrderFilled:        `OrderFilled`,
	OrderClosed:        `OrderClosed`,
	FeedCreate:         `FeedCreate`,
	FeedUpdate:         `FeedUpdate`,
	FileCreate:         `FileCreate`,
	FileDelete:         `FileDelete`,
	ValidatorPropose:   `ValidatorPropose`,
	ValidatorElect:     `ValidatorElect`,
	ValidatorRemove:    `ValidatorRemove`,
	ValidatorSwitch:    `ValidatorSwitch`,
	PackedNFT:          `PackedNFT`,
	ValueCreate:        `ValueCreate`,
	ValueUpdate:        `ValueUpdate`,
	PollCreated:        `PollCreated`,
	PollClosed:         `PollClosed`,
	PollVote:           `PollVote`,
	ChannelCreate:      `ChannelCreate`,
	ChannelRefill:      `ChannelRefill`,
	ChannelSettle:      `ChannelSettle`,
	LeaderboardCreate:  `LeaderboardCreate`,
	LeaderboardInsert:  `LeaderboardInsert`,
	LeaderboardReset:   `LeaderboardReset`,
	PlatformCreate:     `PlatformCreate`,
	ChainSwap:          `ChainSwap`,
	ContractRegister:   `ContractRegister`,
	ContractDeploy:     `ContractDeploy`,
	AddressMigration:   `AddressMigration`,
	ContractUpgrade:    `ContractUpgrade`,
	Log:                `Log`,
	Inflation:          `Inflation`,
	OwnerAdded:         `OwnerAdded`,
	OwnerRemoved:       `OwnerRemoved`,
	DomainCreate:       `DomainCreate`,
	DomainDelete:       `DomainDelete`,
	TaskStart:          `TaskStart`,
	TaskStop:           `TaskStop`,
	CrownRewards:       `CrownRewards`,
	Infusion:           `Infusion`,
	Crowdsale:          `Crowdsale`,
	OrderBid:           `OrderBid`,
	Custom:             `Custom`,
}

func (k EventKind) String() string {
	return eventLookup[k]
}

func (k *EventKind) SetString(eventKind string) {
	for k1, s := range eventLookup {
		if s == eventKind {
			*k = k1
			return
		}
	}
}

func (k EventKind) IsTokenEvent() bool {
	return k == TokenBurn || k == TokenClaim || k == TokenMint || k == TokenReceive || k == TokenSend || k == TokenStake
}

func (k EventKind) IsMarketEvent() bool {
	return k == OrderCreated || k == OrderCancelled || k == OrderFilled || k == OrderClosed || k == OrderBid
}

const (
	Fixed   TypeAuction = 0
	Classic TypeAuction = 1
	Reserve TypeAuction = 2
	Dutch   TypeAuction = 3
)

type OrganizationEventData struct {
	Organization  string
	MemberAddress crypto.Address
}

type TokenEventData struct {
	Symbol    string
	Value     *big.Int
	ChainName string
}

// Serialize implements ther Serializable interface
func (te *TokenEventData) Serialize(writer *io.BinWriter) {
	writer.WriteString(te.Symbol)
	writer.WriteBigInteger(te.Value)
	writer.WriteString(te.ChainName)
}

// Deserialize implements ther Serializable interface
func (te *TokenEventData) Deserialize(reader *io.BinReader) {
	te.Symbol = reader.ReadString()
	te.Value = reader.ReadBigInteger()
	te.ChainName = reader.ReadString()
}

type InfusionEventData struct {
	BaseSymbol    string
	TokenID       big.Int
	InfusedSymbol string
	InfusedValue  big.Int
	ChainName     string
}
type MarketEventData struct {
	BaseSymbol  string
	QuoteSymbol string
	ID          big.Int
	Price       big.Int
	EndPrice    big.Int
	Type        TypeAuction
}

// Serialize implements ther Serializable interface
func (d *MarketEventData) Serialize(writer *io.BinWriter) {
	writer.WriteString(d.BaseSymbol)
	writer.WriteString(d.QuoteSymbol)
	writer.WriteBigInteger(&d.ID)
	writer.WriteBigInteger(&d.Price)
	writer.WriteBigInteger(&d.EndPrice)
	writer.WriteB(byte(d.Type))
}

// Deserialize implements ther Serializable interface
func (d *MarketEventData) Deserialize(reader *io.BinReader) {
	d.BaseSymbol = reader.ReadString()
	d.QuoteSymbol = reader.ReadString()
	d.ID = *reader.ReadBigInteger()
	d.Price = *reader.ReadBigInteger()
	d.EndPrice = *reader.ReadBigInteger()
	d.Type = TypeAuction(reader.ReadU32LE())
}

type ChainValueEventData struct {
	Name  string
	Value big.Int
}

type TransactionSettleEventData struct {
	Hash     crypto.Hash
	Platform string
	Chain    string
}

type GasEventData struct {
	Address crypto.Address
	price   big.Int
	amount  big.Int
}

type Event struct {
	Kind     EventKind
	Address  crypto.Address
	Contract string
	Data     []byte
}

func (e Event) String() string {
	return eventLookup[e.Kind] + "/" + e.Contract + "@" + e.Address.String() + ":" + hex.EncodeToString(e.Data)
}

// Serialize implements ther Serializable interface
func (e *Event) Serialize(writer *io.BinWriter) {
	writer.WriteB(byte(e.Kind))
	e.Address.Serialize(writer)
	writer.WriteString(e.Contract)
	writer.WriteVarBytes(e.Data)
}

// Deserialize implements ther Serializable interface
func (e *Event) Deserialize(reader *io.BinReader) {
	e.Kind = EventKind(reader.ReadU32LE())
	e.Address.Deserialize(reader)
	e.Contract = reader.ReadString()
	e.Data = reader.ReadVarBytes()
}
