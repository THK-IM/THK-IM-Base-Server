package pool

type Factory interface {
	NewMatchPool(key string) MatchPool
	NewRecommendPool(key string) RecommendPool
}
