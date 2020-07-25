package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type todo struct {
	Task     string
	Complete bool
}

func main() {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	fileName := homeDir + "/tasks.json"
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		_, _ = os.Create(fileName)
		file, err = ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	var tasks []todo
	_ = json.Unmarshal(file, &tasks)
	loop := true
	for loop {
		fmt.Print("$ ")
		cmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd = strings.ReplaceAll(cmd, "\n", "")
		cmds := strings.Fields(cmd)
		if cmds[0] == "task" {
			if len(cmds) == 1 {
				taskInfo()
			} else {
				switch cmds[1] {
				case "add":
					if len(cmds) == 2 {
						fmt.Println("Entered an empty task.")
						break
					}
					tasks = append(tasks, addTask(cmds[2:]))
				case "list":
					listTasks(tasks)
				case "close":
					fmt.Println("Task manager closed.")
					loop = false
				default:
					fmt.Printf("Command \"%s\" not found.\n", strings.Join(cmds, " "))
				}
			}
		} else {
			fmt.Printf("Command \"%s\" not found.\n", cmds[0])
		}
	}
}

func taskInfo() {
	fmt.Println("task is a CLI for managing your TODOs")
	fmt.Println("\nUsage:\n\ttask [command]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("\tadd [task]\tAdd a new task to your TODO list")
	fmt.Println("\tdo\t\tMark a task on your TODO list as complete")
	fmt.Println("\tlist\t\tList all of your incomplete tasks")
	fmt.Println("\trm \t\tDelete task")
	fmt.Println("\tcompleted\tList out any tasks completed")
	fmt.Println("\tclose\t\tClose task manager")
}

func addTask(s []string) todo {
	task := todo{
		Task:     strings.Join(s, " "),
		Complete: false,
	}
	fmt.Printf("Added \"%s\" to your task list.\n", task.Task)
	return task
}

func listTasks(tasks []todo) {
	incomplete := false
	for i := 0; i < len(tasks); i++ {
		if !tasks[i].Complete {
			incomplete = true
			break
		}
	}
	if !incomplete {
		fmt.Println("Your task list is empty")
		fmt.Println("Use \"task add [task]\" to add a new task")
		return
	}
	fmt.Println("You have the following tasks:")
	for i, g := 0, 0; i < len(tasks); i++ {
		if !tasks[i].Complete {
			fmt.Printf(" %d. %s\n", g+1, tasks[i].Task)
			g++
		}
	}
}
