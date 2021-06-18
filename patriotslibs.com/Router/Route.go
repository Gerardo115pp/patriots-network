package patriot_router

import (
	"regexp"
)

type Route struct {
	route_str   string
	route_regex *regexp.Regexp
	is_string   bool
}

func (self *Route) String() string {
	if self.is_string {
		return self.route_str
	} else {
		return self.route_regex.String()
	}
}

func (self *Route) Match(candidate string) (is_match bool) {
	if self.is_string {
		is_match = self.route_str == candidate
	} else {
		is_match = self.route_regex.MatchString(candidate)
	}
	return
}

func NewRoute(route_exp string, is_string bool) *Route {
	var new_route *Route = new(Route)
	new_route.is_string = is_string
	new_route.route_str = route_exp
	if is_string {
		new_route.route_regex = regexp.MustCompile(".^")
	} else {
		new_route.route_regex = regexp.MustCompile(route_exp)
	}
	return new_route
}
