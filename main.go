﻿// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// https://github.com/line/line-bot-sdk-go/tree/master/linebot

package main

import (
	"strconv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Constants
var timeFormat = "01/02 PM03:04:05"
var user_zchien = "U696bcb700dfc9254b27605374b86968b"
var user_yaoming = "U3aaab6c6248bb38f194134948c60f757"
var user_jackal = "U3effab06ddf5bcf0b46c1c60bcd39ef5"
var user_shane = "U2ade7ac4456cb3ca99ffdf9d7257329a"
var user_tenshi = "U1afa490979eb023cf3ed0db8c233f61e"
// Global Settings
var channelSecret = os.Getenv("CHANNEL_SECRET")
var channelToken = os.Getenv("CHANNEL_TOKEN")
//var baseURL = os.Getenv("APP_BASE_URL")
var baseURL = "https://line-talking-bot-go.herokuapp.com"
var endpointBase = os.Getenv("ENDPOINT_BASE")
var tellTimeInterval int = 15
var answers_TextMessage = []string{
		"人被殺，就會死。",
		"凡是每天喝水的人，有高機率在100年內死去",
		"今年中秋節剛好是滿月、今年七夕恰逢鬼月、今年母親節正好是星期日",
		"只要每天省下買一杯奶茶的錢，十天後就能買十杯奶茶",
		"台灣人在睡覺時，大多數的美國人都在工作",
		"台灣競爭力低落，在美國就連小學生都會說流利的英語",
		"在非洲，每六十秒，就有一分鐘過去",
		"每呼吸60秒，就減少一分鐘的壽命",
		"身高170cm的女生看起來和身高1米7的女生一樣高",
		"英國研究證實，全世界的人口中，減去瘦子的人口數後，剩下來的都是胖子。",
		"張開你的眼睛！否則，你將什麼都看不見。",
		"嗯嗯，呵呵，我要去洗澡了",
		"當一個便當吃不飽時.你可以吃兩個",
		"當你吃下吃下廿碗白飯，換算竟相當於吃下了二十碗白飯的熱量",
		"當你的左臉被人打，那你的左臉就會痛",
		"當蝴蝶在南半球拍了兩下翅膀，牠就會稍微飛高一點點",
		"誰能想的到，這名16歲少女，在四年前，只是一名12歲少女",
		"據統計，未婚生子的人數中有高機率為女性",
		"在非洲，每一分鐘，就有六十秒過去。",
		"在你的面前閉氣的話，就會不能呼吸喔。",
		"跟你在一起時，回憶一天前的事，就像回想昨天的事情。",
		"你不在的這十二個月，對我來說如同一年般長。",
		"不知道為什麼，把眼睛矇上後什麼都看不到。",
		"出生時，大家都是裸體的喔。",
		"英國研究 生日過越多的人就越老",
		"歲數越長活的越久",
		"當別人贏過你時，你就輸了！",
		"研究指出日本人的母語是日語",
		"你知道嗎 當你背對太陽 你就看不見金星",
		"當你失眠的時候，你就會睡不著",
		"今天是昨天的明天。",
		"吃得苦中苦，那一口特別苦",
	        "我愛你",
	        "天使是無敵的",
	        "本機器人作者是天使lineid：lovemountain2",
	}
var answers_ImageMessage = []string{
		"傳這甚麼廢圖？你有認真在分享嗎？快發色圖啦",
	}
var answers_StickerMessage = []string{
		"腳踏實地打字好嗎？傳這甚麼貼圖？又不是色圖",
	}
var answers_VideoMessage = []string{
		"看甚麼影片，不知道我的流量快用光了嗎？真是的又不是av",
	}
var answers_AudioMessage = []string{
		"說的比唱的好聽，唱得鬼哭神號，是要嚇唬誰？",
	}
var answers_LocationMessage = []string{
		"這是哪裡啊？火星嗎？",
	}
var answers_ReplyCurseMessage = []string{
		"真的無恥",
		"有夠無恥",
		"超級無恥",
		"就是無恥",
	        "超級無敵無恥",
	        "和天使一樣無恥（跑",
                "就是啊對你太失望了死變態",
	}
var answers_Replyyansheng = []string{
	       "又在說燕神的事了",
	       "來啊喂又是燕神啊",
        }
var answers_Replyxiang = []string{
	       "祥這個名字讓我想吐",
	       "請不要再提祥這個讓人想吐的名字了",
        }
var answers_ReplyQQ = []string{
	      "Q什麼Q",
	      "再Q試試看",
        }      
var silentMap = make(map[string]bool) // [UserID/GroupID/RoomID]:bool
var washMap = make(map[string]bool)
var AwashMap = make(map[string]bool)
var clearMap = make(map[string]bool)
var highCMap = make(map[string]bool)
//var echoMap = make(map[string]bool)

var loc, _ = time.LoadLocation("Asia/Taipei")
var bot *linebot.Client


func tellTime(replyToken string, doTell bool){
	var silent = false
	now := time.Now().In(loc)
	nowString := now.Format(timeFormat)
	
	if doTell {
		log.Println("現在時間(台北): " + nowString)
		bot.ReplyMessage(replyToken, linebot.NewTextMessage("現在時間(台北): " + nowString)).Do()
	} else if silent != true{
		log.Println("自動報時(台北): " + nowString)
		bot.PushMessage(replyToken, linebot.NewTextMessage("自動報時(台北): " + nowString)).Do()
	} else {
		log.Println("tell time misfired")
	}
}

func tellTimeJob(sourceId string) {
	for {
		time.Sleep(time.Duration(tellTimeInterval) * time.Minute)
		now := time.Now().In(loc)
		log.Println("time to tell time to : " + sourceId + ", " + now.Format(timeFormat))
		tellTime(sourceId, false)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	go func() {
		tellTimeJob(user_zchien);
	}()
	go func() {
		for {
			now := time.Now().In(loc)
			log.Println("keep alive at : " + now.Format(timeFormat))
			//http.Get("https://line-talking-bot-go.herokuapp.com")
			time.Sleep(time.Duration(rand.Int31n(29)) * time.Minute)
		}
	}()
	
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	log.Print("URL:"  + r.URL.String())
	
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		var replyToken = event.ReplyToken

		var source = event.Source //EventSource		
		var userId = source.UserID
		var groupId = source.GroupID
		var roomId = source.RoomID
		log.Print("callbackHandler to source UserID/GroupID/RoomID: " + userId + "/" + groupId + "/" + roomId)
		
		var sourceId = roomId
		if sourceId == "" {
			sourceId = groupId
			if sourceId == "" {
				sourceId = userId
			}
		}
		
		if event.Type == linebot.EventTypeMessage {
			_, silent := silentMap[sourceId]
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				log.Print("ReplyToken[" + replyToken + "] TextMessage: ID(" + message.ID + "), Text(" + message.Text  + "), current silent status=" + strconv.FormatBool(silent) )
				//if _, err = bot.ReplyMessage(replyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				//	log.Print(err)
				//}
				
				if source.UserID != "" && source.UserID != user_zchien {
					profile, err := bot.GetProfile(source.UserID).Do()
					if err != nil {
						log.Print(err)
					} else if _, err := bot.PushMessage(user_zchien, linebot.NewTextMessage(profile.DisplayName + ": "+message.Text)).Do(); err != nil {
							log.Print(err)
					}
				}
				
				if strings.Contains(message.Text, "你閉嘴") {
					silentMap[sourceId] = true
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("QQ")).Do()
				} else if strings.Contains(message.Text, "現在幾點") {
					tellTime(replyToken, true)
				}else if strings.Contains(message.Text, "try") {
					bot.findAndAddContactsByMid("Ub01e3901e840a9dd76a1bf6904a93cee")
				} else if "說吧" == message.Text {
					silentMap[sourceId] = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("麥克風測試，1、2、3... OK")).Do()
				}else if "洗版開始密碼010220" == message.Text {
					washMap[sourceId] = true
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("準備完畢")).Do()
				}else if "洗版結束密碼010220" == message.Text {
					washMap[sourceId] = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("已關閉")).Do()
				}else if "全使用權關閉" == message.Text {
					if source.UserID != user_tenshi {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("你滾，這不是你可以用的")).Do()
					}
					if source.UserID == user_tenshi{
			                AwashMap[sourceId] = false
					clearMap[sourceId] = false
					washMap[sourceId] = false
					silentMap[sourceId] = true
					highCMap[sourceId] = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("全功能已經鎖上")).Do()
					}
			         }else if "終結開始使用者確認" == message.Text {
					if source.UserID != user_tenshi {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("你不是我的主人，無法使用此功能")).Do()
					}
					if source.UserID == user_tenshi{
					clearMap[sourceId] = true
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("使用者確認，終結模式開始(只不過還在寫 沒完成w）")).Do()
					}
				}else if "終結結束使用者確認" == message.Text {
					if source.UserID != user_tenshi {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("你不是我的主人，無法使用此功能")).Do()
					}
					if source.UserID == user_tenshi{
					clearMap[sourceId] = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("使用者確認，終結模式結束(只不過還在寫 沒完成w）")).Do()
					}
				}else if "使用者確認" == message.Text {
					if source.UserID != user_tenshi {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("你沒有通過認證")).Do()
					}
					if source.UserID == user_tenshi{
					highCMap[sourceId] = true
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("使用者確認，歡迎回家我的主人")).Do()
					}
				}else if "結束使用權" == message.Text {
					if source.UserID != user_tenshi {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("你又沒獲得過")).Do()
					}
					if source.UserID == user_tenshi{
					highCMap[sourceId] = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("使用者確認，一路小心我的主人")).Do()
					}
				}else if "全體洗版允許" == message.Text {
					if highCMap[sourceId] == false {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("還沒通過使用者認證，請說‘使用者確認’")).Do()
					}
					if highCMap[sourceId] == true {
					AwashMap[sourceId] = true
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("主人確認，已允許全體使用洗版功能")).Do()
					}
				}else if "全體洗版關閉" == message.Text {
					if highCMap[sourceId] == false {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("還沒通過使用者認證，請說‘使用者確認’")).Do()
					}
					if highCMap[sourceId] == true {
					AwashMap[sourceId] = false
					bot.ReplyMessage(replyToken, linebot.NewTextMessage("主人確認，已關閉全體使用洗版功能")).Do()
					}
				}else if "profile" == message.Text {
					if source.UserID != "" {
						profile, err := bot.GetProfile(source.UserID).Do()
						if err != nil {
							log.Print(err)
						} else if _, err := bot.ReplyMessage(
							replyToken,
							linebot.NewTextMessage("Display name: "+profile.DisplayName + ", Status message: "+profile.StatusMessage)).Do(); err != nil {
								log.Print(err)
						}
					} else {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("很抱歉你並不是天使本人無法使用此功能")).Do()
					}
				} else if "深夜選擇" == message.Text {
					imageURL := baseURL + "/static/buttons/1040.jpg"
					//log.Print(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> "+imageURL)
					template := linebot.NewButtonsTemplate(
						imageURL, "歡迎你喔死變態", "對就是你死變態",
						linebot.NewURITemplateAction("thisav的傳送門", "http://thisav.com/"),
						linebot.NewMessageTemplateAction("按這個代表你不是變態", "我是死變態"),
						linebot.NewMessageTemplateAction("按這個代表你是變態", "我是無敵大變態"),
						linebot.NewMessageTemplateAction("不選擇", "我是個變態但是不敢承認"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Buttons alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if strings.Contains(message.Text, "變態") {
					template := linebot.NewConfirmTemplate(
						"你是變態嗎?",
						linebot.NewMessageTemplateAction("是", "我是個誠實的大變態"),
						linebot.NewMessageTemplateAction("不是", "我是個偷偷侵犯妹妹的死變態"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Confirm alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if strings.Contains(message.Text, "fgo")  {
					template := linebot.NewConfirmTemplate(
						"你誠實嗎?",
						linebot.NewMessageTemplateAction("我很誠實", "我很誠實告訴你我是個變態"),
						linebot.NewMessageTemplateAction("我不誠實", "我是個偷偷侵犯妹妹的死變態"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Confirm alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if strings.Contains(message.Text, "貪心")  {
					template := linebot.NewConfirmTemplate(
						"你貪心嗎?",
						linebot.NewMessageTemplateAction("我很貪心", "貪心到連我妹妹的同學都搞"),
						linebot.NewMessageTemplateAction("我不貪心", "我只搞我妹妹沒搞他同學"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Confirm alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if strings.Contains(message.Text, "騙人") {
					template := linebot.NewConfirmTemplate(
						"你騙人嗎?",
						linebot.NewMessageTemplateAction("我沒騙人", "是的我騙你我不是變態"),
						linebot.NewMessageTemplateAction("我在騙人", "我是個偷偷侵犯妹妹的死變態"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Confirm alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "洗版功能" == message.Text {
					template := linebot.NewConfirmTemplate(
						"你是變態嗎?",
						linebot.NewMessageTemplateAction("啟動洗版一次", "wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww笑死"),
						linebot.NewMessageTemplateAction("不啟動洗版", "好的已變回可愛的天使bot"),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Confirm alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if "av推薦" == message.Text {
					imageURL := baseURL + "/static/buttons/1040.jpg"
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(
							imageURL, "thisav", "屬於變態的天堂",
							linebot.NewURITemplateAction("thisav傳送門", "http://thisav.com/"),
							linebot.NewMessageTemplateAction("你不喜歡看av請按這裡", "我超愛看av的"),
						),
						linebot.NewCarouselColumn(
							imageURL, "h-game", "給喜歡動漫的你",
							linebot.NewURITemplateAction("h-game傳送門", "https://h-game123.com/"),
							linebot.NewMessageTemplateAction("按這裡代表你你不看色漫", "我超愛蘿莉的他們的鮑魚都好禁"),
						),
						linebot.NewCarouselColumn(
							imageURL, "xnxx", "給喜歡多元化的你",
							linebot.NewURITemplateAction("xnxx傳送門", "http://www.xnxx.com/"),
							linebot.NewMessageTemplateAction("按這裡代表你沒有妹妹", "我沒有妹妹但我有爽手"),
						),
						linebot.NewCarouselColumn(
							imageURL, "xhamster", "給喜歡外國人的你",
							linebot.NewURITemplateAction("xhamster傳送門", "https://xhamster.com/"),
							linebot.NewMessageTemplateAction("按這裡代表你不打飛機", "我每天3餐各打3次飛機"),
						),
					)
					if _, err := bot.ReplyMessage(
						replyToken,
						linebot.NewTemplateMessage("Carousel alt text", template),
					).Do(); err != nil {
						log.Print(err)
					}
				
				} else if strings.Contains(message.Text, "滾") {
					if source.UserID != user_tenshi {
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("你不是我的主人沒資格叫我走")).Do()
					}
					if source.UserID == user_tenshi {
							switch source.Type {
						case linebot.EventSourceTypeUser:
							bot.ReplyMessage(replyToken, linebot.NewTextMessage("我想走, 但是我走不了...")).Do()
						case linebot.EventSourceTypeGroup:
							bot.ReplyMessage(replyToken, linebot.NewTextMessage("我揮一揮衣袖 不帶走一片雲彩")).Do()
							bot.LeaveGroup(source.GroupID).Do()
						case linebot.EventSourceTypeRoom:
							bot.ReplyMessage(replyToken, linebot.NewTextMessage("我揮一揮衣袖 不帶走一片雲彩")).Do()
							bot.LeaveRoom(source.RoomID).Do()
						}
					}
						
					
				} else if strings.Contains(message.Text, "無恥") {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_ReplyCurseMessage[rand.Intn(len(answers_ReplyCurseMessage))])).Do()
				} else if strings.Contains(message.Text, "祥") {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_Replyxiang[rand.Intn(len(answers_Replyxiang))])).Do()
				} else if strings.Contains(message.Text, "燕") {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_Replyyansheng[rand.Intn(len(answers_Replyyansheng))])).Do()
				} else if strings.Contains(message.Text, "Q") {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_ReplyQQ[rand.Intn(len(answers_ReplyQQ))])).Do()
				} else if AwashMap[sourceId] == true {
					if message.Text == message.Text{
						bot.ReplyMessage(replyToken, linebot.NewTextMessage(message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text)).Do()
					}
				}else if washMap[sourceId] == true {
					if highCMap[sourceId] == false{
						bot.ReplyMessage(replyToken, linebot.NewTextMessage("還沒通過使用者認證，請說‘使用者確認’")).Do()
					}
					if highCMap[sourceId] == true{
					     if  message.Text == message.Text{
						  if source.UserID != user_tenshi{
							 washMap[sourceId] = false
							 highCMap[sourceId] = false
							 bot.ReplyMessage(replyToken, linebot.NewTextMessage("檢測到非主人命令，結束一切權限")).Do()
						                                }
						       if source.UserID == user_tenshi {
						       bot.ReplyMessage(replyToken, linebot.NewTextMessage(message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+
												    message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text+message.Text)).Do()
					                                               }
					                                       }
					                                }
					}else if silentMap[sourceId] != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_TextMessage[rand.Intn(len(answers_TextMessage))])).Do()
				}
			case *linebot.ImageMessage :
				log.Print("ReplyToken[" + replyToken + "] ImageMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), PreviewImageURL(" + message.PreviewImageURL + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_ImageMessage[rand.Intn(len(answers_ImageMessage))])).Do()
				}
			case *linebot.VideoMessage :
				log.Print("ReplyToken[" + replyToken + "] VideoMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), PreviewImageURL(" + message.PreviewImageURL + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_VideoMessage[rand.Intn(len(answers_VideoMessage))])).Do()
				}
			case *linebot.AudioMessage :
				log.Print("ReplyToken[" + replyToken + "] AudioMessage[" + message.ID + "] OriginalContentURL(" + message.OriginalContentURL + "), Duration(" + strconv.Itoa(message.Duration) + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_AudioMessage[rand.Intn(len(answers_AudioMessage))])).Do()
				}
			case *linebot.LocationMessage:
				log.Print("ReplyToken[" + replyToken + "] LocationMessage[" + message.ID + "] Title(" + message.Title  + "), Address(" + message.Address + "), Latitude(" + strconv.FormatFloat(message.Latitude, 'f', -1, 64) + "), Longitude(" + strconv.FormatFloat(message.Longitude, 'f', -1, 64) + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_LocationMessage[rand.Intn(len(answers_LocationMessage))])).Do()
				}
			case *linebot.StickerMessage :
				log.Print("ReplyToken[" + replyToken + "] StickerMessage[" + message.ID + "] PackageID(" + message.PackageID + "), StickerID(" + message.StickerID + ")" )
				if silent != true {
					bot.ReplyMessage(replyToken, linebot.NewTextMessage(answers_StickerMessage[rand.Intn(len(answers_StickerMessage))])).Do()
				}
			}
		} else if event.Type == linebot.EventTypePostback {
		} else if event.Type == linebot.EventTypeBeacon {
		}
	}
	
}
