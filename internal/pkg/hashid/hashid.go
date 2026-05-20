package hashid

import (
    `fmt`
    
    `cms/internal/config`
    
    `github.com/speps/go-hashids`
)

type HashID struct {
    *hashids.HashID
}

func New(config *config.Config) (*HashID, error) {
    hashData := hashids.NewData()
    hashData.Salt = config.AppKey
    hashData.MinLength = 16
    hd, err := hashids.NewWithData(hashData)
    if err != nil {
        return nil, err
    }
    return &HashID{
        HashID: hd,
    }, nil
}

func (h *HashID) EncodeUint64(n uint64) (string, error) {
    ints := []int64{int64(n >> 32), int64(n & 0xffffffff)}
    return h.HashID.EncodeInt64(ints)
}

func (h *HashID) DecodeUint64(hash string) (uint64, error) {
    decResult, err := h.HashID.DecodeInt64WithError(hash)
    if err != nil {
        return 0, err
    }
    if len(decResult) != 2 {
        return 0, fmt.Errorf("invalid hash %s", hash)
    }
    result := (uint64(decResult[0]) << 32) | uint64(decResult[1])
    return result, nil
}
