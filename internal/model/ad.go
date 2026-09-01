package model

import "time"

// Card contains data available in the category listing.
type Card struct {
	ID             int64
	URL            string
	Title          string
	Price          *int
	Rooms          *int
	Floor          *int
	Promoted       bool
	PromotionLabel string
	ImageURL       string
	AgeText        string
	Position       int
}

// Ad is the enriched representation read from an individual advertisement page.
type Ad struct {
	Card
	Description string
	SellerName  string
	SellerSince string
	SellerAds   *int
}

// RuntimeStatus is a small snapshot shown by /status.
type RuntimeStatus struct {
	Mode               string
	LastSuccessfulPoll time.Time
	NextPoll           time.Time
	BackoffUntil       time.Time
	LastError          string
	SeenCount          int64
	LastCardCount      int
	LastNewCount       int
	LastSentCount      int
}
