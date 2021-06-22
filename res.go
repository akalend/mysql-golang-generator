package main

import "database/sql"
import "fmt"
import _ "github.com/go-sql-driver/mysql"

func sql2String(str sql.NullString) string {
    if str.Valid {
        return str.String
    }
    return ""
}

func sql2Int(str sql.NullInt64) int64 {
    if str.Valid {
        return str.Int64
    }
    return 0
}

func sql2Float(str sql.NullFloat64) float64 {
    if str.Valid {
        return str.Float64
    }
    return 0.0
}

func errorCheck(err error) {
    if err != nil {
        panic(err.Error())
    }
}

type Review struct {
    id       int64
    model    string
    url      string
    rate     int64
    positive string
    negative string
    review   string
    created  int64
    title    string
}

func getReview(db *sql.DB, id int64) Review {
    var (
        ret      Review
        model    sql.NullString
        url      sql.NullString
        rate     sql.NullInt64
        positive sql.NullString
        negative sql.NullString
        review   sql.NullString
        created  sql.NullInt64
        title    sql.NullString
    )
    sql_s := " SELECT model,url,rate,positive,negative,review,created,title FROM review WHERE id=? "
    rows, err := db.Query(sql_s, id)
    errorCheck(err)
    defer rows.Close()
    for rows.Next() {
        err = rows.Scan(&model, &url, &rate, &positive, &negative, &review, &created, &title)
        errorCheck(err)
        ret.id = id
        ret.model = sql2String(model)
        ret.url = sql2String(url)
        ret.rate = sql2Int(rate)
        ret.positive = sql2String(positive)
        ret.negative = sql2String(negative)
        ret.review = sql2String(review)
        ret.created = sql2Int(created)
        ret.title = sql2String(title)
    }
    return ret
}

func initDb(cnn_str string) *sql.DB {

    Db, err := sql.Open("mysql", cnn_str)
    errorCheck(err)
    return Db
}
func main() {

    db := initDb("client:client@(localhost:3306)/nlp")
    res := getReview(db, 4)
    fmt.Println(res)
}
