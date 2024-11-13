package dao

import (
	"context"
	models "cursos/models"
	"errors"
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCourseDAO struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoCourseDAO(client *mongo.Client, database, collection string) *MongoCourseDAO {
	return &MongoCourseDAO{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (dao *MongoCourseDAO) getNextID(ctx context.Context) (uint, error) {
	counterCollection := dao.client.Database(dao.database).Collection("counter")
	filter := bson.M{"_id": "course_id"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		Seq uint `bson:"seq"`
	}

	log.Println("Ejecutando FindOneAndUpdate con filtro:", filter, "y actualización:", update)

	err := counterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("error generating ID: %w", err)
	}

	log.Println("Nuevo ID generado:", result.Seq)

	return result.Seq, nil
}


func (dao *MongoCourseDAO) CreateCourse(ctx context.Context, course *models.Course) error {

	nextID, err := dao.getNextID(ctx)
    if err != nil {
        return fmt.Errorf("error generating course ID: %w", err)
    }
    course.ID = nextID

	collection := dao.client.Database(dao.database).Collection(dao.collection)
	_, err = collection.InsertOne(ctx, course)
	if err != nil {
		return fmt.Errorf("error creating course: %w", err)
	}
	return nil
}

func (dao *MongoCourseDAO) GetCourseByName(ctx context.Context, courseName string) ([]models.Course, error) {
	collection := dao.client.Database(dao.database).Collection(dao.collection)

	filter := bson.M{
		"$or": []bson.M{
			{"course_name": bson.M{"$regex": courseName, "$options": "i"}},
			{"category": bson.M{"$regex": courseName, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding courses: %w", err)
	}
	defer cursor.Close(ctx)

	var courses []models.Course
	if err = cursor.All(ctx, &courses); err != nil {
		return nil, fmt.Errorf("error decoding courses: %w", err)
	}

	return courses, nil
}

func (dao *MongoCourseDAO) GetAllCourses(ctx context.Context)([]models.Course, error){
	collection := dao.client.Database(dao.database).Collection(dao.collection)

	filter := bson.M{} 

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding courses: %w", err)
	}
	defer cursor.Close(ctx)

	var courses []models.Course
	if err = cursor.All(ctx, &courses); err != nil {
		return nil, fmt.Errorf("error decoding courses: %w", err)
	}

	return courses, nil
}
func (dao *MongoCourseDAO) GetCourseByID(ctx context.Context, id string) (*models.Course, error) {
	collection := dao.client.Database(dao.database).Collection(dao.collection)
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("error converting id to int: %w", err)
	}
	var course models.Course
	if err := collection.FindOne(ctx, bson.M{"_id": idNum}).Decode(&course); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding course: %w", err)
	}

	return &course, nil
}

func (dao *MongoCourseDAO) UpdateCourse(ctx context.Context, request *models.UpdateCourseRequest) (*models.UpdateCourseResponse, error) {
	collection := dao.client.Database(dao.database).Collection(dao.collection)

	update := bson.M{
		"$set": bson.M{
			"course_name":  request.CourseName,
			"category":     request.Category,
			"description":  request.Description,
			"length":       request.Length,
			"last_updated": request.LastUpdated,
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": request.ID}, update)
	if err != nil {
		return nil, fmt.Errorf("error updating course: %w", err)
	}

	response := &models.UpdateCourseResponse{
		ID:          request.ID,
		CourseName:  request.CourseName,
		Category:    request.Category,
		Description: request.Description,
		Length:      request.Length,
	}
	return response, nil
}


func (dao *MongoCourseDAO) DeleteCourse(ctx context.Context, request *models.DeleteCourseRequest) (*models.DeleteCourseResponse, error) {
	collection := dao.client.Database(dao.database).Collection(dao.collection)

	idString := fmt.Sprintf("%d", request.ID)

	objectID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return nil, fmt.Errorf("error converting id to object ID: %w", err)
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, fmt.Errorf("error deleting course: %w", err)
	}

	response := &models.DeleteCourseResponse{
		Message: "Se eliminó el curso de manera exitosa",
	}

	return response, nil
}



