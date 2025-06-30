package services

import (
	"admin/models"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Mapa para almacenar las instancias activas
var instanciasActivas = make(map[string]models.Instancia)

func CrearInstancia(instancia models.Instancia) error {
	// Generar un ID único si no se proporciona
	if instancia.ID == "" {
		instancia.ID = fmt.Sprintf("inst_%s_%d", instancia.Nombre, len(instanciasActivas)+1)
	}
	
	// Verificar si la instancia ya existe
	if _, existe := instanciasActivas[instancia.ID]; existe {
		return fmt.Errorf("la instancia con ID %s ya existe", instancia.ID)
	}
	
	// Verificar el estado de la instancia
	estado, err := ProbarInstancia(instancia.URL)
	if err != nil {
		estado = "Inactivo"
	}
	instancia.Estado = estado
	
	// Agregar la instancia al mapa
	instanciasActivas[instancia.ID] = instancia
	
	return nil
}

func EliminarInstancia(instancia models.Instancia) error {
	// Verificar si la instancia existe
	if _, existe := instanciasActivas[instancia.ID]; !existe {
		return fmt.Errorf("la instancia con ID %s no existe", instancia.ID)
	}
	
	// Eliminar la instancia del mapa
	delete(instanciasActivas, instancia.ID)
	
	return nil
}

func ObtenerInstancias() ([]models.Instancia, error) {
	// Crear cliente Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error al crear cliente Docker: %v", err)
	}
	defer cli.Close()

	// Obtener lista de contenedores
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar contenedores: %v", err)
	}

	var resultado []models.Instancia

	// Filtrar solo los contenedores de nuestro proyecto (que empiecen con "proyecto-arq-soft-2" o nombres específicos)
	proyectoContainers := []string{
		"backend_users", "backend_users_2", "backend_courses", "backend_subscriptions",
		"search_cursos", "admin", "frontend", "nginx", "mysql", "mongo", "solr", "rabbitmq", "memcached-container",
	}

	for _, container := range containers {
		// Verificar si el contenedor pertenece a nuestro proyecto
		esContenedorProyecto := false
		for _, nombreProyecto := range proyectoContainers {
			if strings.Contains(container.Names[0], nombreProyecto) || 
			   strings.Contains(container.Names[0], "proyecto-arq-soft-2") {
				esContenedorProyecto = true
				break
			}
		}

		if !esContenedorProyecto {
			continue
		}

		// Determinar el estado del contenedor
		estado := "Inactivo"
		if container.State == "running" {
			estado = "Activo"
		}

		// Obtener el nombre limpio del contenedor
		nombre := strings.TrimPrefix(container.Names[0], "/")
		
		// Determinar la URL basada en el nombre del contenedor
		var url string
		switch nombre {
		case "backend_users":
			url = "http://backend_users:8082/users/1"
		case "backend_users_2":
			url = "http://backend_users_2:8088/users/1"
		case "backend_courses":
			url = "http://backend_courses:8083/cursos/all"
		case "backend_subscriptions":
			url = "http://backend_subscriptions:8084/subscriptions/get/1"
		case "search_cursos":
			url = "http://search_cursos:8085/search"
		case "admin":
			url = "http://admin:8087/admin/services"
		case "frontend":
			url = "http://frontend:80"
		case "nginx":
			url = "http://nginx:80"
		case "mysql":
			url = "mysql://mysql:3306"
		case "mongo":
			url = "mongodb://mongo:27017"
		case "solr":
			url = "http://solr:8983/solr"
		case "rabbitmq":
			url = "http://rabbitmq:15672"
		case "memcached-container":
			url = "memcached://memcached-container:11211"
		default:
			url = ""
		}

		// Crear la instancia
		instancia := models.Instancia{
			ID:     container.ID[:12], // Usar los primeros 12 caracteres del ID del contenedor
			Nombre: nombre,
			Estado: estado,
			URL:    url,
		}

		resultado = append(resultado, instancia)
	}

	// Agregar las instancias dinámicas creadas por el usuario
	for _, instancia := range instanciasActivas {
		resultado = append(resultado, instancia)
	}

	return resultado, nil
}

// ProbarInstancia realiza una solicitud HTTP para verificar si la instancia está activa.
func ProbarInstancia(url string) (string, error) {
	// Si la URL está vacía o no es HTTP, no probar
	if url == "" || !strings.HasPrefix(url, "http") {
		return "Inactivo", fmt.Errorf("URL no válida para probar")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "Inactivo", fmt.Errorf("error al hacer la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Inactivo", nil
	}

	return "Activo", nil
}
