package task

import (
	"fmt"
	"time"
	"twitter-to-lark/config"
	"twitter-to-lark/notify"
	"twitter-to-lark/tool"
)

type BalanceSample struct {
	Banlance uint64
	Time     time.Time
}

func (t *Task) doTwitter() (err error) {
	addrList := config.Cfg.Twitter.TwitterList
	for _, v := range addrList {
		twitterUsername := v.Username
		larkKey := v.LarkKey
		if twitterUsername == "" || larkKey == "" {
			return fmt.Errorf("config err ")
		}
		err = t.twitter2lark(twitterUsername, larkKey)
		if err != nil {
			return fmt.Errorf("t.twitter2lark err: %s", err.Error())
		}
	}
	return
}
func (t *Task) twitter2lark(twitterName, larkKey string) (err error) {
	token := config.Cfg.Twitter.BearerToken
	if token == "" {
		return fmt.Errorf("BearerToken is empty")
	}
	resp, err := tool.TweetsSearch(token, twitterName)
	if err != nil {
		return fmt.Errorf("tool.TweetsSearch err: %s", err.Error())
	}
	if resp.Meta.ResultCount == 0 || len(resp.Data) == 0 {
		return
	}
	for _, v := range resp.Data {
		text := v.Text
		if len(v.ReferencedTweets) > 0 {
			if v.ReferencedTweets[0].Type == "retweeted" {
				//search retweeted in includes
				for _, vv := range resp.Includes.Tweets {
					if vv.Id == v.ReferencedTweets[0].Id {
						text = vv.Text
					}
				}
			}
		}
		if isSend := t.Rc.GetTweets2lark("ccc"); !isSend {
			notify.SendLarkTextNotify(larkKey, twitterName, text)
			if err := t.Rc.SetTweets2lark("ccc"); err != nil {
				return fmt.Errorf("t.Rc.SetTweets2lark err: %s", err.Error())
			}
			continue
		}

	}
	return
}
