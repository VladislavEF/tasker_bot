package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MemoryStorage struct {
	// UserName on UserId
	userIds map[string]int64
	// UserId on TaskIds
	taskList map[int64][]string
	// TaskId on TaskInfo
	tasks map[string]TaskInfo
	// UserId on current Line
	lines map[int64]Line
	path  string
}

type Database interface {
	GetUserId(userName string) int64
	GetUserTasks(userId int64) []string
	GetTaskInfo(taskId string) TaskInfo
	GetLine(userId int64) *Line

	AddUser(userName string, userId int64)
	AddUserTask(userId int64, taskId string)
	AddNewTask(task TaskInfo)

	IsOpenLine(userId int64) bool
	ChangeLine(userId int64, line Line)

	Save() error
}

// func GetDatabase() Database {
// 	return Database{}
// }

func (this *MemoryStorage) GetUserId(id string) int64 {
	return 0
}
func (this *MemoryStorage) GetUserTasks(userId int64) []string {
	return []string{""}
}
func (this *MemoryStorage) GetTaskInfo(taskId string) TaskInfo {
	return TaskInfo{}
}
func (this *MemoryStorage) GetLine(userId int64) *Line {
	if _, ok := this.lines[userId]; ok {
		line := this.lines[userId]
		return &line
	} else {
		return nil
	}
}

func (this *MemoryStorage) AddUser(userName string, userId int64) {
	if this.userIds == nil {
		this.userIds = make(map[string]int64)
	}
}
func (this *MemoryStorage) AddUserTask(userId int64, taskId string) {
	if this.taskList == nil {
		this.taskList = make(map[int64][]string)
	}
}
func (this *MemoryStorage) AddNewTask(task TaskInfo) {
	if this.tasks == nil {
		this.tasks = make(map[string]TaskInfo)
	}
}
func (this *MemoryStorage) IsOpenLine(userId int64) bool {
	line, ok := this.lines[userId]
	return ok && line.command != ""
}
func (this *MemoryStorage) ChangeLine(userId int64, line Line) {
	if this.lines == nil {
		this.lines = make(map[int64]Line)
	}
	if _, ok := this.lines[userId]; !ok {
		this.lines[userId] = line
	} else {
		delete(this.lines, userId)
	}
}

func (this *MemoryStorage) Save() error {
	data, err := json.MarshalIndent(this, "", "	")
	if err != nil {
		return err
	}
	fmt.Println(this)
	return os.WriteFile(this.path, data, 0644)
}

func GetLocalDatabase(path string) (Database, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	rawFile := []byte{}
	_, err = file.Read(rawFile)
	if err != nil {
		return nil, err
	}
	db := &MemoryStorage{}
	if len(rawFile) == 0 {
		db.path = path
	} else {
		err = json.Unmarshal(rawFile, db)
	}

	return Database(db), err
}
