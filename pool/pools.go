package pool

const (
	StrategyRandom    = "random"
	StrategyScoreAsc  = "score-asc"
	StrategyScoreDesc = "score-desc"
)

type (
	Element struct {
		Id    int64 `json:"id"`
		Score int64 `json:"score"`
	}

	Strategy struct {
		Type  string `json:"type"`
		Count uint16 `json:"count"`
	}

	MatchPool interface {
		Add(elements ...*Element) error
		Remove(elements ...*Element) error
		MemberCount() (int64, error)
		FetchMembers(uId int64, strategies []*Strategy) ([]*Element, error)
	}

	RecommendPool interface {
		Add(elements ...*Element) error
		Remove(elements ...*Element) error
		MemberCount() (int64, error)
		FetchMembers(uId int64, strategies []*Strategy) ([]*Element, error)
	}
)
