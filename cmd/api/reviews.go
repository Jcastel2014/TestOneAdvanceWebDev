package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jcastel2014/test1/internal/data"
	"github.com/jcastel2014/test1/internal/validator"
)

func (a *appDependencies) createReview(w http.ResponseWriter, r *http.Request) {
	id, _, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	var incomingData struct {
		Rating  float64 `json:"rating"`
		Comment string  `json:"comment"`
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	review := &data.Reviews{
		Rating:  incomingData.Rating,
		Comment: incomingData.Comment,
	}

	v := validator.New()
	data.ValidateReview(v, review, 1)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.productModel.InsertReview(review, id)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/review/%d", review.ID))

	data := envelope{
		"review": review,
	}

	err = a.writeJSON(w, http.StatusCreated, data, headers)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)

}

func (a *appDependencies) getReview(w http.ResponseWriter, r *http.Request) {

	id, rid, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	review, err := a.productModel.GetReview(id, rid)

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
		"review": review,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

}

func (a *appDependencies) updateReview(w http.ResponseWriter, r *http.Request) {

	id, rid, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	review, err := a.productModel.GetReview(id, rid)

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
		Rating  *float64 `json:"rating"`
		Comment *string  `json:"comment"`
	}

	err = a.readJSON(w, r, &incomingData)

	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.Rating != nil {
		review.Rating = *incomingData.Rating
	}

	if incomingData.Comment != nil {
		review.Comment = *incomingData.Comment
	}

	v := validator.New()

	data.ValidateReview(v, review, 1)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)

		return
	}

	err = a.productModel.UpdateReview(review, rid)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	data := envelope{
		"review": review,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

}

func (a *appDependencies) deleteReview(w http.ResponseWriter, r *http.Request) {

	id, rid, err := a.readIDParam(r)

	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.productModel.DeleteReview(id, rid)

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

func (a *appDependencies) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	var queryParametersData struct {
		Product string
		data.Filters
	}

	queryParameters := r.URL.Query()
	queryParametersData.Product = a.getSingleQueryParameters(queryParameters, "product", "")

	queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "id")
	queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "rating")
	queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "helpful_count")
	queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "created_at")
	queryParametersData.Filters.Sort = a.getSingleQueryParameters(queryParameters, "sort", "updated_at")

	queryParametersData.Filters.SortSafeList = []string{"id", "rating", "helpful_count", "created_at", "updated_at", "-id", "-rating", "-helpful_count", "-created_at", "-updated_at"}
	v := validator.New()

	queryParametersData.Filters.Page = a.getSingleIntegerParameters(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameters(queryParameters, "page_size", 10, v)

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	product_id, err := toInt(queryParametersData.Product)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	review, err := a.productModel.GetAllReviews(product_id, queryParametersData.Filters)

	if err != nil {
		a.serverErrResponse(w, r, err)
		return
	}

	data := envelope{
		"review": review,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		a.serverErrResponse(w, r, err)
	}
}

// func (a *appDependencies) SortReviews(w http.ResponseWriter, r *http.Request) {

// 	var queryParametersData struct {
// 		Name           string
// 		Description    string
// 		Category       string
// 		Average_rating string
// 		Price          string
// 		data.Filters
// 	}
// }
