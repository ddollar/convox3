package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/convox/convox/pkg/structs"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Task struct {
	Kind   string            `json:"kind"`
	Logs   string            `json:"logs"`
	Params map[string]string `json:"params"`
	Status string            `json:"status"`
	Title  string            `json:"title"`

	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
}

type Tasks []Task

type TaskWriter struct {
	m      *Model
	key    string
	output []byte
}

func (m *Model) NewTask(jid, kind, title string) (int, error) {
	j, err := m.JobGet(jid)
	if err != nil {
		return 0, err
	}

	t := Task{
		Kind:   kind,
		Logs:   fmt.Sprintf("jobs/%s/tasks/%s", j.ID, uuid.New().String()),
		Status: "running",
		Title:  title,
	}

	j.Tasks = append(j.Tasks, t)

	if err := m.JobSave(j); err != nil {
		return 0, errors.WithStack(err)
	}

	return len(j.Tasks) - 1, nil
}

func (m *Model) TaskOutput(t Task) (string, error) {
	r, err := m.rack.ObjectFetch(os.Getenv("APP"), t.Logs)
	if err != nil {
		return "", nil
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(data), nil
}

func (m *Model) TaskWriter(t Task) (*TaskWriter, error) {
	if t.Logs == "" {
		return nil, errors.WithStack(fmt.Errorf("no logs key for task"))
	}

	tw := &TaskWriter{
		m:      m,
		key:    t.Logs,
		output: []byte{},
	}

	return tw, nil
}

func (tw *TaskWriter) Write(data []byte) (int, error) {
	tw.output = append(tw.output, data...)

	_, err := tw.m.rack.ObjectStore(os.Getenv("APP"), tw.key, bytes.NewReader(tw.output), structs.ObjectStoreOptions{})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return len(data), nil
}
