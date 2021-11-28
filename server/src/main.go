package main

import (
    "fmt"
    "net/http"
    "github.com/PuerkitoBio/goquery"
    "github.com/djimenez/iconv-go"
)

func main() {
    // TOHOシネマの劇場一覧サイト
    url := "https://www.tohotheater.jp/theater/find.html"

    // GETリクエスト
    res, _ := http.Get(url)

    // 呼び出し元の関数がreturnされるまで接続を切らない
    defer res.Body.Close()

    // 対象サイトのBody部分の読み取り
    utfBody, err := iconv.NewReader(res.Body, "Shift_JIS", "utf-8")

    if err != nil {
        fmt.Println("エンコーディングに失敗しました。")
    }

    // HTMLパース
    doc, _ := goquery.NewDocumentFromReader(utfBody)

    // // titleを抜き出し
    result := doc.Find("title").Text()
    fmt.Println(result)
}