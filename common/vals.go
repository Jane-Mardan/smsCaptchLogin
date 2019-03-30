package common

import (
	"regexp"
	"time"
)

var (
	RE_PHONE_NUM                = regexp.MustCompile(`^\+?[0-9]{8,14}$`)
	RE_CHINA_PHONE_NUM          = regexp.MustCompile(`^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`)
	RE_SMS_CAPTCHA_CODE         = regexp.MustCompile(`^[0-9]{4}$`)
	SmsExpiredSeconds           = 120
	SmsValidResendPeriodSeconds = 30
)

// 隐式导入
var ConstVals = &struct {
	Player struct {
		CaptchaExpire    time.Duration
		CaptchaMaxTTL    time.Duration
		CaptchaExpireMin int
	}
}{}

func constantsPost() {
	ConstVals.Player.CaptchaExpire = time.Duration(SmsExpiredSeconds) * time.Second
	ConstVals.Player.CaptchaMaxTTL = ConstVals.Player.CaptchaExpire -
		time.Duration(SmsValidResendPeriodSeconds)*time.Second
	ConstVals.Player.CaptchaExpireMin = SmsExpiredSeconds / 60
}
