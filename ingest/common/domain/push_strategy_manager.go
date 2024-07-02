package commondomain

// PushStrategyManager is an interface for managing the strategy of pushing the blocks.
// Either all block data or only the block update are the possible options
type PushStrategyManager interface {
	// ShouldPushAllData returns true if all data should be pushed.
	ShouldPushAllData() bool

	// MarkInitialDataIngested marks the initial data as ingested.
	// After calling this function, ShouldPushAllData should return false.
	MarkInitialDataIngested()

	// MarkErrorObserved marks that an error has been observed.
	MarkErrorObserved()
}

type coldStartManager struct {
	shouldPushAllData bool
}

var _ PushStrategyManager = &coldStartManager{}

// NewPushStrategyManager creates a new push strategy manager.
// It is initialized with the strategy of pushing all data.
func NewPushStrategyManager() PushStrategyManager {
	return &coldStartManager{
		shouldPushAllData: true,
	}
}

func (c *coldStartManager) ShouldPushAllData() bool {
	return c.shouldPushAllData
}

func (c *coldStartManager) MarkInitialDataIngested() {
	c.shouldPushAllData = false
}

func (c *coldStartManager) MarkErrorObserved() {
	c.shouldPushAllData = true
}
