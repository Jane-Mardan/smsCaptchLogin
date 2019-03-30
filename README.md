配置好golang环境后，go run main.go进行项目启动



接口说明

获取验证码

contect-type:application/x-www-form-urlencoded, method: get

发送内容：{
  phoneNum: String,
  phoneCountryCode: String
}

http://localhost:9992/api/player/v1/SmsCaptcha/get

返回： {
  ret: success
}





使用验证码登录

contect-type:application/x-www-form-urlencoded, method: post

发送内容：{
  phoneNum: String,
  phoneCountryCode: String，
  smsLoginCaptcha： string
}

http://localhost:9992/api/player/v1/SmsCaptcha/login

返回： {
  ret: success
}
