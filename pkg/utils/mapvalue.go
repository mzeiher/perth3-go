package utils

func MapValue(value float64, fromMin float64, fromMax float64, toMin float64, toMax float64) float64 {
	// Calculate the range of values for both the source and target ranges
	fromRange := fromMax - fromMin
	toRange := toMax - toMin

	// Map the value from the source range to the target range
	mappedValue := (value-fromMin)*(toRange/fromRange) + toMin

	return mappedValue
}
