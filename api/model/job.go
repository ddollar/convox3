package model

import (
	"fmt"
	"time"

	"github.com/convox/console/pkg/integration"
	"github.com/convox/console/pkg/logger"
	"github.com/convox/console/pkg/settings"
	"github.com/pkg/errors"
)

type Job struct {
	ID string `dynamo:"id"`

	Kind     string            `dynamo:"kind"`
	Name     string            `dynamo:"name"`
	Params   map[string]string `dynamo:"params"`
	Pid      string            `dynamo:"pid"`
	Services map[string]string `dynamo:"services"`
	Status   string            `dynamo:"status"`
	Tasks    Tasks             `dynamo:"tasks"`
	TTL      int64             `dynamo:"ttl"`

	OrganizationID string `dynamo:"organization-id"`
	WorkflowID     string `dynamo:"workflow-id"`
	Workflow       Workflow

	Created  time.Time `dynamo:"created"`
	Started  time.Time `dynamo:"started"`
	Finished time.Time `dynamo:"finished"`
}

func (m *Model) JobFail(id string, failure error) error {
	logger.New("ns=models model=Job at=Fail").ErrorBacktrace(failure)

	j, err := m.JobGet(id)
	if err != nil {
		return err
	}

	w, err := m.WorkflowGet(j.WorkflowID)
	if err != nil {
		return errors.WithStack(err)
	}

	is, err := m.OrganizationIntegrations(w.OrganizationId)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, n := range is.ByKind("notification") {
		if ni, err := n.Notification(); err == nil {
			ni.EventSend(w.Name, n.Attributes, integration.NotificationEvent{
				Action: "workflow:complete",
				Data: map[string]string{
					"kind": w.Kind,
				},
				Status:    "error",
				Timestamp: time.Now().UTC(),
			})
		}
	}

	for i, t := range j.Tasks {
		if t.Status == "running" {
			j.Tasks[i].Status = "failed"
		}
	}

	s, err := m.WorkflowSource(w.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	switch w.Kind {
	case "review":
		r, err := m.WorkflowRepository(w.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		s.StatusUpdate(r, j.Params["ref"], "failure", failure.Error(), j.URL())
	}

	j.Finished = time.Now().UTC()
	j.Status = "failed"

	if err := m.JobSave(j); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Model) JobGet(id string) (*Job, error) {
	j := &Job{}

	if err := m.storage.Get("jobs", id, j); err != nil {
		return nil, errors.WithStack(err)
	}

	if w, _ := m.WorkflowGet(j.WorkflowID); w != nil {
		j.Workflow = *w
	}

	return j, nil
}

func (m *Model) JobListByStatus(status string) (Jobs, error) {
	js := Jobs{}

	if err := m.storage.GetIndex("jobs", "status-created-index", map[string]string{"status": status}, &js); err != nil {
		return nil, errors.WithStack(err)
	}

	return js, nil
}

func (m *Model) JobSave(j *Job) error {
	if err := m.storage.Put("jobs", j); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Model) JobSourceStatus(id, status string) error {
	j, err := m.JobGet(id)
	if err != nil {
		return err
	}

	w, err := m.WorkflowGet(j.WorkflowID)
	if err != nil {
		return errors.WithStack(err)
	}

	i, err := m.IntegrationGet(w.IntegrationID)
	if err != nil {
		return errors.WithStack(err)
	}

	si, err := i.Source()
	if err != nil {
		return errors.WithStack(err)
	}

	switch j.Workflow.Kind {
	case "review":
		r, err := m.WorkflowRepository(w.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		si.StatusUpdate(r, j.Params["ref"], "pending", status, j.URL())
	}

	return nil
}

func (m *Model) JobSucceed(id string) error {
	j, err := m.JobGet(id)
	if err != nil {
		return err
	}

	w, err := m.WorkflowGet(j.WorkflowID)
	if err != nil {
		return errors.WithStack(err)
	}

	is, err := m.OrganizationIntegrations(w.OrganizationId)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, n := range is.ByKind("notification") {
		if ni, err := n.Notification(); err == nil {
			ni.EventSend(w.Name, n.Attributes, integration.NotificationEvent{
				Action: "workflow:complete",
				Data: map[string]string{
					"kind": w.Kind,
				},
				Status:    "success",
				Timestamp: time.Now().UTC(),
			})
		}
	}

	s, err := m.WorkflowSource(w.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	switch w.Kind {
	case "review":
		r, err := m.WorkflowRepository(w.ID)
		if err != nil {
			return errors.WithStack(err)
		}

		s.StatusUpdate(r, j.Params["ref"], "success", "Workflow completed successfully", j.URL())
	}

	j.Finished = time.Now().UTC()
	j.Status = "complete"

	if err := m.JobSave(j); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (j *Job) URL() string {
	return fmt.Sprintf("https://%s/organizations/%s/jobs/%s", settings.ExternalHost, j.OrganizationID, j.ID)
}
