package stores

// StatsStore defines storage behavior (memory / Redis interchangeable)
type StatsStore interface {
	Increment(userID string) error
	GetStats() (map[string]int, error)
}
