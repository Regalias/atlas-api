package apiserver

import (
	"errors"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type taskop int

const maxRetryCount = 3

const (
	setLink taskop = iota
	removeLink
)

type cacheTask struct {
	operation taskop
	linkpath  string
	linkdest  string
}

type cacheTaskHandler struct {
	taskQueue chan *cacheTask
	logger    *zerolog.Logger
	cache     *cacheProvider
}

// newTaskQueue creates a new task queue with the specified queue size
func newTaskQueue(queueSize int, logger *zerolog.Logger, cacheprov *cacheProvider) *cacheTaskHandler {
	ch := make(chan *cacheTask, queueSize)
	return &cacheTaskHandler{
		taskQueue: ch,
		logger:    logger,
		cache:     cacheprov,
	}
}

func (th *cacheTaskHandler) submitTask(ct *cacheTask) error {

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
func (th *cacheTaskHandler) runWorker() {
	for {
		t := <-th.taskQueue
		// th.logger.Debug().Msg("Got task: " + strconv.Itoa(int(t.operation)) + t.linkpath + t.linkdest)
		switch t.operation {
		case setLink:
			// TODO: query the link to check if it exists before setting?
			if err := th.cache.upsertLink(t.linkpath, t.linkdest); err != nil {
				th.logError("Couldn't set link in cache: " + err.Error())
			}
		case removeLink:
			// remove the link from cache
			if err := th.cache.deleteLink(t.linkpath); err != nil {
				th.logError("Couldn't delete link in cache: " + err.Error())
			}
		default:
			// unknown op
			th.logError("Unknown operation: " + strconv.Itoa(int(t.operation)))
		}
	}
}

func (th *cacheTaskHandler) logError(errmsg string) {
	th.logger.Error().Str("Segment", "TaskWorker").Msg(errmsg)
}
