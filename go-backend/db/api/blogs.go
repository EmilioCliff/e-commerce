package api

import (
	"database/sql"
	"net/http"

	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/gin-gonic/gin"
)

// TODO: How to pass content of a blog and what format to save it as
type addBlogRequestQuery struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type addBlogRequestUri struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) addBlog(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"unathorized": "user unauthorized to create blog"})
		return
	}

	var uri addBlogRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req addBlogRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	blog, err := server.store.CreateBlog(ctx, db.CreateBlogParams{
		Author:  uri.ID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

type editBlogUri struct {
	ID int64 `uri:"id" binding:"required"`
}

type editBlogQuery struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (server *Server) editBlog(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"unathorized": "user unauthorized to edit blog"})
		return
	}

	var uri editBlogUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req editBlogQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	blogToEdit, err := server.store.GetBlog(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payloadAssert.UserId != blogToEdit.Author {
		ctx.JSON(http.StatusUnauthorized, gin.H{"unathorized": "admin unauthorized to edit another admins blog"})
		return
	}

	if req.Content == "" {
		req.Content = blogToEdit.Content
	}

	if req.Title == "" {
		req.Title = blogToEdit.Title
	}

	blog, err := server.store.EditBlog(ctx, db.EditBlogParams{
		ID:      uri.ID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

func (server *Server) listBlogs(ctx *gin.Context) {
	blogs, err := server.store.ListBlogs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

type deleteBlogRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) deleteBlog(ctx *gin.Context) {
	payload := ctx.MustGet(PayloadKey)
	payloadAssert := payload.(token.Payload)
	if !payloadAssert.IsAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"unathorized": "user unauthorized to delete blog"})
		return
	}

	var req deleteBlogRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	blogToDelete, err := server.store.GetBlog(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payloadAssert.UserId != blogToDelete.Author {
		ctx.JSON(http.StatusUnauthorized, gin.H{"unathorized": "admin unauthorized to delete another admins blog"})
		return
	}

	if err := server.store.DeleteBlog(ctx, req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"delete": "blog deleted successful"})
}

type getAdminsBlogRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getAdminsBlogs(ctx *gin.Context) {
	var req getAdminsBlogRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	blogs, err := server.store.GetAdminsBlog(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, blogs)
}
