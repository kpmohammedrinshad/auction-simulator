package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// outputDir is the folder where per-auction result files are written.
const outputDir = "auction_outputs"

// NewSimulator applies resource constraints and initialises TotalBidders bidders.
func NewSimulator(cfg ResourceConfig) *Simulator {
	cfg.Apply()

	bidders := make([]Bidder, TotalBidders)
	for i := range bidders {
		bidders[i] = Bidder{ID: i + 1}
	}
	return &Simulator{Bidders: bidders, Config: cfg}
}

// Run launches all TotalAuctions concurrently, waits for all to finish,
// writes individual output files, then prints the summary and total elapsed time.
func (s *Simulator) Run() {
	// Create output directory if it doesn't exist
	os.MkdirAll(outputDir, os.ModePerm)

	results := make([]AuctionResult, TotalAuctions)
	var wg sync.WaitGroup

	globalStart := time.Now() // ← wall-clock start: before first auction fires

	for i := 0; i < TotalAuctions; i++ {
		wg.Add(1)
		go func(auctionID int) {
			defer wg.Done()
			results[auctionID-1] = RunAuction(auctionID, s.Bidders)
		}(i + 1)
	}

	wg.Wait()
	totalElapsed := time.Since(globalStart) // ← wall-clock end: after last auction done

	// Write one file per auction, then print console summary
	s.writeOutputFiles(results)
	s.printResults(results, totalElapsed)
}

// writeOutputFiles writes a dedicated .txt file for each auction result
// into the auction_outputs/ directory.
func (s *Simulator) writeOutputFiles(results []AuctionResult) {
	for _, r := range results {
		path := fmt.Sprintf("%s/auction_%02d.txt", outputDir, r.AuctionID)
		content := formatAuctionResult(r)

		// Write file — truncate if it already exists
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Printf("[Error] Could not write file for Auction #%02d: %v\n", r.AuctionID, err)
		}
	}
	fmt.Printf("[Output] %d auction files written to ./%s/\n\n", len(results), outputDir)
}

// formatAuctionResult builds the text content for a single auction output file.
func formatAuctionResult(r AuctionResult) string {
	separator := "======================================\n"
	content := separator
	content += fmt.Sprintf("  AUCTION #%02d — RESULT\n", r.AuctionID)
	content += separator
	content += fmt.Sprintf("  Auction ID   : %d\n", r.AuctionID)

	if r.WinnerID == -1 {
		content += "  Winner       : No winner (no bids received)\n"
		content += "  Winning Bid  : N/A\n"
	} else {
		content += fmt.Sprintf("  Winner       : Bidder #%03d\n", r.WinnerID)
		content += fmt.Sprintf("  Winning Bid  : $%.2f\n", r.WinBid)
	}

	content += fmt.Sprintf("  Total Bids   : %d / %d bidders\n", r.TotalBids, TotalBidders)
	content += fmt.Sprintf("  Duration     : %s\n", r.Duration.Round(time.Millisecond))
	content += separator
	return content
}

// printResults prints each auction outcome to the console with overall timing.
func (s *Simulator) printResults(results []AuctionResult, total time.Duration) {
	fmt.Println("========== AUCTION RESULTS ==========")
	for _, r := range results {
		if r.WinnerID == -1 {
			fmt.Printf("Auction #%02d | No bids received   | Bids: %d  | Duration: %s\n",
				r.AuctionID, r.TotalBids, r.Duration.Round(time.Millisecond))
		} else {
			fmt.Printf("Auction #%02d | Winner: Bidder #%03d | WinBid: $%.2f | Bids: %d | Duration: %s\n",
				r.AuctionID, r.WinnerID, r.WinBid, r.TotalBids, r.Duration.Round(time.Millisecond))
		}
	}
	fmt.Println("=====================================")
	fmt.Printf("[Timing] First auction started → last auction completed: %s\n\n",
		total.Round(time.Millisecond))
}
