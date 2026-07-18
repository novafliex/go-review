package main

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func createStudent(c *gin.Context) {
	var req StudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "invalid request body")
		return
	}

	req.StudentNo = strings.TrimSpace(req.StudentNo)
	req.Name = strings.TrimSpace(req.Name)
	req.Status = strings.TrimSpace(req.Status)

	// 检验字段
	if req.StudentNo == "" {
		ErrorResponse(c, 400, "student_no is required")
		return
	}
	if req.Name == "" {
		ErrorResponse(c, 400, "name is required")
		return
	}
	if len([]rune(req.Name)) > 50 {
		ErrorResponse(c, 400, "name cannot exceed 50 characters")
		return
	}
	if req.Age <= 0 || req.Age > 150 {
		ErrorResponse(c, 400, "age must be between 1 and 150")
		return
	}

	if req.Status != "enabled" && req.Status != "disabled" {
		ErrorResponse(c, 400, "status must be enabled or disabled")
		return
	}

	// 查询数据库 看是否会重复
	var existingID int
	err := DB.QueryRow("select id from students where student_no = ?", req.StudentNo).Scan(&existingID)
	if err == nil {
		ErrorResponse(c, 409, "student_no already exists")
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 500, "check student_no failed")
		return
	}

	// 插入
	result, err := DB.Exec("insert into students (student_no, name, age, status) values (?, ?, ?, ?)", req.StudentNo, req.Name, req.Age, req.Status)
	if err != nil {
		ErrorResponse(c, 500, "insert failed")
		return
	}

	// 得到新加的id
	id, err := result.LastInsertId()
	if err != nil {
		ErrorResponse(c, 500, "get last insert id failed")
		return
	}

	student := Student{
		ID:        int(id),
		StudentNo: req.StudentNo,
		Name:      req.Name,
		Age:       req.Age,
		Status:    req.Status,
	}

	JSONResponse(c, 201, "student created successfully", student)
}

func getStudents(c *gin.Context) {
	rows, err := DB.Query("select id, student_no, name, age, status from students")
	if err != nil {
		ErrorResponse(c, 500, "query students failed")
		return
	}

	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		var student Student
		err := rows.Scan(
			&student.ID,
			&student.StudentNo,
			&student.Name,
			&student.Age,
			&student.Status,
		)
		if err != nil {
			ErrorResponse(c, 500, "scan student failed")
			return
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		ErrorResponse(c, 500, "iterate students failed")
		return
	}

	Success(c, students)
}

func getStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ErrorResponse(c, 400, "invalid id")
		return
	}

	var student Student

	err = DB.QueryRow("select id, student_no, name, age, status from students where id = ?", id).Scan(
		&student.ID,
		&student.StudentNo,
		&student.Name,
		&student.Age,
		&student.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 404, "student not found")
		return
	}
	if err != nil {
		ErrorResponse(c, 500, "query student failed")
		return
	}

	Success(c, student)
}

func updateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ErrorResponse(c, 400, "invalid id")
		return
	}

	var req StudentRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "request body failed")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.StudentNo = strings.TrimSpace(req.StudentNo)
	req.Status = strings.TrimSpace(req.Status)

	if req.StudentNo == "" {
		ErrorResponse(c, 400, "student_no is required")
		return
	}

	if req.Name == "" {
		ErrorResponse(c, 400, "name is required")
		return
	}

	if len([]rune(req.Name)) > 50 {
		ErrorResponse(c, 400, "name cannot exceed 50 characters")
		return
	}

	if req.Age <= 0 || req.Age > 150 {
		ErrorResponse(c, 400, "age must be between 1 and 150")
		return
	}

	if req.Status != "enabled" && req.Status != "disabled" {
		ErrorResponse(c, 400, "status must be enabled or disabled")
		return
	}
	// 看要修改的id是否在数据库存在
	var exctingId int
	err = DB.QueryRow("select id from students where id = ?", id).Scan(&exctingId)
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 404, "student not found")
		return
	}
	if err != nil {
		ErrorResponse(c, 500, "query student failed")
		return
	}
	//判断新学号是否被其他学生占用
	var duplicateID int

	err = DB.QueryRow(
		"SELECT id FROM students WHERE student_no = ? AND id != ?",
		req.StudentNo,
		id,
	).Scan(&duplicateID)

	if err == nil {
		ErrorResponse(c, 409, "student_no already exists")
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 500, "check student_no failed")
		return
	}

	result, err := DB.Exec("update students set student_no = ?, name = ?, age = ?, status = ? where id = ?",
		req.StudentNo, req.Name, req.Age, req.Status, id)
	if err != nil {
		ErrorResponse(c, 500, "update failed")
		return
	}

	_, err = result.RowsAffected()
	if err != nil {
		ErrorResponse(c, 500, "get affected rows failed")
		return
	}
	// // affected == 0 说明可能没有变化，仍然继续查询并返回完整数据
	// _ = affected

	var student Student

	err = DB.QueryRow("select id, student_no, name, age, status from students where id = ?", id).Scan(
		&student.ID,
		&student.StudentNo,
		&student.Name,
		&student.Age,
		&student.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 404, "student not found")
		return
	}
	if err != nil {
		ErrorResponse(c, 500, "query student failed")
		return
	}

	Success(c, student)

}

func deleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ErrorResponse(c, 400, "invalid id")
		return
	}

	var student Student

	err = DB.QueryRow("select id, student_no, name, age, status from students where id = ?", id).Scan(
		&student.ID,
		&student.StudentNo,
		&student.Name,
		&student.Age,
		&student.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 404, "student not found")
		return
	}
	if err != nil {
		ErrorResponse(c, 500, "query student failed")
		return
	}

	result, err := DB.Exec("delete from students where id = ?", id)
	if err != nil {
		ErrorResponse(c, 500, "delete student failed")
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		ErrorResponse(c, 500, "get affected rows failed")
		return
	}
	if affected == 0 {
		ErrorResponse(c, 404, "not found")
		return
	}
	JSONResponse(c, 200, "student deleted successfully", student)
}

func updateStudentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ErrorResponse(c, 400, "invalid id")
		return
	}

	var req StudentStatusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "request body failed")
		return
	}

	req.Status = strings.TrimSpace(req.Status)

	if req.Status == "" {
		ErrorResponse(c, 400, "status is required")
		return
	}
	if req.Status != "enabled" && req.Status != "disabled" {
		ErrorResponse(c, 400, "status must be enabled or disabled")
		return
	}

	result, err := DB.Exec(
		"UPDATE students SET status = ? WHERE id = ?",
		req.Status,
		id,
	)
	if err != nil {
		ErrorResponse(c, 500, "update student status failed")
		return
	}

	// 5. 获取影响行数
	affected, err := result.RowsAffected()
	if err != nil {
		ErrorResponse(c, 500, "get affected rows failed")
		return
	}

	// 6. 影响 0 行时，区分“不存在”和“状态没变化”
	if affected == 0 {
		var existingID int

		err = DB.QueryRow(
			"SELECT id FROM students WHERE id = ?",
			id,
		).Scan(&existingID)

		if errors.Is(err, sql.ErrNoRows) {
			ErrorResponse(c, 404, "student not found")
			return
		}

		if err != nil {
			ErrorResponse(c, 500, "query student failed")
			return
		}

		// 运行到这里说明学生存在，只是状态原本就一样
	}

	// 7. 查询修改后的完整学生数据
	var student Student

	err = DB.QueryRow(
		`SELECT id, student_no, name, age, status
				 FROM students
				 WHERE id = ?`,
		id,
	).Scan(
		&student.ID,
		&student.StudentNo,
		&student.Name,
		&student.Age,
		&student.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 404, "student not found")
		return
	}

	if err != nil {
		ErrorResponse(c, 500, "query updated student failed")
		return
	}

	// 8. 返回修改后的数据
	JSONResponse(c, 200, "student status updated successfully", student)
}
