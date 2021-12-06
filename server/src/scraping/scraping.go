package scraping

import (
    "fmt"
    "strings"
    "net/http"
    "golang.org/x/text/transform"
    "github.com/PuerkitoBio/goquery"
    "github.com/djimenez/iconv-go"
    "github.com/labstack/echo"
    "bufio"
    "golang.org/x/text/encoding/japanese"
)

func GetMovieTheater(c echo.Context) error {

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
        return c.JSON(200, map[string]interface{}{"error": "エンコーディングに失敗しました。"})
    }

    // HTMLパース
    doc, err := goquery.NewDocumentFromReader(utfBody)

    if err != nil {
        return c.JSON(200, map[string]interface{}{"error": "HTMLのパースに失敗しました。"})
    }

    place := c.FormValue("prefectures")
    
    theaterList := []string{}
    // // titleを抜き出し
    doc.Find(".section > h1").Each(func(_ int, page *goquery.Selection) {
        // HTMLの中から劇場一覧のDOMを指定
        if strings.Index(page.Text(), "劇場一覧") != -1 {
            // 劇場一覧のDOMの中から都道府県を指定
            page.Next().Find(".theater-list-area > h4").Each(func(_ int, prefectures *goquery.Selection) {
                if strings.Index(prefectures.Text(), place) != -1 {
                    // 指定した都道府県の映画館のリンクを取得
                    prefectures.Next().Find(".item > a").Each(func(_ int, theaterLink *goquery.Selection) {
                        href, _ := theaterLink.Attr("href")
                        theaterList = append(theaterList, href)
                
                    })
                    // node, _ := prefectures.Next().Find(".item > a").Attr("href")

                    // for i := 0; i < len(node); i++ {
                    //     // theaterList = append(theaterList, node[i].FirstChild.Data)
                    //     theaterList = append(theaterList, node[i].FirstChild.Data)
                    // }
                }
            })
        }
    })

    result := []string{}

    // 都道府県に対応する劇場のサイトに入り上映状況を取得
    for _, theater := range theaterList {
        // 各劇場サイト
        url2 := "https://hlo.tohotheater.jp/" + theater

        // GETリクエスト
        res2, _ := http.Get(url2)

        // 対象サイトのBody部分の読み取り
        utfBody2 := transform.NewReader(bufio.NewReader(res2.Body), japanese.ShiftJIS.NewDecoder())

        if err != nil {
            fmt.Println("エンコーディングに失敗しました2")
            fmt.Errorf("Some context: %v", err)
            return c.JSON(200, map[string]interface{}{"error": "エンコーディングに失敗しました2"})
        }

        // HTMLパース
        doc2, err := goquery.NewDocumentFromReader(utfBody2)

        if err != nil {
            return c.JSON(200, map[string]interface{}{"error": err})
        }

        title2 := doc2.Find("title").Text()
        result = append(result, title2)
    }


    return c.JSON(200, map[string]interface{}{"hello": result})
}