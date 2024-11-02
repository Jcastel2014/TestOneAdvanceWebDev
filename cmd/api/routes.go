package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *appDependencies) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.notAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/createProduct", a.createProduct)
	router.HandlerFunc(http.MethodGet, "/displayProduct/:id", a.displayProduct)
	router.HandlerFunc(http.MethodDelete, "/deleteProduct/:id", a.deleteProduct)
	router.HandlerFunc(http.MethodGet, "/displayAllProducts", a.displayAllProducts)
	router.HandlerFunc(http.MethodPatch, "/updateProduct/:id", a.updateProduct)

	router.HandlerFunc(http.MethodPost, "/product/:id/createReview", a.createReview)
	router.HandlerFunc(http.MethodGet, "/product/:id/getReview/:rid", a.getReview)
	router.HandlerFunc(http.MethodPatch, "/product/:id/updateReview/:rid", a.updateReview)

	//update
	// router.HandlerFunc(http.MethodPatch, "/updateProduct/:id", a.updateProduct)

	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	// router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.updateCommentHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)

	// router.HandlerFunc(http.MethodGet, "/v1/comments", a.listCommentsHandler)
	return a.recoverPanic(router)
}
