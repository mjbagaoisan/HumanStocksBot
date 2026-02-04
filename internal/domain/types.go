package domain

import "time"

// MarketStatus represents the state of a user's market
type MarketStatus string

const (
	MarketStatusActive     MarketStatus = "ACTIVE"
	MarketStatusSunsetting MarketStatus = "SUNSETTING"
	MarketStatusClosed     MarketStatus = "CLOSED"
)

// GuildConfig stores per-guild configuration
type GuildConfig struct {
	GuildID       string
	StartingCash  int64 // Amount given to new users (in cents)
	BasePrice     int64 // Base price for bonding curve
	Slope         int64 // Slope for bonding curve
	TradeFeeBps   int64 // Total fee in basis points (e.g., 200 = 2%)
	SubjectFeeBps int64 // Portion of fee to subject (e.g., 100 = 1%)
	TradingPaused bool
	SunsetDays    int // Days before market closes after optout
}

// GuildMember represents a user's membership in a guild
type GuildMember struct {
	GuildID    string
	UserID     string
	OptedIn    bool
	OptedInAt  *time.Time
	OptedOutAt *time.Time
}

// Wallet holds a user's cash balance
type Wallet struct {
	GuildID   string
	UserID    string
	Cash      int64 // Balance in cents
}

// Market represents a tradable user's market
type Market struct {
	GuildID           string
	SubjectUserID     string
	Status            MarketStatus
	SharesOutstanding int64
	ReserveBalance    int64 // Money held in bonding curve
	LastPrice         int64 // Price of most recent trade
	SunsetAt          *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Holding represents shares owned by a user
type Holding struct {
	GuildID       string
	OwnerUserID   string // Who owns the shares
	SubjectUserID string // Whose shares they are
	Shares        int64
	AvgCost       int64 // Average cost basis per share (for P&L)
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Trade records a completed trade
type Trade struct {
	ID                     string
	GuildID                string
	TraderUserID           string
	SubjectUserID          string
	Side                   string // "BUY", "SELL", or "LIQUIDATION"
	Shares                 int64
	GrossAmount            int64
	FeeAmount              int64
	SubjectFeeAmount       int64
	SystemFeeAmount        int64
	NetAmount              int64
	PriceBefore            int64
	PriceAfter             int64
	SharesOutstandingAfter int64
	ReserveBalanceAfter    int64
	IdempotencyKey         string
	CreatedAt              time.Time
}

// Treasury tracks system fees collected
type Treasury struct {
	GuildID    string
	SystemFees int64 // Total fees collected
	UpdatedAt  time.Time
}
