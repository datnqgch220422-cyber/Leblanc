package models

const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCompleted = "completed"
	BookingStatusCancelled = "cancelled"
)

func NormalizeBookingStatus(status string) string {
	switch status {
	case BookingStatusConfirmed, BookingStatusCompleted, BookingStatusCancelled:
		return status
	default:
		return BookingStatusPending
	}
}
