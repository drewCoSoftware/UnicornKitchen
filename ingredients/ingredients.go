package ingredients

func HasIngredient(name string) bool {
	allIngredients := []string{"fish", "potato", "rice"}

	return contains(name, allIngredients)

}

func contains(value string, allValues []string) bool {
	for _, a := range allValues {
		if a == value {
			return true
		}
	}
	return false
}
