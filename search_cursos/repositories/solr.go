package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	cursos "search_cursos/dao"
	"strconv"

	"github.com/stevenferrer/solr-go"
)

type SolrConfig struct {
	Host       string
	Port       string
	Collection string
	BaseURL string
}

type Solr struct {
	Client     *solr.JSONClient
	Collection string
	BaseURL string
}

func NewSolr(config SolrConfig) Solr {
	baseURL := fmt.Sprintf("http://%s:%s/solr/%s",config.Host, config.Port, config.Collection)

	client := solr.NewJSONClient(baseURL)
	return Solr{
		Client:     client,
		Collection: config.Collection,
		BaseURL: baseURL,
	}
}

func (searchEngine Solr) Index(ctx context.Context, curso cursos.Curso) (string, error) {
	// Convertir el ID a string antes de indexarlo
	doc := map[string]interface{}{
		"id":          fmt.Sprintf("%d", curso.ID),  // Convertir ID a string
		"course_name": curso.CourseName,
		"description": curso.Description,
		"category":    curso.Category,
		"length":      curso.Length,
	}

	indexRequest := map[string]interface{}{
		"add": []interface{}{doc},
	}

	body, err := json.Marshal(indexRequest)
	if err != nil {
		return "", fmt.Errorf("error serializando documento del curso: %w", err)
	}
	
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error al indexar curso: %w", err)
	}
	if resp.Error != nil {
		return "", fmt.Errorf("error en la respuesta al indexar curso: %v", resp.Error)
	}

	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return "", fmt.Errorf("error al confirmar cambios en Solr: %w", err)
	}

	// Retornar el ID como string también
	return fmt.Sprintf("%d", curso.ID), nil
}

func (searchEngine Solr) Search(ctx context.Context, query string, limit int, offset int) ([]cursos.Curso, error) {
    // Construir manualmente la URL para /select
	selectURL := fmt.Sprintf("%s/select?q=course_name:%s&rows=%d&start=%d&wt=json",
	searchEngine.BaseURL, query, limit, offset)
	fmt.Println("SELECT URL: ", selectURL)
    // Crear la solicitud HTTP
    req, err := http.NewRequestWithContext(ctx, "GET", selectURL, nil)
    if err != nil {
        return nil, fmt.Errorf("error construyendo solicitud: %w", err)
    }

    // Ejecutar la solicitud
	fmt.Println("SOLICITUD: ", req)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error al realizar la solicitud: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("error en la respuesta: %s", resp.Status)
    }

    // Leer y parsear la respuesta
    var result struct {
        Response struct {
            NumFound int                `json:"numFound"`
            Docs     []map[string]interface{} `json:"docs"`
        } `json:"response"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("error parseando respuesta: %w", err)
    }

    // Mapear documentos a la estructura Curso
    var cursosList []cursos.Curso
    for _, doc := range result.Response.Docs {
        curso := cursos.Curso{
            ID:          uint(getIntField(doc, "id")),
            CourseName:  getStringField(doc, "course_name"),
            Description: getStringField(doc, "description"),
            Category:    getStringField(doc, "category"),
            Length:      getIntField(doc, "length"),
        }
        cursosList = append(cursosList, curso)
    }

    return cursosList, nil
}



func getStringField(doc map[string]interface{}, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	return ""
}

func getIntField(doc map[string]interface{}, field string) int {
	if val, ok := doc[field].(float64); ok {
		return int(val)
	}
	if valStr, ok := doc[field].(string); ok {
        if val, err := strconv.Atoi(valStr); err == nil {
            return val
        }
    }
	return 0
}

func (s Solr) Update(ctx context.Context, curso cursos.Curso) error {
	// Convertir el curso a JSON para Solr
	updateRequest := map[string]interface{}{
		"add": []interface{}{
			map[string]interface{}{
				"id":          curso.ID,
				"course_name": curso.CourseName,
				"description": curso.Description,
				"category":    curso.Category,
				"length":      curso.Length,
			},
		},
	}

	// Convertir la solicitud a JSON
	body, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("error serializando solicitud de actualización: %w", err)
	}

	// Enviar la solicitud de actualización a Solr
	resp, err := s.Client.Update(ctx, s.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error al actualizar curso en Solr: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("error en la respuesta al actualizar curso en Solr: %v", resp.Error)
	}

	// Confirmar cambios en Solr
	if err := s.Client.Commit(ctx, s.Collection); err != nil {
		return fmt.Errorf("error al confirmar actualización en Solr: %w", err)
	}

	return nil
}

func (s Solr) Delete(ctx context.Context, id string) error {
	// Crear la solicitud de eliminación en formato JSON
	deleteRequest := map[string]interface{}{
		"delete": []interface{}{
			map[string]interface{}{"id": id},
		},
	}

	// Convertir la solicitud a JSON
	body, err := json.Marshal(deleteRequest)
	if err != nil {
		return fmt.Errorf("error serializando solicitud de eliminación: %w", err)
	}

	// Enviar la solicitud de eliminación a Solr
	resp, err := s.Client.Update(ctx, s.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error al eliminar curso en Solr: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("error en la respuesta al eliminar curso en Solr: %v", resp.Error)
	}

	// Confirmar cambios en Solr
	if err := s.Client.Commit(ctx, s.Collection); err != nil {
		return fmt.Errorf("error al confirmar eliminación en Solr: %w", err)
	}

	return nil
}
