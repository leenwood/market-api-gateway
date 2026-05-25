package proxy

import (
	"regexp"
	"strings"

	"market-api-gateway/internal/core/domain"
)

type compiledRoute struct {
	route      domain.Route
	re         *regexp.Regexp
	paramNames []string
}

// Router matches requests against the static routing table.
type Router struct {
	routes []compiledRoute
}

func NewRouter() *Router {
	r := &Router{}
	for _, route := range domain.Routes {
		r.routes = append(r.routes, compile(route))
	}
	return r
}

// GetRoute returns the matching route and extracted path params, or false if no match.
func (r *Router) GetRoute(path, method string) (*domain.Route, map[string]string, bool) {
	for i := range r.routes {
		cr := &r.routes[i]
		if !strings.EqualFold(cr.route.Method, method) {
			continue
		}
		m := cr.re.FindStringSubmatch(path)
		if m == nil {
			continue
		}
		params := make(map[string]string, len(cr.paramNames))
		for j, name := range cr.paramNames {
			if j+1 < len(m) {
				params[name] = m[j+1]
			}
		}
		route := cr.route
		return &route, params, true
	}
	return nil, nil, false
}

// compile converts a route pattern like /api/v1/cart/:userID to a regexp.
func compile(route domain.Route) compiledRoute {
	pattern := route.Pattern
	var paramNames []string

	// Wildcard suffix (e.g. /api/v1/analytics/)
	if strings.HasSuffix(pattern, "/") {
		re := regexp.MustCompile("^" + regexp.QuoteMeta(strings.TrimSuffix(pattern, "/")) + "(/.*)?$")
		return compiledRoute{route: route, re: re, paramNames: nil}
	}

	// Replace :param segments with named capture groups.
	parts := strings.Split(pattern, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			name := part[1:]
			paramNames = append(paramNames, name)
			parts[i] = "([^/]+)"
		} else {
			parts[i] = regexp.QuoteMeta(part)
		}
	}
	re := regexp.MustCompile("^" + strings.Join(parts, "/") + "$")
	return compiledRoute{route: route, re: re, paramNames: paramNames}
}
