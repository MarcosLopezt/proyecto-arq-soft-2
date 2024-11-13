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
