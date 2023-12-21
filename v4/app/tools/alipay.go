package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go_code/lms/app/model"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AlipayConfig struct {
	AppID           string
	AlipayPublicKey string
	PrivateKey      string
}
type Book struct {
	Id    int64
	Title string
	Price float64
}

func HandlePayment(c *gin.Context) {
	config := &AlipayConfig{
		AppID:           "9021000133605979",
		AlipayPublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnyToBCwftdRDzT2ZiMkEm2PfAgF3ttHFaKoGfpoZ6O0pwyYXgdHNVxO/7R2xAVSSPQRSC335KK2gj3Nev/PRZZ65RK+7luIV+vqaacEuPOB2KaLfQjheIWoUGbamw1NyJB4faKeFNDIxvWOVkXyucaxwOF6dQk8Lr1hBvyCom8hUtWrw9p6aS539psByEE2Z+VoTROZhoSIJzAWcpc0hglgVGxE1qsnMB/5phzbMEaOy2hCPrlEYstLo4PKURucJBOIkW/NMXioUp0AVAUJLkNmcNjdSMgTGOwxgsL134mygVVQAwzRps/rRcE9P08RP2yRk5KAEuvfCSqn9lHn/iQIDAQAB",
		PrivateKey:      "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCfJOgELB+11EPNPZmIyQSbY98CAXe20cVoqgZ+mhno7SnDJheB0c1XE7/tHbEBVJI9BFILffkoraCPc16/89FlnrlEr7uW4hX6+pppwS484HYpot9COF4hahQZtqbDU3IkHh9op4U0MjG9Y5WRfK5xrHA4Xp1CTwuvWEG/IKibyFS1avD2nppLnf2mwHIQTZn5WhNE5mGhIgnMBZylzSGCWBUbETWqycwH/mmHNswRo7LaEI+uURiy0ujg8pRG5wkE4iRb80xeKhSnQBUBQkuQ2Zw2N1IyBMY7DGCwvXfibKBVVADDNGmz+tFwT0/TxE/bJGTkoAS698JKqf2Uef+JAgMBAAECggEAOZl6AEiYEY+KGra44zEeYb9775XoZlr2QDOJtjjAN/Xer6sRxwLQdzvGs2OTtQ+O/laZ+17U10xAWKtBF+h/WXBhTeLs2mdp3TTmvnAU9COpoNg6RhVwvFdQfx4ErZ1+KmRcqdw38fRY3Fs2vzmJSctHsp9L+7vwVr5yzWJcpEUDHc7H6ofOceamhARxt7c+Ci+mW7Usmkn4n/yY/dsiOHEkIP2jig1mZVtHW+l3oZGwADafDx//14ZHbQGmHL6xBJfdfL8Qe6Q0bR8lTb98QXvCNvjWBeLiiuEKn8hu5CgQ1+dDUhCzCMaS1r85A2qp4JmZgPSksrzlQaPdXab1IQKBgQDZmz3mIk8pE/QKjUYAOpQJ5tahZ7q+G2sJOxUV9eQmAATCfd+NumFg9s01kziVEcOAS5nQmQQ/Nh4M/IxVStQAuYwR08Nbv252ohhEyBIWFJu1O9lyFRjrk672IrkOTR34MFVXpOb0QXH5vVQ5xuaFZ72Sjj1fY49WQ3mvNViU9wKBgQC7ORC/RGIq6vA59NaCG7yb4rSnL6yKv5QW0igecpBPTqhRpCshrsdjWmQX0JWeOuslU2v2MV7+A61b0EU68algJ1KYp9M9bxvMzW1gtG5aQYeu//InneeITxYKGnXg20EiOwKvD8EDjFVCGYf/JDmZohMMasgfFqKdmnlGU1dvfwKBgF/yzJNJs3/YVXXFnwUAzz+icibPFw82BbcXPw/k0QlUXsTnPYg/kypvCELEPGG2aG6MZzEMF0xL72oofTQqf1omdjVyGyS8Pte+V1cUpKzpv1JlJlbgKJHPF4ld/BzrMfi6TxcLFe7DYJ8OtTGYmxJdkMArSbJistQFgoUXz4w3AoGASWDmVXCh4RxPpwd1A22HkLlcAKLIx2Cq0/7uRnC0asDza3wig4MTFreYv1S7L1b8TpbRm55iEsCyM2f/mkiZD9yZnMc6Hbvsc2qYUeyly64fVdFuwWQ8GMqWYsNYLBcDAj2Kob5U8eUQjWWgTbmO8c36FUXMxZKDja27lnqXrFUCgYEA0LgGh4SUzOkFRWrr1OpvPv6gDRZOgVgOxnV6eyUNcTwdtsnyOC9FbI0WK13II+jujL34/t587ud78hZ6xf5y9UFvqXJOJmoQBHUTk60ZGqe4bTzHIUkSmfbnOgMBKMBn5x++7rX3JiW0trS97HtbUUczyHvd7b0Iz4p5jPxcy2c=",
	}
	// 创建支付宝客户端
	client, err := alipay.New(config.AppID, config.PrivateKey, false)
	// 从全局中间件中获取支付宝客户端
	c.Set("alipay", client)
	client, _ = c.MustGet("alipay").(*alipay.Client)
	var id int64
	idStr := c.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetBookName(strconv.FormatInt(id, 10))
	// 创建支付请求参数
	p := alipay.TradePagePay{}
	//p.NotifyURL = "https://www.bilibili.com" // 设置支付宝回调通知URL
	//p.ReturnURL = "https://www.bilibili.com" // 设置支付成功后跳转的URL
	p.Subject = ret.Title
	p.OutTradeNo = strconv.FormatInt(time.Now().Unix(), 10)
	p.TotalAmount = strconv.FormatFloat(ret.Price, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	// 发起支付请求
	result, err := client.TradePagePay(p)
	if err != nil {
		// 处理支付请求错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Max-Age", "86400")
	c.Redirect(http.StatusFound, result.String())

	go func() {
		// 假设订单超时时间为30分钟
		timeout := time.NewTimer(30 * time.Second)
		select {
		case <-timeout.C:
			// 关闭订单的逻辑
			CloseOrder(p.OutTradeNo)
		}
	}()
}
func CloseOrder(orderNo string) {
	client, _ := alipay.New("9021000133605979", "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCfJOgELB+11EPNPZmIyQSbY98CAXe20cVoqgZ+mhno7SnDJheB0c1XE7/tHbEBVJI9BFILffkoraCPc16/89FlnrlEr7uW4hX6+pppwS484HYpot9COF4hahQZtqbDU3IkHh9op4U0MjG9Y5WRfK5xrHA4Xp1CTwuvWEG/IKibyFS1avD2nppLnf2mwHIQTZn5WhNE5mGhIgnMBZylzSGCWBUbETWqycwH/mmHNswRo7LaEI+uURiy0ujg8pRG5wkE4iRb80xeKhSnQBUBQkuQ2Zw2N1IyBMY7DGCwvXfibKBVVADDNGmz+tFwT0/TxE/bJGTkoAS698JKqf2Uef+JAgMBAAECggEAOZl6AEiYEY+KGra44zEeYb9775XoZlr2QDOJtjjAN/Xer6sRxwLQdzvGs2OTtQ+O/laZ+17U10xAWKtBF+h/WXBhTeLs2mdp3TTmvnAU9COpoNg6RhVwvFdQfx4ErZ1+KmRcqdw38fRY3Fs2vzmJSctHsp9L+7vwVr5yzWJcpEUDHc7H6ofOceamhARxt7c+Ci+mW7Usmkn4n/yY/dsiOHEkIP2jig1mZVtHW+l3oZGwADafDx//14ZHbQGmHL6xBJfdfL8Qe6Q0bR8lTb98QXvCNvjWBeLiiuEKn8hu5CgQ1+dDUhCzCMaS1r85A2qp4JmZgPSksrzlQaPdXab1IQKBgQDZmz3mIk8pE/QKjUYAOpQJ5tahZ7q+G2sJOxUV9eQmAATCfd+NumFg9s01kziVEcOAS5nQmQQ/Nh4M/IxVStQAuYwR08Nbv252ohhEyBIWFJu1O9lyFRjrk672IrkOTR34MFVXpOb0QXH5vVQ5xuaFZ72Sjj1fY49WQ3mvNViU9wKBgQC7ORC/RGIq6vA59NaCG7yb4rSnL6yKv5QW0igecpBPTqhRpCshrsdjWmQX0JWeOuslU2v2MV7+A61b0EU68algJ1KYp9M9bxvMzW1gtG5aQYeu//InneeITxYKGnXg20EiOwKvD8EDjFVCGYf/JDmZohMMasgfFqKdmnlGU1dvfwKBgF/yzJNJs3/YVXXFnwUAzz+icibPFw82BbcXPw/k0QlUXsTnPYg/kypvCELEPGG2aG6MZzEMF0xL72oofTQqf1omdjVyGyS8Pte+V1cUpKzpv1JlJlbgKJHPF4ld/BzrMfi6TxcLFe7DYJ8OtTGYmxJdkMArSbJistQFgoUXz4w3AoGASWDmVXCh4RxPpwd1A22HkLlcAKLIx2Cq0/7uRnC0asDza3wig4MTFreYv1S7L1b8TpbRm55iEsCyM2f/mkiZD9yZnMc6Hbvsc2qYUeyly64fVdFuwWQ8GMqWYsNYLBcDAj2Kob5U8eUQjWWgTbmO8c36FUXMxZKDja27lnqXrFUCgYEA0LgGh4SUzOkFRWrr1OpvPv6gDRZOgVgOxnV6eyUNcTwdtsnyOC9FbI0WK13II+jujL34/t587ud78hZ6xf5y9UFvqXJOJmoQBHUTk60ZGqe4bTzHIUkSmfbnOgMBKMBn5x++7rX3JiW0trS97HtbUUczyHvd7b0Iz4p5jPxcy2c=", false)

	// 创建交易关闭请求参数
	closeReq := alipay.TradeClose{
		OutTradeNo: orderNo, // 要关闭的订单号
	}
	fmt.Println("Closing order:", orderNo)

	// 发起交易关闭请求
	closeRes, err := client.TradeClose(closeReq)
	if err != nil {
		// 处理交易关闭请求错误
		return
	}

	// 处理交易关闭结果
	if closeRes.Code != "10000" {
		// 交易关闭失败
		return
	}
	return
}
func HandleCallback(c *gin.Context) {
	// 获取请求中的所有参数
	params := make(map[string]string)
	c.Request.ParseForm()
	for key, values := range c.Request.Form {
		params[key] = values[0]
	}

	// 验证签名
	if VerifySign(params) {
		// 签名验证通过，处理业务逻辑
		// TODO: 在这里写下你的业务逻辑代码

		// 返回成功响应
		c.String(http.StatusOK, "success")
	} else {
		// 签名验证失败，返回错误响应
		c.String(http.StatusOK, "error")
	}
}

func VerifySign(params map[string]string) bool {
	// 将参数按照键名进行升序排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 拼接排序后的参数键值对
	var signStr string
	for _, key := range keys {
		if key == "sign" || key == "sign_type" {
			continue
		}
		value := params[key]
		signStr += key + "=" + value + "&"
	}
	signStr = strings.TrimRight(signStr, "&")

	publicKey := `MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCfJOgELB+11EPNPZmIyQSbY98CAXe20cVoqgZ+mhno7SnDJheB0c1XE7/tHbEBVJI9BFILffkoraCPc16/89FlnrlEr7uW4hX6+pppwS484HYpot9COF4hahQZtqbDU3IkHh9op4U0MjG9Y5WRfK5xrHA4Xp1CTwuvWEG/IKibyFS1avD2nppLnf2mwHIQTZn5WhNE5mGhIgnMBZylzSGCWBUbETWqycwH/mmHNswRo7LaEI+uURiy0ujg8pRG5wkE4iRb80xeKhSnQBUBQkuQ2Zw2N1IyBMY7DGCwvXfibKBVVADDNGmz+tFwT0/TxE/bJGTkoAS698JKqf2Uef+JAgMBAAECggEAOZl6AEiYEY+KGra44zEeYb9775XoZlr2QDOJtjjAN/Xer6sRxwLQdzvGs2OTtQ+O/laZ+17U10xAWKtBF+h/WXBhTeLs2mdp3TTmvnAU9COpoNg6RhVwvFdQfx4ErZ1+KmRcqdw38fRY3Fs2vzmJSctHsp9L+7vwVr5yzWJcpEUDHc7H6ofOceamhARxt7c+Ci+mW7Usmkn4n/yY/dsiOHEkIP2jig1mZVtHW+l3oZGwADafDx//14ZHbQGmHL6xBJfdfL8Qe6Q0bR8lTb98QXvCNvjWBeLiiuEKn8hu5CgQ1+dDUhCzCMaS1r85A2qp4JmZgPSksrzlQaPdXab1IQKBgQDZmz3mIk8pE/QKjUYAOpQJ5tahZ7q+G2sJOxUV9eQmAATCfd+NumFg9s01kziVEcOAS5nQmQQ/Nh4M/IxVStQAuYwR08Nbv252ohhEyBIWFJu1O9lyFRjrk672IrkOTR34MFVXpOb0QXH5vVQ5xuaFZ72Sjj1fY49WQ3mvNViU9wKBgQC7ORC/RGIq6vA59NaCG7yb4rSnL6yKv5QW0igecpBPTqhRpCshrsdjWmQX0JWeOuslU2v2MV7+A61b0EU68algJ1KYp9M9bxvMzW1gtG5aQYeu//InneeITxYKGnXg20EiOwKvD8EDjFVCGYf/JDmZohMMasgfFqKdmnlGU1dvfwKBgF/yzJNJs3/YVXXFnwUAzz+icibPFw82BbcXPw/k0QlUXsTnPYg/kypvCELEPGG2aG6MZzEMF0xL72oofTQqf1omdjVyGyS8Pte+V1cUpKzpv1JlJlbgKJHPF4ld/BzrMfi6TxcLFe7DYJ8OtTGYmxJdkMArSbJistQFgoUXz4w3AoGASWDmVXCh4RxPpwd1A22HkLlcAKLIx2Cq0/7uRnC0asDza3wig4MTFreYv1S7L1b8TpbRm55iEsCyM2f/mkiZD9yZnMc6Hbvsc2qYUeyly64fVdFuwWQ8GMqWYsNYLBcDAj2Kob5U8eUQjWWgTbmO8c36FUXMxZKDja27lnqXrFUCgYEA0LgGh4SUzOkFRWrr1OpvPv6gDRZOgVgOxnV6eyUNcTwdtsnyOC9FbI0WK13II+jujL34/t587ud78hZ6xf5y9UFvqXJOJmoQBHUTk60ZGqe4bTzHIUkSmfbnOgMBKMBn5x++7rX3JiW0trS97HtbUUczyHvd7b0Iz4p5jPxcy2c=`

	sign := params["sign"]
	// TODO: 进行签名验证的代码
	// 这里需要使用你自己的验签方法，示例中的 Verify 方法仅供参考
	valid := Verify(signStr, sign, publicKey)

	return valid
}

func Verify(signStr, sign, publicKey string) bool {
	// TODO: 实现你的验签逻辑
	// 这里使用支付宝提供的验签 SDK 进行验证，示例中的 Verify 方法仅供参考
	// 假设这里的 Verify 方法是支付宝提供的验签方法
	valid := Verify(signStr, sign, publicKey)

	return valid
}
