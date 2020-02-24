package kecho

import "strings"

// RouteName erturns the route name by provided string values.
// e.g RouteName("auth","register") returns "auth::register"
func RouteName(a ...string) string {
	return strings.Join(a, "::")
}
