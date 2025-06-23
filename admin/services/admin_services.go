package services

import (
	"admin/models"
	"fmt"
	"net/http"
)



func CrearInstancia(instancia models.Instancia) error {
	// cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	// if err != nil {
	// 	return err
	// }

	// // Definir la configuración del contenedor
	// containerConfig := &container.Config{
	// 	Image: "my-docker-image", // Cambiar por la imagen que desees usar
	// }

	// // Crear el contenedor
	// ctx := context.Background()
	// containerName := instancia.Nombre // Usamos el nombre de la instancia como nombre del contenedor
	// createdContainer, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, containerName)
	// if err != nil {
	// 	return err
	// }

	// // Iniciar el contenedor
	// err = cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})
	// if err != nil {
	// 	return err
	// }

	 return nil
}

func EliminarInstancia(instancia models.Instancia) error {
	// cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	// if err != nil {
	// 	return err
	// }

	// // Eliminar el contenedor de forma forzada
	// ctx := context.Background()
	// err = cli.ContainerRemove(ctx, instancia.ID, types.ContainerRemoveOptions{Force: true})
	// if err != nil {
	// 	return err
	// }

	 return nil
}

func ObtenerInstancias() ([]models.Instancia, error) {
	// Lista de URLs y nombres de las instancias
	instancias := []struct {
		Nombre string
		URL    string
	}{
		{"backend_users", "http://backend_users:8082/users/1"},
		{"backend_cursos", "http://backend_courses:8083/cursos/all"},
		{"backend_subscriptions", "http://backend_subscriptions:8084/subscriptions/get/1"},
	}

	var resultado []models.Instancia

	// Iterar sobre cada instancia y verificar su estado
	for _, instancia := range instancias {
		estado, err := ProbarInstancia(instancia.URL)
		if err != nil {
			// Si hay error, marcar como "Inactivo" y continuar
			fmt.Printf("Error probando %s: %v\n", instancia.Nombre, err)
			estado = "Inactivo"
		}

		// Agregar la instancia con su estado al resultado
		resultado = append(resultado, models.Instancia{
			Nombre: instancia.Nombre,
			Estado: estado,
		})
	}

	return resultado, nil
}

// ProbarInstancia realiza una solicitud HTTP para verificar si la instancia está activa.
func ProbarInstancia(url string) (string, error) {
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
