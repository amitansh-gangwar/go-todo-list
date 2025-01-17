package todo

import (
	"database/sql"
	"fmt"
	"os"

	"testing"
	"time"

	"github.com/amitansh-gangwar/go-todo-list/config"
	"github.com/amitansh-gangwar/go-todo-list/db"
	"github.com/stretchr/testify/assert"
)

const (
	task         = "some random testing task"
	taskID       = 1
	updateTask   = "updated testing task"
	updateTaskID = 1
	deleteTaskID = 1
)

func TestCreate(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	task_id, err := Create(dbIns, task)
	assert.NoError(t, err)

	var task_check string
	s, err := dbIns.Prepare("SELECT task AS task_check FROM tasks WHERE task_id = $1")
	rows := s.QueryRow(task_id)
	rows.Scan(&task_check)

	dbIns.Exec("truncate table tasks;")
	assert.Equal(t, task, task_check)
	assert.NoError(t, err)
}

func TestCannotCreateEmptyTask(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	task_id, err := Create(dbIns, "")
	assert.EqualError(t, err, "Cannot Create an empty task")
	assert.Equal(t, -1, task_id)
}
func TestReadForExistingTask(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	var task_id int
	statement, err := dbIns.Prepare("INSERT INTO tasks(task, created_at, status) VALUES($1, $2, $3) RETURNING task_id;")
	rows := statement.QueryRow("read existing test task", fmt.Sprintf("%v-%d-%v", time.Now().Year(), int(time.Now().Month()), time.Now().Day()), false)
	rows.Scan(&task_id)

	ts, err := Read(dbIns, task_id)

	dbIns.Exec("truncate table tasks;")
	assert.NoError(t, err)
	assert.Equal(t, "read existing test task", ts.Task)
}

func TestReadForNoTask(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	_, err := Read(dbIns, -10000000)
	assert.EqualError(t, err, "Task Id is non-existent")
}

func TestShowAll(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	_, err := ShowAll(dbIns)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	var task_id int
	var task string
	statement, err := dbIns.Prepare("INSERT INTO tasks(task, created_at, status) VALUES($1, $2, $3) RETURNING task_id;")
	rows := statement.QueryRow("update test task", fmt.Sprintf("%v-%d-%v", time.Now().Year(), int(time.Now().Month()), time.Now().Day()), false)
	rows.Scan(&task_id)

	err = Update(dbIns, task_id, "updated task")
	assert.NoError(t, err)

	statement, err = dbIns.Prepare("SELECT task from tasks where task_id=$1;")
	row := statement.QueryRow(task_id)
	row.Scan(&task)
	assert.Equal(t, "updated task", task)
	assert.NoError(t, err)
	dbIns.Exec("truncate table tasks;")
}

func TestCannotUpdateWithEmptyTask(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	var task_id int
	statement, err := dbIns.Prepare("INSERT INTO tasks(task, created_at, status) VALUES($1, $2, $3) RETURNING task_id;")
	rows := statement.QueryRow("update test task", fmt.Sprintf("%v-%d-%v", time.Now().Year(), int(time.Now().Month()), time.Now().Day()), false)
	rows.Scan(&task_id)

	err = Update(dbIns, task_id, "")
	assert.EqualError(t, err, "Cannot update with an empty task")
	dbIns.Exec("truncate table tasks;")
}

func TestMarkDone(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	var task_id int
	statement, err := dbIns.Prepare("INSERT INTO tasks(task, created_at, status) VALUES($1, $2, $3) RETURNING task_id;")
	rows := statement.QueryRow("update test task", fmt.Sprintf("%v-%d-%v", time.Now().Year(), int(time.Now().Month()), time.Now().Day()), false)
	rows.Scan(&task_id)

	err = MarkDone(dbIns, task_id)
	assert.NoError(t, err)

	var status bool
	statement, err = dbIns.Prepare("SELECT status from tasks where task_id=$1;")
	row := statement.QueryRow(task_id)
	row.Scan(&status)
	assert.Equal(t, true, status)
	assert.NoError(t, err)
	dbIns.Exec("truncate table tasks;")
}

func TestMarkDoneNotExistingTask(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	err := MarkDone(dbIns, -100)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestDelete(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	config.Load()

	dbIns := db.InitDB()
	defer dbIns.Close()

	var task_id int
	statement, err := dbIns.Prepare("INSERT INTO tasks(task, created_at, status) VALUES($1, $2, $3) RETURNING task_id;")
	rows := statement.QueryRow("delete test task", fmt.Sprintf("%v-%d-%v", time.Now().Year(), int(time.Now().Month()), time.Now().Day()), false)
	rows.Scan(&task_id)

	err = Delete(dbIns, task_id)
	var counter int
	statement, err = dbIns.Prepare("SELECT COUNT(*) from tasks where task_id=$1;")
	rows = statement.QueryRow(task_id)
	rows.Scan(&counter)
	dbIns.Exec("truncate table tasks;")
	assert.Zero(t, counter)
	assert.NoError(t, err)
}
