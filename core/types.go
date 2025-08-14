package core

import (
	"fmt"
	"math/rand"
)

type Task struct {
	Id          int
	Description string
	Done        bool
}

type ListInfo struct {
	Name       string
	NumDone    int
	NumPending int
	NumTasks   int
}

type List struct {
	Info    ListInfo
	TaskIds []int
	Tasks   map[int]*Task
	UsedIds map[int]struct{}
}

func NewList(name string) List {
	return List{
		Info: ListInfo{
			Name: name,
		},
		TaskIds: []int{},
		Tasks:   make(map[int]*Task),
		UsedIds: make(map[int]struct{}),
	}
}

func (l *List) generateTaskId() (int, error) {
	numAttempts := 0
	for {
		id := rand.Int()
		if _, ok := l.UsedIds[id]; !ok {
			l.UsedIds[id] = struct{}{}
			return id, nil
		}
		numAttempts++
		if numAttempts > 100 {
			return -1, fmt.Errorf("failed to generate unique task id - make sure that they are being deleted from UsedIds correctly")
		}
	}
}

// make new task with the given description and completion status without adding it to the list.
func (l *List) NewTask(description string, done bool) (Task, error) {
	id, err := l.generateTaskId()
	if err != nil {
		return Task{}, err
	}
	task := Task{
		Id:          id,
		Description: description,
		Done:        done,
	}

	return task, nil
}

// add a task to the list.
func (l *List) AddTask(task Task) error {
	if _, ok := l.Tasks[task.Id]; ok {
		return fmt.Errorf("task with id %d already exists", task.Id)
	}
	l.Tasks[task.Id] = &task
	l.TaskIds = append(l.TaskIds, task.Id)
	l.Info.NumTasks++
	if task.Done {
		l.Info.NumDone++
	} else {
		l.Info.NumPending++
	}
	return nil
}

// insert a task after a certain index of the list.
func (l *List) Insert(task Task, index int) error {
	if index < 0 || index > len(l.TaskIds) {
		return fmt.Errorf("invalid index %d", index)
	}
	if _, ok := l.Tasks[task.Id]; ok {
		return fmt.Errorf("task with id %d already exists", task.Id)
	}
	l.Tasks[task.Id] = &task
	l.TaskIds = append(l.TaskIds[:index], append([]int{task.Id}, l.TaskIds[index:]...)...)
	l.Info.NumTasks++
	if task.Done {
		l.Info.NumDone++
	} else {
		l.Info.NumPending++
	}
	return nil
}

// add a new task to the list and update meta data
// returns the id of the newly created task
func (l *List) AddNewTask(description string, done bool) (int, error) {
	task, err := l.NewTask(description, done)
	if err != nil {
		return -1, err
	}
	err = l.AddTask(task)
	if err != nil {
		return -1, err
	}
	return task.Id, nil
}

// insert a new task into the list at the specified index and update meta data
func (l *List) InsertNewTask(description string, index int) (int, error) {
	task, err := l.NewTask(description, false)
	if err != nil {
		return -1, err
	}
	err = l.Insert(task, index)
	if err != nil {
		return -1, err
	}
	return task.Id, nil
}

// remove the task with the given id from the list and update meta data
func (l *List) RemoveTask(taskId int) error {
	if task, ok := l.Tasks[taskId]; ok {
		delete(l.Tasks, taskId)
		delete(l.UsedIds, taskId)
		l.TaskIds = RemoveIntFromSlice(l.TaskIds, taskId)
		l.Info.NumTasks--
		if task.Done {
			l.Info.NumDone--
		} else {
			l.Info.NumPending--
		}
	} else {
		return fmt.Errorf("tried removing non-existent task id %d from list %s", taskId, l.Info.Name)
	}
	return nil
}

// update the description of the task with the given id
func (l *List) EditTaskDescription(taskId int, newDescription string) error {
	if task, ok := l.Tasks[taskId]; ok {
		task.Description = newDescription
	} else {
		return fmt.Errorf("tried editing non-existent task id %d in list %s", taskId, l.Info.Name)
	}
	return nil
}

// toggle the completion status of the task with the given id
func (l *List) ToggleCompletion(taskId int) error {
	if task, ok := l.Tasks[taskId]; ok {
		task.Done = !task.Done
		if task.Done {
			l.Info.NumDone++
			l.Info.NumPending--
		} else {
			l.Info.NumDone--
			l.Info.NumPending++
		}
	} else {
		return fmt.Errorf("tried toggling non-existent task id %d in list %s", taskId, l.Info.Name)
	}
	return nil
}
