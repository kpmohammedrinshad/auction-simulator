package main

import "time"

// Item represents the object being auctioned, described by TotalAttributes attributes.
type Item struct {
	ID         int
	Attributes [TotalAttributes]float64 // 20 numeric attributes broadcast to all bidders
}

// Bid represents a single offer submitted by a bidder for an auction item.
type Bid struct {
	BidderID int
	Amount   float64
}

// Bidder represents a participant who evaluates items and optionally places bids.
type Bidder struct {
	ID int
}

// AuctionResult captures the final outcome of a completed auction.
type AuctionResult struct {
	AuctionID int
	WinnerID  int           // -1 when no bids were received
	WinBid    float64       // Winning bid amount (0 if no winner)
	TotalBids int           // Count of bids received before timeout
	Duration  time.Duration // Wall-clock time from auction start to close
}

// ResourceConfig specifies the vCPU and RAM constraints for the simulator.
type ResourceConfig struct {
	VCPUs  int // Logical CPU cores (sets GOMAXPROCS)
	RAMMiB int // Informational RAM cap in MiB (for deployment tooling)
}

// Simulator orchestrates all concurrent auctions.
type Simulator struct {
	Bidders []Bidder
	Config  ResourceConfig
}
