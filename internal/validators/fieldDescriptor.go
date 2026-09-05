package validators

import "github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"

type fieldDescriptor[T any] struct {
	name  string
	value func(state *editor_state_model.EditorState) *T
}

type rangedFieldDescriptor[T comparable] struct {
	fieldDescriptor[T]

	lowest  T
	highest T
}

func newRangedFieldDescriptor[T comparable](
	name string,
	lowest, highest T,
	value func(state *editor_state_model.EditorState) *T) rangedFieldDescriptor[T] {
	return rangedFieldDescriptor[T]{
		name:    name,
		value:   value,
		lowest:  lowest,
		highest: highest,
	}
}
