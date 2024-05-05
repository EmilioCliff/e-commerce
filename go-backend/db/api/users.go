package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/EmilioCliff/e-commerce/db/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userResponse struct {
	Username         string       `json:"username"`
	Email            string       `json:"email"`
	Subscription     bool         `json:"subscription"`
	UserCartProducts []db.Product `json:"user_cart"`
	Role             string       `json:"role"`
}

func (server *Server) newUserResponse(user db.User, ctx *gin.Context) (userResponse, error) {
	var products []db.Product
	for _, productId := range user.UserCart {
		product, err := server.store.GetProduct(ctx, productId)
		if err != nil {
			return userResponse{}, err
		}

		products = append(products, product)
	}

	return userResponse{
		Username:         user.Username,
		Email:            user.Email,
		Subscription:     user.Subscription,
		UserCartProducts: products,
		Role:             user.Role,
	}, nil
}

type createUserRequest struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Subscription bool   `json:"subscription"`
	Role         string `json:"role"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		Subscription: req.Subscription,
		Role:         req.Role,
		Password:     hashPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, err := server.newUserResponse(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAT time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.VerifyPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.maker.CreateToken(user.ID, server.config.AccessTokenDuration, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.maker.CreateToken(user.ID, server.config.RefreshTokenDuration, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       refreshPayload.UserID,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		UserAgent:    ctx.Request.UserAgent(),
		UserIp:       ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpireAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userRsp, err := server.newUserResponse(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpireAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAT: refreshPayload.ExpireAt,
		User:                  userRsp,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type refreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) refreshAccessToken(ctx *gin.Context) {
	var req refreshAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, _, err := server.maker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("refresh token is blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatch in refresh token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("refresh token is expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var role string
	if payload.IsAdmin {
		role = "admin"
	} else {
		role = "user"
	}

	accessToken, accessTokenPayload, err := server.maker.CreateToken(session.UserID, server.config.AccessTokenDuration, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := refreshAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpireAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type blockRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) blockRefreshToken(ctx *gin.Context) {
	var req blockRefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, _, err := server.maker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, err = server.store.BlockSession(ctx, db.BlockSessionParams{
		IsBlocked: true,
		ID:        payload.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "refresh token is blocked"})
}

func (server *Server) listUsers(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized to view all users"})
		return
	}

	users, err := server.store.ListUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []userResponse
	for _, user := range users {
		structuredUser, err := server.newUserResponse(user, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp = append(rsp, structuredUser)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var uri getUserRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, uri.Id)
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
	if !payloadAssert.IsAdmin && payloadAssert.UserID != user.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized to view another users"})
		return
	}

	rsp, err := server.newUserResponse(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

type updateUserRequestUri struct {
	Id int64 `uri:"id"`
}

type updateUserRequestQuery struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var uri updateUserRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateUserRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserForUpdate(ctx, uri.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.Username == "" {
		req.Username = user.Username
	}

	if req.Password == "" {
		req.Password = user.Password
	}

	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if payloadAssert.IsAdmin {
		if req.Role == "" {
			req.Role = user.Role
		}

	} else {
		req.Role = user.Role
	}

	// TODO:	update the updated time
	updatedUser, err := server.store.UpdateUserCredentials(ctx, db.UpdateUserCredentialsParams{
		Username:  req.Username,
		Password:  req.Password,
		Role:      req.Role,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, err := server.newUserResponse(updatedUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)

}

// Not implemente as it will cascadily delete all data with the userId
func (server *Server) deleteUser(ctx *gin.Context) {

}

type changeSubscriptionRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type changeSubscriptionRequestQuery struct {
	Subscription bool `json:"subscription" binding:"required"`
}

func (server *Server) changeSubscription(ctx *gin.Context) {
	var uri changeSubscriptionRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req changeSubscriptionRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserForUpdate(ctx, uri.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// arg := db.UpdateUserSubscriptionParams{
	// 	ID:        user.ID,
	// 	UpdatedAt: time.Now(),
	// }

	// if req.Subscription == true {
	// 	arg.Subscription = true
	// } else {
	// 	arg.Subscription = false
	// }

	updatedUser, err := server.store.UpdateUserSubscription(ctx, db.UpdateUserSubscriptionParams{
		ID:           user.ID,
		Subscription: req.Subscription,
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, err := server.newUserResponse(updatedUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getSubscribeUsers(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user can not view subscribes users"})
		return
	}

	users, err := server.store.GetSubscribedUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []userResponse
	for _, user := range users {
		structuredUser, err := server.newUserResponse(user, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp = append(rsp, structuredUser)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listUserCartProductRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

func (server *Server) listUserCartProduct(ctx *gin.Context) {
	var uri listUserCartProductRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, uri.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []db.Product
	for _, productId := range user.UserCart {
		product, err := server.store.GetProduct(ctx, productId)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		rsp = append(rsp, product)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type addToUserCartRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type addToUserCartRequestQuery struct {
	ProductId int64 `json:"product_id" binding:"product_id"`
}

func (server *Server) addToUserCart(ctx *gin.Context) {
	var uri addToUserCartRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var req addToUserCartRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUserForUpdate(ctx, uri.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// proves that the product id exist in the db
	_, err = server.store.GetProduct(ctx, req.ProductId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user.UserCart = append(user.UserCart, req.ProductId)

	updatedUser, err := server.store.UpdateUserCart(ctx, db.UpdateUserCartParams{
		ID:       user.ID,
		UserCart: user.UserCart,
	})

	rsp, err := server.newUserResponse(updatedUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

type removeToUserCartRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type removeToUserCartRequestQuery struct {
	ProductId int64 `json:"product_id" binding:"product_id"`
}

func (server *Server) removeToUserCart(ctx *gin.Context) {
	var uri removeToUserCartRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var req removeToUserCartRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUserForUpdate(ctx, uri.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// proves that the product id exist in the db
	_, err = server.store.GetProduct(ctx, req.ProductId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for idx, productID := range user.UserCart {
		if productID == req.ProductId {
			user.UserCart = append(user.UserCart[:idx], user.UserCart[idx+1:]...)
			break
		}
	}

	updatedUser, err := server.store.UpdateUserCart(ctx, db.UpdateUserCartParams{
		ID:       user.ID,
		UserCart: user.UserCart,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, err := server.newUserResponse(updatedUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
