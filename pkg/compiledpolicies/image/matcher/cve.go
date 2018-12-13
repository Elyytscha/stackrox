package matcher

import (
	"fmt"
	"regexp"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/compiledpolicies/utils"
)

func init() {
	compilers = append(compilers, newCVEMatcher)
}

func newCVEMatcher(policy *storage.Policy) (Matcher, error) {
	cve := policy.GetFields().GetCve()
	if cve == "" {
		return nil, nil
	}

	cveRegex, err := utils.CompileStringRegex(cve)
	if err != nil {
		return nil, err
	}
	matcher := &cveMatcherImpl{cveRegex}
	return matcher.match, nil
}

type cveMatcherImpl struct {
	cveRegex *regexp.Regexp
}

func (p *cveMatcherImpl) match(image *storage.Image) (violations []*storage.Alert_Violation) {
	for _, component := range image.GetScan().GetComponents() {
		for _, vuln := range component.GetVulns() {
			if p.cveRegex.MatchString(vuln.GetCve()) {
				violations = append(violations, &storage.Alert_Violation{
					Message: fmt.Sprintf("'%v' in Component '%v' matches the regex '%+v'", vuln.GetCve(), component.GetName(), p.cveRegex),
					Link:    vuln.GetLink(),
				})
			}
		}
	}
	return
}
