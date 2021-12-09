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

func GetAllEvents(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var events []models.Event

	filter := bson.M{}
	findOptions := options.Find()

	if s := c.Query("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"title": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
				{
					"link": bson.M{
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
		var event models.Event
		cursor.Decode(&event)
		events = append(events, event)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      events,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetEvent(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var event models.Event
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := eventCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Event Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&event)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Event Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    event,
		"success": true,
	})
}

func AddEvent(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	event := new(models.Event)

	if err := c.BodyParser(event); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := eventCollection.InsertOne(ctx, event)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Event failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Event inserted successfully",
	})

}

func UpdateEvent(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	event := new(models.Event)

	if err := c.BodyParser(event); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Event not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": event,
	}
	_, err = eventCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Event failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Event updated successfully",
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Event not found",
			"error":   err,
		})
	}
	_, err = eventCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Event failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Event deleted successfully",
	})
}

func JoinEvent(c *fiber.Ctx) error {
	eventCollection := config.MI.DB.Collection("events")
	userCollection := config.MI.DB.Collection("users")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	var event models.Event
	var userActivity models.UserActivity
	var eventActivity models.EventActivity

	eventObjId, err := primitive.ObjectIDFromHex(c.Params("eventId"))
	userObjId, err := primitive.ObjectIDFromHex(c.Params("userId"))

	findEventResult := eventCollection.FindOne(ctx, bson.M{"_id": eventObjId})
	findUserResult := userCollection.FindOne(ctx, bson.M{"_id": userObjId})

	err = findEventResult.Decode(&event)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get Event by Id",
			"error":   err,
		})
	}
	err = findUserResult.Decode(&user)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get User by Id",
			"error":   err,
		})
	}
	userActivity = models.UserActivity{EventId: event.ID, EventTitle: event.Title, Attende: "no"}
	eventActivity = models.EventActivity{UserId: user.ID, Email: user.Email, Attende: "no"}

	event.Participant = append(event.Participant, eventActivity)
	user.Activities = append(user.Activities, userActivity)

	updateEvent := bson.M{
		"$set": event,
	}

	updateUser := bson.M{
		"$set": user,
	}

	_, err = eventCollection.UpdateOne(ctx, bson.M{"_id": eventObjId}, updateEvent)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to add participant",
			"error":   err.Error(),
		})
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": userObjId}, updateUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to add userActivity",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User joined successfully",
	})
}

// func JoinEvent(c *fiber.Ctx) error {
// 	eventCollection := config.MI.DB.Collection("events")
// 	userCollection := config.MI.DB.Collection("events")
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

// 	var participant models.User
// 	var event models.Event
// 	objId, err := primitive.ObjectIDFromHex(c.Params("event-id"))
// 	findResult := eventCollection.FindOne(ctx, bson.M{"_id": objId})
// 	err = findResult.Decode(&event)
// 	if err := c.BodyParser(&participant); err != nil {
// 		log.Println(err)
// 		return c.Status(400).JSON(fiber.Map{
// 			"success": false,
// 			"message": "Failed to parse body",
// 			"error":   err,
// 		})
// 	}
// 	event.Participant = append(event.Participant, participant)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"success": false,
// 			"message": "Failed to get event ID",
// 			"error":   err.Error(),
// 		})
// 	}
// 	updateEvent := bson.M{
// 		"$set": event,
// 	}

// 	_, err = eventCollection.UpdateOne(ctx, bson.M{"_id": objId}, updateEvent)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"success": false,
// 			"message": "Failed to add participant",
// 			"error":   err.Error(),
// 		})
// 	}

// 	updateUser := bson.M{
// 		"$set": participant,
// 	}
// 	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objId}, updateUser)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"success": false,
// 			"message": "Failed to add participant",
// 			"error":   err.Error(),
// 		})
// 	}
// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"success": true,
// 		"message": "User joined successfully",
// 	})
// }
