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

var kategoriCollection *mongo.Collection = configs.GetCollection(configs.DB, "Main_Kategori")
var validasi_kategori = validator.New()

// membuat/memasukan kategori prodok baru
func KategoriCreate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var kategori models.MainCategory
	defer cancel()

	//val the request body
	if err := c.Bind(&kategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  400,
			"Message": err.Error(),
		})
		return
	}

	//use the validator library to val required fields
	if validationErr := validasi_kategori.Struct(&kategori); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  400,
			"Message": validationErr.Error(),
		})
		return
	}

	//membuat data banner baru
	newKategori := models.MainCategory{
		Id:              primitive.NewObjectID(),
		Kategori_Produk: kategori.Kategori_Produk,
		Image:           kategori.Image,
	}

	result, err := kategoriCollection.InsertOne(ctx, newKategori)
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

// mengambil satu data kategori dengan filter by ID
func GetAKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	kategoriid := c.Param("kategoriId")
	var kategori models.MainCategory
	defer cancel()

	//merubah data id dari string menjadi Integer
	fmt.Println(kategoriid)
	objId, _ := primitive.ObjectIDFromHex(kategoriid)

	err := kategoriCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&kategori)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  500,
			"Message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"Data": kategori,
	})
}

// mengubah data kategori dengan filter by ID
func EditKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	kategoriid := c.Param("kategoriID")
	var kategori models.MainCategory
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(kategoriid)

	//val the request body
	if err := c.Bind(&kategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  500,
			"Message": err.Error(),
		})
	}

	//menggunakan validasi untuk digunakan validasi required
	if validationErr := validasi_kategori.Struct(&kategori); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  500,
			"Message": validationErr.Error(),
		})
	}
	// update func
	update := bson.M{"Kategori Produk": kategori.Kategori_Produk, "Image": kategori.Image}

	result, err := kategoriCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  400,
			"Message": err.Error(),
		})
	}
	//mengambil data yang sudah diubah
	var updatedKategori models.MainCategory
	if result.MatchedCount == 1 {
		err := kategoriCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedKategori)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  400,
				"Message": err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "success",
		"Data":    updatedKategori,
	})
}

// menghapus data kategori
func DeleteKategori(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	kategoriid := c.Param("kategoriID")
	defer cancel()
	fmt.Println(kategoriid)
	objId, _ := primitive.ObjectIDFromHex(kategoriid)

	result, err := kategoriCollection.DeleteOne(ctx, bson.M{"id": objId})

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

// mengambil seluruh data kategori
func GetAllKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var kategori []models.MainCategory

	defer cancel()

	results, err := kategoriCollection.Find(ctx, bson.M{})

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
		var singleKategori models.MainCategory
		if err = results.Decode(&singleKategori); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  500,
				"Message": err.Error(),
			})
			return
		}

		kategori = append(kategori, singleKategori)
	}

	c.JSON(http.StatusOK, gin.H{
		"Data":    kategori,
		"Status":  200,
		"Message": "success",
	})

}

//get all mainkategori and child kategori

func Getloadkategori(ctx *gin.Context) {
	//lookupStage := bson.D{{"$lookup", bson.D{{"from", "child_Kategori"}, {"localField", "Main_Kategori"}, {"foreignField", "_id"}, {"as", "child_Kategori"}}}}
	//unwindStage := bson.D{{"$unwind", bson.D{{"path", "$child_Kategori"}, {"preserveNullAndEmptyArrays", false}}}}
	qry := []bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":         "child_Kategori",
				"localField":   "idmaincategory",
				"foreignField": "idmaincategory",
				"as":           "result",
			},
		},
	}

	showLoadedCursor, err := kategoriCollection.Aggregate(context.Background(), qry)
	if err != nil {
		panic(err)
	}
	var showsLoaded []bson.M
	if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
		panic(err)
	}

	showsLoaded = append(showsLoaded)
	ctx.JSON(http.StatusOK, gin.H{
		"Data":    showsLoaded,
		"Status":  200,
		"Message": "success",
	})

}