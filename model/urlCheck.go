package model

import (
	"github.com/jinzhu/gorm"
)

type (
	// URLCheckModel is model
	URLCheckModel struct {
		gorm.Model
		URL1             string `json:"url1"`
		CrawlTime        int    `json:"crawlTime"`
		Frequency        int    `json:"frequency"`
		FailureThreshold int    `json:"failureThreshold"`
	}

	// TransformedurlCheck is model
	TransformedurlCheck struct {
		ID               uint   `json:"id"`
		URL1             string `json:"url1"`
		CrawlTime        int    `json:"crawlTime"`
		Frequency        int    `json:"frequency"`
		FailureThreshold int    `json:"failureThreshold"`
	}
)

type (
	// URLStatusModel is model
	URLStatusModel struct {
		gorm.Model
		URL1      string `json:"url1"`
		AttemptNo int    `json:"attemptNo"`
		Stat      bool   `json:"stat"`
	}

	// TransformedurlStatus is model
	TransformedurlStatus struct {
		ID        uint   `json:"id"`
		URL1      string `json:"url1"`
		AttemptNo int    `json:"attemptNo"`
		Stat      bool   `json:"stat"`
	}
)
