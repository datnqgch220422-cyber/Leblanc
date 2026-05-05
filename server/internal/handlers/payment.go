package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"
	"leblanc/server/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func VnpayIPN(c *gin.Context) {
    params := make(map[string]string)
    _ = c.Request.ParseForm()
    for key, values := range c.Request.Form {
        if key == "vnp_SecureHash" || key == "vnp_SecureHashType" {
            continue
        }
        if len(values) == 0 {
            continue
        }
        if value := strings.TrimSpace(values[0]); value != "" {
            params[key] = value
        }
    }

    secureHash := strings.TrimSpace(c.Query("vnp_SecureHash"))
    if secureHash == "" {
        secureHash = strings.TrimSpace(c.Request.FormValue("vnp_SecureHash"))
    }

    if !services.VerifyVNPaySignature(params, secureHash) {
        c.JSON(http.StatusBadRequest, gin.H{"RspCode": "97", "Message": "Invalid signature"})
        return
    }

    orderID := strings.TrimSpace(c.Query("vnp_TxnRef"))
    responseCode := strings.TrimSpace(c.Query("vnp_ResponseCode"))
    transactionNo := strings.TrimSpace(c.Query("vnp_TransactionNo"))
    message := strings.TrimSpace(c.Query("vnp_OrderInfo"))

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var booking models.Booking
    if err := db.DB.Collection("bookings").FindOne(ctx, bson.M{"paymentOrderId": orderID}).Decode(&booking); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"RspCode": "01", "Message": "Order not found"})
        return
    }

    setFields := bson.M{
        "paymentMessage": message,
        "paymentTransId": transactionNo,
    }

    if responseCode == "00" {
        setFields["paymentStatus"] = models.PaymentStatusPaid
        setFields["status"] = models.BookingStatusConfirmed

        for _, item := range booking.Items {
            _, _ = db.DB.Collection("drinks").UpdateOne(
                ctx,
                bson.M{"_id": item.DrinkID},
                bson.M{"$inc": bson.M{"stock": -item.Qty}},
            )
        }
        _, _ = db.DB.Collection("carts").DeleteOne(ctx, bson.M{"email": booking.Email})
    } else {
        setFields["paymentStatus"] = models.PaymentStatusFailed
    }

    _, _ = db.DB.Collection("bookings").UpdateOne(
        ctx,
        bson.M{"_id": booking.ID},
        bson.M{"$set": setFields},
    )

    if responseCode == "00" {
        c.JSON(http.StatusOK, gin.H{"RspCode": "00", "Message": "Confirm Success"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"RspCode": "99", "Message": "Payment Failed"})
}