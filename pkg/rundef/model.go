package rundef

type FunctionArg struct {
	ValuesOverride			map[string]interface{}	`json:"valuesOverride"`
	FieldTags				map[string]string		`json:"fieldTags"`
}


type FunctionDef struct {
	Name 		string						`json:"name"`
	Args		map[string]FunctionArg		`json:"args"`
}


type RunDef struct {
	Name            string                 `json:"name"`
	ClientClassName string                 `json:"clientClassName"`
	Function        FunctionDef            `json:"function"`
	RunValues       map[string]interface{} `json:"runValues"`
	RunTags         map[string]string      `json:"runTags"`
	Validator		string				   `json:"validator"`
}


type RunDefSuite struct {
	Tests			[]RunDef              	`json:"runDefinitions"`
	GlobalValues	map[string]interface{} 	`json:"globalValues"`
	GlobalTags		map[string]string    	`json:"globalTags"`
}


type RunResult struct {
	Name		string
	Passed		bool
	Error		error
}