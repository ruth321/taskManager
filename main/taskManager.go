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

type todos []todo

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
	var tasks todos
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
					tasks.addTask(cmds[2:])
				case "do":
					if len(cmds) > 3 {
						fmt.Printf("Unknown arguments - \"%s\".\n", strings.Join(cmds[3:], " "))
						break
					}
					tasks = doTask(cmds[2], tasks)
				case "list":
					tasks.listTasks()
				case "close":
					fmt.Println("Task manager closed.")
					loop = false
				default:
					fmt.Printf("Unknown command - \"%s\".\n", strings.Join(cmds, " "))
				}
			}
		} else {
			fmt.Printf("Unknown command - \"%s\".\n", strings.Join(cmds, " "))
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

func (tasks *todos) addTask(s []string) {
	task := todo{
		Task:     strings.Join(s, " "),
		Complete: false,
	}
	*tasks = append(*tasks, task)
	fmt.Printf("Added \"%s\" to your task list.\n", task.Task)
}

func (tasks *todos) listTasks() {
	incomplete := false
	for i := 0; i < len(*tasks); i++ {
		if !(*tasks)[i].Complete {
			incomplete = true
			break
		}
	}
	if !incomplete {
		fmt.Println("Your task list is empty.")
		fmt.Println("Use \"task add [task]\" to add a new task.")
		return
	}
	fmt.Println("You have the following tasks:")
	for i, g := 0, 0; i < len(*tasks); i++ {
		if !(*tasks)[i].Complete {
			fmt.Printf(" %d. %s\n", g+1, (*tasks)[i].Task)
			g++
		}
	}
}

func doTask(n string, tasks []todo) []todo {

	return tasks
}
