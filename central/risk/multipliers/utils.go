package multipliers

// NormalizeScore normalizes score into a value between 1 and maxAllowedValue.
func NormalizeScore(score, saturation, maxAllowedValue float32) float32 {
	if score > saturation {
		return float32(maxAllowedValue)
	}
	return 1 + (score/saturation)*(maxAllowedValue-1)
}
