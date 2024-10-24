package db

import (
	//"fmt"

	"log"

	// "backend/models/comentarios"
	// "backend/models/cursos"
	// "backend/models/files"
	// "backend/models/subscripciones"
	// "backend/models/users"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	// Reemplaza usuario, contraseña, host, puerto y nombre_basedatos con tus propios valores
	DBName := "arq-soft"         //Nombre de la base de datos local
	DBUser := "root"             //usuario de la base de datos, habitualmente root
	DBPass := "marcoslopez1719$" //password del root en la instalacion
	DBHost := "127.0.0.1"
	// Conecta a la base de datos
	dsn := DBUser + ":" + DBPass + "@tcp(" + DBHost + ":3306)/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Connection Failed to Open")
		log.Fatal(err)
	} else {
		log.Println("Connection Established")
	}
	// Asigna la instancia de la base de datos a la variable DB
	DB = db

	return nil
}





// import (
// 	"backend/models/comentarios"
// 	"backend/models/cursos"
// 	"backend/models/files"
// 	"backend/models/subscripciones"
// 	"backend/models/users"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func Connect() error {
// 	// Reemplaza usuario, contraseña, host, puerto y nombre_basedatos con tus propios valores
// 	// DBName := "arq-soft"         //Nombre de la base de datos local de ustedes
// 	// DBUser := "root"             //usuario de la base de datos, habitualmente root
// 	// DBPass := "marcoslopez1719$" //password del root en la instalacion
// 	// DBHost := "127.0.0.1"

// 	dbHost := os.Getenv("DB_HOST")
//     dbPort := os.Getenv("DB_PORT")
//     dbUser := os.Getenv("DB_USER")
//     dbPassword := os.Getenv("DB_PASS")
//     dbName := os.Getenv("DB_NAME")
// 	// Conecta a la base de datos
// 	//dsn := DBUser + ":" + DBPass + "@tcp(" + DBHost + ":3306)/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
//         dbUser, dbPassword, dbHost, dbPort, dbName)

	
// 	var err error
// 	var db *gorm.DB
// 	for i:= 0; i < 10 ; i++{
// 		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 		if err == nil {
//             fmt.Println("Connection Established")
//             break
//         }
//         fmt.Printf("Failed to connect to database. Retrying in 5 seconds... (%d/10)\n", i+1)
//         time.Sleep(5 * time.Second)
// 	}

// 	if err != nil {
// 		log.Println("Connection Failed to Open")
// 		log.Fatal(err)
// 	} else {
// 		log.Println("Connection Established")
// 	}
// 	// Asigna la instancia de la base de datos a la variable DB
// 	DB = db

// 	//userClient.Db = db
// 	err = AutoMigrate()
//     if err != nil {
//         log.Println("Failed to perform auto migration")
//         log.Fatal(err)
//     } else {
//         log.Println("Auto migration completed")
//     }


// 	return nil
// }

// func AutoMigrate() error {
// 	err := DB.AutoMigrate(&users.User{}, &cursos.Course{}, &files.File{}, &comentarios.Comentario{}, &subscripciones.Subscription{})
//     if err != nil {
// 		log.Printf("No se creo la tabla")
//         return err
//     }

// 	log.Println("Auto migration completed")

// 	return nil
// }

