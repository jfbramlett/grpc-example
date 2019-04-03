package testdef

type FunctionArg struct {
	ValuesOverrideJson		string					`json:"valuesOverrideJson"`
	FieldTags				map[string]string		`json:"fieldTags"`
}

type Function struct {
	Name 		string						`json:"name"`
	Args		map[string]FunctionArg		`json:"args"`
}


type TestDef struct {
	Name				string			`json:"name"`
	ClientClassName		string			`json:"clientClassName"`
	Function			Function		`json:"function"`
}



type TestSuiteDef struct {
	Tests			[]TestDef		`json:"tests"`
}