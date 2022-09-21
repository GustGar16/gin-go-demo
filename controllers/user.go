package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Traemos la cleccion que utilizaremos
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// Nueva validacion de la estructura
var validate = validator.New()

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		var users []models.User
		defer cancel()

		result, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GeneralResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer result.Close(ctx)
		for result.Next(ctx) {
			var singleUser models.User
			if err = result.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.GeneralResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}
			users = append(users, singleUser)
		}
		c.JSON(http.StatusOK, responses.GeneralResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}})

	}
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		var user models.User
		defer cancel()

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.GeneralResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Se valida la estructura del usuario
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GeneralResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		//Creacion de un nuevo usuario
		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Location: user.Location,
			Title:    user.Title,
		}
		//Almacenamiento del usuario en la bd
		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GeneralResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.GeneralResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		//Obtenemos el id que viene como parametro
		userId := c.Param("userId")
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)
		cond := bson.M{
			"_id": objId,
		}
		err := userCollection.FindOne(ctx, cond).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GeneralResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.GeneralResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.User
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(userId)

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.GeneralResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GeneralResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"name":     user.Name,
			"location": user.Location,
			"title":    user.Title,
		}
		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GeneralResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.GeneralResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}
		c.JSON(http.StatusOK, responses.GeneralResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})

	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)
		result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GeneralResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.GeneralResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}})
			return
		}
		c.JSON(http.StatusOK, responses.GeneralResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "user successfully deleted"}})
	}
}
