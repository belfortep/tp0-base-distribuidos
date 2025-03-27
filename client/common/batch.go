package common

import "strings"

const MAX_MEMORY = 8192

// Batch entity, encapsulates how to serialize many bets
type Batch struct {
	bets        []Bet
	batchSize   int
	batchMemory int
}

// Creates a new batch
func NewBatch(batchSize int) Batch {
	return Batch{
		bets:        []Bet{},
		batchSize:   batchSize,
		batchMemory: MAX_MEMORY,
	}
}

// Verify if we can append a new bet to the batch,
// Verifying the number of bets we have and the memory size needed for the bet
func (batch *Batch) CantAppend(bet Bet) bool {
	return batch.batchSize-1 < 0 || batch.batchMemory-bet.MemorySize() <= 1
}

// Append a new bet to the batch if we can
func (batch *Batch) Append(bet Bet) {
	if batch.CantAppend(bet) {
		return
	}

	batch.bets = append(batch.bets, bet)
	batch.batchSize--
	batch.batchMemory -= bet.MemorySize()
}

// Serialize the batch, this means
// to serialize all the individual bets
func (batch *Batch) Serialize() string {
	serialized := ""
	for _, bet := range batch.bets {
		serialized += bet.Serialize()
	}
	serialized = strings.TrimSuffix(serialized, "\n")

	return serialized
}

// Returns if we have 0 bets in this batch
func (batch *Batch) isEmpty() bool {
	return len(batch.bets) == 0
}
