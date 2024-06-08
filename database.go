package main

import (
	"encoding/json"
	"os"
)

type MemoryStorage struct {
	// UserName on UserId
	userIds map[string]string
	// UserId on TaskIds
	taskList map[string][]string
	// TaskId on TaskInfo
	tasks map[string]TaskInfo
}

type Database interface {
	GetUserId(string) string
	GetUserTasks(string) []string
	GetTaskInfo(string) TaskInfo

	AddUser(string)
	AddUserTask(string, string)
	AddNewTask(TaskInfo)

	Save() error
}

func GetDatabase() *Database {
	return new(Database)
}

func (this *MemoryStorage) GetUserId(id string) string {
	return ""
}
func (this *MemoryStorage) GetUserTasks(id string) []string {
	return []string{""}
}
func (this *MemoryStorage) GetTaskInfo(id string) TaskInfo {
	return TaskInfo{}
}

func (this *MemoryStorage) AddUser(name string) {
	if this.userIds == nil {
		this.userIds = make(map[string]string)
	}
}
func (this *MemoryStorage) AddUserTask(userId, taskId string) {
	if this.taskList == nil {
		this.taskList = make(map[string][]string)
	}
}
func (this *MemoryStorage) AddNewTask(task TaskInfo) {
	if this.tasks == nil {
		this.tasks = make(map[string]TaskInfo)
	}
}

func (this *MemoryStorage) Save() error {
	data, err := json.MarshalIndent(this, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile("database.json", data, 0644)
}

func GetLocalDatabase() (Database, error) {
	path := "database.json"
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
	err = json.Unmarshal(rawFile, db)
	return Database(db), err
}
