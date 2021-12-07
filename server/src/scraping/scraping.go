package scraping

import (
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
                }
            })
        }
    })

    result := []string{}
    title := c.FormValue("title")

    // Chromeを利用することを宣言
    agoutiDriver := agouti.ChromeDriver(
        agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox"}),
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

    // 都道府県に対応する劇場のサイトに入り上映状況を取得
    for _, theater := range theaterList {
        // 各劇場サイト
        url2 := "https://hlo.tohotheater.jp/" + theater
        page.Navigate(url2)

        dom, err := page.HTML()

        if err != nil {
            fmt.Println("page.HTML()でエラー発生")
            fmt.Println("Some context: %v", err)
        }

        contents := strings.NewReader(dom)

        doc, err := goquery.NewDocumentFromReader(contents)
        if err != nil {
            fmt.Println("page.HTML()でエラー発生")
            fmt.Println("Some context: %v", err)
        }

        // // titleを抜き出し
        doc.Find(".schedule-body-section-item").Each(func(_ int, page *goquery.Selection) {
            if strings.Contains(title, page.Find(".schedule-body-title").Text()) {
                page.Find(".schedule-item").Each(func(_ int, element *goquery.Selection){
                    text := element.Find(".start").Text() + "〜" + element.Find(".end").Text() + "：" + element.Find(".status").Text()
                    result = append(result, text)
                })
            } 
        })

        
    }


    return c.JSON(200, map[string]interface{}{"hello": result})
}