package config

import (
	"fmt"
)

// EnvType is a string type alias used for application environment names. The
// application environment indicates the level of application stability and
// the target audience for clients.
type EnvType string

// EnvDev indicates that the application code is unstable. And should only
// be accessed by the developers
const EnvDevelop EnvType = "develop"

// EnvTest indicates that the application code is ustable. And that tests are
// being run on it. Only the test cases should access the application.
const EnvTest EnvType = "test"

// EnvProd indicates that the application code is stable. And should be viewed
// by actual clients.
const EnvProd EnvType = "production"

// EnvErr indicates that an error occurred while parsing a raw value into a
// EnvType
const EnvErr EnvType = "err"

// NewEnvType creates a EnvType with the provided raw value. An error is
// returned if one occurs parsing this raw value into an EnvType. Nil on
// success.
func NewEnvType(raw string) (EnvType, error) {
	// Check each type
	if raw == string(EnvDevelop) {
		return EnvDevelop, nil
	} else if raw == string(EnvTest) {
		return EnvTest, nil
	} else if raw == string(EnvProd) {
		return EnvProd, nil
	} else {
		// If no match
		return EnvErr, fmt.Errorf("no EnvType with raw value: %s", raw)
	}
}
