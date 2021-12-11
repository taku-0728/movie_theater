package scraping

import (
    // "encoding/json"
    "fmt"
    "strings"
    "net/http"
    "github.com/PuerkitoBio/goquery"
    "github.com/djimenez/iconv-go"
    "github.com/labstack/echo"
    "github.com/sclevine/agouti"
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

    // 入力された都道府県を取得
    place := c.FormValue("prefectures")
    
    theaterList := []string{}

    // 劇場一覧のHTMLから情報を取得
    doc.Find(".section > h1").Each(func(_ int, page *goquery.Selection) {
        // HTMLの中から劇場一覧の部分を指定
        if strings.Index(page.Text(), "劇場一覧") != -1 {
            // 劇場一覧の中から都道府県を指定
            page.Next().Find(".theater-list-area > h4").Each(func(_ int, prefectures *goquery.Selection) {
                // 入力された都道府県と一致する都道府県を探す
                if strings.Index(prefectures.Text(), place) != -1 {
                    // 指定した都道府県の劇場一覧のリンクを取得
                    prefectures.Next().Find(".item > a").Each(func(_ int, theaterLink *goquery.Selection) {
                        href, _ := theaterLink.Attr("href")
                        theaterList = append(theaterList, href)
                    })
                }
            })
        }
    })

    // JavaScriptを使った動的ページのHTMLを取得するためにChromeを利用することを宣言
    agoutiDriver := agouti.ChromeDriver(
        agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox", "--disable-dev-shm-usage"}),
    )
    
    if err := agoutiDriver.Start(); err != nil {
        fmt.Println("Failed to start driver:%v", err)
        return c.JSON(200, map[string]interface{}{"error": "agoutiDriver.Start()でエラー発生"})
    }

    defer agoutiDriver.Stop()
    page, err := agoutiDriver.NewPage()

    if err != nil {
        fmt.Println("NewPage()でエラー発生")
        fmt.Println("Some context: %v", err)
        return c.JSON(200, map[string]interface{}{"error": "NewPage()でエラー発生"})
    }

    // 入力された映画タイトルを取得
    title := c.FormValue("title")

    var theaterName []string
    var schedule [][]string

    // 都道府県に対応する劇場のサイトに入り上映状況を取得
    for _, theater := range theaterList {
        // 各劇場サイト
        url2 := "https://hlo.tohotheater.jp/" + theater

        // 各劇場サイトに入る
        page.Navigate(url2)

        // 動的ページのHTMLを格納
        dom, err := page.HTML()

        if err != nil {
            fmt.Println("page.HTML()でエラー発生")
            fmt.Println("Some context: %v", err)
        }

        // スクレイピングするためにDOMを読み込みなおす
        contents := strings.NewReader(dom)

        doc, err := goquery.NewDocumentFromReader(contents)
        if err != nil {
            fmt.Println("page.HTML()でエラー発生")
            fmt.Println("Some context: %v", err)
        }

        // 各劇場一覧ページのタイトル(映画館名)の「：」より前を格納
        theaterName = append(theaterName, doc.Find("title").Text()[:strings.Index(doc.Find("title").Text(), "：")])
        var time []string

        doc.Find(".schedule-body-section-item").EachWithBreak(func(_ int, page *goquery.Selection) bool {
            // 映画タイトルの中から入力されたタイトルと一致しているかを判別
            if strings.Contains(title, page.Find(".schedule-body-title").Text()) {
                // 入力されたタイトルと一致していれば上映スケジュールを取得
                page.Find(".schedule-item").Each(func(_ int, element *goquery.Selection){
                    text := element.Find(".start").Text() + "〜" + element.Find(".end").Text() + " " + element.Find(".status").Text()
                    time = append(time, text)
                })
                return false
            }
            return true
        })

        schedule = append(schedule, time)
        time = nil
    }

    result := map[int]map[string]interface{}{}

    for i := 0; i < len(theaterName); i++ {
        result[i] = map[string]interface{}{}
        result[i]["theaterName"] = theaterName[i]
        result[i]["schedule"] = schedule[i]
    }

    return c.JSON(200, result)
}