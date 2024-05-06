package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/gin-gonic/gin"
)

type viewProductResponse struct {
	Product db.Product  `json:"product"`
	Reviews []db.Review `json:"reviews"`
}

type createProductRequest struct {
	ProductName  string   `json:"product_name"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	Quantity     int32    `json:"quantity"`
	Discount     float64  `json:"discount"`
	SizeOption   []string `json:"size_options"`
	ColorOptions []string `json:"color_options"`
	Category     string   `json:"category"`
	Brand        string   `json:"brand"`
	ImageUrl     []string `json:"image_url"`
}

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user cannot add a product"})
		return
	}

	product, err := server.store.CreateProduct(ctx, db.CreateProductParams{
		ProductName:  req.ProductName,
		Description:  req.Description,
		Price:        req.Price,
		Quantity:     int64(req.Quantity),
		Discount:     &req.Discount,
		SizeOptions:  req.SizeOption,
		ColorOptions: req.ColorOptions,
		Category:     req.Category,
		Brand:        &req.Brand,
		ImageUrl:     req.ImageUrl,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, product)
}

func (server *Server) listProducts(ctx *gin.Context) {
	products, err := server.store.GetAllProducts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, products)
}

type getProductRequest struct {
	ID int64 `uri:"id"`
}

func (server *Server) getProduct(ctx *gin.Context) {
	var uri getProductRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reviews, err := server.store.GetProductReviews(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := viewProductResponse{
		Product: product,
		Reviews: reviews,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getProductReviewsRequest struct {
	ID int64 `uri:"id"`
}

func (server *Server) getProductReviews(ctx *gin.Context) {
	var uri getProductReviewsRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	reviews, err := server.store.GetProductReviews(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}

type updateProductRequestUri struct {
	ID int64 `uri:"id" binding:"required"`
}

type updateProductRequestQuery struct {
	ProductName  string   `json:"product_name"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	Discount     float64  `json:"discount"`
	SizeOption   []string `json:"size_options"`
	ColorOptions []string `json:"color_options"`
	Category     string   `json:"category"`
	Brand        string   `json:"brand"`
	ImageUrl     []string `json:"image_url"`
}

func (server *Server) updateProduct(ctx *gin.Context) {
	var uri updateProductRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateProductRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.GetProductForUpdate(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.ProductName == "" {
		req.ProductName = product.ProductName
	}

	if req.Description == "" {
		req.Description = product.Description
	}

	if req.Price == 0.0 {
		req.Price = product.Price
	}

	if req.Discount == 0.0 {
		req.Discount = *product.Discount
	}

	if len(req.SizeOption) == 0 {
		req.SizeOption = product.SizeOptions
	}

	if len(req.ColorOptions) == 0 {
		req.ColorOptions = product.ColorOptions
	}

	if req.Category == "" {
		req.Category = product.Category
	}

	if req.Brand == "" {
		req.Brand = *product.Brand
	}

	if len(req.ImageUrl) == 0 {
		req.ImageUrl = product.ImageUrl
	}

	arg := db.UpdateProductParams{}
	arg.ID = product.ID
	arg.ProductName = req.ProductName
	arg.Description = req.Description
	arg.Price = req.Price
	arg.Discount = &req.Discount
	arg.SizeOptions = append(arg.SizeOptions, req.SizeOption...)
	arg.ColorOptions = append(arg.ColorOptions, req.ColorOptions...)
	arg.Category = req.Category
	arg.Brand = &req.Brand
	arg.ImageUrl = append(arg.ImageUrl, req.ImageUrl...)
	arg.UpdatedAt = time.Now()

	updatedProduct, err := server.store.UpdateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}

type addProductQuantityRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type addProductQuantityQuery struct {
	Quantity int32 `json:"quanity" binding:"required"`
}

func (server *Server) addProductQuantity(ctx *gin.Context) {
	var uri addProductQuantityRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req addProductQuantityQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.GetProductForUpdate(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updatedProduct, err := server.store.AddProductQuantity(ctx, db.AddProductQuantityParams{
		ID:       product.ID,
		Quantity: int64(req.Quantity),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}

// not implementing it will remove all the entities with this product ID
func (server *Server) deleteProduct(ctx *gin.Context) {

}

func (server *Server) getCollectionProducts(ctx *gin.Context) {

}
