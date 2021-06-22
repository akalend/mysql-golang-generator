package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)
type Field struct {
	name  string
	ftype string
	sqltype string
	fn_conv string
}

type DbGen struct {
	Db *sql.DB
	fields []Field
	tablename string
	pk string
	type_pk string
}

func (dg *DbGen) getSchema(tabName string) {
	dg.tablename = tabName
	sql_1 := "DESCRIBE " + tabName
	rows, err := dg.Db.Query(sql_1)
	errorCheck(err)
	defer rows.Close()
	var fieldName string
	var fieldType string
	var fieldIsNull sql.NullString
	var fieldDefault sql.NullString
	var fieldComment string
	var isKey string
	for rows.Next() {
		err = rows.Scan(&fieldName, &fieldType, &fieldIsNull, &isKey, &fieldDefault, &fieldComment)
		errorCheck(err)
		type_out := "string"
		sql_type := "sql.NullString"
		fn_conv := "sql2String"
		if strings.Index(fieldType, "int") >= 0 {
			type_out = "int64"
			sql_type = "sql.NullInt64"
			fn_conv  = "sql2Int"
		} else if strings.Index(fieldType, "char") >= 0 {
			type_out = "string"
			sql_type = "sql.NullString"
			fn_conv = "sql2String"
		} else if strings.Index(fieldType, "date") >= 0 {
			type_out = "string"
			sql_type = "sql.NullString"
			fn_conv = "sql2String"
		} else if strings.Index(fieldType, "double") >= 0 {
			type_out = "float"
			sql_type = "sql.NullFloat64"
			fn_conv  = "sql2Float"
		} else if strings.Index(fieldType, "text") >= 0 {
			type_out = "string"
			sql_type = "sql.NullString"
			fn_conv = "sql2String"
		}
		//fmt.Println(fieldName, fieldType , fieldIsNull, isKey)
		dg.fields = append(dg.fields, Field{fieldName, type_out, sql_type, fn_conv} )

		if isKey == "PRI" {
			dg.pk = fieldName
			dg.type_pk = type_out
		}
	}
}

func errorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (dg *DbGen)  initDb(cnn_str string) {
	var err error
	dg.Db, err = sql.Open("mysql", cnn_str)
	errorCheck(err)
}

func (dg *DbGen) generate() {
	var fieldList []string

	fmt.Printf("type %s struct {\n`",strings.Title(strings.ToLower(dg.tablename )))
	for _,field := range dg.fields {
		fmt.Printf("\t\t%s %s\n", field.name, field.ftype)
	}

	fmt.Println("}")

	fmt.Printf("\tfunc get%s(db *sql.DB, %s %s) %s {\n", strings.Title(strings.ToLower(dg.tablename)),
		dg.pk, dg.type_pk, strings.Title(strings.ToLower(dg.tablename )))
	fmt.Println("\t\tvar(")
	fmt.Printf("\t\t\tret %s\n", strings.Title(strings.ToLower(dg.tablename )))

	for _,field := range dg.fields {
		if field.name == dg.pk {
			continue
		}
		fieldList = append(fieldList, field.name)
		fmt.Printf("\t\t\t%s %s\n", field.name, field.sqltype)
	}
	fmt.Println("\t\t)")

	out_fieldList := strings.Join(fieldList, ",")
	out_vars := strings.Join(fieldList, ", &")


	sql_txt := fmt.Sprintf( "SELECT %s FROM %s WHERE %s=?", out_fieldList, dg.tablename, dg.pk)
	fmt.Printf("\t\tsql_s := \" %s \"\n" , sql_txt)
	fmt.Printf("\t\trows, err := db.Query(sql_s, %s)\n ", dg.pk)
	fmt.Println("\t\terrorCheck(err)")
	fmt.Println("\t\tdefer rows.Close()")

	fmt.Println("\t\tfor rows.Next() {")

	fmt.Printf("\t\t\terr = rows.Scan(&%s)\n", out_vars)
	fmt.Println("\t\t\terrorCheck(err)")

	for _,field := range dg.fields {
		if field.name == dg.pk {
			fmt.Printf("\t\t\tret.%s=%s\n", field.name, field.name)

		} else {
			fmt.Printf("\t\t\tret.%s=%s(%s)\n", field.name, field.fn_conv, field.name)
		}
	}

	fmt.Println("\t\t}")

	fmt.Println("\t\treturn ret")
	fmt.Println("\t}")


}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Printf("\t%s [connection string] [tablrname]\n", os.Args[0])
		os.Exit(0)
	}


	dbGen := DbGen{}
	cnn_str := os.Args[1]

	dbGen.initDb(cnn_str)
	dbGen.getSchema(os.Args[2])
	dbGen.generate()

	defer dbGen.Db.Close()
	fmt.Println("// ---------------- funush ------------------")
}


