package todo

type TodoReader interface {
	List() []Todo
}

func CountPending(reader TodoReader) int {
	todos := reader.List()
	count := 0
	for _, t := range todos {
		if !t.Done {
			count++
		}
	}

	return count
}

type StatsService struct {
	reader TodoReader
}

func NewStatsService(reader TodoReader) *StatsService {
	return &StatsService{reader: reader}
}

func (s *StatsService) PendingCount() int {
	return CountPending(s.reader)
}
