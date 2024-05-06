package api

import (
	"database/sql"
	"net/http"

	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/gin-gonic/gin"
)

type responseOrder struct {
	Username        string              `json:"username"`
	OrderItems      []responseOrderItem `json:"order_items"`
	Amount          float64             `json:"amount"`
	Status          string              `json:"status"`
	ShippingAddress string              `json:"shipping_address"`
}

func (server *Server) newOrderResponse(order db.Order, ctx *gin.Context) (responseOrder, error) {
	user, err := server.store.GetUser(ctx, order.UserID)
	if err != nil {
		return responseOrder{}, err
	}

	orderItems, err := server.store.GetOrderOrderItems(ctx, order.ID)
	if err != nil {
		return responseOrder{}, err
	}

	var structuredOrderItems []responseOrderItem
	var amount float64
	for _, orderItem := range orderItems {
		structuredOrderItem, returnedAmount, err := server.newOrderItemResponse(orderItem, ctx)
		if err != nil {
			return responseOrder{}, err
		}

		amount += returnedAmount
		structuredOrderItems = append(structuredOrderItems, structuredOrderItem)

	}

	return responseOrder{
		Username:        user.Username,
		OrderItems:      structuredOrderItems,
		Amount:          amount,
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
	}, nil

}

type createOrderRequestUri struct {
	ID int64 `uri:"id" binding:"required"`
}

type createOrderRequestQuery struct {
	ProductIDs      []int64  `json:"product_ids" binding:"required"`
	Quantity        []int32  `json:"quantities" binding:"required"`
	Colors          []string `json:"colors"`
	Size            []string `json:"size"`
	ShippingAddress string   `json:"shipping_address" binding:"required"`
}

// Create order and empty the cart list and send receipts
func (server *Server) createOrder(ctx *gin.Context) {
	var uri createOrderRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	var req createOrderRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var amount float64
	for idx, productID := range req.ProductIDs {
		product, err := server.store.GetProduct(ctx, productID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if *product.Discount != 0.0 {
			calculatedPrice := product.Price * (100 - *product.Discount)
			amount += calculatedPrice * float64(req.Quantity[idx])
			continue
		}

		amount += product.Price * float64(req.Quantity[idx])

	}
	// ACHIEVE THESE THROUGH A TRANSACTION if it fails both the order and items created are reverted
	// TODO: calculate the total shipping fee then distribute task to asynq to send_stk_push then processit and create a transactions
	//		 wait for the mpesa callback url then process it and update the transactions if successful if not do nothing
	//		 create the order and its order_items then empty the cart list.
	//       Send receipt order to user if successfull transactions.

	order, err := server.store.CreateOrder(ctx, db.CreateOrderParams{
		UserID:          uri.ID,
		Amount:          amount,
		ShippingAddress: req.ShippingAddress,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for idx, productID := range req.ProductIDs {
		_, err := server.store.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: productID,
			Quantity:  req.Quantity[idx],
			Color:     &req.Colors[idx],
			Size:      &req.Size[idx],
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	rsp, err := server.newOrderResponse(order, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	// ctx.JSON(http.StatusOK, gin.H{"order": "creating order"})
}

type getOrderRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getOrder(ctx *gin.Context) {
	var req getOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.ID)
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized to view another user order"})
		return
	}

	rsp, err := server.newOrderResponse(order, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listOrders(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user cannot list all orders"})
		return
	}

	orders, err := server.store.ListOrders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := make(map[string][]responseOrder)
	rsp["orders"] = make([]responseOrder, len(orders))
	for _, order := range orders {
		structureOrder, err := server.newOrderResponse(order, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp["orders"] = append(rsp["orders"], structureOrder)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listUsersOrdersRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) listUsersOrders(ctx *gin.Context) {
	var req listUsersOrdersRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin && payloadAssert.UserID != req.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user can not view another users orders"})
		return
	}

	orders, err := server.store.GetUserOrders(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := make(map[string][]responseOrder)
	rsp["orders"] = make([]responseOrder, len(orders))
	for _, order := range orders {
		structuredOrder, err := server.newOrderResponse(order, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp["orders"] = append(rsp["orders"], structuredOrder)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type adminChangeOrderStatusRequestUri struct {
	ID int64 `uri:"id" binding:"required"`
}

// Implement the oneof where the status can be pending, inprocess and delivered
type adminChangeOrderStatusRequestQuery struct {
	Status string `json:"status" binding:"required"`
}

func (server *Server) adminChangeOrderStatus(ctx *gin.Context) {
	var uri adminChangeOrderStatusRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req adminChangeOrderStatusRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user can not update an order status"})
		return
	}

	order, err := server.store.GetOrderForUpdate(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	editedOrder, err := server.store.UpdateOrder(ctx, db.UpdateOrderParams{
		Status: req.Status,
		ID:     order.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, err := server.newOrderResponse(editedOrder, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
