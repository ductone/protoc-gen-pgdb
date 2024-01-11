package v1

import (
	"fmt"
	"strings"
)

// In order to insert a vector type it needs to be of the form '[1.0,2.0,3.0,...]'.
func FloatArrayToVectorString(in []float32) string {
	return strings.Join(strings.Fields(fmt.Sprint(in)), ",")
}
