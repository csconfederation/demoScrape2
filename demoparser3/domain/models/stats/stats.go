package stats

type Stats struct {
	StatsByPlayer       map[uint64]*PlayerStats
	CTSideStatsByPlayer map[uint64]*PlayerStats
	TSideStatsByPlayer  map[uint64]*PlayerStats
	StatsByTeam         map[string]*TeamStats
}
