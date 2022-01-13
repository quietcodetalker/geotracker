package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/util/pagination"
	"regexp"
)

var usernameRegexp = regexp.MustCompile("^[a-zA-Z0-9]{4,16}$")

func ValidateUsername(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return usernameRegexp.Match([]byte(value))
}

func ValidateLongitude(fl validator.FieldLevel) bool {
	value := fl.Field().Float()
	return -180.0 <= value && value <= 180.0
}
func ValidateLatitude(fl validator.FieldLevel) bool {
	value := fl.Field().Float()
	return -90.0 <= value && value <= 90.0
}

func ValidateGeoPoint(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(domain.Point); ok {
		long := value.Longitude()
		lat := value.Latitude()
		if long < -180 || long > 180 {
			return false
		}
		if lat < -90 || lat > 90 {
			return false
		}
		return true
	}

	return false
}

func ValidatePageToken(fl validator.FieldLevel) bool {
	cursor := fl.Field().String()
	if cursor == "" {
		return true
	}
	pageToken, pageSize, err := pagination.DecodeCursor(cursor)
	if pageToken < 0 || pageSize < 0 || err != nil {
		return false
	}
	return true
}

func ValidateListMethod(sl validator.StructLevel) {
	switch v := sl.Current().Interface().(type) {
	case port.UserServiceListUsersInRadiusRequest:
		if (v.PageToken == "" && v.PageSize == 0) ||
			(v.PageToken != "" && v.PageSize != 0) {
			sl.ReportError(v.PageToken, "page_token", "PageToken", "pagesize_or_pagetoken", v.PageToken)
			sl.ReportError(v.PageSize, "page_size", "PageSize", "pagesize_or_pagetoken", fmt.Sprint(v.PageSize))
		}
	}
}
