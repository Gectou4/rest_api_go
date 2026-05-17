package router

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/gectou4/rest_api_go/internal/config"
)

type compiledRoute struct {
	Methods    []string
	Regex      *regexp.Regexp
	ParamNames []string
	Controller string
	Action     string
}

type Router struct {
	routes []compiledRoute
}

func NewRouter(routes []config.Route) *Router {
	r := &Router{}
	for _, route := range routes {
		methods := strings.Split(route.Method, "|")
		paramNames := extractParamNames(route.Pattern)
		regex := regexp.MustCompile(route.Pattern)

		parts := strings.Split(route.Handler, ":")
		controller := ""
		action := ""
		if len(parts) == 2 {
			controller = parts[0]
			action = parts[1]
		}

		r.routes = append(r.routes, compiledRoute{
			Methods:    methods,
			Regex:      regex,
			ParamNames: paramNames,
			Controller: controller,
			Action:     action,
		})
	}
	return r
}

func extractParamNames(pattern string) []string {
	var names []string
	re := regexp.MustCompile(`\([^)]+\)`)
	matches := re.FindAllStringIndex(pattern, -1)
	for i := range matches {
		names = append(names, string(rune('0'+i+1)))
	}
	return names
}

func (r *Router) Match(method, path string) (string, string, map[string]string, bool) {
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

		return route.Controller, route.Action, params, true
	}
	return "", "", nil, false
}

type contextKey string

const (
	paramsKey     contextKey = "router_params"
	controllerKey contextKey = "router_controller"
	actionKey     contextKey = "router_action"
)

func (r *Router) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		method := req.Method

		if method == http.MethodPost && req.Header.Get("X-HTTP-Method-Override") != "" {
			method = req.Header.Get("X-HTTP-Method-Override")
		}

		controller, action, params, found := r.Match(method, path)
		if !found {
			http.NotFound(w, req)
			return
		}

		ctx := context.WithValue(req.Context(), paramsKey, params)
		ctx = context.WithValue(ctx, controllerKey, controller)
		ctx = context.WithValue(ctx, actionKey, action)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func GetParams(ctx context.Context) map[string]string {
	if params, ok := ctx.Value(paramsKey).(map[string]string); ok {
		return params
	}
	return nil
}

func GetController(ctx context.Context) string {
	if c, ok := ctx.Value(controllerKey).(string); ok {
		return c
	}
	return ""
}

func GetAction(ctx context.Context) string {
	if a, ok := ctx.Value(actionKey).(string); ok {
		return a
	}
	return ""
}
