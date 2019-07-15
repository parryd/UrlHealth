package crudfunc

import (
	"UrlHealth/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

//Init to initialize
func Init() {
	var err error
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/urlCheck?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()

	db.AutoMigrate(&model.URLCheckModel{})
	db.AutoMigrate(&model.URLStatusModel{})
}

//CreateurlCheck is a function to store in db
func CreateurlCheck(c *gin.Context) {

	var urlCheck []model.URLCheckModel
	c.BindJSON(&urlCheck)

	for item := range urlCheck {
		db.Create(&urlCheck[item])
		c.JSON(200, urlCheck[item])
	}
}

// CheckAllurl is a func to check all the url
func CheckAllurl(url1 string, crawlTime int, frequency int, failureThreshold int) {

	timeout := time.Duration(crawlTime) * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	req, _ := http.NewRequest("GET", url1, nil)

	i := 1
	for ; i <= failureThreshold; i++ {
		fmt.Println("yo", i)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("ao", i)
			urlStat := model.URLStatusModel{URL1: url1, AttemptNo: i, Stat: false}
			db.Save(&urlStat)
			//c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": err})
		} else {
			fmt.Println("bo", i)
			urlStat := model.URLStatusModel{URL1: url1, AttemptNo: i, Stat: true}
			db.Save(&urlStat)
			defer resp.Body.Close()
			//c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "HTTP Response Status": resp.StatusCode, "error": err})
			break
		}
		//fmt.Println("yo ", i)
		time.Sleep(time.Duration(frequency) * time.Second)
	}

}

// UpdateurlCheck is func to update
func UpdateurlCheck(c *gin.Context) {
	var urlCheck model.URLCheckModel
	urlCheckID := c.Param("id")
	db.First(&urlCheck, urlCheckID)
	if urlCheck.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No url found!"})
		return
	}
	db.Model(&urlCheck).Update("url1", c.PostForm("url1"))
	crawlTime, _ := strconv.Atoi(c.PostForm("crawlTime"))
	frequency, _ := strconv.Atoi(c.PostForm("frequency"))
	failureThreshold, _ := strconv.Atoi(c.PostForm("failureThreshold"))
	db.Model(&urlCheck).Update("crawlTime", crawlTime)
	db.Model(&urlCheck).Update("frequency", frequency)
	db.Model(&urlCheck).Update("failureThreshold", failureThreshold)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "url updated successfully!"})
}

//DeleteurlCheck is to delete entry
func DeleteurlCheck(c *gin.Context) {
	var urlCheck model.URLCheckModel
	urlCheckID := c.Param("id")
	db.First(&urlCheck, urlCheckID)
	if urlCheck.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No url found!"})
		return
	}
	db.Delete(&urlCheck)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "url deleted successfully!"})
}

//FetchSingleurlCheck is to fetch and check one url
func FetchSingleurlCheck(c *gin.Context) {
	var urlCheck model.URLCheckModel
	urlCheckID := c.Param("id")
	db.First(&urlCheck, urlCheckID)
	if urlCheck.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Wrong entry"})
		return
	}

	timeout := time.Duration(urlCheck.CrawlTime) * time.Second

	client := http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest("GET", urlCheck.URL1, nil)

	i := 1
	for ; i <= urlCheck.FailureThreshold; i++ {
		resp, err := client.Do(req)
		if err != nil {
			//fmt.Println("yo2")
			urlStat := model.URLStatusModel{URL1: urlCheck.URL1, AttemptNo: i, Stat: false}
			db.Save(&urlStat)
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": err})
		} else {
			urlStat := model.URLStatusModel{URL1: urlCheck.URL1, AttemptNo: i, Stat: true}
			db.Save(&urlStat)
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "HTTP Response Status": resp.StatusCode, "error": err})
			break
		}

		//fmt.Println("yo ", i)

		time.Sleep(time.Duration(urlCheck.Frequency) * time.Second)

	}

	//defer resp.Body.Close()
	//fmt.Println("yo1")

	//_urlCheck := transformedurlCheck{ID: urlCheck.ID, URL1: urlCheck.URL1, CrawlTime: urlCheck.CrawlTime, Frequency: urlCheck.Frequency, FailureThreshold: urlCheck.FailureThreshold}

	//createChecked(urlCheck.URL1, resp.StatusCode)
}

//FetchAllurlCheck is to fetch and check all urls
func FetchAllurlCheck(c *gin.Context) {
	var urlChecks []model.URLCheckModel
	//var _urlChecks []transformedurlCheck
	db.Find(&urlChecks)
	if len(urlChecks) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No url found!"})
		return
	}

	for _, item := range urlChecks {

		go CheckAllurl(item.URL1, item.CrawlTime, item.Frequency, item.FailureThreshold)
		//_urlChecks = append(_urlChecks, transformedurlCheck{ID: item.ID, URL1: item.URL1, CrawlTime: item.CrawlTime, Frequency: item.Frequency, FailureThreshold: item.FailureThreshold})
		// timeout := time.Duration(item.CrawlTime) * time.Second

		// client := http.Client{
		// 	Timeout: timeout,
		// }

		// req, _ := http.NewRequest("GET", item.URL1, nil)

		// i := 1
		// for ; i <= item.FailureThreshold; i++ {
		// 	resp, err := client.Do(req)
		// 	if err != nil {
		// 		//fmt.Println("yo2")
		// 		urlStat := urlStatusModel{URL1: item.URL1, AttemptNo: i, Stat: false}
		// 		db.Save(&urlStat)
		// 		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": err})
		// 	} else {
		// 		urlStat := urlStatusModel{URL1: item.URL1, AttemptNo: i, Stat: true}
		// 		db.Save(&urlStat)
		// 		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "HTTP Response Status": resp.StatusCode, "error": err})
		// 		break
		// 	}
		// 	//fmt.Println("yo ", i)
		// 	time.Sleep(time.Duration(item.Frequency) * time.Second)
		// }
	}
	//c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _urlChecks})
}
