package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jcastel2014/test1/internal/data"
	"github.com/jcastel2014/test1/internal/validator"
)

func (a *appDependencies) createProduct(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Image_url   string  `json:"image_url"`
		Price       float64 `json:"price"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	product := &data.Product{
		Name:        incomingData.Name,
		Description: incomingData.Description,
		Category:    incomingData.Category,
		Image_url:   incomingData.Image_url,
		Price:       incomingData.Price,
	}

	v := validator.New()

	// one sent to identify handler
	data.ValidateProduct(v, product, 1)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.productModel.Insert(product)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/product/%d", product.ID))

	data := envelope{
		"product": product,
	}

	err = a.writeJSON(w, http.StatusCreated, data, headers)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
}

func (a *appDependencies) displayProduct(w http.ResponseWriter, r *http.Request) {
	id, _, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	product, err := a.productModel.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrResponse(w, r, err)
		}

		return
	}

	data := envelope{
		"product": product,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}
}

func (a *appDependencies) updateProduct(w http.ResponseWriter, r *http.Request) {
	id, _, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	product, err := a.productModel.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrResponse(w, r, err)
		}

		return
	}

	var incomingData struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Category    *string  `json:"category"`
		Image_url   *string  `json:"image_url"`
		Price       *float64 `json:"price"`
	}

	err = a.readJSON(w, r, &incomingData)

	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.Name != nil {
		product.Name = *incomingData.Name
	}

	if incomingData.Description != nil {
		product.Description = *incomingData.Description
	}

	if incomingData.Category != nil {
		product.Category = *incomingData.Category
	}

	if incomingData.Image_url != nil {
		product.Image_url = *incomingData.Image_url
	}

	if incomingData.Price != nil {
		product.Price = *incomingData.Price
	}

	v := validator.New()

	data.ValidateProduct(v, product, 1)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)

		return
	}

	err = a.productModel.Update(product)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	data := envelope{
		"product": product,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}
}

func (a *appDependencies) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.productModel.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrResponse(w, r, err)
		}

		return
	}

	data := envelope{
		"message": "comment successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
	}
}

func (a *appDependencies) displayAllProducts(w http.ResponseWriter, r *http.Request) {

	// due to image url being too big, to make the queries more readable i havnt included it into the query
	// the images will show as a number, which indicates the ID of
	var queryParametersData struct {
		Name           string
		Description    string
		Category       string
		Average_rating string
		Price          string
		data.Filters
	}

	queryParameters := r.URL.Query()
	queryParametersData.Name = a.getSingleQueryParameters(queryParameters, "name", "")
	queryParametersData.Description = a.getSingleQueryParameters(queryParameters, "description", "")
	queryParametersData.Category = a.getSingleQueryParameters(queryParameters, "category", "")
	queryParametersData.Average_rating = a.getSingleQueryParameters(queryParameters, "average_rating", "")
	queryParametersData.Price = a.getSingleQueryParameters(queryParameters, "price", "")

	v := validator.New()

	queryParametersData.Filters.Page = a.getSingleIntegerParameters(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameters(queryParameters, "page_size", 10, v)

	queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "id")

	queryParametersData.Filters.SortSafeList = []string{"id", "name", "-id", "-name"}

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	product, metadata, err := a.productModel.GetAll(queryParametersData.Name, queryParametersData.Description, queryParametersData.Category, queryParametersData.Average_rating, queryParametersData.Filters)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	data := envelope{
		"product":  product,
		"metadata": metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		a.serverErrResponse(w, r, err)
	}
}
