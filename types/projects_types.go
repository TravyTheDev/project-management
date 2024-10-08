package types

import "time"

type ProjectsStore interface {
	CreateProject(*Project) error
	DeleteProject(int) error
	UpdateProject(*Project) error
	GetAllProjects() ([]*Project, error)
	GetProjectByID(int) (*ProjectRes, error)
	GetProjectsByParentID(int) ([]*Project, error)
	GetProjectsByAssigneeID(int) ([]*Project, error)
	GetProjectsByStatus(int) ([]*Project, error)
}

type Project struct {
	ID          int       `json:"id"`
	ParentID    int       `json:"parentID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	AssigneeID  int       `json:"assigneeID"`
	Urgency     int       `json:"urgency"`
	Notes       string    `json:"notes"`
	StartDate   string    `json:"startDate"`
	EndDate     string    `json:"endDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProjectRes struct {
	Project *Project `json:"project"`
	User    *UserRes `json:"user,omitempty"`
}
