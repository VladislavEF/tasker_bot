package main

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"
)

type TaskInfo struct {
	Name         string
	Id           string
	Status       TaskStatus
	Dependens    []string
	CreationDate time.Time
	FinishDate   time.Time
	// daysInWork int
}

type TaskStatus int

const (
	Backlog = iota
	Done
	Cancelled
)

func MakeId(name string) string {
	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}

func NewTask(name string) *TaskInfo {
	task := TaskInfo{
		Name:         name,
		Id:           MakeId(strings.ToUpper(name)),
		Status:       Backlog,
		Dependens:    make([]string, 0),
		CreationDate: time.Now(),
	}

	return &task
}

func (this *TaskInfo) Done() {
	this.Status = Done
	this.FinishDate = time.Now()
}

func (this *TaskInfo) Cancelled() {
	this.Status = Cancelled
	this.FinishDate = time.Now()
}

func (this *TaskInfo) Backlog() {
	this.Status = Backlog
	this.FinishDate = time.Time{}
}

func (this *TaskInfo) BoundWith(id string) {
	this.Dependens = append(this.Dependens, id)
}

func (this *TaskInfo) GetId() (id string) {
	return this.Id
}

func (this *TaskInfo) GetName() (id string) {
	return this.Name
}
