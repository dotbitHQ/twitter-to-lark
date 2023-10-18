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
		if isSend := t.Rc.GetTweets2lark(v.Id); !isSend {

			inputTime, err := time.Parse(time.RFC3339, v.CreatedAt)
			if err != nil {
				fmt.Println("无法解析输入时间字符串:", err)
				return
			}
			newTime := inputTime.Add(8 * time.Hour)
			resultTimeStr := newTime.Format("2006-01-02 15:04:05")
			title := fmt.Sprintf("%s     %s", twitterName, resultTimeStr)
			notify.SendLarkTextNotify(larkKey, title, text)
			if err := t.Rc.SetTweets2lark(v.Id); err != nil {
				return fmt.Errorf("t.Rc.SetTweets2lark err: %s", err.Error())
			}
			continue
		}

	}
	return
}
