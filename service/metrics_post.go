package service

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-team/bitrise-step-analytics/metrics"
	"github.com/bitrise-team/bitrise-step-analytics/models"
	"github.com/pkg/errors"
)

// MetricsPostHandler ...
func MetricsPostHandler(w http.ResponseWriter, r *http.Request) error {
	var buildAnalytics models.BuildAnalytics
	defer httpresponse.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&buildAnalytics); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	if reflect.DeepEqual(buildAnalytics, models.BuildAnalytics{}) {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide metrics data")
	}
	dogstatsd, err := GetDogStatsDMetricsFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	dogstatsd.Track(buildAnalytics, metrics.DogStatsDBuildCounterMetricName)
	for _, aStepAnalytic := range buildAnalytics.StepAnalytics {
		dogstatsd.Track(aStepAnalytic, metrics.DogStatsDStepCounterMetricName)
	}

	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "ok"})
}
