package commands

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/config"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/execmanager"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/k8s"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/providers"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/providers/packet"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/providers/shell"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/reporting"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/runners"
	shell_mgr "github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/shell"
	"github.com/networkservicemesh/networkservicemesh/test/cloudtest/pkg/utils"
)

const (
	defaultConfigFile string = ".cloudtest.yaml"
)

// Arguments - command line arguments
type Arguments struct {
	clusters        []string // A list of enabled clusters from configuration.
	providerConfig  string   // A folder to start scaning for tests inside
	count           int      // Limit number of tests to be run per every cloud
	instanceOptions providers.InstanceOptions
	onlyEnabled     bool // Disable all clusters and enable only enabled in command line.
}

type clusterState byte

const (
	clusterAdded clusterState = iota
	clusterReady
	clusterBusy
	clusterStarting
	clusterStopping
	clusterCrashed
	clusterNotAvailable
	clusterShutdown
)

// Cluster operation record, to be added as testcase into report.
type clusterOperationRecord struct {
	time     time.Time
	duration time.Duration
	status   clusterState
	attempt  int
	logFile  string
	errMsg   error
}

type clusterInstance struct {
	instance      providers.ClusterInstance
	state         clusterState
	group         *clustersGroup
	startCount    int
	id            string
	taskCancel    context.CancelFunc
	cancelMonitor context.CancelFunc
	startTime     time.Time

	executions    []*clusterOperationRecord
	retestCounter int // If test is requesting retest on this cluster instance, we count how many times it is happening, it will be set to 0 if test is not request retest.
}
type clustersGroup struct {
	instances []*clusterInstance
	provider  providers.ClusterProvider
	config    *config.ClusterProviderConfig
	tasks     map[string]*testTask // All tasks assigned to this cluster.
	completed map[string]*testTask
}

type testTask struct {
	taskID           string
	test             *model.TestEntry
	clusters         []*clustersGroup
	clusterInstances []*clusterInstance
	clusterTaskID    string
}

type eventKind byte

const (
	eventTaskUpdate eventKind = iota
	eventClusterUpdate
)

type operationEvent struct {
	kind            eventKind
	clusterInstance *clusterInstance
	task            *testTask
}

type executionContext struct {
	sync.RWMutex
	manager          execmanager.ExecutionManager
	clusters         []*clustersGroup
	operationChannel chan operationEvent
	tests            []*model.TestEntry
	tasks            []*testTask
	running          map[string]*testTask
	completed        []*testTask
	skipped          []*testTask
	cloudTestConfig  *config.CloudTestConfig
	report           *reporting.JUnitFile
	startTime        time.Time
	clusterReadyTime time.Time
	factory          k8s.ValidationFactory
	arguments        *Arguments
	clusterWaitGroup sync.WaitGroup // Wait group for clusters destroying
}

// CloudTestRun - CloudTestRun
func CloudTestRun(cmd *cloudTestCmd) {
	var configFileContent []byte
	var err error

	if cmd.cmdArguments.providerConfig == "" {
		cmd.cmdArguments.providerConfig = defaultConfigFile
	}

	configFileContent, err = ioutil.ReadFile(cmd.cmdArguments.providerConfig)
	if err != nil {
		logrus.Errorf("Failed to read config file %v", err)
		return
	}

	// Root config
	testConfig := &config.CloudTestConfig{}
	err = parseConfig(testConfig, configFileContent)
	if err != nil {
		logrus.Errorf("Failed to parse config %v", err)
		os.Exit(1)
	}

	// Process config imports
	err = performImport(testConfig)
	if err != nil {
		logrus.Errorf("Failed to process config imports %v", err)
		os.Exit(1)
	}

	_, err = PerformTesting(testConfig, k8s.CreateFactory(), cmd.cmdArguments)
	if err != nil {
		logrus.Errorf("Failed to process tests %v", err)
		os.Exit(1)
	}
}

func performImport(testConfig *config.CloudTestConfig) error {
	for _, imp := range testConfig.Imports {
		importConfig := &config.CloudTestConfig{}

		configFileContent, err := ioutil.ReadFile(imp)
		if err != nil {
			logrus.Errorf("аailed to read config file %v", err)
			return err
		}
		if err = parseConfig(importConfig, configFileContent); err != nil {
			return err
		}

		// Do add imported items
		testConfig.Executions = append(testConfig.Executions, importConfig.Executions...)
		testConfig.Providers = append(testConfig.Providers, importConfig.Providers...)
	}
	return nil
}

// PerformTesting - PerformTesting
func PerformTesting(config *config.CloudTestConfig, factory k8s.ValidationFactory, arguments *Arguments) (*reporting.JUnitFile, error) {
	ctx := &executionContext{
		cloudTestConfig:  config,
		operationChannel: make(chan operationEvent, 100),
		tasks:            []*testTask{},
		running:          map[string]*testTask{},
		completed:        []*testTask{},
		tests:            []*model.TestEntry{},
		factory:          factory,
		arguments:        arguments,
	}
	ctx.manager = execmanager.NewExecutionManager(ctx.cloudTestConfig.ConfigRoot)
	// Create cluster instance handles
	if err := ctx.createClusters(); err != nil {
		return nil, err
	}
	// Collect tests
	if err := ctx.findTests(); err != nil {
		logrus.Errorf("Error finding tests %v", err)
		return nil, err
	}
	// We need to be sure all clusters will be deleted on end of execution.
	defer ctx.performShutdown()
	// Fill tasks to be executed..
	ctx.createTasks()

	err := ctx.performExecution()
	reportfile, err2 := ctx.generateJUnitReportFile()
	if err2 != nil {
		logrus.Errorf("Error during generation of report: %v", err2)
	}
	if err != nil {
		return reportfile, err
	}
	return reportfile, err2
}

func parseConfig(cloudTestConfig *config.CloudTestConfig, configFileContent []byte) error {
	err := yaml.Unmarshal(configFileContent, cloudTestConfig)
	if err != nil {
		err = errors.Wrap(err, "failed to parse configuration file")
		logrus.Errorf(err.Error())
		return err
	}
	logrus.Infof("configuration file loaded successfully...")
	return nil
}

func (ctx *executionContext) performShutdown() {
	// We need to stop all clusters we started
	if !ctx.arguments.instanceOptions.NoStop {
		for _, clG := range ctx.clusters {
			group := clG
			for _, cInst := range group.instances {
				curInst := cInst
				if curInst.taskCancel != nil {
					logrus.Infof("Canceling currently running task")
					curInst.taskCancel()
				}
				logrus.Infof("Schedule Closing cluster %v %v", group.config.Name, curInst.id)
				ctx.clusterWaitGroup.Add(1)

				go func() {
					defer ctx.clusterWaitGroup.Done()
					logrus.Infof("Closing cluster %v %v", group.config.Name, curInst.id)
					ctx.destroyCluster(curInst, false, false)
				}()
			}
		}
		ctx.clusterWaitGroup.Wait()
	}
	logrus.Infof("All clusters destroyed")
}

func (ctx *executionContext) performExecution() error {
	logrus.Infof("Starting test execution")
	ctx.startTime = time.Now()
	ctx.clusterReadyTime = ctx.startTime

	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(ctx.cloudTestConfig.Timeout)*time.Second)
	defer cancelFunc()

	termChannel := tools.NewOSSignalChannel()
	defer ctx.printStatistics()

	statTicker := time.NewTicker(60 * time.Second)
	defer statTicker.Stop()

	for len(ctx.tasks) > 0 || len(ctx.running) > 0 {
		// WE take 1 test task from list and do execution.
		ctx.assignTasks()
		ctx.checkClustersUsage()

		if len(ctx.tasks) == 0 && len(ctx.running) == 0 {
			// No tasks left to execute.
			break
		}

		select {
		case event := <-ctx.operationChannel:
			switch event.kind {
			case eventClusterUpdate:
				ctx.performClusterUpdate(event)
			case eventTaskUpdate:
				// Remove from running onces.
				ctx.processTaskUpdate(event)
			}
		case <-termChannel:
			return errors.New("termination request is received")
		case <-timeoutCtx.Done():
			return errors.Errorf("global timeout elapsed: %v seconds", ctx.cloudTestConfig.Timeout)
		case <-statTicker.C:
			ctx.printStatistics()
		}
	}

	logrus.Info("Finished test execution")
	return nil
}

func (ctx *executionContext) assignTasks() {
	if len(ctx.tasks) > 0 {
		// Lets check if we have cluster required and start it
		// Check if we have cluster we could assign.
		newTasks := []*testTask{}
		for _, task := range ctx.tasks {
			if task.test.Status == model.StatusSkipped {
				logrus.Infof("Ignoring skipped task:  %s", task.test.Name)
				continue
			}
			clustersAvailable, clustersToUse, assigned := ctx.selectClusterForTask(task)
			if assigned {
				// Start task execution.
				err := ctx.startTask(task, clustersToUse)
				if err != nil {
					logrus.Errorf("Error starting task  %s %v", task.test.Name, err)
					assigned = false
				} else {
					ctx.running[task.taskID] = task
				}
			}
			// If we finally not assigned.
			if !assigned {
				if clustersAvailable < len(task.clusters) {
					// We move task to skipped since, no clusters could execute it, all attempts for clusters to recover are finished.
					logrus.Errorf("Move task to skipped since no clusters could execute it: %s", task.test.Name)
					task.test.Status = model.StatusSkippedSinceNoClusters
					for _, cl := range task.clusters {
						delete(cl.tasks, task.test.Key)
						cl.completed[task.test.Key] = task
					}
					ctx.completed = append(ctx.completed, task)
				} else {
					newTasks = append(newTasks, task)
				}
			}
		}
		ctx.tasks = newTasks
	}
}

func (ctx *executionContext) performClusterUpdate(event operationEvent) {
	ctx.Lock()
	defer ctx.Unlock()
	logrus.Infof("Cluster instance %s is updated: %v", event.clusterInstance.id, event.clusterInstance)
	if event.clusterInstance.taskCancel != nil && event.clusterInstance.state == clusterCrashed {
		// We have task running on cluster
		event.clusterInstance.taskCancel()
	}
	if event.clusterInstance.state == clusterReady {
		if ctx.clusterReadyTime == ctx.startTime {
			ctx.clusterReadyTime = time.Now()

		}
	}

}

func (ctx *executionContext) processTaskUpdate(event operationEvent) {
	delete(ctx.running, event.task.taskID)

	if event.task.test.Status == model.StatusRerunRequest && ctx.cloudTestConfig.RetestConfig.WarmupTimeout > 0 {
		go func() {
			ids := []string{}
			for _, ci := range event.task.clusterInstances {
				ids = append(ids, ci.id)
			}
			wtime := time.Second * time.Duration(ctx.cloudTestConfig.RetestConfig.WarmupTimeout)
			logrus.Infof("Warmup cluster operations: %v timeout: %v", ids, wtime)
			<-time.After(wtime)
			// Make cluster as ready
			ctx.makeInstancesReady(event.task.clusterInstances)

			for _, ci := range event.task.clusterInstances {
				ctx.operationChannel <- operationEvent{
					kind:            eventClusterUpdate,
					clusterInstance: ci,
				}
			}
		}()
	} else {
		// Make cluster as ready
		ctx.makeInstancesReady(event.task.clusterInstances)
	}
	if event.task.test.Status == model.StatusSuccess || event.task.test.Status == model.StatusFailed {
		ctx.Lock()
		ctx.completed = append(ctx.completed, event.task)
		ctx.Unlock()
		logrus.Infof("Completed %s on %s, %s, runtime: %v",
			event.task.test.Name,
			event.task.clusterTaskID,
			statusName(event.task.test.Status),
			event.task.test.Duration.Round(time.Second))

		for ind, cl := range event.task.clusters {
			delete(cl.tasks, event.task.test.Key)

			// Add test only to first cluster
			if ind == 0 {
				cl.completed[event.task.test.Key] = event.task
			}
		}
	} else {
		logrus.Infof("Re schedule task %v reason: %v", event.task.test.Name, statusName(event.task.test.Status))
		ctx.tasks = append(ctx.tasks, event.task)
	}
}

func (ctx *executionContext) makeInstancesReady(instances []*clusterInstance) {
	for _, inst := range instances {
		ctx.Lock()
		if inst.state == clusterBusy {
			inst.state = clusterReady
		}
		inst.taskCancel = nil
		ctx.Unlock()
	}
}

func statusName(status model.Status) interface{} {
	switch status {
	case model.StatusAdded:
		return "added"
	case model.StatusFailed:
		return "failed"
	case model.StatusSkipped:
		return "skipped"
	case model.StatusSuccess:
		return "success"
	case model.StatusTimeout:
		return "timeout"
	case model.StatusRerunRequest:
		return "rerun-request"
	}
	return fmt.Sprintf("code: %v", status)
}

func (ctx *executionContext) selectClusterForTask(task *testTask) (int, []*clusterInstance, bool) {
	var clustersToUse []*clusterInstance
	clustersAvailable := 0
	for _, cluster := range task.clusters {
		groupAssigned := false
		groupAvailable := false
		for _, ci := range cluster.instances {
			var state clusterState
			ctx.Lock()
			state = ci.state
			ctx.Unlock()
			// No task is assigned for cluster.
			switch state {
			case clusterAdded, clusterCrashed:
				// Try starting cluster
				if ctx.startCluster(ci) {
					groupAvailable = true
				}
			case clusterReady:
				groupAvailable = true
				if groupAssigned {
					continue
				}
				// Check if we match requirements.
				// We could assign task and start it running.
				clustersToUse = append(clustersToUse, ci)
				// We need to remove task from list
				groupAssigned = true
			case clusterBusy, clusterStarting, clusterStopping:
				groupAvailable = true
			}
		}
		if groupAvailable {
			clustersAvailable++
		}
	}
	return clustersAvailable, clustersToUse, len(clustersToUse) == len(task.clusters)
}

func (ctx *executionContext) printStatistics() {
	elapsed := time.Since(ctx.startTime)
	var elapsedRunning time.Duration
	ctx.RLock()
	elapsedRunning = time.Since(ctx.clusterReadyTime)
	running := ""
	for _, r := range ctx.running {
		running += fmt.Sprintf("\t\t%s on %v, %v\n", r.test.Name, r.clusterTaskID, time.Since(r.test.Started).Round(time.Second))
	}
	ctx.RUnlock()

	if len(running) > 0 {
		running = "\n\tRunning:\n" + running
	}
	clustersMsg := strings.Builder{}
	if len(ctx.clusters) > 0 {
		_, _ = clustersMsg.WriteString("\n\tClusters:\n")
	}
	for _, cl := range ctx.clusters {
		_, _ = clustersMsg.WriteString(fmt.Sprintf("\t\tCluster: %v Tasks left: %v\n", cl.config.Name, len(cl.tasks)))
		ctx.RLock()
		for _, inst := range cl.instances {
			_, _ = clustersMsg.WriteString(fmt.Sprintf("\t\t\t%s: %v, uptime: %v\n", inst.id, fromClusterState(inst.state),
				time.Since(inst.startTime).Round(time.Second)))
		}
		ctx.RUnlock()
	}

	remaining := ""
	if len(ctx.completed) > 0 {
		oneTask := elapsed / time.Duration(len(ctx.completed))
		remaining = fmt.Sprintf("%v", (time.Duration(len(ctx.tasks)+len(ctx.running)) * oneTask).Round(time.Second))
	}

	successTests := 0
	failedTests := 0
	skippedTests := 0
	timeoutTests := 0

	failedNames := ""

	for _, t := range ctx.completed {
		switch t.test.Status {
		case model.StatusSuccess:
			successTests++
		case model.StatusTimeout:
			timeoutTests++
		case model.StatusSkipped:
			skippedTests++
		case model.StatusFailed:
			failedTests++
			failedNames += "\n\t\t" + t.test.Name
		case model.StatusSkippedSinceNoClusters:
			skippedTests++
		}
	}

	logrus.Infof("Statistics:" +
		fmt.Sprintf("\n\tElapsed total: %v", elapsed.Round(time.Second)) +
		fmt.Sprintf("\n\tTests time: %v", elapsedRunning.Round(time.Second)) +
		fmt.Sprintf("\n\tTasks  Completed: %d", len(ctx.completed)) +
		fmt.Sprintf("\n\t       Remaining: %d (~%v)\n", len(ctx.running)+len(ctx.tasks), remaining) +
		fmt.Sprintf("%s%s", running, clustersMsg.String()) +
		fmt.Sprintf("\n\tStatus  Passed: %d"+
			"\n\tStatus  Failed: %d%v"+
			"\n\tStatus  Timeout: %d"+
			"\n\tStatus  Skipped: %d", successTests, failedTests, failedNames, timeoutTests, skippedTests))
}

func fromClusterState(state clusterState) string {
	switch state {
	case clusterReady:
		return "ready"
	case clusterAdded:
		return "added"
	case clusterBusy:
		return "running test"
	case clusterCrashed:
		return "crashed"
	case clusterNotAvailable:
		return "not available"
	case clusterStarting:
		return "starting"
	case clusterStopping:
		return "stopping"
	case clusterShutdown:
		return "shutdown"
	}
	return fmt.Sprintf("unknown state: %v", state)
}

func (ctx *executionContext) createTasks() {
	taskIndex := 0
	for i, test := range ctx.tests {
		selector := test.ExecutionConfig.ClusterSelector
		// In case of one cluster, we create task copies and execute on every cloud.

		var task *testTask
		if test.ExecutionConfig.ClusterCount > 1 {
			for _, clusterName := range selector {
				for _, cluster := range ctx.clusters {
					if clusterName == cluster.config.Name {
						if task == nil {
							task = ctx.createSingleTask(taskIndex, test, cluster, i)
							taskIndex++
						} else {
							task.clusters = append(task.clusters, cluster)
							cluster.tasks[task.test.Key] = task
						}
						break
					}
				}
			}
		} else {
			for _, cluster := range ctx.clusters {
				if len(selector) > 0 && utils.Contains(selector, cluster.config.Name) ||
					len(selector) == 0 {
					task = ctx.createSingleTask(taskIndex, test, cluster, i)
					taskIndex++
				}
			}
		}

		if task == nil {
			logrus.Errorf("%s: no clusters defined of required %v", test.Name, selector)
		} else if len(task.clusters) < test.ExecutionConfig.ClusterCount {
			logrus.Errorf("%s: not all clusters defined of required %v", test.Name, selector)
			task.test.Status = model.StatusSkipped
		}
	}
}

func (ctx *executionContext) createSingleTask(taskIndex int, test *model.TestEntry, cluster *clustersGroup, i int) *testTask {
	task := &testTask{
		taskID: fmt.Sprintf("%d", taskIndex),
		test: &model.TestEntry{
			Kind:            test.Kind,
			Name:            test.Name,
			Tags:            test.Tags,
			Status:          test.Status,
			ExecutionConfig: test.ExecutionConfig,
			Executions:      []model.TestEntryExecution{},
			RunScript:       test.RunScript,
		},
		clusters: []*clustersGroup{cluster},
	}

	// Generate task key to avoid crossing in cluster tasks map
	testKey := ""
	for _, clusterName := range test.ExecutionConfig.ClusterSelector {
		testKey = fmt.Sprintf("%s_%s", testKey, clusterName)
	}
	task.test.Key = fmt.Sprintf("%s_%s", testKey, test.Name)

	// To track cluster task executions.
	cluster.tasks[task.test.Key] = task
	if ctx.arguments.count > 0 && i >= ctx.arguments.count {
		logrus.Infof("Limit of tests for execution:: %v is reached. Skipping test %s", ctx.arguments.count, test.Name)
		test.Status = model.StatusSkipped
		ctx.skipped = append(ctx.skipped, task)
	} else {
		ctx.tasks = append(ctx.tasks, task)
	}
	return task
}

func (ctx *executionContext) startTask(task *testTask, instances []*clusterInstance) error {
	ids := ""
	for _, ci := range instances {
		if len(ids) > 0 {
			ids += "_"
		}
		ids += ci.id
		ctx.Lock()
		ci.state = clusterBusy
		ctx.Unlock()
	}
	fileName, file, err := ctx.manager.OpenFileTest(ids, task.test.Name, "run")
	if err != nil {
		return err
	}

	clusterConfigs := []string{}

	for _, inst := range instances {
		var clusterConfig string
		clusterConfig, err = inst.instance.GetClusterConfig()
		if err != nil {
			return err
		}
		clusterConfigs = append(clusterConfigs, clusterConfig)
	}

	task.clusterInstances = instances
	task.clusterTaskID = ids

	timeout := ctx.getTestTimeout(task)

	var runner runners.TestRunner
	switch task.test.Kind {
	case model.TestEntryKindShellTest:
		runner = runners.NewShellTestRunner(ids, task.test)
	case model.TestEntryKindGoTest:
		runner = runners.NewGoTestRunner(ids, task.test, timeout)
	default:
		return errors.New("invalid task runner")
	}

	ctx.executeTask(task, clusterConfigs, file, runner, timeout, instances, err, fileName)
	return nil
}

func (ctx *executionContext) executeTask(task *testTask, clusterConfigs []string, file io.Writer, runner runners.TestRunner, timeout time.Duration, instances []*clusterInstance, err error, fileName string) {
	go func() {
		testDelay := func() int {
			first := true
			ctx.RLock()
			for _, tt := range ctx.completed {
				if tt.clusterTaskID == task.clusterTaskID {
					first = false
					break
				}
			}
			ctx.RUnlock()
			delay := 0
			if !first {
				for _, cl := range task.clusters {
					if cl.config.TestDelay > delay {
						delay = cl.config.TestDelay
					}
				}
			}
			return delay
		}()
		if testDelay != 0 {
			logrus.Infof("Cluster %v requires %v seconds delay between tests", task.clusterTaskID, testDelay)
			<-time.After(time.Duration(testDelay) * time.Second)
			logrus.Infof("Cluster %v: %v seconds delay between tests completed", task.clusterTaskID, testDelay)
		}

		st := time.Now()
		env := []string{}
		// Fill Kubernetes environment variables.

		if len(task.test.ExecutionConfig.KubernetesEnv) > 0 {
			for ind, envV := range task.test.ExecutionConfig.KubernetesEnv {
				env = append(env, fmt.Sprintf("%s=%s", envV, clusterConfigs[ind]))
			}
		} else {
			for idx, cfg := range clusterConfigs {
				if idx == 0 {
					env = append(env, fmt.Sprintf("KUBECONFIG=%s", cfg))
				} else {
					env = append(env, fmt.Sprintf("KUBECONFIG%d=%s", idx, cfg))
				}
			}
		}

		writer := bufio.NewWriter(file)

		msg := fmt.Sprintf("Starting %s on %v\n", task.test.Name, task.clusterTaskID)
		logrus.Info(msg)
		_, _ = writer.WriteString(msg)
		_, _ = writer.WriteString(fmt.Sprintf("Command line %v\nenv==%v \n\n", runner.GetCmdLine(), env))
		_ = writer.Flush()

		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ctx.Lock()
		for _, inst := range instances {
			inst.taskCancel = cancel
		}
		task.test.Started = time.Now()
		ctx.Unlock()

		errCode := runner.Run(timeoutCtx, env, writer)

		_ = writer.Flush()

		// Check if test ask us restart it, and have few executions left
		if errCode != nil && len(ctx.cloudTestConfig.RetestConfig.Patterns) > 0 && ctx.cloudTestConfig.RetestConfig.RestartCount > 0 {
			if ctx.matchRestartRequest(fileName) {
				if len(task.test.Executions) < ctx.cloudTestConfig.RetestConfig.RestartCount {
					// Let's check if we have same cluster instance fail few times one after another with this error.
					for _, cinst := range task.clusterInstances {
						cinst.retestCounter++
						if cinst.retestCounter == ctx.cloudTestConfig.RetestConfig.AllowedRetests { // We it happened again, we need to re-start this cluster and give test one more attempt.
							// If cluster failed with network error most of time, let's re-create it.
							logrus.Errorf("Reached a limit of re-tests per cluster instance: %v %v %v", task.test.Name, cinst.id, ctx.cloudTestConfig.RetestConfig.AllowedRetests)
							cinst.retestCounter = 0
							_ = ctx.destroyCluster(cinst, true, false)
						}
						ctx.Lock()
						cinst.taskCancel = nil
						ctx.Unlock()
					}

					ctx.updateTestExecution(task, fileName, model.StatusRerunRequest)
				} else {
					msg := fmt.Sprintf("Test %v retry count %v exceed: err: %v", task.test.Name, ctx.cloudTestConfig.RetestConfig.RestartCount, errCode.Error())
					logrus.Errorf(msg)
					_, _ = writer.WriteString(errCode.Error())
					_ = writer.Flush()
					taskStatus := model.StatusFailed
					if ctx.cloudTestConfig.RetestConfig.RetestFailResult == "skip" {
						taskStatus = model.StatusSkipped
						task.test.SkipMessage = msg
					}
					ctx.updateTestExecution(task, fileName, taskStatus)
				}
				return
			}
		}

		// Update retestCounter if test is not retesting.
		for _, cinst := range task.clusterInstances {
			ctx.Lock()
			cinst.retestCounter = 0
			ctx.Unlock()
		}

		if errCode != nil {
			onFailErr := ctx.handleOnFailTask(task, env, writer)
			if err != nil {
				errCode = errors.Wrap(errCode, onFailErr.Error())
			}
		}

		task.test.Duration = time.Since(st)

		if errCode != nil {
			// Check if cluster is alive.
			clusterNotAvailable := false
			for _, inst := range instances {
				err = inst.instance.CheckIsAlive()
				if err != nil {
					logrus.Errorf("Task failed because cluster is not valid: %v %v %v", task.test.Name, inst.id, err)
					clusterNotAvailable = true
					ctx.destroyCluster(inst, true, false)
				}
				ctx.Lock()
				inst.taskCancel = nil
				ctx.Unlock()
			}
			if clusterNotAvailable {
				logrus.Errorf("Test is canceled due timeout and cluster error.. Will be re-run")
				ctx.updateTestExecution(task, fileName, model.StatusTimeout)
			} else {
				logrus.Errorf(errCode.Error())
				_, _ = writer.WriteString(errCode.Error())
				_ = writer.Flush()
				ctx.updateTestExecution(task, fileName, model.StatusFailed)
			}
		} else {
			ctx.updateTestExecution(task, fileName, model.StatusSuccess)
		}
	}()
}

func (ctx *executionContext) matchRestartRequest(fileName string) bool {
	// Check if output file contains restart request marker
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()

	reader := bufio.NewReader(f)
	for {
		r, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if utils.MatchRetestPattern(ctx.cloudTestConfig.RetestConfig.Patterns, r) {
			return true
		}
	}
	return false
}

func (ctx *executionContext) handleOnFailTask(task *testTask, env []string, writer *bufio.Writer) error {
	config := task.test.ExecutionConfig
	if config == nil {
		logrus.Warnf("OnFail: no execution config for %v", task.test.Name)
		return nil
	}
	if strings.TrimSpace(config.OnFail) == "" {
		logrus.Infof("OnFail: not provided OnFail script for config %v", config.Name)
		return nil
	}
	mgr := shell_mgr.NewEnvironmentManager()
	if err := mgr.ProcessEnvironment(task.clusterTaskID, "shellrun", os.TempDir(), append(env, config.Env...), nil); err != nil {
		logrus.Errorf("OnFail: an error during process env: %v", err)
		return err
	}
	timeout := ctx.getTestTimeout(task)
	context, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return runOnFailScript(context, config.OnFail, mgr.GetProcessedEnv(), writer)
}

func runOnFailScript(ctx context.Context, script string, env []string, writer *bufio.Writer) error {
	logger := func(s string) {
	}
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, cmd := range utils.ParseScript(script) {
		_, err := utils.RunCommand(ctx, cmd, root, logger, writer, env, map[string]string{}, false)
		if err != nil {
			logrus.Errorf("OnFail: an error during run cmd: %v, err: %v", cmd, err.Error())
			return err
		}
	}
	return nil
}

func (ctx *executionContext) getTestTimeout(task *testTask) time.Duration {
	timeout := time.Second * time.Duration(task.test.ExecutionConfig.Timeout) * 2
	if timeout == 0 {
		logrus.Infof("test timeout is not specified, use default value, 3min")
		timeout = time.Minute * 3
	}
	return timeout
}

func (ctx *executionContext) updateTestExecution(task *testTask, fileName string, status model.Status) {
	task.test.Status = status
	task.test.Executions = append(task.test.Executions, model.TestEntryExecution{
		Status:     status,
		Retry:      len(task.test.Executions) + 1,
		OutputFile: fileName,
	})
	ctx.operationChannel <- operationEvent{
		task: task,
		kind: eventTaskUpdate,
	}
}

func (ctx *executionContext) startCluster(ci *clusterInstance) bool {
	if ci.state != clusterAdded && ci.state != clusterCrashed {
		// no need to start
		return true
	}

	if ci.startCount > ci.group.config.RetryCount {
		logrus.Infof("Marking cluster %v as not available attempts reached: %v", ci.id, ci.group.config.RetryCount)
		ci.state = clusterNotAvailable
		return false
	}

	ci.state = clusterStarting
	execution := &clusterOperationRecord{
		time: time.Now(),
	}
	ci.executions = append(ci.executions, execution)
	go func() {
		timeout := ctx.getClusterTimeout(ci.group)
		ctx.Lock()
		ci.startCount++
		execution.attempt = ci.startCount
		ctx.Unlock()
		errFile, err := ci.instance.Start(timeout)
		if err != nil {
			execution.logFile = errFile
			execution.errMsg = err
			execution.status = clusterCrashed
			destroyErr := ctx.destroyCluster(ci, true, false)
			if destroyErr != nil {
				logrus.Errorf("Both start and destroy of cluster returned errors, stop retrying operations with this cluster %v", ci.instance)
				ci.startCount = ci.group.config.RetryCount + 1
				execution.status = clusterNotAvailable
			}
		} else {
			execution.status = clusterReady
		}
		execution.duration = time.Since(execution.time)
		// Starting cloud monitoring thread
		var state clusterState
		ctx.Lock()
		state = ci.state
		ctx.Unlock()
		if state != clusterCrashed {
			monitorContext, monitorCancel := context.WithCancel(context.Background())
			ci.cancelMonitor = monitorCancel
			ctx.monitorCluster(monitorContext, ci)
		} else {
			ctx.operationChannel <- operationEvent{
				kind:            eventClusterUpdate,
				clusterInstance: ci,
			}
		}
	}()
	return true
}

func (ctx *executionContext) getClusterTimeout(group *clustersGroup) time.Duration {
	timeout := time.Duration(group.config.Timeout) * time.Second
	if group.config.Timeout == 0 {
		logrus.Infof("cluster timeout is not specified, use default value 15min")
		timeout = 15 * time.Minute
	}
	return timeout
}

func (ctx *executionContext) monitorCluster(context context.Context, ci *clusterInstance) {
	checks := 0
	for {
		err := ci.instance.CheckIsAlive()
		if err != nil {
			logrus.Errorf("Failed to interact with %s: %v", ci.id, err)
			ctx.destroyCluster(ci, true, false)
			break
		}

		if checks == 0 {
			// Initial check performed, we need to make cluster ready.
			ctx.Lock()
			ci.state = clusterReady
			ci.startTime = time.Now()
			ctx.Unlock()
			ctx.operationChannel <- operationEvent{
				kind:            eventClusterUpdate,
				clusterInstance: ci,
			}
			logrus.Infof("Cluster instance started: %s", ci.id)
		}
		checks++
		select {
		case <-time.After(5 * time.Second):
			// Just pass
		case <-context.Done():
			logrus.Infof("cluster monitoring is canceled: %s. Uptime: %v seconds", ci.id, checks*5)
			return
		}
	}
}

func (ctx *executionContext) destroyCluster(ci *clusterInstance, sendUpdate, fork bool) error {

	if ci.state == clusterCrashed || ci.state == clusterNotAvailable || ci.state == clusterShutdown {
		// It is already destroyed or not available.
		return nil
	}
	ctx.Lock()
	ci.state = clusterStopping
	if ci.cancelMonitor != nil {
		ci.cancelMonitor()
	}
	ctx.Unlock()
	timeout := ctx.getClusterTimeout(ci.group)
	if fork {
		ctx.clusterWaitGroup.Add(1)
		go func() {
			defer ctx.clusterWaitGroup.Done()
			err := ci.instance.Destroy(timeout)
			if err != nil {
				logrus.Errorf("Failed to destroy cluster")
			}
		}()
		return nil
	}
	err := ci.instance.Destroy(timeout)
	if err != nil {
		logrus.Errorf("Failed to destroy cluster")
	}

	if ci.group.config.StopDelay != 0 {
		logrus.Infof("Cluster stop warm-up timeout specified %v", ci.group.config.StopDelay)
		<-time.After(time.Duration(ci.group.config.StopDelay) * time.Second)
	}
	ctx.Lock()
	ci.state = clusterCrashed
	ctx.Unlock()
	if sendUpdate {
		ctx.operationChannel <- operationEvent{
			clusterInstance: ci,
			kind:            eventClusterUpdate,
		}
	}
	return err
}

func (ctx *executionContext) createClusters() error {
	ctx.clusters = []*clustersGroup{}
	clusterProviders, err := createClusterProviders(ctx.manager)
	if err != nil {
		return err
	}

	for _, cl := range ctx.cloudTestConfig.Providers {
		if ctx.arguments.onlyEnabled {
			logrus.Infof("Disable cluster config:: %v since onlyEnabled is passed...", cl.Name)
			cl.Enabled = false
		}
		for _, cc := range ctx.arguments.clusters {
			if cl.Name == cc {
				if !cl.Enabled {
					logrus.Infof("Enabling config:: %v", cl.Name)
				}
				cl.Enabled = true
			}
		}
		if cl.Enabled {
			logrus.Infof("Initialize provider for config:: %v %v", cl.Name, cl.Kind)
			provider, ok := clusterProviders[cl.Kind]
			if !ok {
				msg := fmt.Sprintf("Cluster provider %s are not found...", cl.Kind)
				logrus.Errorf(msg)
				return errors.New(msg)
			}
			instances := []*clusterInstance{}
			group := &clustersGroup{
				provider:  provider,
				config:    cl,
				tasks:     map[string]*testTask{},
				completed: map[string]*testTask{},
			}
			for i := 0; i < cl.Instances; i++ {
				cluster, err := provider.CreateCluster(cl, ctx.factory, ctx.manager, ctx.arguments.instanceOptions)
				if err != nil {
					msg := fmt.Sprintf("Failed to create cluster instance. Error %v", err)
					logrus.Errorf(msg)
					return errors.New(msg)
				}
				instances = append(instances, &clusterInstance{
					instance:  cluster,
					startTime: time.Now(),
					state:     clusterAdded,
					id:        cluster.GetID(),
					group:     group,
				})
			}
			group.instances = instances
			if len(instances) == 0 {
				msg := fmt.Sprintf("No instances are specified for %s.", cl.Name)
				logrus.Errorf(msg)
				return errors.New(msg)
			}
			ctx.clusters = append(ctx.clusters, group)
		}
	}
	if len(ctx.clusters) == 0 {
		msg := "there is no clusters defined. Exiting"
		logrus.Errorf(msg)
		return errors.New(msg)
	}
	return nil
}

func (ctx *executionContext) findTests() error {
	logrus.Infof("Finding tests")
	for _, exec := range ctx.cloudTestConfig.Executions {
		if exec.Name == "" {
			return errors.New("execution name should be specified")
		}
		if exec.Kind == "" || exec.Kind == "gotest" {
			tests, err := ctx.findGoTest(exec)
			if err != nil {
				return err
			}
			if len(tests) > 0 {
				ctx.tests = append(ctx.tests, tests...)
			}
		} else if exec.Kind == "shell" {
			tests := ctx.findShellTest(exec)
			if len(tests) > 0 {
				ctx.tests = append(ctx.tests, tests...)
			}
		} else {
			return errors.Errorf("unknown executon kind %v", exec.Kind)
		}
	}
	// If we have execution without tags, we need to remove all tests from it from tagged executions.
	logrus.Infof("Total tests found: %v", len(ctx.tests))
	if len(ctx.tests) == 0 {
		return errors.New("there is no tests defined")
	}
	return nil
}

func (ctx *executionContext) findShellTest(exec *config.ExecutionConfig) []*model.TestEntry {
	return []*model.TestEntry{
		{
			Name:            exec.Name,
			Kind:            model.TestEntryKindShellTest,
			Tags:            "",
			ExecutionConfig: exec,
			Status:          model.StatusAdded,
			RunScript:       exec.Run,
		},
	}
}

func (ctx *executionContext) findGoTest(executionConfig *config.ExecutionConfig) ([]*model.TestEntry, error) {
	st := time.Now()
	logrus.Infof("Starting finding tests by tags %v", executionConfig.Tags)
	execTests, err := model.GetTestConfiguration(ctx.manager, executionConfig.PackageRoot, executionConfig.Tags)
	if err != nil {
		logrus.Errorf("Failed during test lookup %v", err)
		return nil, err
	}
	logrus.Infof("Tests found: %v Elapsed: %v", len(execTests), time.Since(st))
	result := []*model.TestEntry{}
	for _, t := range execTests {
		t.Kind = model.TestEntryKindGoTest
		t.ExecutionConfig = executionConfig
		match := true
		for _, v := range executionConfig.OnlyRun {
			match = false
			if v == t.Name {
				match = true
				break
			}
		}
		if match {
			result = append(result, t)
		}
	}
	return result, nil
}

func (ctx *executionContext) generateJUnitReportFile() (*reporting.JUnitFile, error) {
	// generate and write report
	ctx.report = &reporting.JUnitFile{}

	totalFailures := 0
	for _, cluster := range ctx.clusters {
		failures := 0
		totalTests := 0
		totalTime := time.Duration(0)
		suite := &reporting.Suite{
			Name: cluster.config.Name,
		}

		for _, test := range cluster.tasks {
			totalTests, totalTime, failures = ctx.generateTestCaseReport(test, totalTests, totalTime, failures, suite)
		}

		for _, test := range cluster.completed {
			totalTests, totalTime, failures = ctx.generateTestCaseReport(test, totalTests, totalTime, failures, suite)
		}

		// Check cluster executions.
		availableClusters := 0
		for _, inst := range cluster.instances {
			if inst.state != clusterNotAvailable {
				availableClusters++
			}
		}
		if availableClusters == 0 {
			// No clusters available let's mark this as error.
			for _, inst := range cluster.instances {
				if inst.state == clusterNotAvailable {
					for _, exec := range inst.executions {
						ctx.generateClusterFailedReportEntry(inst, exec, suite)
						failures++
						totalTests++
						break
					}
				}
			}
		}

		suite.Tests = totalTests
		suite.Failures = failures
		suite.Time = fmt.Sprintf("%v", totalTime)
		totalFailures += failures

		ctx.report.Suites = append(ctx.report.Suites, suite)
	}

	output, err := xml.MarshalIndent(ctx.report, "  ", "    ")
	if err != nil {
		logrus.Errorf("failed to store JUnit xml report: %v\n", err)
	}
	if ctx.cloudTestConfig.Reporting.JUnitReportFile != "" {
		ctx.manager.AddFile(ctx.cloudTestConfig.Reporting.JUnitReportFile, output)
	}
	if totalFailures > 0 {
		return ctx.report, errors.Errorf("there is failed tests %v", totalFailures)
	}
	return ctx.report, nil
}

func (ctx *executionContext) generateClusterFailedReportEntry(inst *clusterInstance, exec *clusterOperationRecord, suite *reporting.Suite) {
	message := fmt.Sprintf("Cluster start failed %v", inst.id)
	result := fmt.Sprintf("Error: %v\n", exec.errMsg)
	if exec.logFile != "" {
		lines, err := utils.ReadFile(exec.logFile)
		if err == nil {
			// No file
			result += strings.Join(lines, "\n")
		}
	}
	startCase := &reporting.TestCase{
		Name: fmt.Sprintf("Startup-%v", inst.id),
		Time: fmt.Sprintf("%v", exec.duration),
	}
	startCase.Failure = &reporting.Failure{
		Type:     "ERROR",
		Contents: result,
		Message:  message,
	}
	suite.TestCases = append(suite.TestCases, startCase)
}

func (ctx *executionContext) generateTestCaseReport(test *testTask, totalTests int, totalTime time.Duration, failures int, suite *reporting.Suite) (int, time.Duration, int) {
	testCase := &reporting.TestCase{
		Name: test.test.Key,
		Time: test.test.Duration.String(),
	}

	totalTests++
	totalTime += test.test.Duration
	switch test.test.Status {
	case model.StatusFailed, model.StatusTimeout:
		failures++

		message := fmt.Sprintf("Test execution failed %v", test.test.Name)
		result := strings.Builder{}
		for idx, ex := range test.test.Executions {
			lines, err := utils.ReadFile(ex.OutputFile)
			if err != nil {
				logrus.Errorf("Failed to read stored output %v", ex.OutputFile)
				lines = []string{"Failed to read stored output:", ex.OutputFile, err.Error()}
			}
			result.WriteString(fmt.Sprintf("Execution attempt: %v Output file: %v", idx, ex.OutputFile))
			result.WriteString(strings.Join(lines, "\n"))
		}
		testCase.Failure = &reporting.Failure{
			Type:     "ERROR",
			Contents: result.String(),
			Message:  message,
		}
	case model.StatusSkipped:
		msg := "By limit of number of tests to run"
		if test.test.SkipMessage != "" {
			msg = test.test.SkipMessage
		}

		testCase.SkipMessage = &reporting.SkipMessage{
			Message: msg,
		}
	case model.StatusSkippedSinceNoClusters:
		testCase.SkipMessage = &reporting.SkipMessage{
			Message: "No clusters are available, all clusters reached restart limits...",
		}
	}
	suite.TestCases = append(suite.TestCases, testCase)
	return totalTests, totalTime, failures
}

func (ctx *executionContext) checkClustersUsage() {
	for _, ci := range ctx.clusters {
		if len(ci.tasks) == 0 {
			up := 0
			for _, inst := range ci.instances {
				if !ctx.isClusterDown(inst) {
					up++
				}
			}
			if up > 0 {
				logrus.Infof("All tasks for cluster group %v are complete. Starting cluster shutdown.", ci.config.Name)
				for _, inst := range ci.instances {
					if !ctx.isClusterDown(inst) && inst.state != clusterBusy {
						ctx.destroyCluster(inst, false, true)
						ctx.Lock()
						inst.state = clusterShutdown
						ctx.Unlock()
					}
				}
			}
		}
	}
}

func (ctx *executionContext) isClusterDown(inst *clusterInstance) bool {
	return inst.state == clusterShutdown || inst.state == clusterCrashed || inst.state == clusterNotAvailable
}

func createClusterProviders(manager execmanager.ExecutionManager) (map[string]providers.ClusterProvider, error) {
	clusterProviders := map[string]providers.ClusterProvider{}

	clusterProviderFactories := map[string]providers.ClusterProviderFunction{
		"packet": packet.NewPacketClusterProvider,
		"shell":  shell.NewShellClusterProvider,
	}

	for key, factory := range clusterProviderFactories {
		if _, ok := clusterProviders[key]; ok {
			msg := fmt.Sprintf("Re-definition of cluster provider... Exiting")
			logrus.Errorf(msg)
			return nil, errors.New(msg)
		}
		root, err := manager.GetRoot(key)
		if err != nil {
			logrus.Errorf("Failed to create cluster provider %v", err)
			return nil, err
		}
		clusterProviders[key] = factory(root)
	}
	return clusterProviders, nil
}

type cloudTestCmd struct {
	cobra.Command

	cmdArguments *Arguments
}

// ExecuteCloudTest - main entry point for command
func ExecuteCloudTest() {
	var rootCmd = &cloudTestCmd{
		cmdArguments: &Arguments{
			providerConfig: defaultConfigFile,
			clusters:       []string{},
		},
	}
	rootCmd.Use = "cloud_test"
	rootCmd.Short = "NSM Cloud Test is cloud helper continuous integration testing tool"
	rootCmd.Long = `Allow to execute all set of individual tests across all clouds provided.`
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		CloudTestRun(rootCmd)
	}
	rootCmd.Args = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	initCmd(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCmd(rootCmd *cloudTestCmd) {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&rootCmd.cmdArguments.providerConfig, "config", "", "", "Config file for providers, default="+defaultConfigFile)
	rootCmd.Flags().StringArrayVarP(&rootCmd.cmdArguments.clusters, "clusters", "c", []string{}, "Enable disable cluster configs, default use from config. Cloud be used to test against selected configuration or locally...")
	rootCmd.Flags().BoolVarP(&rootCmd.cmdArguments.onlyEnabled, "enabled", "e", false, "Use only passed cluster names...")
	rootCmd.Flags().IntVarP(&rootCmd.cmdArguments.count, "count", "", -1, "Execute only count of tests")

	rootCmd.Flags().BoolVarP(&rootCmd.cmdArguments.instanceOptions.NoStop, "noStop", "", false, "Pass to disable stop operations...")
	rootCmd.Flags().BoolVarP(&rootCmd.cmdArguments.instanceOptions.NoInstall, "noInstall", "", false, "Pass to disable do install operations...")
	rootCmd.Flags().BoolVarP(&rootCmd.cmdArguments.instanceOptions.NoPrepare, "noPrepare", "", false, "Pass to disable do prepare operations...")
	rootCmd.Flags().BoolVarP(&rootCmd.cmdArguments.instanceOptions.NoMaskParameters, "noMask", "", false, "Pass to disable masking of environment variables...")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of cloud_test",
		Long:  `All software has versions.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Cloud Test -- HEAD")
		},
	}
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
}
