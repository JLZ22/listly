package core

import (
	"fmt"
	"os"
)

func Abort(msg string) {
	fmt.Println("Abort!", msg)
	os.Exit(1)
}

func Success(msg string) {
	fmt.Println("Success!", msg)
	os.Exit(0)
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
