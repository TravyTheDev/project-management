package projects

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-management/types"
	"strconv"

	"github.com/gorilla/mux"
)

type ProjectsHandler struct {
	projectsStore types.ProjectsStore
}

func NewProjectsHandler(projectsStore types.ProjectsStore) *ProjectsHandler {
	return &ProjectsHandler{
		projectsStore: projectsStore,
	}
}

func (p *ProjectsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/projects/create_project", p.createProject).Methods("POST")
	router.HandleFunc("/projects/delete_project/{id}", p.deleteProject).Methods("DELETE")
	router.HandleFunc("/projects/update_project", p.updateProject).Methods("PUT")
	router.HandleFunc("/projects/all_projects", p.getAllProjects).Methods("GET")
	router.HandleFunc("/projects/project/{id}", p.getProjectByID).Methods("GET")
	router.HandleFunc("/projects/child_projects/{id}", p.getProjectsByParentID).Methods("GET")
	router.HandleFunc("/projects/user_projects/{id}", p.getProjectsByAssigneeID).Methods("GET")
	router.HandleFunc("/projects/project_status/{status}", p.getProjectsByStatus).Methods("GET")
}

func (p *ProjectsHandler) createProject(w http.ResponseWriter, r *http.Request) {
	project := &types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "error decoding project body", http.StatusInternalServerError)
		return
	}
	if err := p.projectsStore.CreateProject(project); err != nil {
		http.Error(w, "error creating project", http.StatusInternalServerError)
		return
	}
}

func (p *ProjectsHandler) deleteProject(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	if err := p.projectsStore.DeleteProject(projectID); err != nil {
		fmt.Println(err)
	}
}

func (p *ProjectsHandler) updateProject(w http.ResponseWriter, r *http.Request) {
	project := &types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		fmt.Println(err)
	}
	if err := p.projectsStore.UpdateProject(project); err != nil {
		fmt.Println(err)
	}
}

func (p *ProjectsHandler) getAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := p.projectsStore.GetAllProjects()
	if err != nil {
		fmt.Println(err)
	}
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}
}

func (p *ProjectsHandler) getProjectByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	project, err := p.projectsStore.GetProjectByID(projectID)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.NewEncoder(w).Encode(project); err != nil {
		http.Error(w, "error getting project", http.StatusInternalServerError)
		return
	}
}

func (p *ProjectsHandler) getProjectsByParentID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	projects, err := p.projectsStore.GetProjectsByParentID(projectID)
	if err != nil {
		fmt.Println(err)
	}
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}
}

func (p *ProjectsHandler) getProjectsByAssigneeID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	projectID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	projects, err := p.projectsStore.GetProjectsByAssigneeID(projectID)
	if err != nil {
		fmt.Println(err)
	}
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}
}

func (p *ProjectsHandler) getProjectsByStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["status"]
	projectID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	projects, err := p.projectsStore.GetProjectsByStatus(projectID)
	if err != nil {
		fmt.Println(err)
	}
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}
}
