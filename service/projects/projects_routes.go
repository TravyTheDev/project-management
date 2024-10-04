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
}

func (p *ProjectsHandler) createProject(w http.ResponseWriter, r *http.Request) {
	project := &types.Project{}
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		fmt.Println(err)
	}
	if err := p.projectsStore.CreateProject(project); err != nil {
		fmt.Println(err)
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
