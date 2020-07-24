package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	fileName := "tasks.json"
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
	for {
		fmt.Print("$ ")
		cmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd = strings.ReplaceAll(cmd, "\n", "")
		cmds := strings.Split(cmd, " ")
		if cmds[0] == "task" {
			if len(cmds) == 1 {
				taskInfo()
			} else {
				switch string(cmd[1]) {
				case "list":

				}
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func taskInfo() {
	fmt.Println("task is a CLI for managing your TODOs.")
	fmt.Println("\nUsage:\n\ttask [command]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("\tadd [task]\tAdd a new task to your TODO list")
	fmt.Println("\tdo\t\t\tMark a task on your TODO list as complete")
	fmt.Println("\tlist\t\tList all of your incomplete tasks")
	fmt.Println("\trm [n]\t\tDelete task n")
	fmt.Println("\tcompleted\tList out any tasks completed")
	fmt.Println("\nclose\t\tClose task manager")
}
