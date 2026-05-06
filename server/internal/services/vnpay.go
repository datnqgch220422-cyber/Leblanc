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

func buildVNPaySignData(params map[string]string) string {
    keys := make([]string, 0, len(params))
    for k, v := range params {
        if strings.TrimSpace(v) == "" {
            continue
        }
        keys = append(keys, k)
    }
    sort.Strings(keys)

    parts := make([]string, 0, len(keys))
    for _, k := range keys {
        // Giữ encode mặc định (space -> +) để khớp sandbox VNPay phổ biến
        parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
    }
    return strings.Join(parts, "&")
}

// Tạo URL thanh toán VNPay
func CreateVNPayURL(orderId string, amount int64, orderInfo string, ipAddr string) (string, error) {
    vnpTmnCode := strings.TrimSpace(os.Getenv("VNPAY_TMN_CODE"))
    vnpHashSecret := strings.TrimSpace(os.Getenv("VNPAY_HASH_SECRET"))
    vnpURL := strings.TrimSpace(os.Getenv("VNPAY_URL"))
    vnpReturnURL := strings.TrimSpace(os.Getenv("VNPAY_RETURN_URL"))

    if vnpTmnCode == "" || vnpHashSecret == "" || vnpURL == "" || vnpReturnURL == "" {
        return "", fmt.Errorf("vnpay credentials are missing")
    }

    if strings.TrimSpace(ipAddr) == "" {
        ipAddr = "127.0.0.1"
    }

    params := map[string]string{
        "vnp_Version":    "2.1.0",
        "vnp_Command":    "pay",
        "vnp_TmnCode":    vnpTmnCode,
        "vnp_Amount":     fmt.Sprintf("%d", amount*100),
        "vnp_CurrCode":   "VND",
        "vnp_TxnRef":     orderId,
        "vnp_OrderInfo":  orderInfo,
        "vnp_OrderType":  "other",
        "vnp_Locale":     "vn",
        "vnp_ReturnUrl":  vnpReturnURL,
        "vnp_IpAddr":     ipAddr,
        "vnp_CreateDate": time.Now().Format("20060102150405"),
    }

    signData := buildVNPaySignData(params)

    h := hmac.New(sha512.New, []byte(vnpHashSecret))
    h.Write([]byte(signData))
    vnpSecureHash := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

    // Thêm HashType để đồng bộ cách verify của phía cổng thanh toán
    fullQuery := signData + "&vnp_SecureHashType=HMACSHA512&vnp_SecureHash=" + url.QueryEscape(vnpSecureHash)

    return vnpURL + "?" + fullQuery, nil
}

func VerifyVNPaySignature(params map[string]string, secureHash string) bool {
    secret := strings.TrimSpace(os.Getenv("VNPAY_HASH_SECRET"))
    secureHash = strings.TrimSpace(secureHash)

    if secret == "" || secureHash == "" {
        return false
    }

    filtered := make(map[string]string, len(params))
    for k, v := range params {
        if k == "vnp_SecureHash" || k == "vnp_SecureHashType" {
            continue
        }
        if strings.TrimSpace(v) == "" {
            continue
        }
        filtered[k] = v
    }

    signData := buildVNPaySignData(filtered)

    h := hmac.New(sha512.New, []byte(secret))
    h.Write([]byte(signData))
    expected := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

    return hmac.Equal([]byte(expected), []byte(strings.ToUpper(secureHash)))
}