package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	baseURL := fmt.Sprintf("http://%s:%s/solr", config.Host, config.Port)


	client := solr.NewJSONClient(baseURL)
	log.Printf("SolR Base URL: %s", baseURL)

	return Solr{
		Client:     client,
		Collection: config.Collection,
		BaseURL: baseURL,
	}
}



func (searchEngine Solr) Search(ctx context.Context, query string, limit int, offset int, availableOnly bool) ([]cursos.Curso, error) {
    if query == "" || query == "*:*" {
        query = "*:*"
    } else {
		query = fmt.Sprintf("%s", url.QueryEscape(query))

    }
    fq := ""
    if availableOnly {
        // No podemos filtrar por disponibilidad directamente en SolR
        // porque necesitamos calcular dinámicamente con el endpoint de subscripciones
        // Por ahora, obtenemos todos los cursos y filtramos después
        fq = ""
    }
	selectURL := fmt.Sprintf("%s/%s/select?q=%s&qf=course_name^2+description^1+category^1&defType=edismax&rows=%d&start=%d&wt=json%s",
    searchEngine.BaseURL, searchEngine.Collection, url.QueryEscape(query), limit, offset, fq)

    log.Printf("SELECT URL: %s", selectURL)
    req, err := http.NewRequestWithContext(ctx, "GET", selectURL, nil)
    if err != nil {
        return nil, fmt.Errorf("error construyendo solicitud: %w", err)
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error al realizar la solicitud: %w", err)
    }
    defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Respuesta Solr: %s", respBody)
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("error en la respuesta: %s", resp.Status)
    }
    var result struct {
        Response struct {
            NumFound int                `json:"numFound"`
            Docs     []map[string]interface{} `json:"docs"`
        } `json:"response"`
    }
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %w", err)
	}
    
    var cursosList []cursos.Curso
    subscriptionsClient := NewSubscriptionsClient("http://backend_subscriptions:8084")
    
    for _, doc := range result.Response.Docs {
        curso := cursos.Curso{
            ID:          uint(getIntField(doc, "id")),
            CourseName:  getStringField(doc, "course_name"),
            Description: getStringField(doc, "description"),
            Category:    getStringField(doc, "category"),
            Length:      getIntField(doc, "length"),
            Cupos:       getIntField(doc, "cupos"),
        }
        
        // Si se requiere filtrar por disponibilidad, calcular dinámicamente
        if availableOnly {
            inscriptos, err := subscriptionsClient.GetSubsByCursoId(ctx, curso.ID)
            if err != nil {
                log.Printf("Warning: No se pudo obtener inscriptos para curso %d: %v", curso.ID, err)
                continue // Saltar este curso si no se puede obtener la información
            }
            
            disponibles := curso.Cupos - inscriptos
            if disponibles <= 0 {
                continue // Saltar cursos sin cupos disponibles
            }
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

func (searchEngine Solr) Index(ctx context.Context, curso cursos.Curso) (string, error) {
	// Obtener la cantidad de inscriptos desde la API de subscripciones
	subscriptionsClient := NewSubscriptionsClient("http://backend_subscriptions:8084")
	inscriptos, err := subscriptionsClient.GetSubsByCursoId(ctx, curso.ID)
	if err != nil {
		log.Printf("Warning: No se pudo obtener inscriptos para curso %d: %v", curso.ID, err)
		inscriptos = 0 // Valor por defecto si no se puede obtener
	}

	// Calcular cupos disponibles
	disponibles := curso.Cupos - inscriptos
	if disponibles < 0 {
		disponibles = 0
	}

	doc := map[string]interface{}{
		"id":          fmt.Sprintf("%d", curso.ID),
		"course_name": curso.CourseName,
		"description": curso.Description,
		"category":    curso.Category,
		"length":      curso.Length,
		"cupos":       curso.Cupos,
	}
	indexRequest := map[string]interface{}{
		"add": []interface{}{doc},
	}
	body, err := json.Marshal(indexRequest)
	if err != nil {
		return "", fmt.Errorf("error serializando documento del curso: %w", err)
	}
	
	updateURL := fmt.Sprintf("%s/%s/update?commit=true", searchEngine.BaseURL, searchEngine.Collection)

	req, err := http.NewRequestWithContext(ctx, "POST", updateURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error creando request manual: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error haciendo POST manual: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("respuesta error: %s - body: %s", resp.Status, respBody)
	}

	return fmt.Sprintf("%d", curso.ID), nil
}


func (s Solr) Update(ctx context.Context, curso cursos.Curso) error {
    // Obtener la cantidad de inscriptos desde la API de subscripciones
    subscriptionsClient := NewSubscriptionsClient("http://backend_subscriptions:8084")
    inscriptos, err := subscriptionsClient.GetSubsByCursoId(ctx, curso.ID)
    if err != nil {
        log.Printf("Warning: No se pudo obtener inscriptos para curso %d: %v", curso.ID, err)
        inscriptos = 0 // Valor por defecto si no se puede obtener
    }

    // Calcular cupos disponibles
    disponibles := curso.Cupos - inscriptos
    if disponibles < 0 {
        disponibles = 0
    }

    // Construir el documento a actualizar
    doc := map[string]interface{}{
        "id":          fmt.Sprintf("%d", curso.ID),
        "course_name": curso.CourseName,
        "description": curso.Description,
        "category":    curso.Category,
        "length":      curso.Length,
        "cupos":       curso.Cupos,
    }
    updateRequest := map[string]interface{}{
        "add": []interface{}{doc},
    }
    body, err := json.Marshal(updateRequest)
    if err != nil {
        return fmt.Errorf("error serializando solicitud de actualización: %w", err)
    }

    // Construir la URL de actualización con commit=true
    updateURL := fmt.Sprintf("%s/%s/update?commit=true", s.BaseURL, s.Collection)

    // Crear la solicitud HTTP
    req, err := http.NewRequestWithContext(ctx, "POST", updateURL, bytes.NewReader(body))
    if err != nil {
        return fmt.Errorf("error creando request de actualización: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la solicitud
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("error haciendo POST de actualización: %w", err)
    }
    defer resp.Body.Close()

    // Verificar el código de estado
    if resp.StatusCode != http.StatusOK {
        respBody, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("respuesta error al actualizar: %s - body: %s", resp.Status, respBody)
    }

    return nil
}
func (s Solr) Delete(ctx context.Context, id string) error {
    // Construir la solicitud de eliminación en formato JSON
    deleteRequest := map[string]interface{}{
        "delete": map[string]interface{}{
            "id": id,
        },
    }
    body, err := json.Marshal(deleteRequest)
    if err != nil {
        return fmt.Errorf("error serializando solicitud de eliminación: %w", err)
    }
    log.Printf("DEBUG: Solicitud de eliminación en Solr: %+v", deleteRequest)

    // Construir la URL de eliminación con commit=true
    deleteURL := fmt.Sprintf("%s/%s/update?commit=true", s.BaseURL, s.Collection)

    // Crear la solicitud HTTP
    req, err := http.NewRequestWithContext(ctx, "POST", deleteURL, bytes.NewReader(body))
    if err != nil {
        return fmt.Errorf("error creando request de eliminación: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la solicitud
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("error haciendo POST de eliminación: %w", err)
    }
    defer resp.Body.Close()

    // Verificar el código de estado
    if resp.StatusCode != http.StatusOK {
        respBody, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("respuesta error al eliminar: %s - body: %s", resp.Status, respBody)
    }

    return nil
}
