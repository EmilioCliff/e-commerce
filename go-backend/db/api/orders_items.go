package api

import (
	"database/sql"
	"net/http"

	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/gin-gonic/gin"
)

type responseOrderItem struct {
	ProductName string   `json:"product_name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Quantity    int32    `json:"quantity"`
	Amount      float64  `json:"amount"`
	Discount    *float64 `json:"discount"`
	Size        *string  `json:"size_options"`
	Color       *string  `json:"color_options"`
	Category    string   `json:"category"`
	Brand       *string  `json:"brand"`
	ImageUrl    []string `json:"image_url"`
}

func (server *Server) newOrderItemResponse(orderItem db.OrderItem, ctx *gin.Context) (responseOrderItem, float64, error) {
	product, err := server.store.GetProduct(ctx, orderItem.ProductID)
	if err != nil {
		return responseOrderItem{}, 0, err
	}

	structuredOrderItem := responseOrderItem{
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    orderItem.Quantity,
		Discount:    product.Discount,
		Size:        orderItem.Size,
		Color:       orderItem.Color,
		Category:    product.Category,
		Brand:       product.Brand,
		ImageUrl:    product.ImageUrl,
	}

	if *structuredOrderItem.Discount != 0 {
		calculatedPrice := product.Price * (100 - *structuredOrderItem.Discount)
		amount := float64(structuredOrderItem.Quantity) * calculatedPrice

		structuredOrderItem.Amount = amount

		return structuredOrderItem, amount, nil
	}

	amount := float64(orderItem.Quantity) * product.Price

	structuredOrderItem.Amount = amount

	return structuredOrderItem, amount, nil
}

type getOrderItemRequestUri struct {
	Id int64 `uri:"id" binding:"required"`
}

func (server *Server) getOrderItem(ctx *gin.Context) {
	var uri getOrderItemRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusOK, errorResponse(err))
		return
	}

	orderItem, err := server.store.GetOrderItem(ctx, uri.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, _, err := server.newOrderItemResponse(orderItem, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getOrderItemsOfAnOrderRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

func (server *Server) getOrderItemsOfAnOrder(ctx *gin.Context) {
	var uri getOrderItemsOfAnOrderRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, uri.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin && payloadAssert.UserID == order.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"unathorized": "user unauthorized to view another user orderItems"})
		return
	}

	orderItems, err := server.store.GetOrderOrderItems(ctx, uri.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp map[string][]responseOrderItem
	rsp["order_items"] = make([]responseOrderItem, len(orderItems))
	for _, orderItem := range orderItems {
		structuredOrderItem, _, err := server.newOrderItemResponse(orderItem, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp["order_items"] = append(rsp["order_items"], structuredOrderItem)
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listOrderItems(ctx *gin.Context) {
	orderItems, err := server.store.ListOrderItems(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp map[int64][]responseOrderItem
	for _, orderItem := range orderItems {
		structuredOrderItem, _, err := server.newOrderItemResponse(orderItem, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if _, ok := rsp[orderItem.OrderID]; !ok {
			rsp[orderItem.OrderID] = make([]responseOrderItem, len(orderItems))
		}

		rsp[orderItem.OrderID] = append(rsp[orderItem.OrderID], structuredOrderItem)
	}

	ctx.JSON(http.StatusOK, rsp)
}
