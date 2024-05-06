package api

import (
	"database/sql"
	"net/http"

	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/gin-gonic/gin"
)

type addReviewRequestUri struct {
	ID int64 `uri:"id" binding:"required"`
}

type addReviewRequestQuery struct {
	UserID int64  `json:"user_id" binding:"required"`
	Rating int32  `json:"rating" binding:"required"`
	Review string `json:"review" binding:"required"`
}

func (server *Server) addReview(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "admin unauthorized to create review"})
		return
	}

	var uri addReviewRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req addReviewRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.CreateReveiw(ctx, db.CreateReveiwParams{
		UserID:    req.UserID,
		ProductID: uri.ID,
		Rating:    req.Rating,
		Review:    req.Review,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	averageRating, err := server.store.StoreCalculateProductRating(ctx, uri.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.UpdateProductRating(ctx, db.UpdateProductRatingParams{
		ID:     uri.ID,
		Rating: &averageRating,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}

type editReviewRequestUri struct {
	ID int64 `uri:"id" binding:"required"`
}

type editReviewRequestQuery struct {
	Rating int32  `json:"rating"`
	Review string `json:"review"`
}

func (server *Server) editReview(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "admin unauthorized to edit review"})
		return
	}

	var uri editReviewRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req editReviewRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.GetReview(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payloadAssert.UserID != review.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized to edit another users review"})
		return
	}

	if req.Rating == 0 {
		req.Rating = review.Rating
	}

	if req.Review == "" {
		req.Review = review.Review
	}

	editedRevies, err := server.store.EditReview(ctx, db.EditReviewParams{
		ID:     uri.ID,
		Rating: req.Rating,
		Review: req.Review,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, editedRevies)
}

type deleteReviewRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) deleteReview(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "admin unauthorized to delete review"})
		return
	}

	var uri deleteReviewRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.GetReview(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payloadAssert.UserID != review.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized to delete another user review"})
		return
	}

	if err := server.store.DeleteReview(ctx, uri.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "review deleted successful"})
}
