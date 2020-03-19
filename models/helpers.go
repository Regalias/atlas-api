package models

// CheckLinkModelsAreEqual compares user controllable properties inside the model
func CheckLinkModelsAreEqual(lm1 *LinkModel, lm2 *LinkModel) bool {
	if (lm1.CanonicalName != lm2.CanonicalName) || (lm1.LinkPath != lm2.LinkPath) || (lm1.TargetURL != lm2.TargetURL || (lm1.Enabled != lm2.Enabled)) {
		return false
	}
	return true
}
