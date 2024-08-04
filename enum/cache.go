package enum

import "strconv"

const (
    CacheUserPrefix    = "user."
    CacheUserBlackList = CacheUserPrefix + "blacklist."
)

func GetCacheUserBlacklistKey(id uint) string {
    return CacheUserBlackList + strconv.FormatUint(uint64(id), 10)
}
