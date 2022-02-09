package regocompile

import (
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator"
	"github.com/stackrox/rox/pkg/booleanpolicy/evaluator/pathutil"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/timeutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TopLevel struct {
	ValA        string `search:"TopLevelA"`
	NestedSlice []Nested
	Base        Base
}

type Base struct {
	// These exist for testing base types.
	ValBaseSlice     []string          `search:"BaseSlice"`
	ValBasePtr       *string           `search:"BasePtr"`
	ValBaseBool      bool              `search:"BaseBool"`
	ValBaseTS        *types.Timestamp  `search:"BaseTS"`
	ValBaseInt       int               `search:"BaseInt"`
	ValBaseUint      uint              `search:"BaseUint"`
	ValBaseFloat     float64           `search:"BaseFloat"`
	ValBaseEnum      storage.Access    `search:"BaseEnum"`
	ValBaseMap       map[string]string `search:"BaseMap"`
	ValBaseStructPtr *random           `search:"BaseStructPtr"`
}

type random struct {
	Val string
}

type Nested struct {
	NestedValA        string `search:"A"`
	NestedValB        string `search:"B"`
	SecondNestedSlice []*SecondNested
}

type SecondNested struct {
	SecondNestedValA string `search:"SecondA"`
	SecondNestedValB string `search:"SecondB"`
}

// Bare versions of top level and nested for the sake of testing augmentation.
// See augmentedCompilerInstance for how these objects are augmented to appear
// like they are their non-bare versions.
type TopLevelBare struct {
	ValA string `search:"TopLevelA"`
	// This struct will be replaced with the other "Base" by the augmentation.
	// We keep it here to make sure the augmentation correctly supplants this object.
	Base struct {
		IrrelevantBaseVal int `search:"BaseInt"`
	}
}

type NestedBare struct {
	NestedValA string `search:"A"`
	NestedValB string `search:"B"`
}

var (
	compilerInstance = MustCreateRegoCompiler(pathutil.NewAugmentedObjMeta((*TopLevel)(nil)))

	augmentedCompilerInstance = MustCreateRegoCompiler(pathutil.NewAugmentedObjMeta((*TopLevelBare)(nil)).
					AddPlainObjectAt([]string{"Base"}, Base{}).
					AddAugmentedObjectAt(
			[]string{"NestedSlice"},
			pathutil.NewAugmentedObjMeta(([]NestedBare)(nil)).AddPlainObjectAt([]string{"SecondNestedSlice"}, ([]*SecondNested)(nil)),
		),
	)

	ts2020Apr01 = protoconv.MustConvertTimeToTimestamp(timeutil.MustParse(time.RFC3339, "2020-04-01T00:00:00Z"))
)

type testCase struct {
	desc           string
	q              *query.Query
	obj            *TopLevel
	expectedResult *evaluator.Result
}

func assertResultsAsExpected(t *testing.T, c testCase, actualRes *evaluator.Result, actualMatched bool) {
	assert.Equal(t, c.expectedResult != nil, actualMatched)
	if c.expectedResult != nil {
		require.NotNil(t, actualRes)
		assert.ElementsMatch(t, c.expectedResult.Matches, actualRes.Matches)
	}
}

func resultWithSingleMatch(fieldName string, values ...string) *evaluator.Result {
	return &evaluator.Result{Matches: []map[string][]string{{fieldName: values}}}
}

func addContextMatchToResult(r *evaluator.Result, fieldName string, values ...string) *evaluator.Result {
	for _, match := range r.Matches {
		match[fieldName] = values
	}
	r.Matches = append(r.Matches, map[string][]string{fieldName: values})
	return r
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, testCase := range testCases {
		c := testCase
		t.Run(c.desc, func(t *testing.T) {
			t.Run("on fully hydrated object", func(t *testing.T) {
				evaluator, err := compilerInstance.CompileRegoBasedEvaluator(c.q)
				require.NoError(t, err)
				res, matched := evaluator.Evaluate(pathutil.NewAugmentedObj(c.obj).Value())
				assertResultsAsExpected(t, c, res, matched)
			})
			t.Run("on augmented object", func(t *testing.T) {
				evaluator, err := augmentedCompilerInstance.CompileRegoBasedEvaluator(c.q)
				require.NoError(t, err)
				topLevelBare := &TopLevelBare{
					ValA: c.obj.ValA,
				}
				base := c.obj.Base
				nestedBare := make([]NestedBare, 0, len(c.obj.NestedSlice))
				for _, elem := range c.obj.NestedSlice {
					nestedBare = append(nestedBare, NestedBare{NestedValA: elem.NestedValA, NestedValB: elem.NestedValB})
				}

				nestedAugmentedObj := pathutil.NewAugmentedObj(nestedBare)
				for i, elem := range c.obj.NestedSlice {
					require.NoError(t, nestedAugmentedObj.AddPlainObjAt(elem.SecondNestedSlice, pathutil.IndexStep(i), pathutil.FieldStep("SecondNestedSlice")))
				}
				topLevelAugmentedObj := pathutil.NewAugmentedObj(topLevelBare)
				require.NoError(t, topLevelAugmentedObj.AddPlainObjAt(base, pathutil.FieldStep("Base")))
				require.NoError(t, topLevelAugmentedObj.AddAugmentedObjAt(nestedAugmentedObj, pathutil.FieldStep("NestedSlice")))
				res, matched := evaluator.Evaluate(topLevelAugmentedObj.Value())
				assertResultsAsExpected(t, c, res, matched)
			})
		})
	}
}

func TestSimpleBase(t *testing.T) {
	qTopLevelAHappy := query.SimpleMatchFieldQuery("TopLevelA", "happy")
	qNestedAHappy := query.SimpleMatchFieldQuery("A", "happy")
	qSecondNestedAHappy := query.SimpleMatchFieldQuery("SecondA", "r/.*ppy")

	runTestCases(t, []testCase{
		{
			desc: "simple one for top level, doesn't pass",
			q:    qTopLevelAHappy,
			obj: &TopLevel{
				ValA: "whatever",
				NestedSlice: []Nested{
					{NestedValA: "blah"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happy"},
					}},
				},
			},
		},
		{
			desc: "simple one for top level, passes",
			q:    qTopLevelAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "blah"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happy"},
					}},
				},
			},
			expectedResult: resultWithSingleMatch("TopLevelA", "happy"),
		},
		{
			desc: "simple one for first level nested, doesn't pass",
			q:    qNestedAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "blah"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happy"},
					}},
				},
			},
		},
		{
			desc: "simple one for first level nested, passes",
			q:    qNestedAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "happy"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happiest"},
					}},
				},
			},
			expectedResult: resultWithSingleMatch("A", "happy"),
		},
		{
			desc: "simple one for second level nested, doesn't pass",
			q:    qSecondNestedAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "happy"},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happiest"},
					}},
				},
			},
		},
		{
			desc: "simple one for second level nested, passes",
			q:    qSecondNestedAHappy,
			obj: &TopLevel{
				ValA: "happy",
				NestedSlice: []Nested{
					{NestedValA: "happy", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "blah"},
						{SecondNestedValA: "blaappy"},
					}},
					{NestedValA: "something else", SecondNestedSlice: []*SecondNested{
						{SecondNestedValA: "happy"},
					}},
				},
			},
			expectedResult: &evaluator.Result{
				Matches: []map[string][]string{
					{"SecondA": {"happy"}},
					{"SecondA": {"blaappy"}},
				},
			},
		},
	})
}

func TestLinked(t *testing.T) {
	runTestCases(t, []testCase{
		{
			desc: "linked, first level of nesting, should match",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A1"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
			expectedResult: &evaluator.Result{
				Matches: []map[string][]string{
					{"A": {"A1"}, "B": {"B1"}},
				},
			},
		},
		{
			desc: "linked, first level of nesting, should not match",
			obj: &TopLevel{
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "A", Values: []string{"A0"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
		},
		{
			desc: "linked, multilevel, should match",
			obj: &TopLevel{
				ValA: "TopLevelValA",
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "TopLevelA", Values: []string{"TopLevelValA"}},
					{Field: "A", Values: []string{"A1"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
			expectedResult: &evaluator.Result{
				Matches: []map[string][]string{
					{"TopLevelA": {"TopLevelValA"}, "A": {"A1"}, "B": {"B1"}},
				},
			},
		},
		{
			desc: "linked, multilevel, should match (group test)",
			obj: &TopLevel{
				ValA: "TopLevelValA",
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
					{NestedValA: "A2", NestedValB: "B2"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "TopLevelA", Values: []string{"TopLevelValA"}},
					{Field: "A", Values: []string{"A1", "A2"}, Operator: query.Or},
					{Field: "B", Values: []string{"B1", "B2"}, Operator: query.Or},
				},
			},
			expectedResult: &evaluator.Result{
				Matches: []map[string][]string{
					{"TopLevelA": {"TopLevelValA"}, "A": {"A1"}, "B": {"B1"}},
					{"TopLevelA": {"TopLevelValA"}, "A": {"A2"}, "B": {"B2"}},
				},
			},
		},
		{
			desc: "linked, multilevel, top doesn't match",
			obj: &TopLevel{
				ValA: "TopLevelValA",
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "TopLevelA", Values: []string{"NONEXISTENT"}},
					{Field: "A", Values: []string{"A1"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
		},
		{
			desc: "linked, multilevel, bottom doesn't match",
			obj: &TopLevel{
				ValA: "TopLevelValA",
				NestedSlice: []Nested{
					{NestedValA: "A0", NestedValB: "B0"},
					{NestedValA: "A1", NestedValB: "B1"},
				},
			},
			q: &query.Query{
				FieldQueries: []*query.FieldQuery{
					{Field: "TopLevelA", Values: []string{"TopLevelValA"}},
					{Field: "A", Values: []string{"A0"}},
					{Field: "B", Values: []string{"B1"}},
				},
			},
		},
	})
}
