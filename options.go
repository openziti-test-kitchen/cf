package cf

import "reflect"

type Instantiator func() interface{}
type Setter func(v interface{}, f reflect.Value) error
type Wiring func(cf interface{}) error
type NameConverter func(f reflect.StructField) string

type Options struct {
	Instantiators map[reflect.Type]Instantiator
	Setters       map[reflect.Type]Setter
	Wirings       map[reflect.Type][]Wiring
	NameConverter NameConverter
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
		NameConverter: PassthroughNameConverter,
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

func (opt *Options) AddWiring(t reflect.Type, w Wiring) *Options {
	if opt.Wirings == nil {
		opt.Wirings = make(map[reflect.Type][]Wiring)
	}
	opt.Wirings[t] = append(opt.Wirings[t], w)
	return opt
}

func (opt *Options) SetNameConverter(nc NameConverter) *Options {
	opt.NameConverter = nc
	return opt
}
