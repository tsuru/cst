package queue

import (
	"time"

	"github.com/tsuru/monsterqueue"
)

// MockQueue implements a monsterqueue.Queue interface for testing purposes.
type MockQueue struct {
	MockRegisterTask func(monsterqueue.Task) error
	MockEnqueue      func(string, monsterqueue.JobParams) (monsterqueue.Job, error)
	MockEnqueueWait  func(string, monsterqueue.JobParams, time.Duration) (monsterqueue.Job, error)
	MockProcessLoop  func()
	MockStop         func()
	MockWait         func()
	MockRetrieveJob  func(string) (monsterqueue.Job, error)
	MockResetStorage func() error
	MockListJobs     func() ([]monsterqueue.Job, error)
	MockDeleteJob    func(string) error
}

// RegisterTask is a mock implementation for testing purposes.
func (mq MockQueue) RegisterTask(task monsterqueue.Task) error {

	if mq.MockRegisterTask != nil {
		return mq.MockRegisterTask(task)
	}

	return nil
}

// Enqueue is a mock implementation for testing purposes.
func (mq MockQueue) Enqueue(tn string, params monsterqueue.JobParams) (monsterqueue.Job, error) {

	if mq.MockEnqueue != nil {
		return mq.MockEnqueue(tn, params)
	}

	return nil, nil
}

// EnqueueWait is a mock implementation for testing purposes.
func (mq MockQueue) EnqueueWait(tn string, params monsterqueue.JobParams, timeout time.Duration) (monsterqueue.Job, error) {

	if mq.MockEnqueueWait != nil {
		return mq.MockEnqueueWait(tn, params, timeout)
	}

	return nil, nil
}

// ProcessLoop is a mock implementation for testing purposes.
func (mq MockQueue) ProcessLoop() {

	if mq.MockProcessLoop != nil {
		mq.MockProcessLoop()
	}
}

// Stop is a mock implementation for testing purposes.
func (mq MockQueue) Stop() {

	if mq.MockStop != nil {
		mq.MockStop()
	}
}

// Wait is a mock implementation for testing purposes.
func (mq MockQueue) Wait() {

	if mq.MockWait != nil {
		mq.MockWait()
	}
}

// RetrieveJob is a mock implementation for testing purposes.
func (mq MockQueue) RetrieveJob(job string) (monsterqueue.Job, error) {

	if mq.MockRetrieveJob != nil {
		return mq.MockRetrieveJob(job)
	}

	return nil, nil
}

// ResetStorage is a mock implementation for testing purposes.
func (mq MockQueue) ResetStorage() error {

	if mq.MockResetStorage != nil {
		return mq.MockResetStorage()
	}

	return nil
}

// ListJobs is a mock implementation for testing purposes.
func (mq MockQueue) ListJobs() ([]monsterqueue.Job, error) {

	if mq.MockListJobs != nil {
		return mq.MockListJobs()
	}

	return []monsterqueue.Job{}, nil
}

// DeleteJob is a mock implementation for testing purposes.
func (mq MockQueue) DeleteJob(job string) error {

	if mq.MockDeleteJob != nil {
		return mq.MockDeleteJob(job)
	}

	return nil
}

// MockJob implements a monsterqueue.Job interface for testing purposes.
type MockJob struct {
	MockSucess       func(monsterqueue.JobResult) (bool, error)
	MockError        func(error) (bool, error)
	MockResult       func() (monsterqueue.JobResult, error)
	MockID           func() string
	MockParameters   func() monsterqueue.JobParams
	MockTaskName     func() string
	MockQueue        func() monsterqueue.Queue
	MockStatus       func() monsterqueue.JobStatus
	MockEnqueueStack func() string
}

// Success is a mock implementation for testing purposes.
func (mj MockJob) Success(result monsterqueue.JobResult) (bool, error) {

	if mj.MockSucess != nil {
		return mj.MockSucess(result)
	}

	return false, nil
}

// Error is a mock implementation for testing purposes.
func (mj MockJob) Error(jobErr error) (bool, error) {

	if mj.MockError != nil {
		return mj.MockError(jobErr)
	}

	return false, nil
}

// Result is a mock implementation for testing purposes.
func (mj MockJob) Result() (monsterqueue.JobResult, error) {

	if mj.MockResult != nil {
		return mj.MockResult()
	}

	return nil, nil
}

// Parameters is a mock implementation for testing purposes.
func (mj MockJob) Parameters() monsterqueue.JobParams {

	if mj.MockParameters != nil {
		return mj.MockParameters()
	}

	return nil
}

// ID is a mock implementation for testing purposes.
func (mj MockJob) ID() string {

	if mj.MockID != nil {
		return mj.MockID()
	}

	return ""
}

// Queue is a mock implementation for testing purposes.
func (mj MockJob) Queue() monsterqueue.Queue {

	if mj.MockQueue != nil {
		return mj.MockQueue()
	}

	return nil
}

// TaskName is a mock implementation for testing purposes.
func (mj MockJob) TaskName() string {

	if mj.MockTaskName != nil {
		return mj.MockTaskName()
	}

	return ""
}

// Status is a mock implementation for testing purposes.
func (mj MockJob) Status() monsterqueue.JobStatus {

	if mj.MockStatus != nil {
		return mj.MockStatus()
	}

	return monsterqueue.JobStatus{}
}

// EnqueueStack is a mock implementation for testing purposes.
func (mj MockJob) EnqueueStack() string {

	if mj.MockEnqueueStack != nil {
		return mj.MockEnqueueStack()
	}

	return ""
}
