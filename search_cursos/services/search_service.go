package services

import (
	"context"
	"fmt"
	cursosDAO "search_cursos/dao"
	cursosDomain "search_cursos/domain"
	"strconv"
)

type Repository interface {
	Index(ctx context.Context, curso cursosDAO.Curso) (string, error)
	Update(ctx context.Context, curso cursosDAO.Curso) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int, offset int) ([]cursosDAO.Curso, error)
}

type ExternalRepository interface {
	GetCursoByID(ctx context.Context, id string) (cursosDomain.Curso, error)
}

type Service struct {
	repository Repository
	cursosAPI  ExternalRepository
}

func NewService(repository Repository, cursosAPI ExternalRepository) Service {
	return Service{
		repository: repository,
		cursosAPI:  cursosAPI,
	}
}

func (service Service) Search(ctx context.Context, query string, offset int, limit int) ([]cursosDomain.Curso, error) {
	// Llamar al método Search del repositorio
	cursosDAOList, err := service.repository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error buscando cursos: %w", err)
	}

	// Convertir los cursos de la capa DAO a la capa del dominio
	cursosDomainList := make([]cursosDomain.Curso, 0)
	for _, curso := range cursosDAOList {
		cursosDomainList = append(cursosDomainList, cursosDomain.Curso{
			ID:          curso.ID,
			CourseName:  curso.CourseName,
			Description: curso.Description,
			Category:    curso.Category,
			Length:      curso.Length,
		})
	}

	return cursosDomainList, nil
}

func (service Service) HandleCursoNew(cursoNew cursosDomain.CursoNew){
	idString := strconv.FormatUint(uint64(cursoNew.CursoID), 10)
    switch cursoNew.Operation {
    case "CREATE", "UPDATE":
        // buscamos los detalles del servicio de cursos
		curso, err := service.cursosAPI.GetCursoByID(context.Background(), idString)
		if err != nil {
			fmt.Printf("Error getting hotel (%s) from API: %v\n", idString, err)
			return
		}

		cursoDAO := cursosDAO.Curso{
			ID: curso.ID,
			CourseName: curso.CourseName,
			Category: curso.Category,
			Length: curso.Length,
			Description: curso.Description,
		}

        // Si la operación es CREATE, indexamos el curso en Solr
        if cursoNew.Operation == "CREATE" {
            if _, err := service.repository.Index(context.Background(), cursoDAO); err != nil {
				fmt.Printf("Error indexing hotel (%s): %v\n", idString, err)
                return 
            }else{
				fmt.Println("Curso indexed successfully:", cursoNew.CursoID)
			}
        } else { // Si la operación es UPDATE, actualizamos el curso en Solr
            if err := service.repository.Update(context.Background(), cursoDAO); err != nil {
				fmt.Printf("Error updating hotel (%s): %v\n", idString, err)
                return 
            }else{
				fmt.Println("Curso updated successfully: ", cursoNew.CursoID)
			}
        }

    case "DELETE":
        if err := service.repository.Delete(context.Background(), idString); err != nil {
			fmt.Printf("Error deleting curso (%s): %v\n", idString, err)
            return 
        }else{
			fmt.Println("Curso deleted successfully: ", cursoNew.CursoID)
		}

    default:
        return 
    }
}
