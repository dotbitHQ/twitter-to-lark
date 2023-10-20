package tool

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/parnurzeal/gorequest"
	"time"
)

const (
	TwitterSearchUrl = "https://api.twitter.com/2/tweets/search/recent?query=from:%s&tweet.fields=created_at&expansions=referenced_tweets.id&user.fields=created_at"
)

type TweetsSearchResp struct {
	Data []struct {
		Id               string `json:"id"`
		Text             string `json:"text"`
		CreatedAt        string `json:"created_at"`
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
		fmt.Println("TweetsSearch curl twitter err:", errs)
		err = fmt.Errorf("gorequest err :", errs[0].Error())
		return
	}

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println("TweetsSearch json.Unmarshal err:", err.Error())
		return
	}
	if len(resp.Data) == 0 {
		log.Info("twitter api response: ", resp)
	}
	return
}
