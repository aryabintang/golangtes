package controllers

import (
	"context"
	"fmt"
	"golang_cms/configs"
	"golang_cms/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var metaCollection *mongo.Collection = configs.GetCollection(configs.DB, "Meta")
var val = validator.New()

// membuat atau memasukan meta  baru
func Createmeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var meta models.Meta
	defer cancel()

	//val the request body
	if err := c.Bind(&meta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  400,
			"Message": err.Error(),
		})
		return
	}

	//menggunakan validasi untuk digunakan validasi required
	if validationErr := val.Struct(&meta); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  400,
			"Message": validationErr.Error(),
		})
		return
	}

	newMeta := models.Meta{
		Id:              primitive.NewObjectID(),
		Meta_Title:      meta.Meta_Title,
		Meta_Url:        meta.Meta_Url,
		Meta_Descrption: meta.Meta_Descrption,
	}
	result, err := metaCollection.InsertOne(ctx, newMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  500,
			"Message": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Data":    result,
		"Status":  200,
		"Message": "Data Berhasil Dibuat",
	})
}

// mengambil satu data meta dengan filter by ID
func GetAmeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	metaId := c.Param("metaId")
	var meta models.Meta
	defer cancel()

	fmt.Println(metaId)
	objId, _ := primitive.ObjectIDFromHex(metaId)

	err := metaCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&meta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  500,
			"Message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"Data": meta,
	})

}

// mengubah data meta dengan filter by ID
func EditAmeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	metaId := c.Param("metaId")
	var meta models.Meta
	defer cancel()

	fmt.Println(metaId)
	objId, _ := primitive.ObjectIDFromHex(metaId)

	//val the request body
	if err := c.Bind(&meta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  500,
			"Message": err.Error(),
		})
	}

	//menggunakan validasi untuk digunakan validasi required
	if validationErr := val.Struct(&meta); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  500,
			"Message": validationErr.Error(),
		})
	}

	// update func
	update := bson.M{"meta_title": meta.Meta_Title, "meta_descrption": meta.Meta_Descrption, "meta_url": meta.Meta_Url}

	result, err := metaCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  400,
			"Message": err.Error(),
		})
	}
	//mengambil data yang sudah diubah
	var updatedUser models.Meta
	if result.MatchedCount == 1 {
		err := metaCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  400,
				"Message": err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "success",
		"Data":    updatedUser,
	})
}

// menghapus data meta
func Deletemeta(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	metaId := c.Param("metaId")
	defer cancel()
	fmt.Println(metaId)
	objId, _ := primitive.ObjectIDFromHex(metaId)

	result, err := metaCollection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  500,
			"Message": err,
		})

	}

	if result.DeletedCount < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"Status":  404,
			"Message": result,
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  200,
		"Message": "Data Berhasil Di Hapus",
	})

}

func GetAllmeta(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var meta []models.Meta
	defer cancel()

	results, err := metaCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  500,
			"Message": err.Error(),
		})
		return
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singlemeta models.Meta
		if err = results.Decode(&singlemeta); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  500,
				"Message": err.Error(),
			})
			return
		}

		meta = append(meta, singlemeta)
	}

	c.JSON(http.StatusOK, gin.H{
		"Data":    meta,
		"Status":  200,
		"Message": "success",
	})

}
