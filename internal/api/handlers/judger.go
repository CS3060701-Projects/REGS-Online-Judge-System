package handlers

import "fmt"

type JudgeJob struct {
	OperatorID string
	Workspace  string
	ProblemID  string
}

var JobQueue chan JudgeJob

func InitJudger(maxWorkers int) {
	JobQueue = make(chan JudgeJob, 1000)

	for i := 1; i <= maxWorkers; i++ {
		go worker(i)
	}
	fmt.Printf("評測系統初始化成功，目前配置 %d 個 Worker\n", maxWorkers)
}

func worker(workerID int) {
	for job := range JobQueue {
		fmt.Printf("[Worker %d] 開始處理任務: %s\n", workerID, job.OperatorID)

		processSubmission(job.OperatorID, job.Workspace, job.ProblemID)

		fmt.Printf("[Worker %d] 完成任務: %s\n", workerID, job.OperatorID)
	}
}
