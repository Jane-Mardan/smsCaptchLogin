package v1

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"server/api"
	. "server/common"
	"server/common/utils"
	"server/storage"
	"strconv"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	RE_CHINA_PHONE_NUM = regexp.MustCompile(`^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`)
)
var Player = playerController{}

type playerController struct {
}

type ExtAuthLoginReq interface {
	extAuthId() string
}

/* SmsCaptchReq[begins]. */
type SmsCaptchReq struct {
	Num         string `json:"phoneNum,omitempty" form:"phoneNum"`
	CountryCode string `json:"phoneCountryCode,omitempty" form:"phoneCountryCode"`
	Captcha     string `json:"smsLoginCaptcha,omitempty" form:"smsLoginCaptcha"`
}

func (req *SmsCaptchReq) extAuthId() string {
	return req.CountryCode + req.Num
}
func (req *SmsCaptchReq) redisKey() string {
	return "/cuisine/sms/captcha/" + req.extAuthId()
}

/* SmsCaptchReq[ends]. */

func (p *playerController) SMSCaptchaGet(c *gin.Context) {
	var req SmsCaptchReq
	err := c.ShouldBindQuery(&req)
	api.CErr(c, err)
	if err != nil || req.Num == "" || req.CountryCode == "" {
		//错误码
		c.Set(api.RET, -1000)
		return
	}
	redisKey := req.redisKey()
	ttl, err := storage.RedisManagerIns.TTL(redisKey).Result()
	api.CErr(c, err)
	if err != nil {
		c.Set(api.RET, -1000)
		return
	}
	// redis剩余时长校验
	if ttl >= ConstVals.Player.CaptchaMaxTTL {
		//错误码
		c.Set(api.RET, -1001)
		return
	}
	succRet := 1000
	resp := struct {
		Ret                        int `json:"ret"`
		GetSmsCaptchaRespErrorCode int `json:"getSmsCaptchaRespErrorCode"`
		SmsCaptchReq
	}{Ret: succRet}
	var captcha string
	if ttl >= 0 {
		// 续验证码时长，重置剩余时长
		storage.RedisManagerIns.Expire(redisKey, ConstVals.Player.CaptchaExpire)
		captcha = storage.RedisManagerIns.Get(redisKey).Val()
		if ttl >= ConstVals.Player.CaptchaExpire/4 {
			if succRet == 1000 {
				getSmsCaptchaRespErrorCode := sendMessage(req.Num, req.CountryCode, captcha)
				if getSmsCaptchaRespErrorCode != 0 {
					//错误码
					resp.Ret = -1002
					resp.GetSmsCaptchaRespErrorCode = getSmsCaptchaRespErrorCode
				}
			}
		}
	} else {
		// 校验通过，进行验证码生成处理
		captcha = strconv.Itoa(utils.Rand.Number(1000, 9999))
		if succRet == 1000 {
			getSmsCaptchaRespErrorCode := sendMessage(req.Num, req.CountryCode, captcha)
			if getSmsCaptchaRespErrorCode != 0 {
				//错误码
				resp.Ret = -1002
				resp.GetSmsCaptchaRespErrorCode = getSmsCaptchaRespErrorCode
			}
		}
		storage.RedisManagerIns.Set(redisKey, captcha, ConstVals.Player.CaptchaExpire)
	}
	resp.SmsCaptchReq.Captcha = captcha
	c.JSON(http.StatusOK, resp)
}

func (p *playerController) SMSCaptchaLogin(c *gin.Context) {
	var req SmsCaptchReq
	err := c.ShouldBindWith(&req, binding.FormPost)
	if err != nil || req.Num == "" || req.CountryCode == "" || req.Captcha == "" {
		c.Set(api.RET, -1001)
		return
	}
	redisKey := req.redisKey()
	captcha := storage.RedisManagerIns.Get(redisKey).Val()
	if captcha != req.Captcha {
		//错误码错误的验证码
		c.Set(api.RET, -1003)
		return
	}

	storage.RedisManagerIns.Del(redisKey)
	resp := struct {
		Ret int `json:"ret"`
	}{1000}

	c.JSON(http.StatusOK, resp)
}

type tel struct {
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
}
type captchaReq struct {
	Ext    string     `json:"ext"`
	Extend string     `json:"extend"`
	Params *[2]string `json:"params"`
	Sig    string     `json:"sig"`
	Sign   string     `json:"sign"`
	Tel    *tel       `json:"tel"`
	Time   int64      `json:"time"`
	Tpl_id int        `json:"tpl_id"`
}

func sendMessage(mobile string, nationcode string, captchaCode string) int {
	tel := &tel{
		Mobile:     mobile,
		Nationcode: nationcode,
	}
	//短信有效期hardcode
	captchaExpireMin := strconv.Itoa(int(ConstVals.Player.CaptchaExpireMin))
	params := [2]string{captchaCode, captchaExpireMin}
	appkey := "TODO"
	rand := strconv.Itoa(utils.Rand.Number(1000, 9999))
	now := utils.UnixtimeSec()

	hash := sha256.New()
	hash.Write([]byte("appkey=" + appkey + "&random=" + rand + "&time=" + strconv.FormatInt(now, 10) + "&mobile=" + mobile))
	md := hash.Sum(nil)
	sig := hex.EncodeToString(md)

	reqData := &captchaReq{
		Ext:    "",
		Extend: "",
		Params: &params,
		Sig:    sig,
		Sign:   "TODO",
		Tel:    tel,
		Time:   now,
		Tpl_id: 0, //TODO
	}
	reqDataString, err := json.Marshal(reqData)
	req := bytes.NewBuffer([]byte(reqDataString))
	if err != nil {
		return -1
	}
	resp, err := http.Post("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=1400150185&random="+rand,
		"application/json",
		req)
	if err != nil {
		fmt.Printf("resp err %v", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("response body err:%v", err)
	}
	type bodyStruct struct {
		Result int `json:"result"`
	}
	var body bodyStruct
	json.Unmarshal(respBody, &body)
	return body.Result
}
