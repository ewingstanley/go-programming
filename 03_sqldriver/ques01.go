package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Student 
type Student struct {
	ID    int
	Name  string
	Age   int
	Grade string
}

var db *sql.DB


func initDB() {
	// 数据库连接
	dsn := "root:123456@tcp(localhost:3306)/dbtest?charset=utf8mb4&parseTime=True&loc=Local"
	
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库 Ping 失败: ", err)
	}

	fmt.Println("数据库连接成功!")
}

// createTable 创建学生表
func createTable() {
	// SQL 语句创建表
	query := `
	CREATE TABLE IF NOT EXISTS students (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		age INT NOT NULL,
		grade VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("创建表失败: ", err)
	}

	fmt.Println("学生表创建成功!")
}

// insertStudent 插入学生记录
func insertStudent(student Student) (int64, error) {
	query := "INSERT INTO students (name, age, grade) VALUES (?, ?, ?)"
	
	result, err := db.Exec(query, student.Name, student.Age, student.Grade)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Printf("插入成功，学生ID: %d\n", id)
	return id, nil
}

// getStudents 查询所有学生
func getStudents() ([]Student, error) {
	query := "SELECT id, name, age, grade FROM students"
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Grade)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, nil
}

// getStudentsByAge 根据年龄查询学生
func getStudentsByAge(minAge int) ([]Student, error) {
	query := "SELECT id, name, age, grade FROM students WHERE age > ?"
	
	rows, err := db.Query(query, minAge)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Grade)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, nil
}

// updateStudentGrade 更新学生年级
func updateStudentGrade(name, newGrade string) (int64, error) {
	query := "UPDATE students SET grade = ? WHERE name = ?"
	
	result, err := db.Exec(query, newGrade, name)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	fmt.Printf("更新了 %d 行数据\n", rowsAffected)
	return rowsAffected, nil
}

// deleteStudentByAge 删除年龄小于指定值的学生
func deleteStudentByAge(maxAge int) (int64, error) {
	query := "DELETE FROM students WHERE age < ?"
	
	result, err := db.Exec(query, maxAge)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	fmt.Printf("删除了 %d 行数据\n", rowsAffected)
	return rowsAffected, nil
}

// getStudentByID 根据ID查询学生
func getStudentByID(id int) (Student, error) {
	query := "SELECT id, name, age, grade FROM students WHERE id = ?"
	
	var student Student
	err := db.QueryRow(query, id).Scan(&student.ID, &student.Name, &student.Age, &student.Grade)
	if err != nil {
		return Student{}, err
	}

	return student, nil
}

func main() {
	// 初始化数据库连接
	initDB()
	defer db.Close()

	// 创建表
	createTable()

	// 1. 插入学生记录
	fmt.Println("\n1. 插入学生记录:")
	student := Student{Name: "张三", Age: 20, Grade: "三年级"}
	insertID, err := insertStudent(student)
	if err != nil {
		log.Fatal("插入失败: ", err)
	}

	// 再插入几个测试数据
	insertStudent(Student{Name: "李四", Age: 19, Grade: "二年级"})
	insertStudent(Student{Name: "王五", Age: 22, Grade: "四年级"})
	insertStudent(Student{Name: "赵六", Age: 14, Grade: "一年级"})

	// 2. 查询所有学生
	fmt.Println("\n2. 所有学生信息:")
	allStudents, err := getStudents()
	if err != nil {
		log.Fatal("查询失败: ", err)
	}
	for _, s := range allStudents {
		fmt.Printf("ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n", 
			s.ID, s.Name, s.Age, s.Grade)
	}

	// 3. 查询年龄大于18的学生
	fmt.Println("\n3. 年龄大于18的学生:")
	adultStudents, err := getStudentsByAge(18)
	if err != nil {
		log.Fatal("查询失败: ", err)
	}
	for _, s := range adultStudents {
		fmt.Printf("ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n", 
			s.ID, s.Name, s.Age, s.Grade)
	}

	// 4. 根据ID查询特定学生
	fmt.Println("\n4. 根据ID查询学生:")
	specificStudent, err := getStudentByID(int(insertID))
	if err != nil {
		log.Fatal("查询失败: ", err)
	}
	fmt.Printf("ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n", 
		specificStudent.ID, specificStudent.Name, specificStudent.Age, specificStudent.Grade)

	// 5. 更新学生年级
	fmt.Println("\n5. 更新张三的年级:")
	_, err = updateStudentGrade("张三", "四年级")
	if err != nil {
		log.Fatal("更新失败: ", err)
	}

	// 6. 删除年龄小于15的学生
	fmt.Println("\n6. 删除年龄小于15的学生:")
	_, err = deleteStudentByAge(15)
	if err != nil {
		log.Fatal("删除失败: ", err)
	}

	// 7. 验证删除结果
	fmt.Println("\n7. 删除后的学生信息:")
	finalStudents, err := getStudents()
	if err != nil {
		log.Fatal("查询失败: ", err)
	}
	for _, s := range finalStudents {
		fmt.Printf("ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n", 
			s.ID, s.Name, s.Age, s.Grade)
	}
}