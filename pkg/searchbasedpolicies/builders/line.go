package builders

import (
	"context"
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/registries/types"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
)

type dockerFileLineFieldQueryBuilder struct {
}

func (c *dockerFileLineFieldQueryBuilder) Name() string {
	return "query builder for docker file lines"
}

func (c *dockerFileLineFieldQueryBuilder) Query(fields *storage.PolicyFields, optionsMap map[search.FieldLabel]*v1.SearchField) (q *v1.Query, v searchbasedpolicies.ViolationPrinter, err error) {
	lineRule := fields.GetLineRule()
	if lineRule == nil {
		return
	}

	instSearchField, err := getSearchField(search.DockerfileInstructionKeyword, optionsMap)
	if err != nil {
		err = errors.Wrapf(err, "%s", c.Name())
	}

	if _, ok := types.DockerfileInstructionSet[lineRule.GetInstruction()]; !ok {
		err = fmt.Errorf("%v is not a valid dockerfile instruction", lineRule.GetInstruction())
		return
	}

	// If no value, then just query for the instruction.
	if lineRule.GetValue() == "" {
		q = search.NewQueryBuilder().AddStringsHighlighted(search.DockerfileInstructionKeyword, lineRule.GetInstruction()).ProtoQuery()
		v = violationPrinterForField(instSearchField.GetFieldPath(), func(match string) string {
			return fmt.Sprintf("Dockerfile instruction %s found", match)
		})
		return
	}

	_, err = regexp.Compile(lineRule.GetValue())
	if err != nil {
		err = errors.Wrapf(err, "invalid line regex %+v", lineRule)
		return
	}

	valSearchField, err := getSearchField(search.DockerfileInstructionValue, optionsMap)
	if err != nil {
		err = errors.Wrapf(err, "%s", c.Name())
	}

	q = search.NewQueryBuilder().AddLinkedFieldsHighlighted(
		[]search.FieldLabel{search.DockerfileInstructionKeyword, search.DockerfileInstructionValue},
		[]string{lineRule.GetInstruction(), search.RegexQueryString(lineRule.GetValue())}).ProtoQuery()

	v = func(_ context.Context, result search.Result) searchbasedpolicies.Violations {
		instMatches := result.Matches[instSearchField.GetFieldPath()]
		valMatches := result.Matches[valSearchField.GetFieldPath()]
		if len(instMatches) == 0 || len(valMatches) == 0 {
			return searchbasedpolicies.Violations{}
		}
		violations := searchbasedpolicies.Violations{
			AlertViolations: make([]*storage.Alert_Violation, 0, len(instMatches)),
		}
		for i, instMatch := range instMatches {
			// This should not happen if search works as expected.
			if i >= len(valMatches) {
				log.Errorf("Matching Dockerfile line rule: %+v, "+
					"instMatches %+v and valMatches %+v not of equal length", lineRule, instMatches, valMatches)
				break
			}
			violations.AlertViolations = append(violations.AlertViolations, &storage.Alert_Violation{
				Message: fmt.Sprintf("Dockerfile Line '%s %s' matches the rule %s %s",
					instMatch, valMatches[i], lineRule.GetInstruction(), lineRule.GetValue()),
			})
		}
		return violations
	}
	return
}

// NewDockerFileLineQueryBuilder returns a query builder that constructs a
// Dockerfile line query.
func NewDockerFileLineQueryBuilder() searchbasedpolicies.PolicyQueryBuilder {
	return &dockerFileLineFieldQueryBuilder{}
}
