package app

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func sendVerificationSMS(phoneNumber, verificationCode string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI5tFn1kUktMJYDpspQNbF", "BWLA7TbKB07isD4jOQNF26M0lZGmdU")
	if err != nil {
		return err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phoneNumber
	request.SignName = "阿里云短信测试"
	request.TemplateCode = "SMS_154950909"
	request.TemplateParam = `{"code":"` + verificationCode + `"}`

	response, err := client.SendSms(request)
	if err != nil {
		return err
	}

	if response.Code != "OK" {
		return errors.New(response.Message)
	}

	return nil
}
