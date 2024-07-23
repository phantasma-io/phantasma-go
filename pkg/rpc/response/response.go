package response

import (
	"encoding/hex"
	"math/big"
	"slices"
	"strings"

	chain "github.com/phantasma-io/phantasma-go/pkg/blockchain"
	"github.com/phantasma-io/phantasma-go/pkg/io"
	"github.com/phantasma-io/phantasma-go/pkg/util"
	"github.com/phantasma-io/phantasma-go/pkg/vm"
)

// ErrorResult comment
type ErrorResult struct {
	Error string `json:"error"`
}

// SingleResult comment
type SingleResult struct {
	Value interface{} `json:"error"`
}

// BalanceResult a
type BalanceResult struct {
	Chain    string   `json:"chain"`
	Amount   string   `json:"amount"`
	Symbol   string   `json:"symbol"`
	Decimals uint     `json:"decimals"`
	Ids      []string `json:"ids"`
}

func (b BalanceResult) ConvertDecimals() string {
	return util.ConvertDecimalsEx(b.Amount, int(b.Decimals), ".")
}

func (b BalanceResult) ConvertDecimalsToFloat() *big.Float {
	f, _ := big.NewFloat(0).SetString(b.ConvertDecimals())
	return f
}

// InteropResult a
type InteropResult struct {
	Logal    string `json:"local"`
	External string `json:"external"`
}

// PlatformResult a
type PlatformResult struct {
	Platform string          `json:"platform"`
	Chain    string          `json:"chain"`
	Fuel     string          `json:"fuel"`
	Tokens   []string        `json:"tokens"`
	Interop  []InteropResult `json:"interop"`
}

// GovernanceResult a
type GovernanceResult struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// OrganizationResult a
type OrganizationResult struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

// CrowdsaleResult comment
type CrowdsaleResult struct {
	Hash          string `json:"hash"`
	Name          string `json:"name"`
	Creator       string `json:"creator"`
	Flags         string `json:"flags"`
	StartDate     uint   `json:"startDate"`
	EndDate       uint   `json:"endDate"`
	SellSymbol    string `json:"sellSymbol"`
	ReceiveSymbol string `json:"receiveSymbol"`
	Price         uint   `json:"price"`
	GlobalSoftCap string `json:"globalSoftCap"`
	GlobalHardCap string `json:"globalHardCap"`
	UserSoftCap   string `json:"userSoftCap"`
	UserHardCap   string `json:"userHardCap"`
}

// NexusResult comment
type NexusResult struct {
	Name          string             `json:"name"`
	Protocol      uint               `json:"protocol"`
	Platforms     []PlatformResult   `json:"platforms"`
	Tokens        []TokenResult      `json:"tokens"`
	Chains        []ChainResult      `json:"chains"`
	Governance    []GovernanceResult `json:"governance"`
	Organizations []string           `json:"organizations"`
}

// StakeResult comment
type StakeResult struct {
	Amount    string `json:"amount"`
	Time      uint   `json:"time"`
	Unclaimed string `json:"unclaimed"`
}

func (s StakeResult) ConvertDecimals() string {
	return util.ConvertDecimalsEx(s.Amount, 8, ".") // Phantasma Stake token (SOUL) has 8 decimals
}

func (s StakeResult) ConvertDecimalsToFloat() *big.Float {
	f, _ := big.NewFloat(0).SetString(s.ConvertDecimals())
	return f
}

// StorageResult comment
type StorageResult struct {
	Available uint            `json:"available"`
	Used      uint            `json:"used"`
	Avatar    string          `json:"avatar"`
	Archives  []ArchiveResult `json:"archives"`
}

// AccountResult comment
type AccountResult struct {
	Address   string          `json:"address"`
	Name      string          `json:"name"`
	Stakes    StakeResult     `json:"stakes"`
	Stake     string          `json:"stake"`
	Unclaimed string          `json:"unclaimed"`
	Relay     string          `json:"relay"`
	Validator string          `json:"validator"`
	Storage   StorageResult   `json:"storage"`
	Balances  []BalanceResult `json:"balances"`
	Txs       []string        `json:"txs"` // Deprecated, returned as an empty array by default. Use GetAddressTransactions() to get transactions for address
}

type AddressTransactionsResult struct {
	Address string              `json:"address"`
	Txs     []TransactionResult `json:"txs"`
}

// LeaderboardRowResult comment
type LeaderboardRowResult struct {
	Address string `json:"address"`
	Value   string `json:"value"`
}

// LeaderboardResult comment
type LeaderboardResult struct {
	Name string                 `json:"name"`
	Rows []LeaderboardRowResult `json:"rows"`
}

// DappResult comment
type DappResult struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Chain   string `json:"chain"`
}

// ChainResult comment
type ChainResult struct {
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Parent       string   `json:"parent"`
	Height       uint     `json:"height"`
	Organization string   `json:"organization"`
	Contracts    []string `json:"contracts"`
	Dapps        []string `json:"dapps"`
}

// EventResult comment
type EventResult struct {
	Address  string `json:"address"`
	Contract string `json:"contract"`
	Kind     string `json:"kind"`
	Data     string `json:"data"`
}

// OracleResult comment
type OracleResult struct {
	URL     string `json:"url"`
	Content string `json:"content"`
}

// SignatureResult comment
type SignatureResult struct {
	Kind string `json:"Kind"`
	Data string `json:"Data"`
}

// TransactionResult comment
type TransactionResult struct {
	Hash         string            `json:"hash"`
	ChainAddress string            `json:"chainAddress"`
	Timestamp    uint              `json:"timestamp"`
	BlockHeight  int               `json:"blockHeight"`
	BlockHash    string            `json:"blockHash"`
	Script       string            `json:"script"`
	Payload      string            `json:"payload"`
	Events       []EventResult     `json:"events"`
	State        string            `json:"state"`
	Result       string            `json:"result"`
	Fee          string            `json:"fee"`
	Signatures   []SignatureResult `json:"signatures"`
	Expiration   uint              `json:"expiration"`
}

func (t TransactionResult) StateIsSuccess() bool {
	return chain.TxStateIsSuccess(t.State)
}

func (t TransactionResult) StateIsFault() bool {
	return chain.TxStateIsFault(t.State)
}

// AccountTransactionsResult comment
type AccountTransactionsResult struct {
	Address string              `json:"address"`
	Txs     []TransactionResult `json:"txs"`
}

// PaginatedResult comment
type PaginatedResult[T any] struct {
	Page       uint `json:"page"`
	PageSize   uint `json:"pageSize"`
	Total      uint `json:"total"`
	TotalPages uint `json:"totalPages"`

	Result T `json:"result"`
}

// BlockResult comment
type BlockResult struct {
	Hash             string              `json:"hash"`
	PreviousHash     string              `json:"previousHash"`
	Timestamp        uint                `json:"timestamp"`
	Height           uint                `json:"height"`
	ChainAddress     string              `json:"chainAddress"`
	Protocol         uint                `json:"protocol"`
	Txs              []TransactionResult `json:"txs"`
	ValidatorAddress string              `json:"validatorAddress"`
	Reward           string              `json:"reward"`
	Events           []EventResult       `json:"events"`
	Oracles          []OracleResult      `json:"oracles"`
}

// TokenExternalResult comment
type TokenExternalResult struct {
	Platform string `json:"platform"`
	Hash     string `json:"hash"`
}

// TokenPriceResult comment
type TokenPriceResult struct {
	Timestamp uint   `json:"Timestamp"`
	Open      string `json:"Open"`
	High      string `json:"High"`
	Low       string `json:"Low"`
	Close     string `json:"Close"`
}

// TokenResult comment
type TokenResult struct {
	Symbol        string                `json:"symbol"`
	Name          string                `json:"name"`
	Decimals      int                   `json:"decimals"`
	CurrentSupply string                `json:"currentSupply"`
	MaxSupply     string                `json:"maxSupply"`
	BurnedSupply  string                `json:"burnedSupply"`
	Address       string                `json:"address"`
	Owner         string                `json:"owner"`
	Flags         string                `json:"flags"`
	Script        string                `json:"script"`
	Series        []TokenSeriesResult   `json:"series"`
	External      []TokenExternalResult `json:"external"`
	Price         []TokenPriceResult    `json:"price"`
}

func (t TokenResult) IsBurnable() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Burnable")
}

func (t TokenResult) IsDivisible() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Divisible")
}

func (t TokenResult) IsFiat() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Fiat")
}

func (t TokenResult) IsFinite() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Finite")
}

func (t TokenResult) IsFuel() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Fuel")
}

func (t TokenResult) IsFungible() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Fungible")
}

func (t TokenResult) IsMintable() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Mintable")
}

func (t TokenResult) IsStakable() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Stakable")
}

func (t TokenResult) IsTransferable() bool {
	return slices.Contains(strings.Split(t.Flags, ", "), "Transferable")
}

// TokenSeriesResult comment
type TokenSeriesResult struct {
	SeriesID      uint              `json:"seriesID"`
	CurrentSupply string            `json:"currentSupply"`
	MaxSupply     string            `json:"maxSupply"`
	BurnedSupply  string            `json:"burnedSupply"`
	Mode          string            `json:"mode"`
	Script        string            `json:"script"`
	Methods       []ABIMethodResult `json:"methods"`
}

// TokenPropertyResult comment
type TokenPropertyResult struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// TokenDataResult comment
type TokenDataResult struct {
	ID             string                `json:"ID"`
	Series         string                `json:"series"`
	Mint           string                `json:"mint"`
	ChainName      string                `json:"chainName"`
	OwnerAddress   string                `json:"ownerAddress"`
	CreatorAddress string                `json:"creatorAddress"`
	RAM            string                `json:"ram"`
	ROM            string                `json:"rom"`
	Status         string                `json:"status"`
	Infusion       []TokenPropertyResult `json:"infusion"`
	Properties     []TokenPropertyResult `json:"properties"`
}

// SendRawTxResult comment
type SendRawTxResult struct {
	Hash  string `json:"hash"`
	Error string `json:"error"`
}

// AuctionResult comment
type AuctionResult struct {
	CreatorAddress  string `json:"creatorAddress"`
	ChainAddress    string `json:"chainAddress"`
	StartDate       uint   `json:"startDate"`
	EndDate         uint   `json:"endDate"`
	BaseSymbol      string `json:"baseSymbol"`
	QuoteSymbol     string `json:"quoteSymbol"`
	TokenID         string `json:"tokenId"`
	Price           string `json:"price"`
	EndPrice        string `json:"endPrice"`
	ExtensionPeriod string `json:"extensionPeriod"`
	Type            string `json:"type"`
	ROM             string `json:"rom"`
	RAM             string `json:"ram"`
	ListingFee      string `json:"listingFee"`
	CurrentWinner   string `json:"currentWinner"`
}

// ScriptResult comment
type ScriptResult struct {
	Events  []EventResult  `json:"events"`
	Result  string         `json:"result"`
	Results []string       `json:"results"`
	Oracles []OracleResult `json:"oracles"`
}

func (s ScriptResult) DecodeResult() *vm.VMObject {
	decoded, _ := hex.DecodeString(s.Result)
	br := io.NewBinReaderFromBuf(decoded)

	var vmObject vm.VMObject
	vmObject.Deserialize(br)
	return &vmObject
}

// ArchiveResult comment
type ArchiveResult struct {
	Name          string   `json:"name"`
	Hash          string   `json:"hash"`
	Time          uint     `json:"time"`
	Size          uint     `json:"size"`
	Encryption    string   `json:"encryption"`
	BlockCount    int      `json:"blockCount"`
	MissingBlocks []int    `json:"missingBlocks"`
	Owners        []string `json:"owners"`
}

// ABIParameterResult comment
type ABIParameterResult struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ABIMethodResult comment
type ABIMethodResult struct {
	Name       string               `json:"name"`
	ReturnType string               `json:"returnType"`
	Parameters []ABIParameterResult `json:"parameters"`
}

// ABIEventResult comment
type ABIEventResult struct {
	Value       int    `json:"value"`
	Name        string `json:"name"`
	ReturnType  string `json:"returnType"`
	Description string `json:"description"`
}

// ContractResult comment
type ContractResult struct {
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Script  string            `json:"script"`
	Methods []ABIMethodResult `json:"methods"`
	Events  []ABIEventResult  `json:"events"`
}

// ChannelResult comment
type ChannelResult struct {
	CreatorAddress string `json:"creatorAddress"`
	TargetAddress  string `json:"targetAddress"`
	Name           string `json:"name"`
	Chain          string `json:"chain"`
	CreationTime   uint   `json:"creationTime"`
	Symbol         string `json:"symbol"`
	Fee            string `json:"fee"`
	Balance        string `json:"balance"`
	Active         bool   `json:"active"`
	Index          int    `json:"index"`
}

// ReceiptResult comment
type ReceiptResult struct {
	Nexus     string `json:"nexus"`
	Channel   string `json:"channel"`
	Index     string `json:"index"`
	Timestamp uint   `json:"timestamp"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Script    string `json:"script"`
}

// PeerResult comment
type PeerResult struct {
	URL     string `json:"url"`
	Version string `json:"version"`
	Flags   string `json:"flags"`
	Fee     string `json:"fee"`
	Pow     uint   `json:"pow"`
}

// ValidatorResult comment
type ValidatorResult struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

// SwapResult comment
type SwapResult struct {
	SourcePlatform      string `json:"sourcePlatform"`
	SourceChain         string `json:"sourceChain"`
	SourceHash          string `json:"sourceHash"`
	SourceAddress       string `json:"sourceAddress"`
	DestinationPlatform string `json:"destinationPlatform"`
	DestinationChain    string `json:"destinationChain"`
	DestinationHash     string `json:"destinationHash"`
	DestinationAddress  string `json:"destinationAddress"`
	Symbol              string `json:"symbol"`
	Value               string `json:"value"`
}
