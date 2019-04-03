package testdef

type FunctionArg struct {
	FuncType				string					`json:"funcType"`
	ValuesOverrideJson		string					`json:"valuesOverrideJson"`
	FieldTags				map[string]string		`json:"fieldTags"`
}

type Function struct {
	Name 		string				`json:"name"`
	Args		[]FunctionArg		`json:"args"`
}


type TestDef struct {
	ClientClassName		string			`json:"clientClassName"`
	Function			Function		`json:"function"`
}



type TestSuiteDef struct {
	Tests			[]TestDef		`json:"tests"`
}