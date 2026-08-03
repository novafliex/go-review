package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		rows, err := DB.Query("SELECT id,name,price FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "internal server error",
				Data:    nil,
			})
			return
		}
		defer rows.Close()
		var products []Product
		for rows.Next() {
			var product Product
			err := rows.Scan(&product.ID, &product.Name, &product.Price)
			if err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Code:    500,
					Message: "scan product failed",
					Data:    nil,
				})
				return
			}
			products = append(products, product)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "iterate products failed",
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			Code:    200,
			Message: "success",
			Data:    products,
		})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	var product Product

	err = DB.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Price)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "not found",
			Data:    nil,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "query product failed",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    product,
	})
}

func updateProduct(c *gin.Context) {
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
	var newProduct Product
	err = c.ShouldBindJSON(&newProduct)
	//判断请求体里的 JSON 能不能解析成 Product。
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
			Message: "name/price is required",
			Data:    nil,
		})
		return
	}
	result, err := DB.Exec("update products set name=?,price=? where id=?", newProduct.Name, newProduct.Price, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "update failed",
			Data:    nil,
		})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "query affected rows failed",
			Data:    nil,
		})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "not found",
			Data:    nil,
		})
		return
	}
	updatedProduct := Product{
		ID:    id,
		Name:  newProduct.Name,
		Price: newProduct.Price,
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "update success",
		Data:    updatedProduct,
	})
}

func deleteProduct(c *gin.Context) {
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

	var deletedProduct Product
	err = DB.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).Scan(&deletedProduct.ID, &deletedProduct.Name, &deletedProduct.Price)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "not found",
			Data:    nil,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "query product failed",
			Data:    nil,
		})
		return
	}

	result, err := DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "delete product failed",
			Data:    nil,
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "get affected rows failed",
			Data:    nil,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "delete success",
		Data:    deletedProduct,
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
	_, err = DB.Exec("insert into products (id, name, price) values (?,?,?)", newProduct.ID, newProduct.Name, newProduct.Price)
	if err != nil {
		c.JSON(http.StatusConflict, Response{
			Code:    409,
			Message: "insert failed",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "create success",
		Data:    newProduct,
	})
}
