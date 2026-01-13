package bootstrap

import (
    `dpcms/internal/config`
    `dpcms/internal/infra/xhashids`
)

func NewXHashIds(cfg *config.Config) (*xhashids.HashID, error) {
    return xhashids.New(xhashids.Config{Salt: cfg.GlobalKey})
}
