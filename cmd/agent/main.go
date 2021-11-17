package main

import (
	"os"

	"github.com/kubeshop/testkube-executor-example/pkg/runner"
	"github.com/kubeshop/testkube/pkg/runner/agent"
)

func main() {
	agent.Run(runner.NewRunner(), os.Args)
}
