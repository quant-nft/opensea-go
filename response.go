package opensea_go

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/xyths/hs/convert"
	"time"
)

type ResponseEvent struct {
	Success     *bool
	AssetEvents []AssetEvent `json:"asset_events"`
}

type AssetEvent struct {
	Asset Asset `json:"asset"`
	// transfer: Mint或者是真实的Transfer
	// created: List
	// bid_entered: Bid
	// bid_withdrawn: Bid Cancel
	// successful: Sale
	// offer_entered: Offer
	EventType string `json:"event_type"`

	// used when EventType = `bid_entered` or `offer_entered`, means bid price
	BidAmount string `json:"bid_amount"`

	// used when EventType = `created`, means list price.
	EndingPrice string `json:"ending_price"`

	// used when EventType = `successful` or `bid_withdrawn`, means sale or cancel offer
	TotalPrice string `json:"total_price"`

	CreatedDate   string   `json:"created_date"`
	FromAccount   *Account `json:"from_account"`
	ToAccount     *Account `json:"to_account"`
	Owner         *Account `json:"owner"`
	Seller        *Account `json:"seller"`
	WinnerAccount *Account `json:"winner_account"`

	PaymentToken PaymentToken `json:"payment_token"`
}

const (
	EventTypeTransfer   = "transfer"
	EventTypeList       = "created"
	EventTypeListCancel = "cancelled"
	EventTypeBid        = "bid_entered"
	EventTypeBidCancel  = "bid_withdrawn"
	EventTypeSale       = "successful"
	EventTypeOffer      = "offer_entered"
)

type Asset struct {
	TokenId       string          `json:"token_id"`
	Name          string          `json:"name"`
	AssetContract AssetContract   `json:"asset_contract"`
	Collection    AssetCollection `json:"collection"`

	ImagePreviewUrl string `json:"image_preview_url"`
}

type AssetContract struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type AssetCollection struct {
	Name string `json:"name"`
}

type Account struct {
	User          User   `json:"user"`
	ProfileImgUrl string `json:"profile_img_url"`
	Address       string `json:"address"`
	Config        string `json:"config"`
}

func (a Account) String() string {
	addr := convert.ShortAddress(a.Address)
	if a.User.Username != "" {
		return fmt.Sprintf("%s(%s)", a.User.Username, addr)
	} else {
		return addr
	}
}

type User struct {
	Username string `json:"username"`
}

type PaymentToken struct {
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

// ResponseCollections is response of `/collection` API
type ResponseCollections struct {
	Collections []RawCollection `json:"collections"`
}

// RawCollection is the `collection` structure in ResponseCollections, the response of `/collection` API.
// It's different from the ResponseEvent.
type RawCollection struct {
	Name                  string
	Description           string
	PrimaryAssetContracts []RawAssetContract `json:"primary_asset_contracts"`
	Stats                 RawStat            `json:"stats"`
	CreatedDate           string             `json:"created_date"`
}

type RawStat struct {
	OneDayVolume       float64 `json:"one_day_volume" bson:"oneDayVolume"`
	OneDayChange       float64 `json:"one_day_change" bson:"oneDayChange"`
	OneDaySales        float64 `json:"one_day_sales" bson:"oneDaySales"`
	OneDayAveragePrice float64 `json:"one_day_average_price" bson:"oneDayAveragePrice"`

	SevenDayVolume       float64 `json:"seven_day_volume" bson:"sevenDayVolume"`
	SevenDayChange       float64 `json:"seven_day_change" bson:"sevenDayChange"`
	SevenDaySales        float64 `json:"seven_day_sales" bson:"sevenDaySales"`
	SevenDayAveragePrice float64 `json:"seven_day_average_price" bson:"sevenDayAveragePrice"`

	ThirtyDayVolume       float64 `json:"thirty_day_volume" bson:"thirtyDayVolume"`
	ThirtyDayChange       float64 `json:"thirty_day_change" bson:"thirtyDayChange"`
	ThirtyDaySales        float64 `json:"thirty_day_sales" bson:"thirtyDaySales"`
	ThirtyDayAveragePrice float64 `json:"thirty_day_average_price" bson:"thirtyDayAveragePrice"`

	TotalVolume float64 `json:"total_volume" bson:"totalVolume"`
	TotalSales  float64 `json:"total_sales" bson:"totalSales"`
	TotalSupply float64 `json:"total_supply" bson:"totalSupply"`

	Count        float64 `json:"count" bson:"count"`
	NumOwners    float64 `json:"num_owners" bson:"numOwners"`
	AveragePrice float64 `json:"average_price" bson:"averagePrice"`
	NumReports   float64 `json:"num_reports" bson:"numReports"`
	MarketCap    float64 `json:"market_cap" bson:"marketCap"`
	FloorPrice   float64 `json:"floor_price" bson:"floorPrice"`
}

type RawAssetContract struct {
	Address string
	Name    string
}

func (ae AssetEvent) ToEvent() Event {
	e := Event{
		Collection:      ae.Asset.Collection.Name,
		Contract:        common.HexToAddress(ae.Asset.AssetContract.Address).Hex(),
		Name:            ae.Asset.Name,
		Id:              ae.Asset.TokenId,
		Date:            toBeijingTime(ae.CreatedDate),
		CreatedAt:       time.Now(),
		ImagePreviewUrl: ae.Asset.ImagePreviewUrl,
	}
	switch ae.EventType {
	case EventTypeTransfer:
		if ae.FromAccount != nil && ae.FromAccount.Address != "0x0000000000000000000000000000000000000000" {
			e.Type = EventTransfer
		} else {
			e.Type = EventMint
		}
		// Mint 和 Transfer 都没有价格需要展示
		if ae.FromAccount != nil {
			e.From = ae.FromAccount.String()
		}
		if ae.ToAccount != nil {
			e.To = ae.ToAccount.String()
		}
	case EventTypeList:
		e.Type = EventList
		if ae.FromAccount != nil {
			e.From = ae.FromAccount.String()
		}
		e.Price = toEther(ae.EndingPrice, ae.PaymentToken)
	case EventTypeListCancel:
		e.Type = EventListCancel
		if ae.Seller != nil {
			e.From = ae.Seller.String()
		}
		e.Price = toEther(ae.EndingPrice, ae.PaymentToken)
	case EventTypeBid:
		e.Type = EventBid
		if ae.FromAccount != nil {
			e.From = ae.FromAccount.String()
		}
		e.Price = toEther(ae.BidAmount, ae.PaymentToken)
	case EventTypeBidCancel:
		e.Type = EventBidCancel
		if ae.FromAccount != nil {
			e.From = ae.FromAccount.String()
		}
		e.Price = toEther(ae.TotalPrice, ae.PaymentToken)
	case EventTypeSale:
		e.Type = EventSale
		e.Price = toEther(ae.TotalPrice, ae.PaymentToken)
		if ae.Seller != nil {
			e.From = ae.Seller.String()
		}
		if ae.WinnerAccount != nil {
			e.To = ae.WinnerAccount.String()
		}
	case EventTypeOffer:
		e.Type = EventOffer
		if ae.FromAccount != nil {
			e.From = ae.FromAccount.String()
		}
		e.Price = toEther(ae.BidAmount, ae.PaymentToken)
	default:
		e.Type = ae.EventType
	}

	return e
}

// "2021-08-28T09:44:43.664713"
func toBeijingTime(date string) string {
	//secondsEastOfUTC := int((8 * time.Hour).Seconds())
	//beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	layout := "2006-01-02T15:04:05.999999"
	t, err := time.Parse(layout, date)
	if err != nil {
		return date
	}
	onlyTime := "15:04:05"
	return t.Local().Format(onlyTime)
}

func toEther(price string, payment PaymentToken) string {
	unit := decimal.New(1, int32(payment.Decimals))
	d, err := decimal.NewFromString(price)
	if err != nil {
		return price
	}
	ret := fmt.Sprintf("%s %s", d.Div(unit).String(), payment.Symbol)
	return ret
}
