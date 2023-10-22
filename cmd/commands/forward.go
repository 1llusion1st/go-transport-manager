package commands

import (
	"context"
	"strings"

	"github.com/1llusion1st/go-transport-manager/manager"
)

type ForwardCmd struct {
	manager.ForwardConfig
	RawHeaders []string `arg:"" name:"headers" help:"extra headers" optional:"true"`
}

func (f *ForwardCmd) Run(ctx *Context) error {

	if len(f.RawHeaders) > 0 {
		for _, header := range f.RawHeaders {
			if strings.Count(header, ":") > 0 {
				splited := strings.SplitN(header, ":", 2)
				if len(splited) != 2 {
					continue
				}
				f.Headers = append(f.Headers, manager.ForwardExtraHeader{
					Name:  splited[0],
					Value: splited[1],
				})
			} else {
				continue
			}
		}
	}

	forwarder, err := manager.NewHTTPForward(f.ForwardConfig)
	if err != nil {
		return err
	}
	err = forwarder.Start(context.Background())
	if err != nil {
		return err
	}
	select {}
}
