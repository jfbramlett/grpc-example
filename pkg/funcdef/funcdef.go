package funcdef

type FunctionArg struct {
	FuncType		string
	ValuesOverride	string
	FieldTags		map[string]string
}

type Function struct {
	Name 		string
	Args		[]FunctionArg
}


type Test struct {
	steps		[]Function
}
