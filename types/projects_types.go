package types

import "time"

type ProjectsStore interface {
	CreateProject(*Project) error
	DeleteProject(int) error
	UpdateProject(*Project) error
	GetAllProjects() ([]*Project, error)
}

type Project struct {
	ID          int       `json:"id"`
	ParentID    int       `json:"parent_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	AssigneeID  int       `json:"assignee_id"`
	Urgency     int       `json:"urgency"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
