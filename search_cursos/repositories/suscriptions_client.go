package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SubscriptionsClient struct {
	BaseURL string
}

func NewSubscriptionsClient(baseURL string) SubscriptionsClient {
	return SubscriptionsClient{
		BaseURL: baseURL,
	}
}

type SubscriptionsResponse struct {
	CursoID   string `json:"curso_id"`
	SubsCount int    `json:"subs_count"`
}

func (c SubscriptionsClient) GetSubsCount(ctx context.Context, cursoID uint) (int, error) {
	url := fmt.Sprintf("%s/subscriptions/get/curso/%d", c.BaseURL, cursoID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creando request: %w", err)
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error haciendo request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status no OK: %s", resp.Status)
	}
	
	var result SubscriptionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error parseando respuesta: %w", err)
	}
	
	return result.SubsCount, nil
}

// GetSubsByCursoId obtiene la cantidad de inscriptos para un curso específico
func (c SubscriptionsClient) GetSubsByCursoId(ctx context.Context, cursoID uint) (int, error) {
	return c.GetSubsCount(ctx, cursoID)
} 