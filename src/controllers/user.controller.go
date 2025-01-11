package controllers

import (
	"context"

	"github.com/DivyanshTiwari001/cjbackend/src/db"
	"github.com/DivyanshTiwari001/cjbackend/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func Register(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate user input
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if user already exists
	count, err := db.UserCollection.CountDocuments(context.TODO(), bson.M{"email": user.Email})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	}

	// Hash password
	err = user.HashPassword()

	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Failed to create user"})
	}

	// generate accessToken
	accessToken,err := user.GenerateAccessToken()

	if err != nil{
		log.Error("access token not created : %v\n",err)
	}

	// Insert user into the database
	insertResult, err := db.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	var createdUser models.User
	err = db.UserCollection.FindOne(context.TODO(),bson.M{"_id": insertResult.InsertedID}).Decode(&createdUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve created user"})
	}

	// Return the created user as a JSON response

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",        // Cookie name
		Value:    accessToken,                  // JWT token value
		HTTPOnly: true,                   // Prevent JavaScript from accessing the cookie
		Secure:   false,                  // Set to true if using HTTPS
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"statusCode":201,"data":createdUser.ToResponse(),"message":"user created successfully"})
}


// Login handles user authentication.
func Login(c *fiber.Ctx) error {
	var input models.LoginInfo

	// Parse the request body into the input struct
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate email and password fields
	if err := input.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "both email and password are required"})
	}

	var user models.User

	// Find the user by email in the database
	err := db.UserCollection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	validPassword,err := user.ComparePassword(&input)

	if !validPassword || err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"Invalid credentials"})
	}

	// generate accessToken
	accessToken,err := user.GenerateAccessToken()

	if err != nil{
		log.Error("access token not created : %v\n",err)
	}

	var loggedInUser models.User
	err = db.UserCollection.FindOne(context.TODO(),bson.M{"_id": user.ID}).Decode(&loggedInUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user"})
	}

	// Return the created user as a JSON response

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",        // Cookie name
		Value:    accessToken,                  // JWT token value
		HTTPOnly: true,                   // Prevent JavaScript from accessing the cookie
		Secure:   false,                  // Set to true if using HTTPS
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"statusCode":200,"data":loggedInUser.ToResponse(),"message":"user loggedIn successfully"})
}

func GetUser(c *fiber.Ctx) error {
    userID := c.Locals("userId").(string)

    if userID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
    }

	objectId, err := primitive.ObjectIDFromHex(userID)

	if err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
	}


	var user models.User
	err = db.UserCollection.FindOne(context.TODO(),bson.M{"_id": objectId}).Decode(&user)

	if err!=nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
	}

	// generate accessToken
	accessToken,err := user.GenerateAccessToken()

	if err != nil{
		log.Error("access token not created : %v\n",err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",        // Cookie name
		Value:    accessToken,                  // JWT token value
		HTTPOnly: true,                   // Prevent JavaScript from accessing the cookie
		Secure:   false,                  // Set to true if using HTTPS
	})

    // Now you can use `userID` in your handler logic
    return c.Status(fiber.StatusOK).JSON(fiber.Map{"statusCode":200,"data":user.ToResponse(),"message":"user fetched successfully"})
}



func Logout(c *fiber.Ctx) error{
	userID := c.Locals("userId").(string)

    if userID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
    }

	objectId, err := primitive.ObjectIDFromHex(userID)

	if err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
	}


	var user models.User
	err = db.UserCollection.FindOne(context.TODO(),bson.M{"_id": objectId}).Decode(&user)

	if err!=nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",        // Cookie name
		Value:    "",                  // JWT token value
		HTTPOnly: true,                   // Prevent JavaScript from accessing the cookie
		Secure:   false,                  // Set to true if using HTTPS
	})

    // Now you can use `userID` in your handler logic
    return c.Status(fiber.StatusOK).JSON(fiber.Map{"statusCode":200,"message":"user loggedOut successfully"})
}