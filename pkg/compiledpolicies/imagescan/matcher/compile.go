package matcher

import (
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/logging"
)

var log = logging.LoggerForModule()

// compiler is a functino that turns a policy into a matcher.
type compiler func(*v1.Policy) (Matcher, error)

// compilers are all of the different Matchers registered.
var compilers []compiler

// Compile creates a new image scan policy matcher.
func Compile(policy *v1.Policy) (Matcher, error) {
	var matcher Matcher
	for _, compiler := range compilers {
		matcherFunction, err := compiler(policy)
		if err != nil {
			return nil, err
		}
		matcher = matcher.MustAlsoViolate(matcherFunction)
	}
	return matcher, nil
}
