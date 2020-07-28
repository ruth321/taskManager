package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type todo struct {
	Task     string
	Complete bool
	Time     time.Time
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
					tasks.addTask(cmds[2:])
				case "do":
					tasks.doTask(cmds[2:])
				case "list":
					tasks.showTaskList(cmds[2:])
				case "rm":
					tasks.removeTask(cmds[2:])
				case "completed":
					tasks.showCompTasks(cmds[2:])
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
	file, _ = json.Marshal(tasks)
	_ = ioutil.WriteFile(fileName, file, 0644)
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
	if len(s) == 0 {
		fmt.Println("Cannot add an empty task.")
		return
	}
	task := todo{
		Task:     strings.Join(s, " "),
		Complete: false,
	}
	*tasks = append(*tasks, task)
	fmt.Printf("Added \"%s\" to your task list.\n", task.Task)
}

func (tasks *todos) showTaskList(cmds []string) {
	if isExcessCmds(cmds) {
		return
	}
	incompTasks := getIncompTasks(*tasks)
	if len(incompTasks) == 0 {
		fmt.Println("Your task list is empty.")
		fmt.Println("Use \"task add [task]\" to add a new task.")
		return
	}
	fmt.Println("You have the following tasks:")
	for i := 0; i < len(incompTasks); i++ {
		fmt.Printf(" %d. %s\n", i+1, incompTasks[i].Task)
	}
}

func (tasks *todos) doTask(s []string) {
	if !isValidNum(s, len(getIncompTasks(*tasks))) {
		return
	}
	n, _ := strconv.Atoi(s[0])
	for i, g := 0, 0; i < len(*tasks); i++ {
		if !(*tasks)[i].Complete {
			g++
			if g == n {
				(*tasks)[i].Complete = true
				(*tasks)[i].Time = time.Now()
				fmt.Printf("You have completed the \"%s\" task.\n", (*tasks)[i].Task)
				break
			}
		}
	}
}

func (tasks *todos) removeTask(s []string) {
	switch s[0] {
	case "incomp":
		(*tasks).rmIncompTask(s[1:])
	case "comp":
		(*tasks).rmCompTask(s[1:])
	default:
		_ = isExcessCmds(s[0:])
	}
}

func (tasks *todos) rmIncompTask(s []string) {
	if !isValidNum(s, len(getIncompTasks(*tasks))) {
		return
	}
	n, _ := strconv.Atoi(s[0])
	for i, g := 0, 0; i < len(*tasks); i++ {
		if !(*tasks)[i].Complete {
			g++
			if g == n {
				n = i
				break
			}
		}
	}
	fmt.Printf("You have deleted the \"%s\" task.\n", (*tasks)[n].Task)
	for i := n; i < len(*tasks)-1; i++ {
		(*tasks)[i] = (*tasks)[i+1]
	}
	*tasks = (*tasks)[:len(*tasks)-1]
}

func (tasks *todos) rmCompTask(s []string) {
	if !isValidNum(s, len(getCompTasks(*tasks))) {
		return
	}
	n, _ := strconv.Atoi(s[0])
	for i, g := 0, 0; i < len(*tasks); i++ {
		if (*tasks)[i].Complete {
			g++
			if g == n {
				n = i
				break
			}
		}
	}
	fmt.Printf("You have deleted the \"%s\" task.\n", (*tasks)[n].Task)
	for i := n; i < len(*tasks)-1; i++ {
		(*tasks)[i] = (*tasks)[i+1]
	}
	*tasks = (*tasks)[:len(*tasks)-1]
}

func (tasks *todos) showCompTasks(cmds []string) {
	var compTasks []todo
	if len(cmds) > 0 {
		switch cmds[0] {
		case "hour":
			if !isValidNum(cmds[1:], math.MaxInt32) {
				return
			}
			compTasks = getCertCompTasks(*tasks, "hour", cmds[1])
		case "day":
			if !isValidNum(cmds[1:], math.MaxInt32) {
				return
			}
			compTasks = getCertCompTasks(*tasks, "day", cmds[1])
		default:
			_ = isExcessCmds(cmds)
			return
		}
	} else {
		compTasks = getCompTasks(*tasks)
	}
	if len(compTasks) == 0 {
		fmt.Println("You do not have completed tasks.")
		return
	}
	fmt.Println("You have finished the following tasks:")
	for i := 0; i < len(compTasks); i++ {
		fmt.Printf(" %d. %s\n", i+1, compTasks[i].Task)
	}
}

func getIncompTasks(tasks []todo) []todo {
	var incompTasks []todo
	for i := 0; i < len(tasks); i++ {
		if !tasks[i].Complete {
			incompTasks = append(incompTasks, tasks[i])
		}
	}
	return incompTasks
}

func getCompTasks(tasks []todo) []todo {
	var compTasks []todo
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Complete {
			compTasks = append(compTasks, tasks[i])
		}
	}
	return compTasks
}

func getCertCompTasks(tasks []todo, t string, ds string) []todo {
	a, _ := strconv.Atoi(ds)
	d := float64(a)
	var compTasks []todo
	now := time.Now()
	if t == "hour" {
		for i := 0; i < len(tasks); i++ {
			if tasks[i].Complete && now.Sub(tasks[i].Time).Hours() < d {
				compTasks = append(compTasks, tasks[i])
			}
		}
	} else {
		for i := 0; i < len(tasks); i++ {
			if tasks[i].Complete && now.Sub(tasks[i].Time).Hours()/24 < d {
				compTasks = append(compTasks, tasks[i])
			}
		}
	}
	return compTasks
}

func isValidNum(s []string, max int) bool {
	if len(s) == 0 {
		fmt.Println("Number is not specified.")
		return false
	}
	if len(s) > 1 {
		fmt.Printf("Unknown args - \"%s\"\n", strings.Join(s, " "))
		return false
	}
	n, err := strconv.Atoi(s[0])
	if err != nil {
		fmt.Printf("Unknown arg - \"%s\"\n", s[0])
		return false
	}
	if n > max || n < 1 {
		fmt.Println("Invalid number.")
		return false
	}
	return true
}

func isExcessCmds(cmds []string) bool {
	if len(cmds) == 1 {
		fmt.Printf("Unknown arg - \"%s\"\n", strings.Join(cmds, " "))
		return true
	}
	if len(cmds) > 1 {
		fmt.Printf("Unknown args - \"%s\"\n", strings.Join(cmds, " "))
		return true
	}
	return false
}
