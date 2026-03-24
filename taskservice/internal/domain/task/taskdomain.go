package taskdomain

import "time"

type TaskDomain struct {
	Id          uint32
	ProjectId   uint32
	Description string
	Deadline    time.Time
}

func NewTaskDomain(projectId uint32, description string, deadline time.Time) (*TaskDomain, error) {
	if err := validateProjectId(projectId); err != nil {
		return nil, err
	}
	if err := validateDescription(description); err != nil {
		return nil, err
	}
	return &TaskDomain{
		Id:          0,
		ProjectId:   projectId,
		Description: description,
		Deadline:    deadline,
	}, nil
}

func RestoreTaskDomain(id uint32, project_id uint32, description string, deadline time.Time) *TaskDomain {
	return &TaskDomain{
		Id:          id,
		ProjectId:   project_id,
		Description: description,
		Deadline:    deadline,
	}
}

func validateProjectId(projectId uint32) error {
	if projectId == 0 {
		return ErrInvalidProjectId
	}
	return nil
}

func validateDescription(description string) error {
	rDesc := []rune(description)
	if len(rDesc) == 0 || len(rDesc) > 255 {
		return ErrInvalidDescription
	}
	return nil
}
