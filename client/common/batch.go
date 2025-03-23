package common

const MAX_MEMORY = 8192

type Batch struct {
	bets        []Bet
	batchSize   int
	batchMemory int
}

func NewBatch(batchSize int) Batch {
	return Batch{
		bets:        []Bet{},
		batchSize:   batchSize,
		batchMemory: MAX_MEMORY,
	}
}

func (batch *Batch) CanAppend(bet Bet) bool {
	return batch.batchSize-1 < 0 || batch.batchMemory-bet.MemorySize() <= 1
}

func (batch *Batch) Append(bet Bet) {
	if !batch.CanAppend(bet) {
		return
	}

	batch.bets = append(batch.bets, bet)
	batch.batchSize--
	batch.batchMemory -= bet.MemorySize()
}

func (batch *Batch) Serialize() string {
	serialized := ""
	for _, bet := range batch.bets {
		serialized += bet.Serialize()
	}
	serialized += "\000"

	return serialized
}
