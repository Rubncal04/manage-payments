package scheduler

import (
	"log"
	"strconv"
	"time"

	"github/Rubncal04/youtube-premium/db"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdatePaymentStatus updates the payment status for users based on their payment date
func UpdatePaymentStatus(mongoRepo *db.MongoRepo) error {
	currentDay := time.Now().Day()
	var startDay, endDay int

	// Determine the date range based on current day
	if currentDay == 13 {
		startDay = 15
		endDay = 20
	} else if currentDay == 25 {
		startDay = 28
		endDay = 30
	} else {
		return nil // Not a scheduled update day
	}

	log.Printf("Updating payment status for users with payment dates between %d and %d", startDay, endDay)

	// Build the date range filter
	var dateFilters []bson.M
	for day := startDay; day <= endDay; day++ {
		dateFilters = append(dateFilters, bson.M{"date_to_pay": strconv.Itoa(day)})
	}

	// Create the filter to find users in the date range
	filter := bson.M{
		"$or": dateFilters,
	}

	// Create the update to set paid=false and status=inactive
	update := bson.M{
		"$set": bson.M{
			"paid":   false,
			"status": "inactive",
		},
	}

	// Update all matching users
	result, err := mongoRepo.UpdateMany("users", filter, update)
	if err != nil {
		return err
	}

	log.Printf("Updated %d users' payment status", result.ModifiedCount)
	return nil
}
