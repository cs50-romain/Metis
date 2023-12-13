package task

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var itasks []Task
var mtasks []Task
var ltasks []Task
var ctasks []Task

type Tasks struct {
	Itasks	[]Task
	Mtasks	[]Task
	Ltasks	[]Task
	Ctasks	[]Task
}

type Task struct {
	Id		int
	Content		string
	IsCompleted	bool
	CreatedAt	time.Time
	Importance	string
}

func GetTask(importance string, task_id int) Task {
	if importance == "important" {
		for idx, task := range itasks {
			if task.Id == task_id {
				return itasks[idx]
			}
		}
	} else if importance == "minor" {
		for idx, task := range mtasks {
			if task.Id == task_id {
				return mtasks[idx]
			}
		}
	} else if importance == "later" {
		for idx, task := range ltasks {
			if task.Id == task_id {
				return ltasks[idx]
			}
		}
	} else {
		log.Printf("Error; Could not find importance type for: %s\n", importance)
		return Task{}
	}
	return Task{} 
}

func (t *Task) DeleteTask(importance string, task_id int) {
	if importance == "important" {
		for idx, task := range itasks {
			if task.Id == task_id {
				log.Print("INFO: Found related task:", task)
				itasks = append(itasks[:idx], itasks[idx+1:]...)
				log.Print("INFO: New itasks:", itasks)
			}
		}
	} else if importance == "minor" {
		for idx, task := range mtasks {
			if task.Id == task_id {
				mtasks = append(mtasks[:idx], mtasks[idx+1:]...)
			}
		}
	} else if importance == "later" {
		for idx, task := range ltasks {
			if task.Id == task_id {
				log.Print("INFO: Found related task:", task)
				ltasks = append(ltasks[:idx], ltasks[idx+1:]...)
				log.Print("INFO: New ltasks:", ltasks)
			}
		}
	} else {
		log.Printf("Error; Could not find importance type for: %s\n", importance)
	}

	itasks = FixId(itasks)
	mtasks = FixId(mtasks)
	ltasks = FixId(ltasks)
}

func (t *Tasks) AddTaskToList(task Task, importance string) {	
	if importance == "important" {
		t.Itasks = append(t.Itasks, task)
	} else if importance == "minor" {
		t.Mtasks = append(t.Mtasks, task)
	} else if importance == "later" {
		t.Ltasks = append(t.Ltasks, task)
	} else if importance == "completed" {
		t.Ctasks = append(t.Ctasks, task)
	}
}

func getFile() []byte {
	b, err := os.ReadFile("./data.json")
	if err != nil {
		log.Println("[ERROR] Error getting tasks file ->", err)
	}
	return b
}

func Prefill(tasks *Tasks) {
	b := getFile()
	err := json.Unmarshal(b, &tasks)
	if err != nil {
		log.Println(err)
	}

	itasks = tasks.Itasks
	mtasks = tasks.Mtasks
	ltasks = tasks.Ltasks
	ctasks = tasks.Ctasks

	itasks = FixId(itasks)
	mtasks = FixId(mtasks)
	ltasks = FixId(ltasks)
	ctasks = FixId(ctasks)

	// AFTER PREFILLING, CHECK DATES AND SEE IF AN OBJECT NEEDS TO SWITCH
	// CHECK EACH TASK LISTS AND MOVE TASK ACCORDINGLY IF NEEDED
	swapTasksBasedOnDate()
}

func swapTasksBasedOnDate() {
	for idx, object := range mtasks {
		//log.Printf("Object %s creation Date:%s\t Current Date:%s", object.Content, object.CreatedAt.String(), time.Now().String())
		if compareDate(object.CreatedAt) >= 7 {
			object.Importance = "important"
			itasks = append(itasks, object)
			mtasks = append(mtasks[:idx], mtasks[idx+1:]...)
		}
	}

	for idx, object := range ltasks {
		//log.Printf("Object %s creation Date:%s\t Current Date:%s", object.Content, object.CreatedAt.String(), time.Now().String())
		if compareDate(object.CreatedAt) >= 4 {
			object.Importance = "minor"
			mtasks = append(mtasks, object)
			ltasks = append(ltasks[:idx], ltasks[idx+1:]...)
		}
	}
}

// compareDate 2 DATES (DATE1 AND TIME.NOW()) AND RETURN AN INTEGER - INTEGER WILL BE HOW MANY DAYS DATE1 IS FROM TIME.NOW()
func compareDate(date time.Time) int{
	if date.Year() == time.Now().Year() {
		if date.Month() == time.Now().Month() {
			return time.Now().Day() - date.Day()
		} else {
			monthDiff := int(time.Now().Month()) - int(date.Month())
			startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
			var daysTillEndOfMonth int = 0 
			if monthDiff <= 1 {
				if date.Month() == time.January || date.Month() == time.March || date.Month() == time.May || date.Month() == time.July || date.Month() == time.August || date.Month() == time.October || date.Month() == time.December {
					daysTillEndOfMonth = 31 - date.Day()
				} else if date.Month() == time.February {
					if date.Year() / 2 == 0 {
						daysTillEndOfMonth = 29 - date.Day() + 1
					} else {
						daysTillEndOfMonth = 28 - date.Day() + 1
					}
					
				} else {
					daysTillEndOfMonth = 30 - date.Day() + 1
				}
				return time.Now().Day() - startOfMonth.Day() + daysTillEndOfMonth 
			} else {
				return -1
			}
		}
	} else {
		return -1
		//return time.Now().Year() - date.Year() + 365 
	}
}

func FixId(array []Task) []Task {
	if array == nil {
		return array
	}

	for idx,_ := range array {
		array[idx].Id = idx
	}
	return array
}
