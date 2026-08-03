-- MySQL 8.0+
-- 本文件只用于保存数据库表结构，不会被 Go 程序自动执行。

CREATE TABLE IF NOT EXISTS students (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '学生 ID',
    student_no VARCHAR(50) NOT NULL COMMENT '学号',
    name VARCHAR(50) NOT NULL COMMENT '姓名',
    age TINYINT UNSIGNED NOT NULL COMMENT '年龄',
    status ENUM('enabled', 'disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态',
    PRIMARY KEY (id),
    UNIQUE KEY uk_students_student_no (student_no),
    CONSTRAINT chk_students_student_no_not_empty
        CHECK (CHAR_LENGTH(TRIM(student_no)) > 0),
    CONSTRAINT chk_students_name_not_empty
        CHECK (CHAR_LENGTH(TRIM(name)) > 0),
    CONSTRAINT chk_students_age
        CHECK (age BETWEEN 1 AND 150)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='学生信息表';
