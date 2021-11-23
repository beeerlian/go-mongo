package controllers

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/beeerlian/go-mongo/config"
	"github.com/beeerlian/go-mongo/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUser(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var users []models.User

	filter := bson.M{}
	findOptions := options.Find()

	if s := c.Query("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"name": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
				{
					"email": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
			},
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
	var limit int64 = int64(limitVal)

	total, _ := eventCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := eventCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Events Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var user models.User
		cursor.Decode(&user)
		users = append(users, user)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      users,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func UserRegistration(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := eventCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "User registration failed",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "User registered successfully",
	})
}

func LoginWithEmail(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var userLoginData models.User
	var user models.User
	if err := c.BodyParser(&userLoginData); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}
	findResult := eventCollection.FindOne(ctx, bson.M{"email": userLoginData.Email})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Email doesn't exist",
			"error":   err,
		})
	}
	err := findResult.Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Failed Decode user data result",
			"error":   err,
		})
	}
	if user.Password != userLoginData.Password {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Email and password doesn't match",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User Loged-in successfully",
	})

}
func DeleteUser(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"error":   err,
		})
	}
	_, err = eventCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "User failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}

func JoinEvent(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var participant models.User
	var event models.Event
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := eventCollection.FindOne(ctx, bson.M{"_id": objId})
	err = findResult.Decode(&event)
	if err := c.BodyParser(&participant); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}
	event.Participant = append(event.Participant, participant)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get event ID",
			"error":   err.Error(),
		})
	}
	update := bson.M{
		"$set": event,
	}
	_, err = eventCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to add participant",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Participant updated successfully",
	})
}
