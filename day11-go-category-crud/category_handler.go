package main

import (
	"database/sql"
	"strconv"

	"strings"

	"github.com/gin-gonic/gin"
)

func getCategory(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		rows, err := DB.Query("select id, name, description, status, created_at, updated_at from categories")
		if err != nil {
			ErrorResponse(c, 500, "internal server error")
			return
		}
		defer rows.Close()

		categories := make([]Category, 0)

		for rows.Next() {
			var category Category
			err := rows.Scan(
				&category.ID,
				&category.Name,
				&category.Description,
				&category.Status,
				&category.CreatedAt,
				&category.UpdatedAt,
			)
			if err != nil {
				ErrorResponse(c, 500, "scan category failed")
				return
			}

			categories = append(categories, category)
		}
		// 检查遍历过程是否发生错误
		if err := rows.Err(); err != nil {
			ErrorResponse(c, 500, "iterate category failed")
			return
		}

		Success(c, categories)
		return
	}
	// 查询单个
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ErrorResponse(c, 400, "invalid category id")
		return
	}

	var category Category
	err = DB.QueryRow("select id, name, description, status, created_at, updated_at from categories where id = ?", id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		ErrorResponse(c, 404, "not found")
		return
	}
	if err != nil {
		ErrorResponse(c, 500, "query category failed")
		return
	}

	Success(c, category)
}

func updateCategory(c *gin.Context) {
	// 解析并且校验id
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ErrorResponse(c, 400, "invalid id")
		return
	}
	//接受并校验请求参数
	var req CategoryRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Status = strings.TrimSpace(req.Status)

	if req.Name == "" {
		ErrorResponse(c, 400, "name is required")
		return
	}
	if req.Status != "enabled" && req.Status != "disabled" {
		ErrorResponse(c, 400, "status must be enabled or disabled")
		return
	}

	// 判断分类是否存在
	var existingID int

	err = DB.QueryRow(
		"SELECT id FROM categories WHERE id = ?",
		id,
	).Scan(&existingID)

	if err == sql.ErrNoRows {
		ErrorResponse(c, 404, "category not found")
		return
	}

	if err != nil {
		ErrorResponse(c, 500, "query category failed")
		return
	}

	// 判断名称是否被其他分类使用
	var duplicateID int

	err = DB.QueryRow(
		"SELECT id FROM categories WHERE name = ? AND id != ?",
		req.Name,
		id,
	).Scan(&duplicateID)

	if err == nil {
		ErrorResponse(c, 409, "category name already exists")
		return
	}

	if err != sql.ErrNoRows {
		ErrorResponse(c, 500, "check category name failed")
		return
	}

	result, err := DB.Exec("update categories set name = ?, description = ?, status = ?  where id = ?", req.Name, req.Description, req.Status, id)
	if err != nil {
		ErrorResponse(c, 500, "update category failed")
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		ErrorResponse(c, 500, "get affected rows failed")
		return
	}
	// 重新查询更新后的完整数据
	var category Category

	err = DB.QueryRow(
		`SELECT id, name, description, status, created_at, updated_at
			 FROM categories
			 WHERE id = ?`,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		ErrorResponse(c, 500, "query updated category failed")
		return
	}

	if affected == 0 {
		JSONResponse(c, 200, "category has no changes", category)
		return
	}

	JSONResponse(c, 200, "category updated successfully", category)

}

func createCategory(c *gin.Context) {
	//接收客户端提交的数据
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "invalid request body")
		return
	}
	// 清除字符串前后的空格
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Status = strings.TrimSpace(req.Status)
	//检验分类名称
	if req.Name == "" {
		ErrorResponse(c, 400, "name is required")
		return
	}
	// 状态没传使用默认值
	if req.Status == "" {
		req.Status = "enabled"
	}

	//校验状态范围
	if req.Status != "enabled" && req.Status != "disabled" {
		ErrorResponse(c, 400, "status must be enabled or disabled")
		return
	}

	//检查分类名称是否重复
	var existingID int

	err := DB.QueryRow("select id from categories where name = ?", req.Name).Scan(&existingID)

	//err 为nil 表示查到了同名分类
	if err == nil {
		ErrorResponse(c, 409, "category name already exists")
		return
	}
	//sql.ErrNoRows 没查到同名分类
	if err != sql.ErrNoRows {
		ErrorResponse(c, 500, "check category name failed")
		return
	}

	//执行新增
	result, err := DB.Exec("insert into categories (name, description, status) values (?, ?, ?)", req.Name, req.Description, req.Status)
	if err != nil {
		ErrorResponse(c, 500, "create category failed")
		return
	}

	//获取数据库自动生成的id
	newID, err := result.LastInsertId()
	if err != nil {
		ErrorResponse(c, 500, "get category id failed")
		return
	}
	//查询创建好的完整数据
	var category Category
	err = DB.QueryRow("select id, name, description, status, created_at, updated_at from categories where id = ?", newID).Scan(&category.ID, &category.Name, &category.Description, &category.Status, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		ErrorResponse(c, 500, "query created category failed")
		return
	}
	JSONResponse(c, 201, "category created successfully", category)
}

func deleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ErrorResponse(c, 400, "invalid id")
		return
	}

	result, err := DB.Exec(
		"DELETE FROM categories WHERE id = ?",
		id,
	)
	if err != nil {
		ErrorResponse(c, 500, "delete category failed")
		return
	}

	// 3. 获取受影响的行数
	affected, err := result.RowsAffected()
	if err != nil {
		ErrorResponse(c, 500, "get affected rows failed")
		return
	}

	// 4. 没有删除任何数据，说明 ID 不存在
	if affected == 0 {
		ErrorResponse(c, 404, "category not found")
		return
	}

	// 5. 返回删除成功
	JSONResponse(
		c,
		200,
		"category deleted successfully",
		nil,
	)
}

func updateCategoryStatus(c *gin.Context) {
	// 1. 获取并校验 ID
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ErrorResponse(c, 400, "invalid category id")
		return
	}

	// 2. 接收状态参数
	var req CategoryStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "invalid request body")
		return
	}

	// 3. 清理并校验状态
	req.Status = strings.TrimSpace(req.Status)

	if req.Status != "enabled" && req.Status != "disabled" {
		ErrorResponse(c, 400, "status must be enabled or disabled")
		return
	}

	// 4. 执行状态更新
	result, err := DB.Exec(
		"UPDATE categories SET status = ? WHERE id = ?",
		req.Status,
		id,
	)
	if err != nil {
		ErrorResponse(c, 500, "update category status failed")
		return
	}

	// 5. 检查分类是否存在
	affected, err := result.RowsAffected()
	if err != nil {
		ErrorResponse(c, 500, "get affected rows failed")
		return
	}

	if affected == 0 {
		var existingID int
		err = DB.QueryRow(
			"SELECT id FROM categories WHERE id = ?",
			id,
		).Scan(&existingID)

		if err == sql.ErrNoRows {
			ErrorResponse(c, 404, "category not found")
			return
		}

		if err != nil {
			ErrorResponse(c, 500, "query category failed")
			return
		}
	}

	// 6. 查询更新后的完整分类数据
	var category Category

	err = DB.QueryRow(
		`SELECT id, name, description, status, created_at, updated_at
		 FROM categories
		 WHERE id = ?`,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		ErrorResponse(c, 500, "query updated category failed")
		return
	}

	// 7. 返回最新分类数据
	JSONResponse(
		c,
		200,
		"category status updated successfully",
		category,
	)
}
