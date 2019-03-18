package caddy_cas

import (
	"github.com/go-cas/cas"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"net/http"
	"net/url"
)

func init() {
	caddy.RegisterPlugin("cas", caddy.Plugin{
		ServerType: "http",
		Action: setup,
	})
}

func setup (c *caddy.Controller) error {
	opt, err := parseConfig(c)
	if err != nil {
		return err
	}

	s := httpserver.GetConfig(c)
	s.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return newCasHandlerWithOptions(next, opt)
	})

	return nil
}

func parseConfig(c *caddy.Controller) (opts cas.Options, err error) {
	if c.Next() {
		v := c.Val()
		opts.URL, err = url.Parse(v)
		if err != nil {
			return
		}
	}
	return
}

type casHandler struct {
	client *cas.Client
	Next httpserver.Handler
}

func (h *casHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return http.StatusFound, nil
	}
	return h.Next.ServeHTTP(w, r)
}

func newCasHandlerWithOptions(next httpserver.Handler, options cas.Options) *casHandler {
	client := cas.NewClient(&options)
	return &casHandler{
		Next: next,
		client: client,
	}
}