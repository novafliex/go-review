package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success get all products!",
		Data:    products,
	})

}

func getProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	for _, product := range products {
		if product.ID == id {
			c.JSON(http.StatusOK, Response{
				Code:    200,
				Message: "success get product",
				Data:    product,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: "product not found",
		Data:    nil,
	})
}

func updateProduct(c *gin.Context) {
	idStr := c.Param("id")
	// 没有 ID 时，请求通常不会匹配该路由。
	// if idStr == "" {
	// 	c.JSON(http.StatusBadRequest, Response{
	// 		Code:    400,
	// 		Message: "id is required",
	// 		Data:    nil,
	// 	})
	// 	return
	// }
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}
	var newProduct Product
	err = c.ShouldBindJSON(&newProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid product",
			Data:    nil,
		})
		return
	}
	if newProduct.Name == "" || newProduct.Price <= 0 {
		c.JSON(http.StatusBadRequest, Response{
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
			c.JSON(http.StatusOK, Response{
				Code:    200,
				Message: "update success",
				Data:    products[index],
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: "not found",
		Data:    nil,
	})
}

func deleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	// 没有 ID 时，请求通常不会匹配该路由。
	// if idStr == "" {
	// 	c.JSON(http.StatusBadRequest, Response{
	// 		Code:    400,
	// 		Message: "id is required",
	// 		Data:    nil,
	// 	})
	// 	return
	// }
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}
	for index, product := range products {
		if product.ID == id {
			products = append(products[:index], products[index+1:]...)
			c.JSON(http.StatusOK, Response{
				Code:    200,
				Message: "delete success",
				Data:    product,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: "not found",
		Data:    nil,
	})
}

func createProduct(c *gin.Context) {
	var newProduct Product
	err := c.ShouldBindJSON(&newProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid product",
			Data:    nil,
		})
		return
	}
	if newProduct.Name == "" || newProduct.ID <= 0 || newProduct.Price <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "product id/name/price is required",
			Data:    nil,
		})
		return
	}
	for _, product := range products {
		if product.ID == newProduct.ID {
			c.JSON(http.StatusConflict, Response{
				Code:    409,
				Message: "product's id already exists",
				Data:    nil,
			})
			return
		}
	}
	products = append(products, newProduct)
	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "create success",
		Data:    newProduct,
	})
}
