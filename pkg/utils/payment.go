package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"github.com/razorpay/razorpay-go"
)

// var (
// 	razorpayKeyID  = os.Getenv("RAZORPAY_KEY_ID")
// 	razorpaySecret = os.Getenv("RAZORPAY_SECRET")
// )

func Executerazorpay(totalPrice float64) (string, error) {
	razorpayKeyID := "rzp_test_MpuH3s9Tj28Inm"
	razorpaySecret := "mOsBkEDi0Z9RbeADEddoR4JV"
	client := razorpay.NewClient(razorpayKeyID, razorpaySecret)
	data := map[string]interface{}{
		"amount":   int(totalPrice) * 100,
		"currency": "INR",
		"receipt":  "101",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", errors.New("payment not initiated")
	}
	razorId, _ := body["id"].(string)
	return razorId, nil
}

func RazorPaymentVerification(sign, orderId, paymentId string) error {
	signature := sign
	secret := "mOsBkEDi0Z9RbeADEddoR4JV"
	data := orderId + "|" + paymentId
	h := hmac.New(sha256.New, []byte(secret))

	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return err
	} else {
		return nil
	}
}
