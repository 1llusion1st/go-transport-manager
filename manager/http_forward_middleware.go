package manager

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func fiberHTTPForwardHandler(config ForwardConfig) fiber.Handler {
	client := &fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return func(c *fiber.Ctx) error {
		req := &fasthttp.Request{}
		c.Request().CopyTo(req)
		resp := c.Response()

		logrus.Debugf("req: %p, resp: %p", req, resp)

		originalURL := utils.CopyString(c.OriginalURL())
		originalPath := string(req.URI().Path())
		logrus.Debugf("originalURL: %s, path: %s", originalURL, originalPath)
		defer req.SetRequestURI(originalURL)

		copiedURL := utils.CopyString(config.Target)
		targetPath, err := url.ParseRequestURI(config.Target)
		if err != nil {
			return err
		}
		logrus.Debugf("new url: %s", copiedURL)
		req.SetRequestURI(copiedURL)
		if scheme := getScheme(utils.UnsafeBytes(copiedURL)); len(scheme) > 0 {
			req.URI().SetSchemeBytes(scheme)
		}
		if string(originalPath) != "/" {
			logrus.Debugf("origpath != /: %s", originalPath)
			remotePath := targetPath.Path + string(originalPath)
			remotePath = strings.Replace(remotePath, "//", "/", 0)
			logrus.Debugf("remote path: %s", remotePath)
			req.URI().SetPath(remotePath)
		}
		if len(config.Headers) > 0 {
			for _, header := range config.Headers {
				req.Header.Set(header.Name, header.Value)
			}
		}
		logrus.Debugf("new url: %s", req.URI().FullURI())

		req.Header.Del(fiber.HeaderConnection)
		if err := client.DoRedirects(req, resp, 1); err != nil {
			return err
		}
		resp.Header.Del(fiber.HeaderConnection)
		c.Status(http.StatusOK)

		return nil
	}
}

func getScheme(uri []byte) []byte {
	i := bytes.IndexByte(uri, '/')
	if i < 1 || uri[i-1] != ':' || i == len(uri)-1 || uri[i+1] != '/' {
		return nil
	}
	return uri[:i-1]
}
