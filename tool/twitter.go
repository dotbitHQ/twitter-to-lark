package tool

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
)

const (
	TwitterSearchUrl = "https://api.twitter.com/2/tweets/search/recent?query=from:%s&expansions=referenced_tweets.id&user.fields=created_at"
)

type TweetsSearchResp struct {
	Data []struct {
		Id               string `json:"id"`
		Text             string `json:"text"`
		ReferencedTweets []struct {
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"referenced_tweets"`
	} `json:"data"`
	Includes struct {
		Tweets []struct {
			Id   string `json:"id"`
			Text string `json:"text"`
		}
	} `json:"includes"`
	Meta struct {
		ResultCount int `json:"result_count"`
	} `json:"meta"`
}

func TweetsSearch(token, twitterName string) (resp TweetsSearchResp, err error) {
	url := fmt.Sprintf(TwitterSearchUrl, twitterName)
	_, body, errs := gorequest.New().Get(url).Set("Authorization", fmt.Sprintf("Bearer %s", token)).Timeout(time.Second * 10).End()
	if len(errs) > 0 {
		fmt.Println("sendLarkTextNotify req err:", errs)
		err = fmt.Errorf("gorequest err :", errs[0].Error())
	}
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err.Error())
		return
	}
	return
}
