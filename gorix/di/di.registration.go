package di

type Registration struct {
	Value      any
	Options    []ProviderOption
	IsInstance bool
}

func Provider(value any, opts ...ProviderOption) Registration {
	return Registration{
		Value:   value,
		Options: opts,
	}
}

func Instance(value any, opts ...ProviderOption) Registration {
	return Registration{
		Value:      value,
		Options:    opts,
		IsInstance: true,
	}
}
