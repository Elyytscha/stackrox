package predicate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/search"
)

// MergeResults merges predicate result into a single result
// If the results slice passed is of length 0, then it will return a nil result
func MergeResults(results ...*search.Result) *search.Result {
	res := search.NewResult()
	for _, r := range results {
		if r == nil {
			continue
		}
		for k, v := range r.Matches {
			res.Matches[k] = append(res.Matches[k], v...)
		}
	}
	return res
}

// Predicate represents a method that accesses data in some interface.
type Predicate func(instance interface{}) (*search.Result, bool)

// Matches is a utility function to just return a bool
func (p Predicate) Matches(instance interface{}) bool {
	_, matches := p(instance)
	return matches
}

// Factory object stores the specs for each when walking the query.
type Factory struct {
	searchFields FieldMap
	searchPaths  wrappedOptionsMap
	exampleObj   interface{}
}

// NewFactory returns a new predicat factory for the type of the given object.
func NewFactory(prefix string, obj interface{}) Factory {
	return Factory{
		searchFields: mapSearchTagsToFieldPaths(obj),
		searchPaths: wrappedOptionsMap{
			optionsMap: search.Walk(v1.SearchCategory(-1), prefix, obj),
			prefix:     fmt.Sprintf("%s.", prefix),
		},
		exampleObj: obj,
	}
}

// GeneratePredicate creates a predicate for the Predicate factories type that returns whether or not the input
// instance matches the query.
func (tb Factory) GeneratePredicate(query *v1.Query) (Predicate, error) {
	ip, err := tb.generatePredicateInternal(query)
	if err != nil {
		return nil, err
	}
	return wrapInternal(ip), nil
}

type internalPredicate func(reflect.Value) (*search.Result, bool)

func wrapInternal(ip internalPredicate) Predicate {
	return func(in interface{}) (*search.Result, bool) {
		val := reflect.ValueOf(in)
		return ip(val)
	}
}

func (tb Factory) generatePredicateInternal(query *v1.Query) (internalPredicate, error) {
	if query.GetQuery() == nil {
		return alwaysTrue, nil
	}
	switch query.GetQuery().(type) {
	case *v1.Query_Disjunction:
		return tb.or(query.GetDisjunction())
	case *v1.Query_Conjunction:
		return tb.and(query.GetConjunction())
	case *v1.Query_BooleanQuery:
		return tb.boolean(query.GetBooleanQuery())
	case *v1.Query_BaseQuery:
		return tb.base(query.GetBaseQuery())
	default:
		return nil, fmt.Errorf("unrecognized query type: %T", query.GetQuery())
	}
}

func (tb Factory) or(q *v1.DisjunctionQuery) (internalPredicate, error) {
	ret := make([]internalPredicate, 0, len(q.GetQueries()))
	for _, dis := range q.GetQueries() {
		next, err := tb.generatePredicateInternal(dis)
		if err != nil {
			return nil, err
		}
		ret = append(ret, next)
	}

	return orOf(ret...), nil
}

func (tb Factory) and(q *v1.ConjunctionQuery) (internalPredicate, error) {
	ret := make([]internalPredicate, 0, len(q.GetQueries()))
	for _, dis := range q.GetQueries() {
		next, err := tb.generatePredicateInternal(dis)
		if err != nil {
			return nil, err
		}
		ret = append(ret, next)
	}

	return andOf(ret...), nil
}

func (tb Factory) boolean(q *v1.BooleanQuery) (internalPredicate, error) {
	must, err := tb.and(q.GetMust())
	if err != nil {
		return nil, err
	}

	mustNot, err := tb.or(q.GetMustNot())
	if err != nil {
		return nil, err
	}

	return func(instance reflect.Value) (*search.Result, bool) {
		mustRes, mustMatch := must(instance)
		if !mustMatch {
			return nil, false
		}
		mustNotRes, mustNotMatch := mustNot(instance)
		if mustNotMatch {
			return nil, false
		}
		return MergeResults(mustRes, mustNotRes), true
	}, nil
}

func (tb Factory) base(q *v1.BaseQuery) (internalPredicate, error) {
	switch q.GetQuery().(type) {
	case *v1.BaseQuery_DocIdQuery:
		return tb.docID(q.GetDocIdQuery())
	case *v1.BaseQuery_MatchNoneQuery:
		return tb.matchNone(q.GetMatchNoneQuery())
	case *v1.BaseQuery_MatchFieldQuery:
		return tb.match(q.GetMatchFieldQuery())
	case *v1.BaseQuery_MatchLinkedFieldsQuery:
		return tb.matchLinked(q.GetMatchLinkedFieldsQuery())
	default:
		return nil, fmt.Errorf("cannot handle base query of type %T", q.GetQuery())
	}
}

func (tb Factory) docID(q *v1.DocIDQuery) (internalPredicate, error) {
	return nil, errors.New("query predicates do not support DocID query types as DocIDs only exist in the index")
}

func (tb Factory) matchNone(q *v1.MatchNoneQuery) (internalPredicate, error) {
	return alwaysFalse, nil
}

func (tb Factory) match(q *v1.MatchFieldQuery) (internalPredicate, error) {
	fp := tb.searchFields.Get(strings.ToLower(q.GetField()))
	if fp == nil {
		return alwaysTrue, nil
	}
	sp, ok := tb.searchPaths.Get(q.GetField())
	if !ok {
		return tb.createPredicate("", fp, q.GetValue())
	}
	return tb.createPredicate(sp.GetFieldPath(), fp, q.GetValue())
}

func (tb Factory) matchLinked(q *v1.MatchLinkedFieldsQuery) (internalPredicate, error) {
	// Find the longest common path with all of the linked fields.
	var commonPath FieldPath
	for _, fieldQuery := range q.GetQuery() {
		path := tb.searchFields.Get(fieldQuery.GetField())
		if path == nil {
			return alwaysTrue, nil
		}
		if commonPath == nil {
			commonPath = path[:len(path)-1]
		} else {
			for idx, field := range path {
				if idx > len(commonPath)-1 {
					break
				}
				if commonPath[idx].Name != field.Name {
					commonPath = commonPath[:idx]
					break
				}
			}
		}
	}

	// Produce a predicate for each of the fields. Use the non common path.
	var preds []internalPredicate
	for _, fieldQuery := range q.GetQuery() {
		path := tb.searchFields.Get(fieldQuery.GetField())
		if path == nil {
			return alwaysTrue, nil
		}
		var fieldPath string
		searchField, ok := tb.searchPaths.Get(fieldQuery.GetField())
		if ok {
			fieldPath = searchField.GetFieldPath()
		}

		pred, err := tb.createPredicate(fieldPath, path[len(commonPath):], fieldQuery.GetValue())
		if err != nil {
			return nil, err
		}
		preds = append(preds, pred)
	}

	// Package all the of predicates as an AND on the common path.
	linked, err := createLinkedNestedPredicate(commonPath[len(commonPath)-1].Type, preds...)
	if err != nil {
		return nil, err
	}
	return createPathPredicate(reflect.TypeOf(tb.exampleObj), commonPath, linked)
}

func (tb Factory) createPredicate(fullPath string, path FieldPath, value string) (internalPredicate, error) {
	// Create the predicate for the search field value.
	pred, err := createBasePredicate(fullPath, path[len(path)-1].Type, value)
	if err != nil {
		return nil, err
	}

	// Create a wrapper predicate which traces the field path down to the value.
	pred, err = createPathPredicate(reflect.TypeOf(tb.exampleObj), path, pred)
	if err != nil {
		return nil, err
	}
	return pred, nil
}

// Combinatorial helpers.
/////////////////////////

func orOf(preds ...internalPredicate) internalPredicate {
	return func(instance reflect.Value) (*search.Result, bool) {
		var results []*search.Result
		for _, pred := range preds {
			if res, match := pred(instance); match {
				results = append(results, res)
			}
		}
		if len(results) > 0 {
			return MergeResults(results...), true
		}
		return nil, false
	}
}

func andOf(preds ...internalPredicate) internalPredicate {
	return func(instance reflect.Value) (*search.Result, bool) {
		var results []*search.Result
		for _, pred := range preds {
			result, ok := pred(instance)
			if !ok {
				return nil, false
			}
			results = append(results, result)
		}
		return MergeResults(results...), true
	}
}

func alwaysFalse(_ reflect.Value) (*search.Result, bool) {
	return nil, false
}

func alwaysTrue(_ reflect.Value) (*search.Result, bool) {
	return &search.Result{}, true
}

// Recursive predicate manufacturing from the input field path.
///////////////////////////////////////////////////////////////

func createPathPredicate(parentType reflect.Type, path FieldPath, pred internalPredicate) (internalPredicate, error) {
	if len(path) == 0 {
		return pred, nil
	}

	// If not, recursively go down to the base.
	child, err := createPathPredicate(path[0].Type, path[1:], pred)
	if err != nil {
		return nil, err
	}

	// Wrap the predicate in field access.
	return createNestedPredicate(parentType, path[0], child)
}

func createNestedPredicate(parentType reflect.Type, field reflect.StructField, pred internalPredicate) (internalPredicate, error) {
	switch parentType.Kind() {
	case reflect.Array:
		return createSliceNestedPredicate(parentType, field, pred)
	case reflect.Slice:
		return createSliceNestedPredicate(parentType, field, pred)
	case reflect.Ptr:
		return createPtrNestedPredicate(parentType, field, pred)
	case reflect.Map:
		return createMapNestedPredicate(parentType, field, pred)
	case reflect.Struct:
		return createStructFieldNestedPredicate(field, pred), nil
	case reflect.Interface:
		return createInterfaceFieldNestedPredicate(field, pred), nil
	default:
		return alwaysFalse, fmt.Errorf("cannot follow: %+v", field)
	}
}

// Complex type predicates.
////////////////////////////

func createSliceNestedPredicate(parentType reflect.Type, field reflect.StructField, pred internalPredicate) (internalPredicate, error) {
	nested, err := createNestedPredicate(parentType.Elem(), field, pred)
	if err != nil {
		return nil, err
	}
	return func(instance reflect.Value) (*search.Result, bool) {
		if instance.IsNil() || instance.IsZero() {
			return nil, false
		}
		var results []*search.Result
		length := instance.Len()
		for i := 0; i < length; i++ {
			idx := instance.Index(i)
			if res, matches := nested(idx); matches {
				results = append(results, res)
			}
		}
		if len(results) > 0 {
			return MergeResults(results...), true
		}
		return nil, false
	}, nil
}

func createPtrNestedPredicate(parentType reflect.Type, field reflect.StructField, pred internalPredicate) (internalPredicate, error) {
	nested, err := createNestedPredicate(parentType.Elem(), field, pred)
	if err != nil {
		return nil, err
	}
	return func(instance reflect.Value) (*search.Result, bool) {
		if instance.IsNil() { // Need to special handle pointers to nil. Good ole typed nils.
			return nil, false
		}
		elem := instance.Elem()
		if res, matches := nested(elem); matches {
			return res, true
		}
		return nil, false
	}, nil
}

func createMapNestedPredicate(parentType reflect.Type, field reflect.StructField, pred internalPredicate) (internalPredicate, error) {
	nestedKey, err := createNestedPredicate(parentType.Key(), field, pred)
	if err != nil {
		return nil, err
	}
	parentTypeElem := parentType.Elem()
	nestedElem, err := createNestedPredicate(parentTypeElem, field, pred)
	if err != nil {
		return nil, err
	}
	return func(instance reflect.Value) (*search.Result, bool) {
		if instance.IsNil() || instance.IsZero() {
			return nil, false
		}
		for _, key := range instance.MapKeys() {
			valueAt := instance.MapIndex(key)
			res, match := nestedKey(key)
			if match {
				return res, true
			}
			res, match = nestedElem(valueAt)
			if match {
				return res, true
			}
		}
		return nil, false
	}, nil
}

func nilCheck(f reflect.Value) bool {
	switch f.Kind() {
	case reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return f.IsNil()
	}
	return false
}

func createStructFieldNestedPredicate(field reflect.StructField, pred internalPredicate) internalPredicate {
	return func(instance reflect.Value) (*search.Result, bool) {
		nextValue := instance.FieldByName(field.Name)
		if nilCheck(nextValue) {
			return nil, false
		}
		return pred(nextValue)
	}
}

func createInterfaceFieldNestedPredicate(field reflect.StructField, pred internalPredicate) internalPredicate {
	return func(instance reflect.Value) (*search.Result, bool) {
		if instance.IsNil() || instance.IsZero() {
			return nil, false
		}
		concrete := instance.Elem()
		if concrete.Type().Kind() == reflect.Ptr {
			concrete = concrete.Elem()
		}
		if concrete.Type().Kind() != reflect.Struct {
			return nil, false
		}
		nextValue := concrete.FieldByName(field.Name)
		if nextValue.IsZero() {
			return nil, false // Field either does not exist, or wasn't populated.
		}
		return pred(nextValue)
	}
}
