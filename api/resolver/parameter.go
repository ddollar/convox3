package resolver

type Parameter struct {
	key   string
	value string
}

func (p *Parameter) Key() string {
	return p.key
}

func (p *Parameter) Value() string {
	return p.value
}

type ParameterArg struct {
	Key   string
	Value string
}
