package crypto

const (
	// Security defines a package level security parameter used to generate
	// instance based parameters.
	Security = 80
)

// GetNISTRecommended returns a pair of corresponding key lengths to a
// given security parameter.
func GetNISTRecommended(secparam int) (modulusLength int, cipherLength int) {
	switch secparam {
	case 80:
		return 1024, 2048
	default:
		return 2048, 4096
	}
}
