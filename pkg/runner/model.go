package runner

type FunctionArg struct {
	ValuesOverride			map[string]interface{}	`json:"valuesOverride"`
	FieldTags				map[string]string		`json:"fieldTags"`
}


type RunDef struct {
	Name            string                 `json:"name"`
	ClientClassName string                 `json:"clientClassName"`
	FunctionName	string				   `json:"functionName"`
	Args			map[string]FunctionArg `json:"args"`
	Validator		string				   `json:"validator"`
}


type RunSuiteDef struct {
	Tests			[]RunDef              	`json:"runDefinitions"`
	GlobalValues	map[string]interface{} 	`json:"globalValues"`
	GlobalTags		map[string]string    	`json:"globalTags"`
}


type RunResult struct {
	Name		string
	Passed		bool
	Error		error
}