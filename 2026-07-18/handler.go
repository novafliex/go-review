package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func product(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProduct(w, r)
	case http.MethodPut:
		updateProduct(w, r)
	case http.MethodPost:
		createProduct(w, r)
	case http.MethodDelete:
		deleteProduct(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, Response{
			Code:    405,
			Message: "请求方法不允许",
			Data:    nil,
		})
	}

}

func getProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusOK, Response{
			Code:    200,
			Message: "success get all products",
			Data:    products,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	for _, product := range products {
		if product.ID == id {
			writeJSON(w, http.StatusOK, Response{
				Code:    200,
				Message: "success get product",
				Data:    product,
			})
			return
		}
	}
	writeJSON(w, http.StatusNotFound, Response{
		Code:    404,
		Message: "product not found",
		Data:    nil,
	})
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	var newProduct Product
	err = json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid product",
			Data:    nil,
		})
		return
	}
	if newProduct.Name == "" || newProduct.Price <= 0 {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "product name/price is required",
			Data:    nil,
		})
		return
	}
	for index, product := range products {
		if product.ID == id {
			products[index].Name = newProduct.Name
			products[index].Price = newProduct.Price
			writeJSON(w, http.StatusOK, Response{
				Code:    200,
				Message: "update success",
				Data:    products[index],
			})
			return
		}
	}
	writeJSON(w, http.StatusNotFound, Response{
		Code:    404,
		Message: "not found",
		Data:    nil,
	})
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid product",
			Data:    nil,
		})
		return
	}

	if newProduct.ID <= 0 || newProduct.Name == "" || newProduct.Price <= 0 {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "product id/name/price is  required.",
			Data:    nil,
		})
		return
	}

	products = append(products, newProduct)
	writeJSON(w, http.StatusCreated, Response{
		Code:    200,
		Message: "create success",
		Data:    newProduct,
	})

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	for index, product := range products {
		if product.ID == id {
			products = append(products[:index], products[index+1:]...)
			writeJSON(w, http.StatusOK, Response{
				Code:    200,
				Message: "delete success",
				Data:    products,
			})
			return
		}
	}
	writeJSON(w, http.StatusNotFound, Response{
		Code:    404,
		Message: "not found.",
		Data:    nil,
	})
}
