package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Banner struct {
	Id     primitive.ObjectID `json:"id,omitempty"`
	Banner string             `json:"banner,omitempty" validate:"required"`
	Alt    string             `json:"alt,omitempty" validate:"required"`
	Link   string             `json:"link,omitempty" validate:"required"`
}

// meta
type Meta struct {
	Id              primitive.ObjectID `json:"id,omitempty"`
	Meta_Title      string             `json:"meta_title,omitempty" validate:"required"`
	Meta_Url        string             `json:"meta_url,omitempty" validate:"required"`
	Meta_Descrption string             `json:"meta_descrption,omitempty" validate:"required"`
}

type MainCategory struct {
	Id              primitive.ObjectID `bson:"idmaincategory"`
	Kategori_Produk string             `json:"kategori_produk,omitempty" validate:"required"`
	Image           string             `json:"image" validate:"required"`
}

type ChildCategory struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	IdMainCategory primitive.ObjectID `json:"idmaincategory"`
	Nama_produk    string             `json:"nama_produk,omitempty" validate:"required"`
	Image          string             `json:"image" validate:"required"`
}