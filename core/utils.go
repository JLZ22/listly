package core

import (
	"fmt"
)

func Success(msg string) {
	fmt.Println(msg) // this used to print a success message, but that was removed, so now this is effectively a no-op. Too lazy to refactor.
}

func ListLists(lists []string, tab string) string {
	out := ""
	for _, list := range lists {
		out += fmt.Sprintf("%s- %s\n", tab, list)
	}
	return out
}

func SplitByCompletion(list List) (completed, pending []*Task) {
	for _, i := range list.TaskIds {
		task, ok := list.Tasks[i]
		if !ok {
			continue
		}
		if task.Done {
			completed = append(completed, task)
		} else {
			pending = append(pending, task)
		}
	}
	return
}

func RemoveIntFromSlice(s []int, val int) []int {
	for i, v := range s {
		if v == val {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}