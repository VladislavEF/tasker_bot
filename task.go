package main

import (
	"crypto/sha1"
	"time"

	"github.com/google/uuid"
)

type TaskInfo struct {
	name         string
	id           string
	status       TaskStatus
	dependens    []string
	creationDate time.Time
	finishDate   time.Time
	// daysInWork int

}

type TaskStatus int

const (
	Backlog = iota
	Done
	Cancelled
)

func MakeId(name string) string {
	md5Hash := sha1.Sum([]byte(name))
	result, _ := uuid.FromBytes(md5Hash[:])
	// if err != nil {
	// 	log.WithError(err).Fatal("Fail to create task")
	// }
	return result.String()
}

func NewTask(name string) *TaskInfo {
	task := TaskInfo{
		name:         name,
		id:           MakeId(name),
		status:       Backlog,
		dependens:    make([]string, 0),
		creationDate: time.Now(),
	}

	return &task
}

func (this *TaskInfo) Done() {
	this.status = Done
	this.finishDate = time.Now()
}
func (this *TaskInfo) Cancelled() {
	this.status = Cancelled
}
func (this *TaskInfo) Backlog() {
	this.status = Backlog
	this.finishDate = time.Time{}
}
func (this *TaskInfo) BoundWith(id string) {
	this.dependens = append(this.dependens, id)
}

func (this *TaskInfo) GetId() (id string) {
	return this.id
}
func (this *TaskInfo) GetName() (id string) {
	return this.name
}
