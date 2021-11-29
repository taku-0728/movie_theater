package main

import (
    "fmt"
    "strings"
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
        fmt.Errorf("Some context: %v", err)
        return
    }

    // HTMLパース
    doc, err := goquery.NewDocumentFromReader(utfBody)

    if err != nil {
        fmt.Println("HTMLのパースに失敗しました。")
        return
    }

    place := "東京都"

    theaterList := []string{}
    // // titleを抜き出し
    doc.Find(".section > h1").Each(func(_ int, page *goquery.Selection) {
        // HTMLの中から劇場一覧のDOMを指定
        if strings.Index(page.Text(), "劇場一覧") != -1 {
            // 劇場一覧のDOMの中から都道府県を指定
            page.Next().Find(".theater-list-area > h4").Each(func(_ int, prefectures *goquery.Selection) {
                if strings.Index(prefectures.Text(), place) != -1 {
                    node := prefectures.Next().Find(".item > a > span").Nodes

                    for i := 0; i < len(node); i++ {
                        theaterList = append(theaterList, node[i].FirstChild.Data)
                    }
                    fmt.Println(theaterList)
                }
            })
        }
    })
}