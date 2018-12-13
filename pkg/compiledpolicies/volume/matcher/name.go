package matcher

import (
	"fmt"
	"regexp"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/compiledpolicies/utils"
)

func init() {
	compilers = append(compilers, newNameMatcher)
}

func newNameMatcher(policy *storage.Policy) (Matcher, error) {
	name := policy.GetFields().GetVolumePolicy().GetName()
	if name == "" {
		return nil, nil
	}

	nameRegex, err := utils.CompileStringRegex(name)
	if err != nil {
		return nil, err
	}
	matcher := &nameMatcherImpl{nameRegex}
	return matcher.match, nil
}

type nameMatcherImpl struct {
	nameRegex *regexp.Regexp
}

func (p *nameMatcherImpl) match(volume *storage.Volume) []*storage.Alert_Violation {
	var violations []*storage.Alert_Violation
	if p.nameRegex.MatchString(volume.GetName()) {
		v := &storage.Alert_Violation{
			Message: fmt.Sprintf("Volume name matched: %s", p.nameRegex),
		}
		violations = append(violations, v)
	}
	return violations
}
