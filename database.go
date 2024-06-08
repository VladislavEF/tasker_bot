package main

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
	AddUserTask(string)
	AddNewTask(TaskInfo)
}

func GetDatabase() *Database {
	return new(Database)
}

func (this *MemoryStorage) GetUserId(string) string {
	return ""
}
func (this *MemoryStorage) GetUserTasks(string) []string {
	return []string{""}
}
func (this *MemoryStorage) GetTaskInfo(string) TaskInfo {
	return TaskInfo{}
}

func (this *MemoryStorage) AddUser(string) {
}
func (this *MemoryStorage) AddUserTask(string) {
}
func (this *MemoryStorage) AddNewTask(TaskInfo) {
}
