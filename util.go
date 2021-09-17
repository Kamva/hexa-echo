package hecho

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/kamva/hexa"
)

// uuidGenerator generate new UUID
func uuidGenerator() string {
	return uuid.New().String()
}

func captureTokens(pattern *regexp.Regexp, input string) *strings.Replacer {
	groups := pattern.FindAllStringSubmatch(input, -1)
	if groups == nil {
		return nil
	}
	values := groups[0][1:]
	replace := make([]string, 2*len(values))
	for i, v := range values {
		j := 2 * i
		replace[j] = "$" + strconv.Itoa(i+1)
		replace[j+1] = v
	}
	return strings.NewReplacer(replace...)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// isInternalErr returns true, if error is an real app error not an hexa Reply or
// Error without server error status code.
// Note:if we decided to do not return hexa.Reply or non-internal errors(e.g, error 4xx,...) as an error
// return param in in echo handler, we can remove this function, because in that situation the
// Echo handler's error return param surely is an app error.
func isInternalErr(err error) bool {
	if err == nil {
		return false
	}

	for err != nil {
		if _, ok := err.(hexa.Reply); ok {
			return false
		}

		if hexaErr, ok := err.(hexa.Error); ok {
			if hexaErr.HTTPStatus() >= 50 { // its an internal error.
				return true
			}
			return false
		}

		err = errors.Unwrap(err)
	}

	// its not hexa Error or Reply, so its an app error.
	return true
}
