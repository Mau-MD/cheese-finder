package db

import "log/slog"

func GetObjectsToRate(limit int) []RateObject {
	var objects []RateObject
	// Need to use RAW SQL to get random objects
	db.Raw("SELECT * FROM rate_objects WHERE liked IS NULL ORDER BY RANDOM() LIMIT ?", limit).Scan(&objects)
	return objects
}

func UpdateObjectWithRating(objectId int, liked bool) error {
	err := db.First(&RateObject{}, objectId).Update("liked", liked).Error
	if err != nil {
		slog.Error("Error updating object", slog.String("error", err.Error()))
		return err
	}
	return nil
}
