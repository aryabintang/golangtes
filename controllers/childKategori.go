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

var ChildCollection *mongo.Collection = configs.GetCollection(configs.DB, "child_Kategori")
var childValidasi = validator.New()

// membuat/memasukan kategori prodok baru
func ChildKategoriCreate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var child_Kategori models.ChildCategory
	//var kategori models.MainCategory
	defer cancel()

	//val the request body
	if err := c.Bind(&child_Kategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  400,
			"Message": err.Error(),
		})
		return
	}

	//use the validator library to val required fields
	if validationErr := childValidasi.Struct(&child_Kategori); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  400,
			"Message": validationErr.Error(),
		})
		return
	}

	//membuat data banner baru
	newChildkategori := models.ChildCategory{
		Id:             primitive.NewObjectID(),
		IdMainCategory: primitive.ObjectID{},
		Nama_produk:    child_Kategori.Nama_produk,
		Image:          child_Kategori.Image,
	}

	result, err := ChildCollection.InsertOne(ctx, newChildkategori)
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
func GetA_childKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id_produk := c.Param("id_produk")
	var child_Kategori models.ChildCategory
	defer cancel()

	//merubah data id dari string menjadi Integer
	fmt.Println(id_produk)
	objId, _ := primitive.ObjectIDFromHex(id_produk)

	err := ChildCollection.FindOne(ctx, bson.M{"id_produk": objId}).Decode(&child_Kategori)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  500,
			"Message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"Data": child_Kategori,
	})
}

// mengubah data kategori dengan filter by ID
func Edit_childKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id_produk := c.Param("id_produk")
	var child_Kategori models.ChildCategory
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id_produk)

	//val the request body
	if err := c.Bind(&child_Kategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  500,
			"Message": err.Error(),
		})
	}

	//menggunakan validasi untuk digunakan validasi required
	if validationErr := childValidasi.Struct(&child_Kategori); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  500,
			"Message": validationErr.Error(),
		})
	}
	// update func
	update := bson.M{"idmaincategory": child_Kategori.IdMainCategory, "nama_produk": child_Kategori.Nama_produk, "image": child_Kategori.Image}

	result, err := ChildCollection.UpdateOne(ctx, bson.M{"id_produk": objId}, bson.M{"$set": update})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  400,
			"Message": err.Error(),
		})
	}
	//mengambil data yang sudah diubah
	var updatedKategori models.ChildCategory
	if result.MatchedCount == 1 {
		err := ChildCollection.FindOne(ctx, bson.M{"id_produk": objId}).Decode(&updatedKategori)

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
func Delete_ChildKategori(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id_produk := c.Param("id_produk")
	defer cancel()
	fmt.Println(id_produk)
	objId, _ := primitive.ObjectIDFromHex(id_produk)

	result, err := ChildCollection.DeleteOne(ctx, bson.M{"id_produk": objId})

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
func GetAll_childKategori(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var kategori []models.ChildCategory

	defer cancel()

	results, err := ChildCollection.Find(ctx, bson.M{})

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
		var singleKategori models.ChildCategory
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