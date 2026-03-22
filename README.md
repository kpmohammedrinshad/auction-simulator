# Auction Simulator

A concurrent auction simulator written in Go that runs multiple auctions simultaneously, collects bids from simulated bidders, measures execution time, and standardizes resource usage.

---

## Overview

- **100 bidders** participate across all auctions
- Each auction item is described by **20 attributes**
- **40 auctions** run concurrently at the same time
- Each auction closes after a **3-second timeout**
- Not every bidder responds — 75% response rate is simulated
- Total wall-clock time is measured from the first auction start to the last auction finish

---

## Project Structure

```
auction-simulator/
├── main.go          # Entry point — seeds RNG, builds config, starts simulator
├── constants.go     # All tuneable parameters in one place
├── models.go        # All struct definitions (Item, Bid, Bidder, AuctionResult, etc.)
├── item.go          # NewItem() — generates auction item with random attributes
├── bidder.go        # Bidder.PlaceBid() — simulates delay, response rate, bid amount
├── auction.go       # RunAuction() — fan-out to bidders, collect bids, enforce timeout
├── resource.go      # ResourceConfig — vCPU and RAM standardization
├── simulator.go     # Simulator — orchestrates auctions, timing, and result reporting
└── go.mod           # Go module definition
```

---

## Getting Started

### Prerequisites

- Go **1.21** or higher
- No external dependencies — standard library only

### Run

```bash
go run .
```

### Build

```bash
go build -o auction-simulator
./auction-simulator
```

---

## Configuration

All parameters are defined in `constants.go`:

| Constant         | Default       | Description                                      |
|------------------|---------------|--------------------------------------------------|
| `TotalBidders`   | `100`         | Number of simulated bidders                      |
| `TotalAttributes`| `20`          | Number of attributes per auction item            |
| `TotalAuctions`  | `40`          | Number of auctions running concurrently          |
| `AuctionTimeout` | `3s`          | Time limit per auction before it closes          |
| `BidderDelay`    | `500` ms      | Max random delay before a bidder responds        |
| `BidResponseRate`| `0.75`        | Probability that a bidder places a bid           |

### Resource Configuration

CPU and RAM constraints are set in `main.go`:

```go
// Auto-detect all available cores (default)
cfg := DefaultResourceConfig()

// Or set manually
cfg := ResourceConfig{VCPUs: 4, RAMMiB: 256}
```

`VCPUs` maps directly to `runtime.GOMAXPROCS(n)`, controlling how many OS threads
run Go goroutines simultaneously.  
`RAMMiB` is informational — enforce it at the container level (Docker/Kubernetes).

---

## How It Works

### Auction Flow

```
Simulator.Run()
    │
    ├── Spawns 40 goroutines (one per auction) ──────────────────────────────┐
    │                                                                         │
    │   RunAuction(id)                                                        │
    │       ├── Generates Item with 20 random attributes                     │
    │       ├── Spawns 100 Bidder goroutines                                 │
    │       │       └── Each bidder: random delay → maybe bid → send to ch  │
    │       ├── Collects bids until timeout OR all bidders done              │
    │       └── Returns AuctionResult (winner, bid amount, count, duration)  │
    │                                                                         │
    └── wg.Wait() — blocks until ALL 40 auctions complete ───────────────────┘
            │
            └── Prints results + total wall-clock time
```

### Timing Measurement

| Marker        | When it is captured                              |
|---------------|--------------------------------------------------|
| `globalStart` | Just before the first auction goroutine spawns   |
| `totalElapsed`| After `wg.Wait()` — all 40 auctions are done    |

### Resource Standardization

| Mechanism       | Implementation                        |
|-----------------|---------------------------------------|
| CPU             | `runtime.GOMAXPROCS(vCPUs)`           |
| RAM (soft cap)  | `ResourceConfig.RAMMiB` (informational)|
| RAM (hard cap)  | Docker `--memory` / k8s `resources.limits.memory` |

---

## Sample Output

```
[Resource] GOMAXPROCS → 8 vCPU(s) | RAM cap: 512 MiB

========== AUCTION RESULTS ==========
Auction #01 | Winner: Bidder #042 | WinBid: $18432.57 | Bids: 76 | Duration: 498ms
Auction #02 | Winner: Bidder #087 | WinBid: $19104.33 | Bids: 74 | Duration: 501ms
Auction #03 | No bids received    | Bids: 0           | Duration: 3000ms
...
Auction #40 | Winner: Bidder #013 | WinBid: $17893.21 | Bids: 71 | Duration: 487ms
=====================================
[Timing] First auction started → last auction completed: 3001ms
```

---

## Memory Considerations

| Source                        | Approx. footprint              |
|-------------------------------|--------------------------------|
| 4,000 goroutines (peak)       | ~32 MB (8 KB stack each)       |
| 40 buffered channels (×100)   | Negligible                     |
| 40 `AuctionResult` structs    | Negligible                     |
| **Total (approx.)**           | **~35–40 MB**                  |

> Splitting code into multiple `.go` files has **zero effect** on memory —
> Go compiles all files in a package into one binary regardless.

---

👤 Author
Mohammed Rinshad K P GitHub: @kpmohammedrinshad
