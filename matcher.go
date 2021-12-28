package pathfinder

func match(patternStr, str string) PathValues {

	pattern := tokenizer(patternStr)
	compiledPattern, err := Compile(pattern)
	if err != nil {
		return nil
	}

	found := compiledPattern.FindStringSubmatch(str)

	if len(found) <= len(pattern) {
		return nil
	}

	var patternValues PathValues

	for i, patternToken := range pattern {
		patternValues = append(patternValues, ValueToken{patternToken, found[i+1]})
	}

	return patternValues
}
