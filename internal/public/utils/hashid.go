package utils

import (
    `fmt`
    
    `github.com/speps/go-hashids`
)

func HashIDEncode(h *hashids.HashID, n uint64) (string, error) {
    ints := []int64{int64(n >> 32), int64(n & 0xffffffff)}
    return h.EncodeInt64(ints)
}

func HashIDDecode(h *hashids.HashID, hash string) (uint64, error) {
    dec, err := h.DecodeInt64WithError(hash)
    if err != nil {
        return 0, err
    }
    if len(dec) != 2 {
        return 0, fmt.Errorf("invalid hash %s", hash)
    }
    result := (uint64(dec[0]) << 32) | uint64(dec[1]&0xffff)
    return result, nil
}
