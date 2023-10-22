package manager

import (
	"context"
	"fmt"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewHTTPForward(conf ForwardConfig) (Forwarder, error) {
	inst := &httpForwarder{
		config:     conf,
		serveError: make(chan error),
		app:        fiber.New(),
	}

	return inst, nil
}

var _ Forwarder = &httpForwarder{}

type httpForwarder struct {
	config     ForwardConfig
	serveError chan error
	app        *fiber.App
}

func (forwarder *httpForwarder) Start(ctx context.Context) (err error) {
	app := forwarder.app
	app.Use(logger.New())

	path := forwarder.config.SourcePathPrefix + "*"

	app.Get(path, fiberHTTPForwardHandler(forwarder.config))
	app.Post(path, fiberHTTPForwardHandler(forwarder.config))
	app.Put(path, fiberHTTPForwardHandler(forwarder.config))
	app.Delete(path, fiberHTTPForwardHandler(forwarder.config))

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", forwarder.config.ListenPort))
	if err != nil {
		return err
	}

	go func() {
		forwarder.serveError <- app.Listener(listener)
	}()
	return nil
}

func (forwarder *httpForwarder) Stop(ctx context.Context) error {
	forwarder.app.ShutdownWithContext(ctx)

	return nil
}

func (forwarder *httpForwarder) Error() error {
	select {
	case err := <-forwarder.serveError:
		return err
	default:
		return nil
	}
}
