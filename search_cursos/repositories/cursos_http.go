package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	cursosDomain "search_cursos/domain"
)

type HTTPConfig struct {
	Host string
	Port string
}

type HTTP struct {
	baseURL func(courseID string) string
	config HTTPConfig
}

func NewHTTP(config HTTPConfig) HTTP {
	return HTTP{
		baseURL: func(courseID string) string {
			return fmt.Sprintf("http://%s:%s/cursos/get/%s", config.Host, config.Port, courseID)
		},
	}
}

func (repository HTTP) GetCursoByID(ctx context.Context, id string) (cursosDomain.Curso, error) {
	resp, err := http.Get(repository.baseURL(id))
	if err != nil {
		return cursosDomain.Curso{}, fmt.Errorf("Error fetching course (%s): %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return cursosDomain.Curso{}, fmt.Errorf("Failed to fetch course (%s): received status code %d", id, resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cursosDomain.Curso{}, fmt.Errorf("Error reading response body for course (%s): %w", id, err)
	}

	// Unmarshal the course details into the course struct
	var course cursosDomain.Curso
	if err := json.Unmarshal(body, &course); err != nil {
		return cursosDomain.Curso{}, fmt.Errorf("Error unmarshaling course data (%s): %w", id, err)
	}

	return course, nil
}

func (repo HTTP) GetAllCursos(ctx context.Context) ([]cursosDomain.Curso, error) {
	url := fmt.Sprintf("http://%s:%s/cursos/all", repo.config.Host, repo.config.Port) // Usamos repo.config para acceder a Host y Port
	fmt.Println("URL get all cursos: ", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al obtener cursos: %w", err)
	}
	defer resp.Body.Close()

	var cursos []cursosDomain.Curso
	if err := json.NewDecoder(resp.Body).Decode(&cursos); err != nil {
		return nil, fmt.Errorf("error al decodificar los cursos: %w", err)
	}

	return cursos, nil
}