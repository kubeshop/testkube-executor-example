package runner

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

func NewRunner() *ExampleRunner {
	return &ExampleRunner{}
}

// ExampleRunner for template - change me to some valid runner
type ExampleRunner struct {
}

func (r *ExampleRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	// ScriptContent will have URI
	uri := execution.ScriptContent
	resp, err := http.Get(uri)
	if err != nil {
		return result, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	// if get is successful return success result
	if resp.StatusCode == 200 {
		return testkube.ExecutionResult{
			Status: testkube.ExecutionStatusSuccess,
			Output: string(b),
		}, nil
	}

	// else we'll return error to simplify example
	err = fmt.Errorf("invalid status code %d, (uri:%s)", resp.StatusCode, uri)
	return result.Err(err), nil
}
