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
	stmt := `INSERT INTO projects (parent_id, title, description, status, assignee_id, urgency, notes, start_date, end_date)` +
		`VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := p.db.Exec(
		stmt,
		project.ParentID,
		project.Title,
		project.Description,
		project.Status,
		project.AssigneeID,
		project.Urgency,
		project.Notes,
		project.StartDate,
		project.EndDate,
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
		`start_date = ?, end_date = ?, updated_at = datetime(current_timestamp, 'localtime') where id = ?`

	_, err := p.db.Exec(
		stmt,
		project.ParentID,
		project.Title,
		project.Description,
		project.Status,
		project.AssigneeID,
		project.Urgency,
		project.Notes,
		project.StartDate,
		project.EndDate,
		project.ID,
	)

	return err
}

func (p *ProjectsStore) GetAllProjects() ([]*types.Project, error) {
	stmt := `SELECT * FROM projects where parent_id = 0`

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

func (p *ProjectsStore) GetProjectByID(id int) (*types.ProjectRes, error) {
	stmt := `SELECT projects.id, parent_id, title, description, status, assignee_id,` +
		`urgency, notes, start_date, end_date, COALESCE(users.id, 0), COALESCE(users.username, ''), COALESCE(users.email, '') ` +
		`FROM projects LEFT JOIN users ON users.id = projects.assignee_id WHERE projects.id = ?`

	rows, err := p.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	projectRes := new(types.ProjectRes)
	for rows.Next() {
		projectRes, err = p.scanRowsIntoProjectRes(rows)
		if err != nil {
			return nil, err
		}
	}

	return projectRes, nil
}

func (p *ProjectsStore) GetProjectsByParentID(id int) ([]*types.Project, error) {
	stmt := `SELECT * FROM projects where parent_id = ?`

	projects, err := p.getProjectsSliceFromIntQuery(stmt, id)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *ProjectsStore) GetProjectsByAssigneeID(id int) ([]*types.Project, error) {
	stmt := `SELECT * FROM projects where assignee_id = ?`

	projects, err := p.getProjectsSliceFromIntQuery(stmt, id)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *ProjectsStore) GetProjectsByStatus(status int) ([]*types.ProjectRes, error) {
	stmt := `SELECT projects.id, parent_id, title, description, status, assignee_id,` +
		`urgency, notes, start_date, end_date, COALESCE(users.id, 0), COALESCE(users.username, ''), COALESCE(users.email, '') ` +
		`FROM projects LEFT JOIN users ON users.id = projects.assignee_id WHERE projects.status = ? AND projects.parent_id = -1`

	rows, err := p.db.Query(stmt, status)
	if err != nil {
		return nil, err
	}
	projects := make([]*types.ProjectRes, 0)

	for rows.Next() {
		p, err := p.scanRowsIntoProjectRes(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (p *ProjectsStore) GetProjectsByUrgency(urgency int) ([]*types.ProjectRes, error) {
	stmt := `SELECT projects.id, parent_id, title, description, status, assignee_id,` +
		`urgency, notes, start_date, end_date, COALESCE(users.id, 0), COALESCE(users.username, ''), COALESCE(users.email, '') ` +
		`FROM projects LEFT JOIN users ON users.id = projects.assignee_id WHERE projects.urgency = ? AND projects.parent_id = -1`

	rows, err := p.db.Query(stmt, urgency)
	if err != nil {
		return nil, err
	}
	projects := make([]*types.ProjectRes, 0)

	for rows.Next() {
		p, err := p.scanRowsIntoProjectRes(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (p *ProjectsStore) SearchProjects(text string) ([]*types.Project, error) {
	projects := make([]*types.Project, 0)
	stmt := `SELECT * FROM projects WHERE title LIKE REPLACE(?, " ", "%")`
	rows, err := p.db.Query(stmt, `%`+text+`%`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		project, err := p.scanRowsIntoProject(rows)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}
	return projects, nil
}

func (p *ProjectsStore) getProjectsSliceFromIntQuery(stmt string, id int) ([]*types.Project, error) {
	rows, err := p.db.Query(stmt, id)
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
		&project.StartDate,
		&project.EndDate,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		fmt.Println(err)
	}
	return project, nil
}

func (p *ProjectsStore) scanRowsIntoProjectRes(rows *sql.Rows) (*types.ProjectRes, error) {
	project := new(types.Project)
	user := new(types.UserRes)
	err := rows.Scan(
		&project.ID,
		&project.ParentID,
		&project.Title,
		&project.Description,
		&project.Status,
		&project.AssigneeID,
		&project.Urgency,
		&project.Notes,
		&project.StartDate,
		&project.EndDate,
		&user.ID,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}
	projectRes := &types.ProjectRes{
		Project: project,
		User:    user,
	}

	return projectRes, nil
}
