package main

import "fmt"

type StudyTask struct {
	ID        int
	Subject   string
	Content   string
	PlanTime  int
	SpentTime int
	Completed bool
}

func main() {
	// Create an empty list named tasks that stores StudyTask values
	tasks := make([]StudyTask, 0)
	// nextID is used to give each new task a unique ID
	nextID := 1

	fmt.Println("Welcome to the Study Task Manager")

	for {
		showMenu()

		var choice int

		fmt.Print("please choose: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			return
		}

		switch choice {
		case 1:
			fmt.Println("the add-task feature is not finished yet.")
		case 2:
			fmt.Println("the show-task feature is finished yet.")
		case 3:
			fmt.Println("the add-task feature is finished yet.")
		case 4:
			fmt.Println("the show-task feature is finished yet.")
		case 0:
			fmt.Println("the program has exited.")
			return
		default:
			fmt.Println("please choose an option.")
		}

		fmt.Println("Current task count:", len(tasks))
		fmt.Println("Next task ID:", nextID)
	}
}

func showMenu() {
	fmt.Println()
	fmt.Println("====== Study Task Manager ======")
	fmt.Println("1. add a task")
	fmt.Println("2. show all tasks")
	fmt.Println("3. complete a task")
	fmt.Println("4. delete a task")
	fmt.Println("0. exit")

}

// inputTask reads task information from the user.
func inputTask(id int) StudyTask {
	var task StudyTask

	task.ID = id

	fmt.Print("Enter the subject: ")
	fmt.Scan(&task.Subject)

	fmt.Print("Enter the task content: ")
	fmt.Scan(&task.Content)

	fmt.Print("Enter the planned study time in minutes: ")
	fmt.Scan(&task.PlanTime)

	task.SpentTime = 0
	task.Completed = false

	return task
}
