package model

type Product struct {
	Type        int        `json:"type" bson:"type"`
	Obj         ProductObj `json:"obj" bson:"obj"`
	Destination string     `json:"destination" bson:"destination"`
}

type ProductObj struct {
	ProductModel ProductModel `json:"productModel" bson:"productModel"`
}

type ProductModel struct {
	ProductId          int    `json:"productId" bson:"productId"`
	ProductName        string `json:"productName" bson:"productName"`
	ProductUnit        int    `json:"productUnit" bson:"productUnit"`
	ProductDescription string `json:"productDescription" bson:"productDescription"`
	ProductDeleted     int    `json:"productDeleted" bson:"productDeleted"`
}
