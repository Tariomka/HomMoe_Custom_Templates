package editor_state_model

import (
	"reflect"
	"slices"
	"strings"
)

// DiffFieldOverrides compares every field the Overrides struct mirrors (by
// name) and returns, as overrides, the edited values that differ from derived,
// plus the JSON names of fields edited leaves nil while derived has a value.
// Fields Overrides does not mirror are structure and never compared. The
// returned overrides share storage with edited.
func DiffFieldOverrides[Item, Overrides any](derived, edited Item) (Overrides, []string) {
	var overrides Overrides
	var unset []string
	overridesValue := reflect.ValueOf(&overrides).Elem()
	derivedValue := reflect.ValueOf(derived)
	editedValue := reflect.ValueOf(edited)
	for index := range overridesValue.NumField() {
		field := overridesValue.Type().Field(index)
		editedField := editedValue.FieldByName(field.Name)
		if reflect.DeepEqual(derivedValue.FieldByName(field.Name).Interface(), editedField.Interface()) {
			continue
		}

		if isNilKind(editedField.Kind()) && editedField.IsNil() {
			unset = append(unset, fieldJSONName(field))
			continue
		}

		target := overridesValue.Field(index)
		if field.Type == editedField.Type() {
			target.Set(editedField)
			continue
		}

		pointer := reflect.New(editedField.Type())
		pointer.Elem().Set(editedField)
		target.Set(pointer)
	}
	return overrides, unset
}

// ApplyFieldOverrides returns derived with every non-nil override written over
// it and every unset field cleared. Unknown unset names are ignored.
func ApplyFieldOverrides[Item, Overrides any](derived Item, overrides Overrides, unset []string) Item {
	result := derived
	resultValue := reflect.ValueOf(&result).Elem()
	overridesValue := reflect.ValueOf(overrides)
	for index := range overridesValue.NumField() {
		field := overridesValue.Type().Field(index)
		target := resultValue.FieldByName(field.Name)
		if slices.Contains(unset, fieldJSONName(field)) {
			target.SetZero()
			continue
		}

		value := overridesValue.Field(index)
		if value.IsNil() {
			continue
		}

		if field.Type == target.Type() {
			target.Set(value)
			continue
		}

		target.Set(value.Elem())
	}
	return result
}

func isNilKind(kind reflect.Kind) bool {
	return kind == reflect.Pointer || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Interface
}

// cloneFieldOverrides deep-copies every pointer and slice reachable from the
// overrides. Interface values (placement-rule args) are immutable scalars and
// are shared.
func cloneFieldOverrides[Overrides any](overrides Overrides) Overrides {
	var clone Overrides
	reflect.ValueOf(&clone).Elem().Set(deepCopy(reflect.ValueOf(overrides)))
	return clone
}

func deepCopy(value reflect.Value) reflect.Value {
	//nolint:exhaustive // only containers need copying; every other kind copies by value.
	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return value
		}

		copied := reflect.New(value.Type().Elem())
		copied.Elem().Set(deepCopy(value.Elem()))
		return copied
	case reflect.Slice:
		if value.IsNil() {
			return value
		}

		copied := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
		for index := range value.Len() {
			copied.Index(index).Set(deepCopy(value.Index(index)))
		}
		return copied
	case reflect.Struct:
		copied := reflect.New(value.Type()).Elem()
		copied.Set(value)
		for index := range value.NumField() {
			copied.Field(index).Set(deepCopy(value.Field(index)))
		}
		return copied
	default:
		return value
	}
}

func fieldJSONName(field reflect.StructField) string {
	name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
	return name
}
