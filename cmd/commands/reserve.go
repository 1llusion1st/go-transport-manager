package commands

import (
	"context"

	"github.com/1llusion1st/go-transport-manager/manager"
)

type ReserveCommand struct {
	manager.ReserveConfig
}

func (r *ReserveCommand) Run(ctx *Context) error {
	reserve, err := manager.NewTCPReserve(r.ReserveConfig)
	if err != nil {
		return err
	}

	err = reserve.Start(context.Background())
	if err != nil {
		return err
	}
	select {}
}
