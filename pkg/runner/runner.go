package runner

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor/output"
)

func NewRunner() *ExampleRunner {
	return &ExampleRunner{}
}

// ExampleRunner for template - change me to some valid runner
type ExampleRunner struct {
}

func (r *ExampleRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {

	// Our test data will be simple string with url inside we can get it as execution.Content.Data
	// but for more sophi

	// execution.Content could have git repo data
	// we're also passing content files/directories as mounted volume in directory
	path := os.Getenv("RUNNER_DATADIR")

	// let's print content of passed volume
	output.PrintEvent("path:", path)
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		d, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		output.PrintEvent("- file: ", path, string(d))
		return nil

	})

	output.PrintLog("some log message from executor")

	// e.g. Cypress test is stored in Git repo so Testkube will checkout it automatically
	// and allow you to use it easily
	// we can create test like below:
	// $ echo "https://httpstat.us/200" | kubectl testkube tests create --name example-google-test --type example/test
	uri := strings.TrimSuffix(execution.Content.Data, "\n") // newline on the end is not needed :)

	// other way to get data could be load it from Git e.g. file in git repo
	// if test will be created using git file content (like this one example with uri stored in file https://github.com/kubeshop/testkube-executor-example/blob/88f6ad4b4317a5a17bf6cd70bb41c4622d3ccfef/examples/test-definition-file)
	// Testkube will automatically map this dir/file to directory defined in the RUNNER_DATADIR env variable

	resp, err := http.Get(uri)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

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
