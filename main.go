package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/mbtamuli/crispy-octo-doodle/database"
)

func CustomGenerator() {
	_ = faker.AddProvider("customID", func(v reflect.Value) (interface{}, error) {
		var str strings.Builder
		str.WriteString(strconv.Itoa(1000 + rand.Intn(9000)))
		str.WriteString("-")
		str.WriteString(strconv.Itoa(1000 + rand.Intn(9000)))
		str.WriteString("-")
		str.WriteString(strconv.Itoa(1000 + rand.Intn(9000)))
		return str.String(), nil
	})
}

type IDCardData struct {
	Name         string `faker:"name" json:"name"`
	Gender       string `faker:"gender" json:"gender"`
	DateOfBirth  string `faker:"date" json:"dateOfBirth"`
	IDNumber     string `faker:"customID" json:"idNumber"`
	AddressLine1 string `faker:"-" json:"addressLine1"`
	AddressLine2 string `faker:"-" json:"addressLine2"`
	Pincode      string `faker:"-" json:"pincode"`
}

var keys = map[string]string{
	"1578859fae": "58896a4d568e44d1843a",
}

func signMessage(msg []byte, key []byte) []byte {
	hmacHasher := hmac.New(sha256.New, key)
	hmacHasher.Write(msg)
	return []byte(base64.StdEncoding.EncodeToString(hmacHasher.Sum(nil)))
}

func HMACValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientHMACdigest := []byte(strings.Split(c.GetHeader("Authorization"), " ")[1])
		clientPublicKey := c.GetHeader("public_key")

		serverHMACdigest := signMessage([]byte(clientPublicKey), []byte(keys[clientPublicKey]))

		if !hmac.Equal(serverHMACdigest, clientHMACdigest) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errorMessage": "invalid access or secret key",
			})
		}

		c.Next()
	}
}

func setupRouter() *gin.Engine {
	// DB setup
	ctx := context.Background()
	db, err := sql.Open("mysql", "dbuser:mySuperSecret123@/ekyc?parseTime=true")
	if err != nil {
		panic(err)
	}

	queries := database.New(db)

	// Gin
	r := gin.Default()
	apiV1 := r.Group("/api/v1")

	apiV1.POST("signup", func(c *gin.Context) {
		var json struct {
			Name  string `json:"name" binding:"required,ascii"`
			Email string `json:"email" binding:"required,email"`
			Plan  string `json:"plan" binding:"required,oneof=basic advanced enterprise"`
		}

		err := c.Bind(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
			return
		}

		plansResult, err := queries.ListPlans(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
			return
		}

		plans := make(map[string]int32)
		for _, plan := range plansResult {
			plans[plan.Name] = plan.ID
		}

		// fmt.Println(plans)

		// create a user
		userID, err := queries.CreateClient(ctx, database.CreateClientParams{
			Name:   json.Name,
			Email:  json.Email,
			PlanID: sql.NullInt32{Int32: plans[json.Plan], Valid: true},
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
			return
		}

		// create keys

		accessKey := strings.ReplaceAll(uuid.NewString(), "-", "")[:10]
		secretKey := strings.ReplaceAll(uuid.NewString(), "-", "")[:20]

		_, err = queries.CreateKey(ctx, database.CreateKeyParams{
			ClientID: int32(userID),
			Access:   accessKey,
			Secret:   secretKey,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"accessKey": accessKey, "secretKey": secretKey})
	})

	authenticated := apiV1.Group("/")
	authenticated.Use(HMACValidator())

	authenticated.POST("image", func(c *gin.Context) {
		var json struct {
			ImageType string `json:"type" binding:"required,oneof=face id_card"`
			File      string `json:"file" binding:"required,base64"`
		}

		err := c.Bind(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
			return
		}

		switch json.ImageType {
		case "face", "id_card":
			c.IndentedJSON(http.StatusOK, gin.H{"id": uuid.NewString()})
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "invalid type, supported types are face or id_card"})
		}
	})

	authenticated.POST("face-match", func(c *gin.Context) {
		var json struct {
			Image1 string `json:"image1" binding:"required,uuid4"`
			Image2 string `json:"image2" binding:"required,uuid4"`
		}

		err := c.Bind(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "invalid or missing image id"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"score": rand.Intn(100)})
	})

	authenticated.POST("ocr", func(c *gin.Context) {
		var json struct {
			Image1 string `json:"image1" binding:"required,uuid4"`
		}

		err := c.Bind(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "invalid or missing image id"})
			return
		}

		fakeAddress := faker.GetRealAddress()

		fmt.Println(fakeAddress.Address)

		data := IDCardData{
			AddressLine1: fakeAddress.Address,
			AddressLine2: fakeAddress.City + ", " + fakeAddress.State,
			Pincode:      fakeAddress.PostalCode,
		}
		err = faker.FakeData(&data)
		if err != nil {
			fmt.Println(err)
		}

		c.IndentedJSON(http.StatusOK, data)
	})
	return r
}

func main() {
	CustomGenerator()
	r := setupRouter()
	r.Run(":8080")
}
