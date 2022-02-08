package pathutil

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	stringType = reflect.TypeOf("")
)

// An augmentTree is a utility class used by AugmentedObj to efficiently store and retrieve
// augmented values that are added to the object by Path.
type augmentTree struct {
	value    *reflect.Value
	children map[Step]*augmentTree
}

func (t *augmentTree) takeStep(key Step) *augmentTree {
	if t == nil {
		return nil
	}
	return t.children[key]
}

func (t *augmentTree) getValue() *reflect.Value {
	if t == nil {
		return nil
	}
	return t.value
}

func convertValue(val reflect.Value) (interface{}, error) {
	kind := val.Kind()
	if kind == reflect.Chan || kind == reflect.Func || kind == reflect.UnsafePointer {
		return nil, fmt.Errorf("kind %s is not supported", kind)
	}
	// built-in type
	if kind <= reflect.Complex128 || kind == reflect.String {
		return val.Interface(), nil
	}
	if kind == reflect.Struct {
		out := make(map[string]interface{})
		err := convertStruct(val, &out)
		if err != nil {
			return nil, err
		}
		return out, nil
	}
	if kind == reflect.Ptr {
		// If it's a nil pointer, explicitly set it to `nil` in the output map.
		if val.IsNil() {
			return nil, nil
		}
		return convertValue(val.Elem())
	}
	if kind == reflect.Array || kind == reflect.Slice {
		out := make([]interface{}, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			converted, err := convertValue(val.Index(i))
			if err != nil {
				return nil, fmt.Errorf("invalid value at index %d: %w", i, err)
			}
			out = append(out, converted)
		}
		return out, nil
	}
	if kind == reflect.Map {
		keyType := val.Type().Key()
		if keyKind := keyType.Kind(); keyKind != reflect.String {
			return nil, fmt.Errorf("unsupported key type for map: %s", keyKind)
		}
		out := make(map[string]interface{})
		mapIter := val.MapRange()
		for mapIter.Next() {
			key := mapIter.Key()
			keyAsString := key.Convert(stringType).Interface().(string)
			mapValue := mapIter.Value()
			convertedMapValue, err := convertValue(mapValue)
			if err != nil {
				return nil, fmt.Errorf("unsupported map value for key %s: %w", keyAsString, err)
			}
			out[keyAsString] = convertedMapValue
		}
		return out, nil
	}
	if kind == reflect.Interface {
		if val.IsNil() {
			return nil, nil
		}
		return convertValue(val.Elem())
	}
	return nil, utils.Should(fmt.Errorf("unsupported kind: %v", kind))
}

func convertStruct(val reflect.Value, out *map[string]interface{}) error {
	if kind := val.Kind(); kind != reflect.Struct {
		return errors.Errorf("value is of type %s, not struct", kind)
	}
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		if field.Anonymous {
			return errors.Errorf("field %s is embedded, which is not currently supported", field.Name)
		}
		if !field.IsExported() {
			continue
		}
		converted, err := convertValue(fieldVal)
		if err != nil {
			return fmt.Errorf("converting field %s: %w", field.Name, err)
		}
		(*out)[field.Name] = converted
	}
	return nil
}

func takeSteps(m *map[string]interface{}, steps []Step) (*map[string]interface{}, error) {
	if len(steps) == 0 {
		return m, nil
	}
	var currentValue interface{} = *m
	for i := 0; i < len(steps)-1; i++ {
		step := steps[i]
		if idx := step.Index(); idx >= 0 {
			asSlice, ok := currentValue.([]interface{})
			if !ok {
				return nil, fmt.Errorf("couldn't take index step %d (among steps %+v): expected a slice", idx, steps)
			}
			if idx >= len(asSlice) {
				return nil, fmt.Errorf("couldn't take index step %d (among steps %+v): slice too short (length %d)", idx, steps, len(asSlice))
			}
			currentValue = asSlice[idx]
			continue
		}
		field := step.Field()
		asMapStringInterface, ok := currentValue.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("couldn't take field step %s (among steps %+v): expected a slice", field, steps)
		}
		currentValue = asMapStringInterface[field]
	}

	// Now, the last step. Here, we must be in a map[string]interface{}, and we will create a new object for the augmented
	// value to populate into.
	field := steps[len(steps)-1].Field()
	if field == "" {
		return nil, fmt.Errorf("invalid augment (after steps %+v): last step should be a field", steps)
	}
	asMapStringInterface, ok := currentValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("steps %+v invalid: value at the end (%+v) is not an object", steps, currentValue)
	}
	childOut := make(map[string]interface{})
	asMapStringInterface[field] = childOut
	return &childOut, nil
}

func (t *augmentTree) populateValue(stepsSoFar []Step, out *map[string]interface{}) error {
	valPtr := t.getValue()
	if valPtr != nil {
		val := *valPtr
		// If it's a pointer to a struct, dereference it.
		if kind := val.Kind(); kind == reflect.Ptr {
			if !val.IsNil() {
				val = val.Elem()
				subOut, err := takeSteps(out, stepsSoFar)
				if err != nil {
					return err
				}
				if err := convertStruct(val, subOut); err != nil {
					return err
				}
			}
		}
	}
	for nextStep, child := range t.children {
		allSteps := append([]Step{}, stepsSoFar...)
		allSteps = append(allSteps, nextStep)
		if err := child.populateValue(allSteps, out); err != nil {
			return fmt.Errorf("failed to populate after child at step %v: %w", nextStep, err)
		}
	}
	return nil
}

func (t *augmentTree) getFullValue() (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := t.populateValue(nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func addAugmentedObjToTreeAtPath(rootTree *augmentTree, path *Path, subObj *AugmentedObj) error {
	currentTree := rootTree
	for _, step := range path.steps {
		if currentTree.children == nil {
			currentTree.children = make(map[Step]*augmentTree)
		}
		subTree := currentTree.children[step]
		if subTree == nil {
			subTree = &augmentTree{}
			currentTree.children[step] = subTree
		}
		currentTree = subTree
	}

	if currentTree.children != nil {
		return errors.Errorf("cannot add subObj %v to tree %v: children exist at this path", subObj, currentTree.children)
	}
	currentTree.value = subObj.augmentTreeRoot.value
	currentTree.children = subObj.augmentTreeRoot.children
	return nil
}

// An AugmentedObj represents an object with some augments.
// Concretely, this means that it effectively consists of two parts:
// -> the core object itself
// -> a mapping of Paths to other (possibly augmented) objects.
// For example, given a struct like
// type A struct {
//    IntVal int
// }
// and an object like A{IntVal: 1},
// you could augment it with "StringVal": "string".
// This makes it possible to treat the Augmented object _as if_
// it was A{IntVal: 1, StringVal: "string"}.
// This is a simple example -- it's possible to augment a value at an
// arbitrary path, traversing struct fields and slice indices, with an
// arbitrary object (which may, in turn, be an augmented object itself).
// It is a concrete realization of the object hierarchy described
// in an AugmentedObjMeta.
// Callers must use NewAugmentedObj to create one.
type AugmentedObj struct {
	augmentTreeRoot augmentTree
}

// NewAugmentedObj returns a ready-to-use instance of AugmentedObj, where the core
// object is the passed object.
// Callers can then call the AddObjAt methods to add augmented objects at various
// paths within this object.
func NewAugmentedObj(actualObj interface{}) *AugmentedObj {
	value := reflect.ValueOf(actualObj)
	return &AugmentedObj{augmentTreeRoot: augmentTree{value: &value}}
}

// AddAugmentedObjAt augments this object with the passed subObj, at the given path.
func (o *AugmentedObj) AddAugmentedObjAt(subObj *AugmentedObj, steps ...Step) error {
	return addAugmentedObjToTreeAtPath(&o.augmentTreeRoot, NewPath(steps...), subObj)
}

// AddPlainObjAt is a convenience wrapper around AddAugmentedObjAt for sub-objects
// that are not augmented.
func (o *AugmentedObj) AddPlainObjAt(subObj interface{}, steps ...Step) error {
	return o.AddAugmentedObjAt(NewAugmentedObj(subObj), steps...)
}

// Value returns an AugmentedValue, which starts off at the "root" of the augmented object.
func (o *AugmentedObj) Value() AugmentedValue {
	return &augmentedValue{underlying: *o.augmentTreeRoot.value, currentNode: &o.augmentTreeRoot}
}

// An AugmentedValue is a wrapper around a reflect.Value which can be traversed in a way
// that is augmentation-aware. It also keeps an internal record of the path traversed so far.
type AugmentedValue interface {
	Underlying() reflect.Value
	TakeStep(step MetaStep) (AugmentedValue, bool)
	// Elem is like calling .Elem on the underlying reflect.Value.
	// It panics if Elem() on the reflect.Value panics.
	Elem() AugmentedValue
	// Index is like calling .Index on the underlying reflect.Value.
	// It panics if Index(i) on the reflect.Value panics.
	Index(int) AugmentedValue
	PathFromRoot() *Path

	GetFullValue() (map[string]interface{}, error)
}

type augmentedValue struct {
	parent       *augmentedValue
	edgeToParent Step
	depth        int

	currentNode *augmentTree
	underlying  reflect.Value
}

func (v *augmentedValue) GetFullValue() (map[string]interface{}, error) {
	return v.currentNode.getFullValue()
}

func (v *augmentedValue) Elem() AugmentedValue {
	return &augmentedValue{underlying: v.underlying.Elem(), currentNode: v.currentNode, parent: v.parent, edgeToParent: v.edgeToParent, depth: v.depth}
}

func (v *augmentedValue) Index(i int) AugmentedValue {
	step := IndexStep(i)
	return v.childValue(v.underlying.Index(i), v.currentNode.takeStep(step), step)
}

func (v *augmentedValue) Underlying() reflect.Value {
	return v.underlying
}

func (v *augmentedValue) TakeStep(metaStep MetaStep) (AugmentedValue, bool) {
	var newUnderlying reflect.Value
	var found bool

	step := FieldStep(metaStep.FieldName)
	nextNode := v.currentNode.takeStep(step)
	if metaStep.StructFieldIndex != nil {
		// This is a "static" struct -- traverse it directly.
		newUnderlying = v.underlying.FieldByIndex(metaStep.StructFieldIndex)
		found = true
	} else {
		// See if this is an augmented path.
		if value := nextNode.getValue(); value != nil {
			newUnderlying = *value
			found = true
		} else {
			// This specific case is hit when the field in the struct is an interface type,
			// in which case StructFieldIndex will not be present.
			if v.underlying.Kind() == reflect.Struct {
				newUnderlying = v.underlying.FieldByName(metaStep.FieldName)
				if newUnderlying.IsValid() {
					found = true
				}
			}
		}
	}
	if !found {
		return nil, false
	}
	return v.childValue(newUnderlying, nextNode, step), true
}

func (v *augmentedValue) childValue(newUnderlying reflect.Value, nextNode *augmentTree, edge Step) *augmentedValue {
	return &augmentedValue{
		parent:       v,
		edgeToParent: edge,
		depth:        v.depth + 1,
		underlying:   newUnderlying,
		currentNode:  nextNode,
	}
}

func (v *augmentedValue) PathFromRoot() *Path {
	p := &Path{steps: make([]Step, v.depth)}
	v.populateIntoSteps(&p.steps)
	return p
}

func (v *augmentedValue) populateIntoSteps(outSlice *[]Step) {
	if v.depth == 0 {
		return
	}
	(*outSlice)[v.depth-1] = v.edgeToParent
	v.parent.populateIntoSteps(outSlice)
}
