package cf

import "reflect"

type Instantiator func() interface{}

type Setter func(v interface{}, f reflect.Value) error

type Options struct {
	Instantiators map[reflect.Type]Instantiator
	Setters       map[reflect.Type]Setter
}

func DefaultOptions() *Options {
	opt := &Options{
		Setters: map[reflect.Type]Setter{
			reflect.TypeOf(0):          intHandler,
			reflect.TypeOf(float64(0)): float64Handler,
			reflect.TypeOf(true):       boolHandler,
			reflect.TypeOf(""):         stringHandler,
			reflect.TypeOf([]string{}): stringArrayHandler,
		},
	}
	return opt
}

func (opt *Options) AddInstantiator(t reflect.Type, i Instantiator) *Options {
	if opt.Instantiators == nil {
		opt.Instantiators = make(map[reflect.Type]Instantiator)
	}
	opt.Instantiators[t] = i
	return opt
}

func (opt *Options) AddSetter(t reflect.Type, s Setter) *Options {
	if opt.Setters == nil {
		opt.Setters = make(map[reflect.Type]Setter)
	}
	opt.Setters[t] = s
	return opt
}