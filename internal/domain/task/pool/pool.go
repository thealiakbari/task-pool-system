package pool

import (
	"context"
	"errors"
	"github.com/thealiakbari/task-pool-system/internal/domain/task/entity"
	"github.com/thealiakbari/task-pool-system/internal/ports/inbound/task"
	"github.com/thealiakbari/task-pool-system/pkg/common/db"
	"log"
	"sync"
	"time"
)

var ErrPoolFull = errors.New("task pool is full")

type Pool struct {
	ctx     context.Context
	cancel  context.CancelFunc
	queue   chan *entity.Task
	workers int
	wg      sync.WaitGroup
}

func New(
	parent context.Context,
	workers int,
	poolSize int,
) *Pool {
	ctx, cancel := context.WithCancel(parent)

	return &Pool{
		ctx:     ctx,
		cancel:  cancel,
		queue:   make(chan *entity.Task, poolSize),
		workers: workers,
	}
}

type WorkerDeps struct {
	TaskService task.TaskService
}

func (p *Pool) Start(deps WorkerDeps) {
	log.Printf("[POOL] starting %d workers", p.workers)

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i+1, deps)
	}
}

func (p *Pool) Submit(task *entity.Task) error {
	select {
	case p.queue <- task:
		log.Printf("[POOL] task submitted: %s", task.Id)
		return nil
	default:
		return ErrPoolFull
	}
}

func (p *Pool) Shutdown() {
	log.Println("[POOL] shutdown initiated")
	p.cancel()
	p.wg.Wait()
	log.Println("[POOL] shutdown completed")
}

func (p *Pool) worker(id int, deps WorkerDeps) {
	defer p.wg.Done()

	log.Printf("[WORKER-%d] started", id)

	for {
		select {
		case <-p.ctx.Done():
			log.Printf("[WORKER-%d] stopping", id)
			return

		case task := <-p.queue:
			p.processTask(id, *task, deps)
		}
	}
}

func (p *Pool) processTask(
	workerID int,
	task entity.Task,
	deps WorkerDeps,
) {
	log.Printf("[WORKER-%d] start task %s", workerID, task.Id)

	taskModel, err := start(task)
	if err != nil {
		return
	}
	_, _ = deps.TaskService.Update(context.Background(), entity.Task{
		UniversalModel: db.UniversalModel{
			Id:        taskModel.Id,
			UpdatedAt: time.Now(),
		},
		Status:   taskModel.Status,
		Duration: taskModel.Duration,
	})

	select {
	case <-p.ctx.Done():
		fail(task)
	case <-time.After(task.Duration):
		complete(task)
	}

	_, _ = deps.TaskService.Update(context.Background(), entity.Task{
		UniversalModel: db.UniversalModel{
			Id:        taskModel.Id,
			UpdatedAt: time.Now(),
		},
		Status:   taskModel.Status,
		Duration: taskModel.Duration,
	})

	log.Printf("[WORKER-%d] finished task %s", workerID, task.Id)
}

var ErrInvalidStateTransition = errors.New("invalid task state transition")

func start(task entity.Task) (entity.Task, error) {
	if task.Status != entity.StatusPending {
		return entity.Task{}, ErrInvalidStateTransition
	}

	task.Status = entity.StatusRunning
	task.UpdatedAt = time.Now()
	return task, nil
}

func complete(task entity.Task) entity.Task {
	task.Status = entity.StatusCompleted
	task.UpdatedAt = time.Now()
	return task
}

func fail(task entity.Task) entity.Task {
	task.Status = entity.StatusFailed
	task.UpdatedAt = time.Now()
	return task
}
