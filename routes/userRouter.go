package routes

import (
	"golang_cms/controllers"
	controller "golang_cms/controllers"
	"golang_cms/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.GET("/user/:Email", controller.GetUserEmail())
	incomingRoutes.PUT("/user/:update", controller.UpdateUser())

	//banner
	incomingRoutes.POST("/banner", controllers.CreateBanner())             //memasukan data banner baru
	incomingRoutes.GET("/banner/:bannerId", controllers.GetABanner())      //mengambil satu data menggunakan filter ID
	incomingRoutes.PUT("/banner/:bannerId", controllers.EditABanner())     //mengedit satu data menggunaakn filter ID
	incomingRoutes.DELETE("/banner/:bannerId", controllers.DeleteBanner()) //menghapus satu data menggunaakn filter ID
	incomingRoutes.GET("/banners", controllers.GetAllBanner())             // mengambil semuah data Banner
	//meta
	incomingRoutes.POST("/meta", controllers.Createmeta)           //memasukan data meta baru
	incomingRoutes.GET("/meta/:metaId", controllers.GetAmeta)      //mengambil satu data meta dengan filter ID
	incomingRoutes.PUT("/meta/:metaId", controllers.EditAmeta)     //mengedit satu data meta dengan filter ID
	incomingRoutes.DELETE("/meta/:metaId", controllers.Deletemeta) //menghapus satu data dengan filter ID
	incomingRoutes.GET("/metas", controllers.GetAllmeta)           //mengambill semuah data meta
	//Kategori Produk main
	incomingRoutes.POST("kategori", controllers.KategoriCreate)               //memasukan data baru pada kategori_produk
	incomingRoutes.GET("kategori/:kategoriid", controllers.GetAKategori)      //memanggil satu data kategori denga filter ID
	incomingRoutes.PUT("kategori/:kategoriid", controllers.EditKategori)      //mengedit satu data kategori dengan filter ID
	incomingRoutes.DELETE("kategori/:kategoriid", controllers.DeleteKategori) //menghaspus satu data kategori dengan filter ID
	incomingRoutes.GET("kategori", controllers.GetAllKategori)                //mengambil semuah data kategori produk

	//Kategori Produk Child
	incomingRoutes.POST("kategori/child", controllers.ChildKategoriCreate)               //memasukan data baru pada kategori_produk
	incomingRoutes.GET("kategori/child/:id_produk", controllers.GetA_childKategori)      //memanggil satu data kategori denga filter ID
	incomingRoutes.PUT("kategori/child/:id_produk", controllers.Edit_childKategori)      //mengedit satu data kategori dengan filter ID
	incomingRoutes.DELETE("kategori/child/:id_produk", controllers.Delete_ChildKategori) //menghaspus satu data kategori dengan filter ID
	incomingRoutes.GET("kategoris/child/", controllers.GetAll_childKategori)             //mengambil semuah data kategori produk

	//Kategori Produk Child

	incomingRoutes.GET("kategoris/load/", controllers.Getloadkategori) //mengambil semuah data kategori produk

}
