package router

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/gectou4/rest_api_go/internal/config"
)

type Route struct {
	Methods      []string
	Regex        *regexp.Regexp
	Handler      string
	ParamNames   []string
}

type Router struct {
	routes []Route
}

func NewRouter(routes []config.Route) *Router {
	r := &Router{}
	for _, route := range routes {
		methods := strings.Split(route.Method, "|")
		pattern := route.Pattern

		paramNames := extractParamNames(pattern)
		regex := regexp.MustCompile(pattern)

		r.routes = append(r.routes, Route{
			Methods:    methods,
			Regex:      regex,
			Handler:    route.Handler,
			ParamNames: paramNames,
		})
	}
	return r
}

func extractParamNames(pattern string) []string {
	var names []string
	re := regexp.MustCompile(`\((?:\?P<([^>]+)>)?[^)]+\)`)
	matches := re.FindAllStringSubmatch(pattern, -1)
	for i, m := range matches {
		if len(m) > 1 && m[1] != "" {
			names = append(names, m[1])
		} else {
			names = append(names, string(rune('0'+i+1)))
		}
	}
	return names
}

func (r *Router) Match(method, path string) (string, map[string]string, bool) {
	for _, route := range r.routes {
		methodMatch := false
		for _, m := range route.Methods {
			if strings.EqualFold(strings.TrimSpace(m), method) {
				methodMatch = true
				break
			}
		}
		if !methodMatch {
			continue
		}

		matches := route.Regex.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		params := make(map[string]string)
		for i, name := range route.ParamNames {
			if i+1 < len(matches) {
				params[name] = matches[i+1]
			}
		}

		parts := strings.Split(route.Handler, ":")
		if len(parts) == 2 {
			params["_controller"] = parts[0]
			params["_action"] = parts[1]
		}

		return route.Handler, params, true
	}
	return "", nil, false
}

type contextKey string

const (
	paramsKey contextKey = "router_params"
)

func GetParams(ctx context.Context) map[string]string {
	if params, ok := ctx.Value(paramsKey).(map[string]string); ok {
		return params
	}
	return nil
}

func (r *Router) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		method := req.Method

		if method == http.MethodPost && req.Header.Get("X-HTTP-Method-Override") != "" {
			method = req.Header.Get("X-HTTP-Method-Override")
		}

		handler, params, found := r.Match(method, path)
		if !found {
			http.NotFound(w, req)
			return
		}

		ctx := context.WithValue(req.Context(), paramsKey, params)
		ctx = context.WithValue(ctx, "handler", handler)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
