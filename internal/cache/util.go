package cache

func makeRedisKey(entity, key string) string {
	return entity + "." + key
}
