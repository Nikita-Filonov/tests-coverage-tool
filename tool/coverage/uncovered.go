package coverage

func enrichWithUncoveredResultParameters(parameter *ResultParameters) bool {
	hasUncoveredChild := false
	for index := range parameter.Parameters {
		if childUncovered := enrichWithUncoveredResultParameters(&parameter.Parameters[index]); childUncovered {
			hasUncoveredChild = true
		}
	}

	if parameter.Covered && hasUncoveredChild {
		parameter.HasUncoveredParameters = true
	}

	return hasUncoveredChild || !parameter.Covered
}

func EnrichSliceWithUncoveredResultParameters(parameters []ResultParameters) {
	for index := range parameters {
		enrichWithUncoveredResultParameters(&parameters[index])
	}
}
