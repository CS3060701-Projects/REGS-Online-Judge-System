package handlers

import (
	"fmt"
)

type JudgeJob struct {
	OperatorID string
	Workspace  string
	ProblemID  string
}

var (
	JobQueue  chan JudgeJob
	semaphore chan struct{}
)

func InitJudger(maxWorkers int) {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	JobQueue = make(chan JudgeJob, 1000)
	semaphore = make(chan struct{}, maxWorkers)

	for i := 1; i <= maxWorkers; i++ {
		go worker(i)
	}
	fmt.Printf("評測系統初始化成功，目前配置 %d 個 Worker（Semaphore 併發上限）\n", maxWorkers)
}

func acquireSemaphore() {
	semaphore <- struct{}{}
}

func releaseSemaphore() {
	<-semaphore
}

func worker(workerID int) {
	for job := range JobQueue {
		acquireSemaphore()
		fmt.Printf("[Worker %d] 開始處理任務: %s\n", workerID, job.OperatorID)
		processSubmission(job.OperatorID, job.Workspace, job.ProblemID)
		fmt.Printf("[Worker %d] 完成任務: %s\n", workerID, job.OperatorID)
		releaseSemaphore()
	}
}
