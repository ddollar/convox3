package job

import (
	"fmt"
	"io"
	"time"

	"github.com/convox/console/api/model"
	"github.com/pkg/errors"
)

func task(j *model.Job, kind, title string, fn func(w io.Writer) error) error {
	m.JobSourceStatus(j.ID, title)

	ti, err := m.NewTask(j.ID, kind, title)
	if err != nil {
		return errors.WithStack(err)
	}

	j.Tasks[ti].Started = time.Now().UTC()

	if err := m.JobSave(j); err != nil {
		return errors.WithStack(err)
	}

	tw, err := m.TaskWriter(j.Tasks[ti])
	if err != nil {
		return err
	}

	err = fn(tw)

	j.Tasks[ti].Ended = time.Now().UTC()

	if err != nil {
		fmt.Printf("ns=job at=task status=fail error=%q\n", err)
		j.Tasks[ti].Status = "failed"
	} else {
		j.Tasks[ti].Status = "complete"
	}

	if err := m.JobSave(j); err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(err)
}
