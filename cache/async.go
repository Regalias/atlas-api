package cache

import (
	"errors"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type taskop int

const maxRetryCount = 3

const (
	SetLink taskop = iota
	RemoveLink
)

// Task containers the cache operation request details
type Task struct {
	Operation taskop
	Linkpath  string
	Linkdest  string
}

// AsyncHandler contains the context for the queue and worker
type AsyncHandler struct {
	taskQueue chan *Task
	logger    *zerolog.Logger
	cache     Provider
}

// NewAsyncQueue creates a new task queue with the specified queue size
func NewAsyncQueue(queueSize int, logger *zerolog.Logger, cacheprov Provider) *AsyncHandler {
	ch := make(chan *Task, queueSize)
	return &AsyncHandler{
		taskQueue: ch,
		logger:    logger,
		cache:     cacheprov,
	}
}

// SubmitTask puts a new cache operation request into the queue
func (th *AsyncHandler) SubmitTask(ct *Task) error {

	for retryCount := 0; retryCount < maxRetryCount; retryCount++ {
		select {
		case th.taskQueue <- ct:
			// do something
			return nil
		default:
			// Queue is full!
			// Retry after a second
			time.Sleep(time.Second)
		}
	}
	return errors.New("TaskSubmitFailed")
}

// TODO: retry?

// RunWorker runs the async worker thread
func (th *AsyncHandler) RunWorker() {
	for {
		t := <-th.taskQueue
		// th.logger.Debug().Msg("Got task: " + strconv.Itoa(int(t.operation)) + t.linkpath + t.linkdest)
		switch t.Operation {
		case SetLink:
			// TODO: query the link to check if it exists before setting?
			if err := th.cache.UpsertLink(t.Linkpath, t.Linkdest); err != nil {
				th.logError("Couldn't set link in cache: " + err.Error())
			}
		case RemoveLink:
			// remove the link from cache
			if err := th.cache.DeleteLink(t.Linkpath); err != nil {
				th.logError("Couldn't delete link in cache: " + err.Error())
			}
		default:
			// unknown op
			th.logError("Unknown operation: " + strconv.Itoa(int(t.Operation)))
		}
	}
}

func (th *AsyncHandler) logError(errmsg string) {
	th.logger.Error().Str("Segment", "TaskWorker").Msg(errmsg)
}
