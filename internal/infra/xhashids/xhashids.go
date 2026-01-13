package xhashids

import `github.com/speps/go-hashids`

type HashID struct {
    *hashids.HashID
}

type Config struct {
    Salt string `json:"salt" yaml:"salt"`
}

func New(config Config) (*HashID, error) {
    hashData := hashids.NewData()
    hashData.Salt = config.Salt
    hashData.MinLength = 16
    hd, err := hashids.NewWithData(hashData)
    if err != nil {
        return nil, err
    }
    return &HashID{
        HashID: hd,
    }, nil
}

func (h *HashID) EncodeUint64(numbers []uint64) (string, error) {
    ints := make([]int64, 0, len(numbers)*2)
    for _, number := range numbers {
        ints = append(ints, int64(number>>32), int64(number&0xffffffff))
    }
    return h.HashID.EncodeInt64(ints)
}

func (h *HashID) DecodeUint64(hash string) ([]uint64, error) {
    decResult, err := h.HashID.DecodeInt64WithError(hash)
    if err != nil {
        return nil, err
    }
    capa := len(decResult) / 2
    result := make([]uint64, 0, capa)
    for i := 0; i < capa; i += 2 {
        hi := uint64(decResult[i] << 32)
        li := uint64(decResult[i+1])
        n := hi | li
        result = append(result, n)
    }
    return result, nil
}
