package runners

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/shell"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/utils"
)

type goTestRunner struct {
	test    *model.TestEntry
	cmdLine string
	envMgr  shell.EnvironmentManager
}

func (runner *goTestRunner) Run(timeoutCtx context.Context, env []string, writer *bufio.Writer) error {
	logger := func(s string) {}
	cmdEnv := append(runner.envMgr.GetProcessedEnv(), env...)
	_, err := utils.RunCommand(timeoutCtx, runner.cmdLine, runner.test.ExecutionConfig.PackageRoot,
		logger, writer, cmdEnv, map[string]string{}, false)
	return err
}

func (runner *goTestRunner) GetCmdLine() string {
	return runner.cmdLine
}

// NewGoTestRunner - creates go test runner
func NewGoTestRunner(ids string, test *model.TestEntry, timeout time.Duration) TestRunner {
	cmdLine := fmt.Sprintf("go test . -test.timeout %v -count 1 --run \"^(%s)$\\\\z\" --tags \"%s\" --test.v",
		timeout, test.Name, test.Tags)

	envMgr := shell.NewEnvironmentManager()
	_ = envMgr.ProcessEnvironment(ids, "gotest", os.TempDir(), test.ExecutionConfig.Env, map[string]string{})

	return &goTestRunner{
		test:    test,
		cmdLine: cmdLine,
		envMgr:  envMgr,
	}
}
