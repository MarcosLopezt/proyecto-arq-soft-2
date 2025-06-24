package db

import (
	"fmt"
	"log"
	"time"
	models "users/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	// Obtiene las variables de entorno
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASS")
	// dbName := os.Getenv("DB_NAME")

	// log.Printf("DB_HOST: %s, DB_PORT: %s, DB_USER: %s, DB_NAME: %s", 
    //     dbHost, dbPort, dbUser, dbName)
	// // Formatea la cadena de conexión
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	dbUser, dbPassword, dbHost, dbPort, dbName)
	dsn := "root:marcoslopez1719$@tcp(localhost:3306)/arq-soft?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	var db *gorm.DB
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Conexión establecida a MySQL")
			break
		}
		fmt.Printf("Error al conectar a la base de datos. Reintentando en 5 segundos... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Println("No se pudo conectar a la base de datos")
		log.Fatal(err)
	} else {
		log.Println("Conexión establecida a MySQL")
	}

	// Asigna la instancia de la base de datos a la variable DB
	DB = db

	// Ejecuta la migración automática para la tabla de usuarios
	err = AutoMigrate()
	if err != nil {
		log.Println("Error en la migración automática")
		log.Fatal(err)
	} else {
		log.Println("Migración automática completada")
	}

	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("No se pudo crear la tabla de usuarios")
		return err
	}

	log.Println("Migración automática completada para la tabla de usuarios")
	return nil
}
