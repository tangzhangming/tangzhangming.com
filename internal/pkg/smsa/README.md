
参考 https://github.com/overtrue/easy-sms
顺序可调整
多短信网关
支持国际号码
可自定义网关
定义短信类

```
import "github.com/tangzhangming/go-smsa"

sms := smsa.New(smsa.cfg{
		Strategy : smsa.OrderStrategy, //调用顺序 顺序、随机、自定义

	})


phone   := 18500001111
content := "您的验证码为: 6379"
template := "templateID"
vars := map[string]string{

}


sms.Send(smsa.PhoneNumber(18500001111), input)

```