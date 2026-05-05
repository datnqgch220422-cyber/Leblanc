package services

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

// Tạo URL thanh toán VNPay
func CreateVNPayURL(orderId string, amount int64, orderInfo string, ipAddr string) (string, error) {
	vnp_TmnCode := os.Getenv("VNPAY_TMN_CODE")
	vnp_HashSecret := os.Getenv("VNPAY_HASH_SECRET")
	vnp_Url := os.Getenv("VNPAY_URL")
	vnp_ReturnUrl := os.Getenv("VNPAY_RETURN_URL")
	if vnp_TmnCode == "" || vnp_HashSecret == "" || vnp_Url == "" || vnp_ReturnUrl == "" {
		return "", fmt.Errorf("vnpay credentials are missing")
	}
	if ipAddr == "" {
		ipAddr = "127.0.0.1"
	}

	vnp_Params := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    vnp_TmnCode,
		"vnp_Amount":     fmt.Sprintf("%d", amount*100), // VNPay yêu cầu nhân 100
		"vnp_CurrCode":   "VND",
		"vnp_TxnRef":     orderId,
		"vnp_OrderInfo":  orderInfo,
		"vnp_OrderType":  "other",
		"vnp_Locale":     "vn",
		"vnp_ReturnUrl":  vnp_ReturnUrl,
		"vnp_IpAddr":     ipAddr,
		"vnp_CreateDate": time.Now().Format("20060102150405"),
	}

	// Sắp xếp keys theo alphabet
	var keys []string
	for k := range vnp_Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Tạo chuỗi query và chuỗi hash
	var signData string
	var queryString string
	for i, k := range keys {
		v := vnp_Params[k]
		if v != "" {
			// Encode URL (Go sẽ biến space thành dấu '+')
			escapedVal := url.QueryEscape(v)
            
			// SỬ DỤNG STRINGS Ở ĐÂY: Đổi dấu '+' thành '%20' theo đúng chuẩn VNPay
			escapedVal = strings.ReplaceAll(escapedVal, "+", "%20")

			if i > 0 {
				signData += "&"
				queryString += "&"
			}
            
			// Lưu ý: VNPay v2.1.0 yêu cầu cả chuỗi tạo chữ ký (signData) cũng phải chứa dữ liệu đã được Encode
			signData += k + "=" + escapedVal
			queryString += k + "=" + escapedVal
		}
	}

	// Tạo chữ ký (Secure Hash) bằng HMAC-SHA512
	h := hmac.New(sha512.New, []byte(vnp_HashSecret))
	h.Write([]byte(signData))
	vnp_SecureHash := hex.EncodeToString(h.Sum(nil))

	// Trả về URL hoàn chỉnh
	return fmt.Sprintf("%s?%s&vnp_SecureHash=%s", vnp_Url, queryString, vnp_SecureHash), nil
}

func VerifyVNPaySignature(params map[string]string, secureHash string) bool {
	secret := os.Getenv("VNPAY_HASH_SECRET")
	if secret == "" || secureHash == "" {
		return false
	}

	var keys []string
	for k := range params {
		if k == "vnp_SecureHash" || k == "vnp_SecureHashType" {
			continue
		}
		if params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signData string
	for i, k := range keys {
		escapedVal := url.QueryEscape(params[k])
		escapedVal = strings.ReplaceAll(escapedVal, "+", "%20")
		if i > 0 {
			signData += "&"
		}
		signData += k + "=" + escapedVal
	}

	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(signData))
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(secureHash))
}