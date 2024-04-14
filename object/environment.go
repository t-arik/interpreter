package object

func NewEnclosedEnvironment(outer *Environment) *Environment{
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store} 
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (env *Environment) Get(name string) (Object, bool) {
	val, ok := env.store[name]
	if !ok && env.outer != nil {
		return env.outer.Get(name)
	}
	return val, ok
}

func (env *Environment) Set(name string, value Object) Object {
	env.store[name] = value
	return value
}
