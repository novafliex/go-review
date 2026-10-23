//所有的接口处理函数

package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 查找商品
func getProduct(c *gin.Context) {
	idStr := c.Param("id")
	// 如果没有传入 id，就返回全部商品
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
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "id must be greater than 0",
			Data:    nil,
		})
		return
	}
	var product Product

	err = DB.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).
		Scan(&product.ID, &product.Name, &product.Price)

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

// 添加商品
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
	//也是需要一个简单的校验 看下传进来的数据是不是合理的
	if newProduct.ID <= 0 || newProduct.Name == "" || newProduct.Price <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "product id/name/price is required",
			Data:    nil,
		})
		return
	}
	_, err = DB.Exec("insert into products (id,name,price) values (?,?,?)", newProduct.ID, newProduct.Name, newProduct.Price)
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

// 修改商品
func updateProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "id is required",
			Data:    nil,
		})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "id must be greater than 0",
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

// 删除商品
func deleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "id is required",
			Data:    nil,
		})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "id must be greater than 0",
			Data:    nil,
		})
		return
	}
	var deletedProduct Product

	err = DB.QueryRow(
		"SELECT id, name, price FROM products WHERE id = ?",
		id,
	).Scan(&deletedProduct.ID, &deletedProduct.Name, &deletedProduct.Price)

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

func getUser(c *gin.Context) {
	idStr := c.Param("id")
	//查询全部的用户
	if idStr == "" {
		rows, err := DB.Query("SELECT id, username, email, phone, role, created_at FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "internal server error",
				Data:    nil,
			})
			return
		}
		defer rows.Close()

		var users []User

		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Role, &user.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Code:    500,
					Message: "scan user failed",
					Data:    nil,
				})
				return
			}

			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "iterate users failed",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Code:    200,
			Message: "query users success!",
			Data:    users,
		})
		return
	}

	//查询单个id的
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return
	}

	var user User
	err = DB.QueryRow("select id, username, email, phone, role, created_at from users where id = ?", id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt)
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
			Message: "query user failed",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "query user success",
		Data:    user,
	})
}

func createUser(c *gin.Context) {
	var newUser User
	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid user",
			Data:    nil,
		})
		return
	}
	//参数校验
	if newUser.Username == "" || newUser.Password == "" || newUser.Email == "" || newUser.Phone == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "username/password/email/phone is required",
			Data:    nil,
		})
		return
	}
	if len(newUser.Password) < 6 || len(newUser.Password) > 72 {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "password must be 6 to 72 bytes"})
		return
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: "hash password failed"})
		return
	}
	newUser.Password = string(hashedPassword)

	//用户的role应该默认设置为普通
	newUser.Role = "user"

	result, err := DB.Exec("insert into users (username, password, email, phone, role) values (?,?,?,?,?)",
		newUser.Username,
		newUser.Password,
		newUser.Email,
		newUser.Phone,
		newUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "create user failed",
			Data:    nil,
		})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "get user id failed",
			Data:    nil,
		})
		return
	}

	newUser.ID = int(id)
	newUser.Password = ""

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "create user success",
		Data:    newUser,
	})

}

func updateUser(c *gin.Context) {
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

	var updatedUser User
	if err = c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid user",
			Data:    nil,
		})
		return
	}
	//校验输入进来的信息是否正确
	if updatedUser.Username == "" || updatedUser.Password == "" || updatedUser.Email == "" || updatedUser.Phone == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "username/password/email/phone is required",
			Data:    nil,
		})
		return
	}

	if len(updatedUser.Password) < 6 || len(updatedUser.Password) > 72 {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "password must be 6 to 72 bytes"})
		return
	}
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: "hash password failed"})
		return
	}
	updatedUser.Password = string(hashedPassword)

	result, err := DB.Exec("update users set username = ?, password = ?, email = ?, phone = ? where id = ?",
		updatedUser.Username,
		updatedUser.Password,
		updatedUser.Email,
		updatedUser.Phone,
		id)
	//Go 在多行参数列表中，最后一个参数后面必须加逗号
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "update user failed",
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

	updatedUser.Password = ""
	updatedUser.Role = "" // 更新接口不修改角色
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "update user success",
		Data:    updatedUser,
	})
}

func deleteUser(c *gin.Context) {
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

	result, err := DB.Exec("delete from users where id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "delete user failed",
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
		Message: "delete user success",
		Data:    nil,
	})
}
