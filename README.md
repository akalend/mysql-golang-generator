# MySQL code generator 

### Installation

    go build mysql_gen.go

### Usage

    mysql_gen conn_string tablename
#### example of generate:

    ./mysql_gen password:user@(localhost)/database review
    type Review struct {
        id int64
        model string
        url string
        rate int64
        positive string
        negative string
        review string
        created int64
        title string
    }

    func getReview(db *sql.DB, id int64) Review {
        var(
            ret Review
            model sql.NullString
            url sql.NullString
            rate sql.NullInt64
            positive sql.NullString
            negative sql.NullString
            review sql.NullString
            created sql.NullInt64
            title sql.NullString
        )
        sql_s := " SELECT model,url,rate,positive,negative,review,created,title FROM review WHERE id=? "
        rows, err := db.Query(sql_s, id)
        errorCheck(err)
        defer rows.Close()
        for rows.Next() {
            err = rows.Scan(&model, &url, &rate, &positive, &negative, &review, &created, &title)
            errorCheck(err)
            ret.id=id
            ret.model=sql2String(model)
            ret.url=sql2String(url)
            ret.rate=sql2Int(rate)
            ret.positive=sql2String(positive)
            ret.negative=sql2String(negative)
            ret.review=sql2String(review)
            ret.created=sql2Int(created)
            ret.title=sql2String(title)
        }
        return ret
    }

 ### Usage  of result

    res.go