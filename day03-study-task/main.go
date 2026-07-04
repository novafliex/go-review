package main

import "fmt"

// StudyTask represents a study task.
type StudyTask struct {
	ID        int
	Subject   string
	Content   string
	PlanTime  int
	SpentTime int
	Completed bool
}

func main() {
	// Create an empty slice that stores StudyTask values.
	tasks := make([]StudyTask, 0)

	// nextID gives each new task a unique ID.
	nextID := 1

	fmt.Println("Welcome to the Study Task Manager")

	for {
		showMenu()

		var choice int

		fmt.Print("Please choose: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input. The program has exited.")
			return
		}

		switch choice {
		case 1:
			task, err := inputTask(nextID)
			if err != nil {
				fmt.Println("Failed to read task information:", err)
				break
			}

			tasks = addTask(tasks, task)
			nextID++

			fmt.Println("Task added successfully.")

		case 2:
			showTasks(tasks)

		case 3:
			var id int

			fmt.Print("Enter the task ID to complete: ")
			_, err := fmt.Scan(&id)
			if err != nil {
				fmt.Println("Invalid task ID.")
				break
			}

			if completeTask(tasks, id) {
				fmt.Println("Task completed successfully.")
			} else {
				fmt.Println("Task not found.")
			}

		case 4:
			var id int

			fmt.Print("Enter the task ID to delete: ")
			_, err := fmt.Scan(&id)
			if err != nil {
				fmt.Println("Invalid task ID.")
				break
			}

			var success bool
			tasks, success = deleteTask(tasks, id)

			if success {
				fmt.Println("Task deleted successfully.")
			} else {
				fmt.Println("Task not found.")
			}

		case 5:
			var id int
			var spentTime int

			fmt.Print("Enter the task ID: ")
			_, err := fmt.Scan(&id)
			if err != nil {
				fmt.Println("Invalid task ID.")
				break
			}

			fmt.Print("Enter the actual study time in minutes: ")
			_, err = fmt.Scan(&spentTime)
			if err != nil {
				fmt.Println("Invalid study time.")
				break
			}

			if updateSpentTime(tasks, id, spentTime) {
				fmt.Println("Spent time updated successfully.")
			} else {
				fmt.Println("Task not found or study time is invalid.")
			}

		case 0:
			fmt.Println("The program has exited.")
			return

		default:
			fmt.Println("Please choose a valid option.")
		}

		fmt.Println("Current task count:", len(tasks))
		fmt.Println("Next task ID:", nextID)
	}
}

// showMenu displays the main menu.
func showMenu() {
	fmt.Println()
	fmt.Println("====== Study Task Manager ======")
	fmt.Println("1. Add a task")
	fmt.Println("2. Show all tasks")
	fmt.Println("3. Complete a task")
	fmt.Println("4. Delete a task")
	fmt.Println("5. Update spent time")
	fmt.Println("0. Exit")
}

// inputTask reads task information from the user.
func inputTask(id int) (StudyTask, error) {
	var task StudyTask

	task.ID = id

	fmt.Print("Enter the subject: ")
	_, err := fmt.Scan(&task.Subject)
	if err != nil {
		return StudyTask{}, err
	}

	fmt.Print("Enter the task content: ")
	_, err = fmt.Scan(&task.Content)
	if err != nil {
		return StudyTask{}, err
	}

	fmt.Print("Enter the planned study time in minutes: ")
	_, err = fmt.Scan(&task.PlanTime)
	if err != nil {
		return StudyTask{}, err
	}

	if task.PlanTime < 0 {
		return StudyTask{}, fmt.Errorf("planned study time cannot be negative")
	}

	task.SpentTime = 0
	task.Completed = false

	return task, nil
}

// addTask adds a task to the task slice.
func addTask(tasks []StudyTask, task StudyTask) []StudyTask {
	return append(tasks, task)
}

// showTasks displays all study tasks.
func showTasks(tasks []StudyTask) {
	if len(tasks) == 0 {
		fmt.Println("There are no tasks.")
		return
	}

	for _, task := range tasks {
		fmt.Println("--------------------")
		fmt.Println("ID:", task.ID)
		fmt.Println("Subject:", task.Subject)
		fmt.Println("Content:", task.Content)
		fmt.Println("Planned time:", task.PlanTime, "minutes")
		fmt.Println("Spent time:", task.SpentTime, "minutes")
		fmt.Println("Completed:", task.Completed)
	}
}

// completeTask marks a task as completed.
func completeTask(tasks []StudyTask, id int) bool {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Completed = true
			return true
		}
	}

	return false
}

// deleteTask deletes a task by ID.
func deleteTask(tasks []StudyTask, id int) ([]StudyTask, bool) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return tasks, true
		}
	}

	return tasks, false
}

// updateSpentTime updates the actual study time of a task.
func updateSpentTime(tasks []StudyTask, id, spentTime int) bool {
	if spentTime < 0 {
		return false
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].SpentTime = spentTime
			return true
		}
	}

	return false
}
