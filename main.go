package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"log"
	"strconv"
)

var templ = template.Must(template.ParseGlob("page/*"))

type Employee struct {
	Id int
	Name string
	City string
}

func dbConn() (db *sql.DB)  {
	dbDriver := "mysql"
	dbUser	 := "root"
	dbPasswd := "root123456"
	dbIP	 := "127.0.0.1"
	dbPort   := "3306"
	dbName   := "db_home"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPasswd+"@tcp("+dbIP+":"+dbPort+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}


func Index(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()
	rowsData, err := db.Query("SELECT * FROM EMPLOYEE ORDER BY ID DESC ")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for rowsData.Next() {
		var id int
		var name, city string
		err = rowsData.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	templ.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()
	keyId := r.URL.Query().Get("id")
	rowsData, err := db.Query("SELECT * FROM EMPLOYEE WHERE id=?", keyId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for rowsData.Next() {
		var id int
		var name, city string
		err := rowsData.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	templ.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request)  {
	templ.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()
	keyId := r.URL.Query().Get("id")
	rowsData, err :=db.Query("select * from employee where id=?", keyId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for rowsData.Next() {
		var id int
		var name, city string
		err := rowsData.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	templ.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		stmt, err:= db.Prepare("insert into employee(name,city) values (?,?)")
		if err != nil {
			panic(err.Error())
		}
		result, err := stmt.Exec(name, city)
		if err != nil {
			panic(err.Error())
		}
		LastInsertId, err := result.LastInsertId()
		log.Println("INSERT: Name: " + name +",City: " + city +",LastInsertId:" + strconv.FormatInt(LastInsertId, 10))
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		result, e := insForm.Exec(name, city, id)
		if e != nil {
			panic(e.Error())
		}
		RowsAffected, e := result.RowsAffected()
		log.Println("UPDATE: Name: " + name + " | City: " + city +",RowsAffected:"+strconv.FormatInt(RowsAffected, 10))
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	empId := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	result, e := delForm.Exec(empId)
	if e != nil {
		panic(e.Error())
	}
	RowsAffected, e := result.RowsAffected()
	log.Println("DELETE, id:", empId , ", RowsAffected:"+strconv.FormatInt(RowsAffected, 10))
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main()  {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
