package opensea_go

import (
	"fmt"
	"strings"
	"time"
)

// Event is from OpenSea event, and reorganized for messaging.
type Event struct {
	Collection string `json:"collection"` // collection name
	Contract   string `json:"contract"`   // collection contract address
	Name       string `json:"name"`       // NFT name
	Id         string `json:"id"`
	Type       string `json:"event"`
	Price      string `json:"price"`
	From       string `json:"from"`
	To         string `json:"to"`
	Date       string `json:"date"`

	ImagePreviewUrl string `json:"imagePreviewUrl"` // for Telegram preview

	CreatedAt time.Time `json:"createdAt"`
}

const (
	EventSale      = "Sale"
	EventOffer     = "Offer"
	EventBid       = "Bid"
	EventBidCancel = "Bid Cancel"
	EventTransfer  = "Transfer"
	EventMint      = "Mint"
	EventList      = "List"
	// EventListCancel is list cancel, OpenSea display "Cancel".
	// We display "List Cancel", to distinguish from "Bid Cancel"
	EventListCancel = "List Cancel"
)

type Layout struct {
	Project      bool // display project
	ImagePreview bool //
}

// FormatDing format text message for Dingding.
func (e Event) FormatDing(layout Layout) string {
	var content string
	if layout.Project {
		content = fmt.Sprintf("项目: %s\n", e.Collection)
	}
	content += fmt.Sprintf("名称: %s\nId: %s", e.Name, e.Id)
	switch e.Type {
	case EventSale:
		content += fmt.Sprintf(
			` 成交(Sale)
买家: %s
卖家: %s
价格: %s`,
			e.To, e.From, e.Price,
		)
	case EventOffer:
		content += fmt.Sprintf(
			` 出价(Offer)
买家: %s
价格: %s`,
			e.From, e.Price,
		)
	case EventBid:
		content += fmt.Sprintf(
			` 出价(Bid)
买家: %s
价格: %s`,
			e.From, e.Price,
		)
	case EventTypeBidCancel:
		content += fmt.Sprintf(
			` 撤销出价(Bid Cancel)
买家: %s
价格: %s`,
			e.From, e.Price,
		)
	case EventTransfer:
		content += fmt.Sprintf(
			` 转让(Transfer)
发送方: %s
接收方: %s`,
			e.From, e.To,
		)
	case EventMint:
		content += fmt.Sprintf(
			` 铸造完成 (Mint)
接收方: %s`,
			e.To,
		)
	case EventListCancel:
		content += fmt.Sprintf(
			` 取消拍卖(List Cancel)
卖家: %s
价格: %s`,
			e.From, e.Price,
		)
	case EventList:
		content += fmt.Sprintf(
			` 拍卖(List)
卖家: %s
价格: %s`,
			e.From, e.Price,
		)
	default:
	}
	content += fmt.Sprintf("\n时间: %s", e.Date)

	content += fmt.Sprintf(
		"\n地址: https://opensea.io/assets/%s/%s",
		strings.ToLower(e.Contract), e.Id,
	)
	if layout.ImagePreview {
		content += fmt.Sprintf("\n预览图片: %s", e.ImagePreviewUrl)
	}

	return content
}

func (e Event) FormatDiscord(layout Layout) string {
	var content string
	if layout.Project {
		content = fmt.Sprintf("**%s**\n", e.Collection)
	}
	content += fmt.Sprintf("**%s** **%s**号 ", e.Name, e.Id)
	switch e.Type {
	case EventSale:
		content += fmt.Sprintf(
			` **成交** (Sale)
  买家: %s
  卖家: %s
  价格: __**%s**__`,
			e.To, e.From, e.Price,
		)
	case EventOffer:
		content += fmt.Sprintf(
			` **出价** (Offer)
  买家: %s
  价格: __**%s**__`,
			e.From, e.Price,
		)
	case EventBid:
		content += fmt.Sprintf(
			` **出价** (Bid)
  买家: %s
  价格: __**%s**__`,
			e.From, e.Price,
		)
	case EventTypeBidCancel:
		content += fmt.Sprintf(
			` **撤销出价** (Bid Cancel)
  买家: %s
  价格: __**%s**__`,
			e.From, e.Price,
		)
	case EventTransfer:
		content += fmt.Sprintf(
			` **转让** (Transfer)
  发送方: %s
  接收方: %s`,
			e.From, e.To,
		)
	case EventMint:
		content += fmt.Sprintf(
			` **铸造完成** (Mint)
  接收方: %s`,
			e.To,
		)
	case EventList:
		content += fmt.Sprintf(
			` **拍卖** (List)
  卖家: %s
  价格: __**%s**__`,
			e.From, e.Price,
		)
	case EventListCancel:
		content += fmt.Sprintf(
			` **取消拍卖** (List Cancel)
  卖家: %s
  价格: __**%s**__`,
			e.From, e.Price,
		)
	default:
	}
	content += fmt.Sprintf("\n  时间: %s", e.Date)
	if layout.ImagePreview {
		content += fmt.Sprintf("\n地址: https://opensea.io/assets/%s/%s\n预览图片: \n%s", strings.ToLower(e.Contract), e.Id, e.ImagePreviewUrl)
	}
	return content
}

func (e Event) FormatTelegram(layout Layout) string {
	var content string
	if layout.Project {
		content = fmt.Sprintf("项目: %s\n", e.Collection)
	}
	content += fmt.Sprintf("名称: %s\nTokenId: %s", e.Name, e.Id)
	switch e.Type {
	case EventSale:
		content += fmt.Sprintf(
			` 成交(Sale)
买家: %s
卖家: %s
价格: %s`,
			e.To, e.From, e.Price,
		)
	case EventOffer:
		content += fmt.Sprintf(
			` 出价(Offer)
买家: %s
价格: %s`,
			e.From, e.Price,
		)
	case EventBid:
		content += fmt.Sprintf(
			` 出价(Bid)
买家: %s
价格: %s`,
			e.From, e.Price,
		)
	case EventTypeBidCancel:
		content += fmt.Sprintf(
			` 撤销出价(Bid Cancel)
买家: %s
价格: %s`,
			e.From, e.Price,
		)
	case EventTransfer:
		content += fmt.Sprintf(
			` 转让(Transfer)
发送方: %s
接收方: %s`,
			e.From, e.To,
		)
	case EventMint:
		content += fmt.Sprintf(
			` 铸造完成 (Mint)
接收方: %s`,
			e.To,
		)
	case EventList:
		content += fmt.Sprintf(
			` 拍卖(List)
卖家: %s
价格: %s`,
			e.From, e.Price,
		)
	case EventListCancel:
		content += fmt.Sprintf(
			` 取消拍卖(List Cancel)
卖家: %s
价格: %s`,
			e.From, e.Price,
		)
	default:
	}
	content += fmt.Sprintf("\n时间: %s", e.Date)
	if layout.ImagePreview {
		content += fmt.Sprintf("\n地址: https://opensea.io/assets/%s/%s\n预览图片: %s", strings.ToLower(e.Contract), e.Id, e.ImagePreviewUrl)
	}
	return content
}
