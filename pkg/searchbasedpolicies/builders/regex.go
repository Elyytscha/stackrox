package builders

import (
	"context"
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/searchbasedpolicies"
)

// A RegexField represents the information required to match a regex-based policy field.
type RegexField struct {
	FieldLabel         search.FieldLabel
	FieldHumanName     string
	AllowSubstrings    bool
	RetrieveFieldValue func(*storage.PolicyFields) string
}

// A RegexQueryBuilder builds a policy query builder from a set of linked regex fields.
type RegexQueryBuilder struct {
	RegexFields []RegexField
}

// Name implements the PolicyQueryBuilder interface.
func (r RegexQueryBuilder) Name() string {
	return fmt.Sprintf("query builder for %+v", r.RegexFields)
}

// Query implements the PolicyQueryBuilder interface.
func (r RegexQueryBuilder) Query(fields *storage.PolicyFields, optionsMap map[search.FieldLabel]*v1.SearchField) (q *v1.Query, v searchbasedpolicies.ViolationPrinter, err error) {
	type presentFieldInfo struct {
		policyVal      string
		searchField    *v1.SearchField
		fieldHumanName string
	}
	var presentFieldValues []presentFieldInfo
	var fieldLabels []search.FieldLabel
	var fieldValues []string

	for _, field := range r.RegexFields {
		policyVal := field.RetrieveFieldValue(fields)
		if policyVal == "" {
			continue
		}
		var searchField *v1.SearchField
		searchField, err = getSearchField(field.FieldLabel, optionsMap)
		if err != nil {
			err = errors.Wrapf(err, "%s", r.Name())
			return
		}

		actualQueriedVal := policyVal
		// If it's a string field, then make a regex query.
		if searchField.GetType() == v1.SearchDataType_SEARCH_STRING {
			if field.AllowSubstrings {
				actualQueriedVal = fmt.Sprintf(".*%s.*", policyVal)
			}
			// Make sure the regex compiles (Bleve will just fail silently.)
			_, err = regexp.Compile(actualQueriedVal)
			if err != nil {
				err = errors.Wrapf(err, "'%s' is an invalid regex", actualQueriedVal)
				return
			}

			actualQueriedVal = search.RegexQueryString(actualQueriedVal)
		}

		presentFieldValues = append(presentFieldValues, presentFieldInfo{
			policyVal:      policyVal,
			searchField:    searchField,
			fieldHumanName: field.FieldHumanName,
		})
		fieldLabels = append(fieldLabels, field.FieldLabel)
		fieldValues = append(fieldValues, actualQueriedVal)
	}
	if len(presentFieldValues) == 0 {
		return
	}

	if len(presentFieldValues) == 1 {
		q = search.NewQueryBuilder().AddStringsHighlighted(fieldLabels[0], fieldValues[0]).ProtoQuery()
	} else {
		q = search.NewQueryBuilder().AddLinkedFieldsHighlighted(fieldLabels, fieldValues).ProtoQuery()
	}

	v = func(_ context.Context, result search.Result) searchbasedpolicies.Violations {
		violations := searchbasedpolicies.Violations{}
		for _, presentFieldValue := range presentFieldValues {
			// Need to fall back to paths from top level image components and CVEs
			matches := result.Matches[presentFieldValue.searchField.GetFieldPath()]
			if len(matches) == 0 {
				matches = result.Matches[swapFieldPaths(presentFieldValue.searchField.GetFieldPath())]
			}
			for _, match := range matches {
				violations.AlertViolations = append(violations.AlertViolations, &storage.Alert_Violation{
					Message: fmt.Sprintf("%s '%s' matched %s", presentFieldValue.fieldHumanName, match, presentFieldValue.policyVal),
				})
			}
		}
		return violations
	}
	return
}
