package exchange

import (
	"errors"
	"github.com/prebid/prebid-server/openrtb_ext"
	"math"
	"strconv"
)

// DEFAULT_PRECISION should be taken care of in openrtb_ext/request.go, but throwing an additional safety check here.
const DEFAULT_PRECISION = 2

// GetCpmStringValue is the externally facing function for computing CPM buckets
func GetCpmStringValue(cpm float64, config openrtb_ext.PriceGranularity) (string, error) {
	cpmStr := ""
	bucketMax := 0.0
	increment := 0.0
	precision := config[0].Precision
	// If we wish to support precision "0", we will need to remove this check
	if precision == 0 {
		precision = DEFAULT_PRECISION
	}
	// calculate max of highest bucket
	for i := 0; i < len(config); i++ {
		if config[i].Max > bucketMax {
			bucketMax = config[i].Max
		}
		if config[i].Precision != precision {
			return "", errors.New("Precision changed within price granularity")
		}
	} // calculate which bucket cpm is in
	if cpm > bucketMax {
		// If we are over max, just return that
		return strconv.FormatFloat(bucketMax, 'f', precision, 64), nil
	}
	for i := 0; i < len(config); i++ {
		if cpm >= config[i].Min && cpm <= config[i].Max {
			increment = config[i].Increment
		}
	}
	if increment > 0 {
		cpmStr = getCpmTarget(cpm, increment, precision)
	}
	return cpmStr, nil
}

func getCpmTarget(cpm float64, increment float64, precision int) string {
	roundedCPM := math.Floor(cpm/increment) * increment
	return strconv.FormatFloat(roundedCPM, 'f', precision, 64)
}
