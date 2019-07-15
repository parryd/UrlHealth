package main

import (
	"UrlHealth/crudfunc"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func doEvery(d time.Duration) {
	for x := range time.Tick(d) {
		fmt.Println(x)
		http.Get("http://127.0.0.1:8080/api/v1/urlChecks/")
	}
}

// func createurlCheck(c *gin.Context) {
// 	crawlTime, _ := strconv.Atoi(c.PostForm("crawlTime"))
// 	frequency, _ := strconv.Atoi(c.PostForm("frequency"))
// 	failureThreshold, _ := strconv.Atoi(c.PostForm("failureThreshold"))
// 	urlCheck := urlCheckModel{URL1: c.PostForm("url1"), CrawlTime: crawlTime, Frequency: frequency, FailureThreshold: failureThreshold}
// 	db.Save(&urlCheck)
// 	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "URLCheck item created successfully!", "resourceId": urlCheck.ID})
// }

// func checkAllurl(url1 string, crawlTime int, frequency int, failureThreshold int) {

// 	timeout := time.Duration(crawlTime) * time.Second
// 	client := http.Client{
// 		Timeout: timeout,
// 	}
// 	req, _ := http.NewRequest("GET", url1, nil)

// 	i := 1
// 	for ; i <= failureThreshold; i++ {
// 		fmt.Println("yo", i)
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println("ao", i)
// 			urlStat := model.URLStatusModel{URL1: url1, AttemptNo: i, Stat: false}
// 			db.Save(&urlStat)
// 			//c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": err})
// 		} else {
// 			fmt.Println("bo", i)
// 			urlStat := model.URLStatusModel{URL1: url1, AttemptNo: i, Stat: true}
// 			db.Save(&urlStat)
// 			defer resp.Body.Close()
// 			//c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "HTTP Response Status": resp.StatusCode, "error": err})
// 			break
// 		}
// 		//fmt.Println("yo ", i)
// 		time.Sleep(time.Duration(frequency) * time.Second)
// 	}

// }

// func checkurl(c *gin.Context) {
// 	var urlCheck urlCheckModel
// 	urlCheckID := c.Param("id")

// 	db.First(&urlCheck, urlCheckID)
// 	if urlCheck.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No url found!"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "HTTP Response Status": resp.StatusCode})
// }

func main() {
	crudfunc.Init()
	go doEvery(10 * time.Minute)
	router := gin.Default()
	v1 := router.Group("/api/v1/urlChecks")
	{
		v1.POST("/", crudfunc.CreateurlCheck)
		v1.GET("/", crudfunc.FetchAllurlCheck)
		v1.GET("/:id", crudfunc.FetchSingleurlCheck)
		v1.PUT("/:id", crudfunc.UpdateurlCheck)
		v1.DELETE("/:id", crudfunc.DeleteurlCheck)
	}
	router.Run()

}
