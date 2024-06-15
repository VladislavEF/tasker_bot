package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MemoryStorage struct {
	// UserName on UserId
	UserIds map[int64]string
	// UserId on TaskIds
	TaskList map[int64][]string
	// TaskId on TaskInfo
	Tasks map[string]TaskInfo
	// UserId on current Line
	Lines map[int64]Line
	path  string
}

type IDatabase interface {
	GetAllUsers() []int64
	GetUserId(userName string) int64
	GetUserTasks(userId int64) []string
	GetTaskInfo(taskId string) TaskInfo
	GetLine(userId int64) *Line

	IsTask(taskId string) bool
	IsOpenLine(userId int64) bool

	AddUser(userName string, userId int64)
	AddUserTask(userId int64, taskId string)
	AddNewTask(task TaskInfo)

	ChangeLine(userId int64, line Line)

	Save() error
	GetDbState()
}

// func GetDatabase() Database {
// 	return Database{}
// }

func (this *MemoryStorage) GetAllUsers() []int64 {
	users := make([]int64, len(this.UserIds))
	for id, _ := range this.UserIds {
		users = append(users, id)
	}
	return users
}
func (this *MemoryStorage) GetUserId(id string) int64 {
	return 0
}
func (this *MemoryStorage) GetUserTasks(userId int64) []string {
	taskIds := this.TaskList[userId]
	tasks := make([]string, 0)
	for _, taskId := range taskIds {
		task := this.Tasks[taskId]
		tasks = append(tasks, task.Name)
	}
	return tasks
}
func (this *MemoryStorage) GetTaskInfo(taskId string) TaskInfo {
	return TaskInfo{}
}
func (this *MemoryStorage) GetLine(userId int64) *Line {
	if _, ok := this.Lines[userId]; ok {
		line := this.Lines[userId]
		return &line
	} else {
		return nil
	}
}

func (this *MemoryStorage) IsTask(taskId string) bool {
	_, ok := this.Tasks[taskId]
	return ok
}
func (this *MemoryStorage) IsOpenLine(userId int64) bool {
	line, ok := this.Lines[userId]
	return ok && line.Command != ""
}

func (this *MemoryStorage) AddUser(userName string, userId int64) {
	if this.UserIds == nil {
		this.UserIds = make(map[int64]string)
	}
	this.UserIds[userId] = userName
}
func (this *MemoryStorage) AddUserTask(userId int64, taskId string) {
	if this.TaskList == nil {
		this.TaskList = make(map[int64][]string)
		this.TaskList[userId] = []string{}
	}
	tasks := append(this.TaskList[userId], taskId)
	this.TaskList[userId] = tasks
}
func (this *MemoryStorage) AddNewTask(task TaskInfo) {
	if this.Tasks == nil {
		this.Tasks = make(map[string]TaskInfo)
	}
	this.Tasks[task.Id] = task
}

func (this *MemoryStorage) ChangeLine(userId int64, line Line) {
	if this.Lines == nil {
		this.Lines = make(map[int64]Line)
	}
	if _, ok := this.Lines[userId]; !ok {
		this.Lines[userId] = line
	} else {
		delete(this.Lines, userId)
	}
}

func (this *MemoryStorage) Save() error {
	data, err := json.MarshalIndent(&this, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(this.path, data, 0644)
}

func (this *MemoryStorage) GetDbState() {
	fmt.Println("DB state:")
	fmt.Println(this)
}

func GetLocalDatabase(path string) (IDatabase, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	file.Close()
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	db := &MemoryStorage{}
	if len(bytes) == 0 {
		db.path = path
	} else {
		err = json.Unmarshal(bytes, db)
		db.path = path
	}

	return IDatabase(db), err
}
