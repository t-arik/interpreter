package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": { Fn: builtinLen },
	"first": { Fn: builtinFirst },
	"last": { Fn: builtinLast },
	"rest": { Fn: builtinRest },
	"push": { Fn: builtinPush },
	"puts": { Fn: builtinPuts },
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got %s", arg.Type())
	}
}

func builtinFirst(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	if len(arr.Elements) == 0 {
		return NULL
	}

	return arr.Elements[0]
}

func builtinLast(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
	}

	if len(arr.Elements) == 0 {
		return NULL
	}

	return arr.Elements[len(arr.Elements) - 1]
}

func builtinRest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	if len(arr.Elements) == 0 {
		return NULL
	}

	newElements := make([]object.Object, len(arr.Elements)-1)
	copy(newElements, arr.Elements[1:])

	return &object.Array{
		Elements: newElements,
	}
}

func builtinPush(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	length := len(arr.Elements)
	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{
		Elements: newElements,
	}
}

func builtinPuts(args ...object.Object) object.Object {
	for _, obj := range args {
		fmt.Println(obj.Inspect())
	}

	return NULL
}

