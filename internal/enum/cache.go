package enum

import "strconv"

const (
    CacheUserPrefix    = "User."
    CacheUserBlackList = CacheUserPrefix + "Blacklist."
)

func GetCacheUserBlacklistKey(id int64) string {
    return CacheUserBlackList + strconv.FormatInt(id, 10)
}
