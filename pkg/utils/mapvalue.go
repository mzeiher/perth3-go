package utils

func MapValue(value float32, fromMin float32, fromMax float32, toMin float32, toMax float32) float32 {
	// Calculate the range of values for both the source and target ranges
	fromRange := fromMax - fromMin
	toRange := toMax - toMin

	// Map the value from the source range to the target range
	mappedValue := (value-fromMin)*(toRange/fromRange) + toMin

	return mappedValue
}
