package rate

import (
	"cheesefinder/db"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"math/rand"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

// Objects to rate will return a list of objects that will be rated in the FE. This list will contain some level of randomness to avoid bias.
func GetObjectsToRate(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())

	limit := chi.URLParam(r, "limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	randomObjectCount := int(math.Ceil(float64(limitInt) * RATE_RANDOMNESS_FACTOR))
	oplog.Info(fmt.Sprintf("Random object count: %d/%d", randomObjectCount, limitInt))

	randomObjects := db.GetObjectsToRate(randomObjectCount)
	if len(randomObjects) < randomObjectCount {
		oplog.Warn(fmt.Sprintf("Not enough random objects were returned from Database. Expected: %d. Got: %d", randomObjectCount, len(randomObjects)))
	}

	predictedObjectCount := limitInt - randomObjectCount
	predictedObjects := make([]db.RateObject, 0, predictedObjectCount) // TODO: Call Prediction Service
	oplog.Info(fmt.Sprintf("Predicted object count: %d/%d", predictedObjectCount, limitInt))

	mixedObjects := mixAndShuffle(randomObjects, predictedObjects)

	b, err := json.Marshal(mixedObjects)
	if err != nil {
		http.Error(w, "Error marshalling objects", http.StatusInternalServerError)
		return
	}

	oplog.Info(fmt.Sprintf("Returning %d objects", len(mixedObjects)))
	w.Write(b)
}

type RateObjectBody struct {
	ObjectID int  `json:"object_id"`
	Liked    bool `json:"liked"`
}

func RateObject(w http.ResponseWriter, r *http.Request) {
	objectBody := RateObjectBody{}
	err := json.NewDecoder(r.Body).Decode(&objectBody)
	if err != nil {
		http.Error(w, "Error decoding body", http.StatusBadRequest)
		return
	}

	err = db.UpdateObjectWithRating(objectBody.ObjectID, objectBody.Liked)
	if err != nil {
		http.Error(w, "Error updating object", http.StatusInternalServerError)
		return
	}

	// TODO: Send signal to the prediction service to update the model

	w.Write([]byte("OK"))
}

func mixAndShuffle(randomObjects []db.RateObject, predictedObjects []db.RateObject) []db.RateObject {
	mixedLength := len(randomObjects) + len(predictedObjects)
	mixedObjects := make([]db.RateObject, 0, mixedLength)
	mixedObjects = append(mixedObjects, randomObjects...)
	mixedObjects = append(mixedObjects, predictedObjects...)

	// This we iterate over the mixedObjects and swap the elements randomly
	for i := range mixedObjects {
		j := rand.Intn(mixedLength)
		mixedObjects[i], mixedObjects[j] = mixedObjects[j], mixedObjects[i]
	}

	return mixedObjects
}
