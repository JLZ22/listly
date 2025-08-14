package core_test

import (
	"fmt"
	"testing"

	"github.com/jlz22/listly/core"
	"github.com/stretchr/testify/require"
)

func TestNewList(t *testing.T) {
	name := "testlist"
	l := core.NewList(name)
	if l.Info.Name != name {
		t.Errorf("expected list name %q, got %q", name, l.Info.Name)
	}
	if l.Info.NumDone != 0 || l.Info.NumPending != 0 || l.Info.NumTasks != 0 {
		t.Errorf("expected empty counters, got %+v", l.Info)
	}
	if len(l.Tasks) != 0 {
		t.Errorf("expected empty Tasks map, got %d items", len(l.Tasks))
	}
}

func TestAddNewTask(t *testing.T) {
	l := core.NewList("test")
	id, err := l.AddNewTask("task 1", false)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}
	if len(l.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(l.Tasks))
	}
	task, ok := l.Tasks[id]
	if !ok {
		t.Fatalf("task with id %d not found", id)
	}
	if task.Description != "task 1" {
		t.Errorf("expected description %q, got %q", "task 1", task.Description)
	}
	if task.Done {
		t.Errorf("expected task done to be false")
	}
	if l.Info.NumTasks != 1 || l.Info.NumDone != 0 || l.Info.NumPending != 1 {
		t.Errorf("unexpected counts %+v", l.Info)
	}

	id2, err := l.AddNewTask("task 2", true)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}
	if l.Info.NumTasks != 2 || l.Info.NumDone != 1 || l.Info.NumPending != 1 {
		t.Errorf("unexpected counts after done task %+v", l.Info)
	}
	_ = id2
}

func TestInsert(t *testing.T) {
	l := core.NewList("test")
	for i := 0; i < 5; i++ {
		_, err := l.AddNewTask(fmt.Sprintf("task %d", i+1), false)
		if err != nil {
			t.Fatalf("unexpected error creating task: %v", err)
		}
	}

	insertedTask, err := l.NewTask("inserted task", true)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}
	err = l.Insert(insertedTask, 2)
	if err != nil {
		t.Fatalf("unexpected error inserting task: %v", err)
	}
	if len(l.Tasks) != 6 {
		t.Fatalf("expected 6 tasks after insert, got %d", len(l.Tasks))
	}
	if l.Info.NumTasks != 6 || l.Info.NumDone != 1 || l.Info.NumPending != 5 {
		t.Errorf("unexpected counts after insert %+v", l.Info)
	}
	require.Equal(t, "inserted task", l.Tasks[l.TaskIds[2]].Description)
}

func TestRemoveTask(t *testing.T) {
	l := core.NewList("test")
	id, err := l.AddNewTask("task 1", false)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}
	_, err = l.AddNewTask("task 2", true)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}

	err = l.RemoveTask(id + 999)
	if err == nil {
		t.Error("expected error removing non-existent task, got nil")
	}

	err = l.RemoveTask(id)
	if err != nil {
		t.Errorf("unexpected error removing existing task: %v", err)
	}
	if len(l.Tasks) != 1 {
		t.Errorf("expected 1 task after removal, got %d", len(l.Tasks))
	}
	if l.Info.NumTasks != 1 || l.Info.NumDone != 1 || l.Info.NumPending != 0 {
		t.Errorf("unexpected counts after removal %+v", l.Info)
	}
}

func TestEditTaskDescription(t *testing.T) {
	l := core.NewList("test")
	id, err := l.AddNewTask("task 1", false)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}

	err = l.EditTaskDescription(id, "new desc")
	if err != nil {
		t.Errorf("unexpected error editing description: %v", err)
	}
	if l.Tasks[id].Description != "new desc" {
		t.Errorf("expected description 'new desc', got %q", l.Tasks[id].Description)
	}

	err = l.EditTaskDescription(id+999, "fail")
	if err == nil {
		t.Error("expected error editing non-existent task, got nil")
	}
}

func TestToggleCompletion(t *testing.T) {
	l := core.NewList("test")
	id, err := l.AddNewTask("task 1", false)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}

	err = l.ToggleCompletion(id)
	if err != nil {
		t.Errorf("unexpected error toggling completion: %v", err)
	}
	if !l.Tasks[id].Done {
		t.Error("expected Done=true after toggle")
	}
	if l.Info.NumDone != 1 || l.Info.NumPending != 0 {
		t.Errorf("unexpected counts %+v", l.Info)
	}

	err = l.ToggleCompletion(id)
	if err != nil {
		t.Errorf("unexpected error toggling completion: %v", err)
	}
	if l.Tasks[id].Done {
		t.Error("expected Done=false after second toggle")
	}
	if l.Info.NumDone != 0 || l.Info.NumPending != 1 {
		t.Errorf("unexpected counts after second toggle %+v", l.Info)
	}

	err = l.ToggleCompletion(id + 999)
	if err == nil {
		t.Error("expected error toggling non-existent task, got nil")
	}
}

func TestUsedIdsUpdate(t *testing.T) {
	list := core.NewList("test")

	// Add a task and check UsedIds
	id, err := list.AddNewTask("task1", false)
	if err != nil {
		t.Fatalf("NewTask failed: %v", err)
	}
	if _, exists := list.UsedIds[id]; !exists {
		t.Errorf("UsedIds missing new task id %d", id)
	}

	// Remove the task and check UsedIds
	err = list.RemoveTask(id)
	if err != nil {
		t.Fatalf("RemoveTask failed: %v", err)
	}
	if _, exists := list.UsedIds[id]; exists {
		t.Errorf("UsedIds still contains removed task id %d", id)
	}
}

func TestRemoveTask_NotFound(t *testing.T) {
	l := core.NewList("test")
	err := l.RemoveTask(999)
	if err == nil {
		t.Error("expected error removing from empty list, got nil")
	}
}

func TestEditTaskDescription_NotFound(t *testing.T) {
	l := core.NewList("test")
	err := l.EditTaskDescription(123, "doesn't exist")
	if err == nil {
		t.Error("expected error editing non-existent task, got nil")
	}
}

func TestToggleCompletion_NotFound(t *testing.T) {
	l := core.NewList("test")
	err := l.ToggleCompletion(123)
	if err == nil {
		t.Error("expected error toggling non-existent task, got nil")
	}
}

func TestListInfoTracking(t *testing.T) {
	l := core.NewList("tracktest")

	// Add three tasks: two pending, one done
	id1, _ := l.AddNewTask("t1", false)
	l.AddNewTask("t2", false)
	id3, _ := l.AddNewTask("t3", true)

	if l.Info.NumTasks != 3 || l.Info.NumDone != 1 || l.Info.NumPending != 2 {
		t.Errorf("unexpected counts after adds: %+v", l.Info)
	}

	// Toggle one pending to done
	_ = l.ToggleCompletion(id1)
	if l.Info.NumDone != 2 || l.Info.NumPending != 1 {
		t.Errorf("unexpected counts after toggle: %+v", l.Info)
	}

	// Remove one done task
	_ = l.RemoveTask(id3)
	if l.Info.NumTasks != 2 || l.Info.NumDone != 1 || l.Info.NumPending != 1 {
		t.Errorf("unexpected counts after removal: %+v", l.Info)
	}

	// Toggle back to pending
	_ = l.ToggleCompletion(id1)
	if l.Info.NumDone != 0 || l.Info.NumPending != 2 {
		t.Errorf("unexpected counts after toggle back: %+v", l.Info)
	}
}
