package projects

import (
	"database/sql"
	"fmt"
	"log"
	"project-management/types"
)

type ProjectsStore struct {
	db *sql.DB
}

func NewProjectsStore(db *sql.DB) *ProjectsStore {
	return &ProjectsStore{
		db: db,
	}
}

func (p *ProjectsStore) CreateProject(project *types.Project) error {
	stmt := `INSERT INTO projects (parent_id, title, description, status, assignee_id, urgency, notes)` +
		`VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := p.db.Exec(
		stmt,
		project.ParentID,
		project.Title,
		project.Description,
		project.Status,
		project.AssigneeID,
		project.Urgency,
		project.Notes,
	)

	return err
}

func (p *ProjectsStore) DeleteProject(id int) error {
	_, err := p.db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	stmt := `DELETE FROM projects WHERE id = ?`

	_, err = p.db.Exec(stmt, id)
	return err
}

func (p *ProjectsStore) UpdateProject(project *types.Project) error {
	stmt := `UPDATE projects SET parent_id = ?, title = ?, description = ?, status = ?, assignee_id = ?, urgency = ?, notes = ?, ` +
		`updated_at = datetime(current_timestamp, 'localtime') where id = ?`

	_, err := p.db.Exec(
		stmt,
		project.ParentID,
		project.Title,
		project.Description,
		project.Status,
		project.AssigneeID,
		project.Urgency,
		project.Notes,
		project.ID,
	)

	return err
}

func (p *ProjectsStore) GetAllProjects() ([]*types.Project, error) {
	stmt := `SELECT * FROM projects`

	rows, err := p.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	projects := make([]*types.Project, 0)
	for rows.Next() {
		p, err := p.scanRowsIntoProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (p *ProjectsStore) scanRowsIntoProject(rows *sql.Rows) (*types.Project, error) {
	project := new(types.Project)

	err := rows.Scan(
		&project.ID,
		&project.ParentID,
		&project.Title,
		&project.Description,
		&project.Status,
		&project.AssigneeID,
		&project.Urgency,
		&project.Notes,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		fmt.Println(err)
	}
	return project, nil
}
