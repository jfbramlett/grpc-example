package testdef

type FunctionArg struct {
	ValuesOverride			map[string]interface{}	`json:"valuesOverride"`
	FieldTags				map[string]string		`json:"fieldTags"`
}

type FunctionDef struct {
	Name 		string						`json:"name"`
	Args		map[string]FunctionArg		`json:"args"`
}


type TestDef struct {
	Name				string					`json:"name"`
	ClientClassName		string					`json:"clientClassName"`
	Function			FunctionDef				`json:"function"`
	TestValues			map[string]interface{}	`json:"testValues"`
	TestTags			map[string]string		`json:"testTags"`
}



type TestSuiteDef struct {
	Tests			[]TestDef					`json:"tests"`
	GlobalValues	map[string]interface{}		`json:"globalValues"`
	GlobalTags		map[string]string			`json:"globalTags"`
}