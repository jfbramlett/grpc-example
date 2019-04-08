package rundef

import "reflect"

type RunResult struct {
	Name		string
	Result		reflect.Value
	Passed		bool
	Error		error
}
