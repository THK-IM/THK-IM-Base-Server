package pool

import "time"

const (
	StrategyRandom = "random"
	StrategyScore  = "score"
)

type (
	Element struct {
		Id    int64   `json:"id"`
		Score float64 `json:"score"`
	}

	RecommendStrategy struct {
		Type             string `json:"type"`
		Count            int    `json:"count"`
		RepeatRetryCount int    `json:"repeat_retry_count"`
	}

	MatchFunction func(id string, candidateId string) (matchedId *string, putBlack bool, err error)

	MatchPool interface {
		Clear() error
		Add(ids ...string) (int64, error)
		Remove(ids ...string) (int64, error)
		Contain(id string) (bool, error)
		Count() (int64, error)
		Match(uId string, f MatchFunction) (matchedId *string, err error)
	}

	RecommendPool interface {
		// Add 添加元素
		Add(elements ...*Element) error
		// Remove 移除元素
		Remove(elements ...*Element) error
		// RemoveByScore 移除<=score的元素
		RemoveByScore(score float64) (int64, error)
		// ElementCount 元素数量
		ElementCount() (int64, error)
		// ElementCountByRange [score1, score2]分数段的元素数量
		ElementCountByRange(score1, score2 float64) (int64, error)
		// FetchElements 根据uId和多个策略取出元素
		FetchElements(uId int64, strategies []*RecommendStrategy) ([]Element, error)
		// AddUserRecord 添加用户推荐记录
		AddUserRecord(uId int64, ex time.Duration, elementId ...int64) error
		// ClearUserRecord 清除用户推荐记录
		ClearUserRecord(uId int64) error
	}
)
