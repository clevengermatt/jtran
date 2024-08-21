package jtran

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// KeywordHandler defines the function signature for keyword handlers.
type KeywordHandler func(value interface{}, context map[string]interface{}, input string) (interface{}, error)

// Registered handlers map for user-registered keywords
var keywordHandlers = make(map[string]KeywordHandler)

// Stock handlers map for default keywords and their handlers
var stockKeywordHandlers = map[string]KeywordHandler{
	"capitalize": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("capitalize keyword expects a string value")
		}

		// Parse input as a comma-separated list of indices
		start, end, err := parseRange(input, len(strVal))
		if err != nil {
			return nil, fmt.Errorf("capitalize keyword: %v", err)
		}

		// Capitalize the first character of each word in the specified range
		runes := []rune(strVal)
		capitalizeNext := true
		for i := start; i < end; i++ {
			if capitalizeNext && unicode.IsLetter(runes[i]) {
				runes[i] = unicode.ToUpper(runes[i])
				capitalizeNext = false
			} else if unicode.IsSpace(runes[i]) {
				capitalizeNext = true
			}
		}

		return string(runes), nil
	},
	"foreach": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		arrayVal, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("foreach keyword expects an array of values")
		}

		var results []string
		for _, item := range arrayVal {
			fieldValue := ResolveField(input, item.(map[string]interface{}))
			if strVal, ok := fieldValue.(string); ok {
				results = append(results, strVal)
			} else {
				return nil, fmt.Errorf("foreach keyword: expected string subfield but got %T", fieldValue)
			}
		}

		return results, nil
	},
	"join": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		sliceVal, ok := value.([]string)
		if !ok {
			return nil, fmt.Errorf("join keyword expects a slice of strings")
		}

		if input == "" {
			return strings.Join(sliceVal, ""), nil
		}

		return strings.Join(sliceVal, input), nil
	},
	"lowercase": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("lowercase keyword expects a string value")
		}

		// Parse input as a comma-separated list of indices
		start, end, err := parseRange(input, len(strVal))
		if err != nil {
			return nil, fmt.Errorf("lowercase keyword: %v", err)
		}

		// Apply lowercase transformation to the specified range
		runes := []rune(strVal)
		for i := start; i < end; i++ {
			runes[i] = rune(strings.ToLower(string(runes[i]))[0])
		}

		return string(runes), nil
	},
	"padleft": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("padleft keyword expects a string value")
		}

		parts := strings.Split(input, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("padleft keyword: invalid input format, expected 'char,length'")
		}

		padChar := parts[0]
		length, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("padleft keyword: invalid length")
		}

		if len(strVal) >= length {
			return strVal, nil
		}

		padding := strings.Repeat(padChar, length-len(strVal))
		return padding + strVal, nil
	},
	"padright": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("padright keyword expects a string value")
		}

		parts := strings.Split(input, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("padright keyword: invalid input format, expected 'char,length'")
		}

		padChar := parts[0]
		length, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("padright keyword: invalid length")
		}

		if len(strVal) >= length {
			return strVal, nil
		}

		padding := strings.Repeat(padChar, length-len(strVal))
		return strVal + padding, nil
	},
	"redact": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("redact keyword expects a string value")
		}

		// Parse input as a comma-separated list of indices
		start, end, err := parseRange(input, len(strVal))
		if err != nil {
			return nil, fmt.Errorf("redact keyword: %v", err)
		}

		// Redact characters in the specified range
		runes := []rune(strVal)
		for i := start; i < end; i++ {
			runes[i] = '*'
		}

		return string(runes), nil
	},
	"replace": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("replace keyword expects a string value")
		}

		// Input format: "old,new"
		parts := strings.SplitN(input, ",", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("replace keyword: invalid input format, expected 'old,new'")
		}

		oldStr := parts[0]
		newStr := parts[1]

		// Replace occurrences of oldStr with newStr in the specified range
		return strings.Replace(strVal, oldStr, newStr, -1), nil
	},
	"reverse": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("reverse keyword expects a string value")
		}

		// Parse input as a comma-separated list of indices
		start, end, err := parseRange(input, len(strVal))
		if err != nil {
			return nil, fmt.Errorf("reverse keyword: %v", err)
		}

		// Reverse characters in the specified range
		runes := []rune(strVal)
		for i, j := start, end-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}

		return string(runes), nil
	},
	"snakecase": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("snakecase keyword expects a string value")
		}

		var result []rune
		for i, r := range strVal {
			if unicode.IsUpper(r) && i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		}

		return string(result), nil
	},
	"split": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("split keyword expects a string value")
		}

		if input == "" {
			return []string{strVal}, nil
		}

		return strings.Split(strVal, input), nil
	},
	"substring": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("substring keyword expects a string value")
		}

		start, end, err := parseRange(input, len(strVal))
		if err != nil {
			return nil, fmt.Errorf("substring keyword: %v", err)
		}

		return strVal[start:end], nil
	},
	"title": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("title keyword expects a string value")
		}

		// Parse input as a comma-separated list of indices
		start, end, err := parseRange(input, len(strVal))
		if err != nil {
			return nil, fmt.Errorf("title keyword: %v", err)
		}

		// Apply lowercase transformation to the specified range
		runes := []rune(strVal)
		for i := start; i < end; i++ {
			runes[i] = rune(strings.ToTitle(string(runes[i]))[0])
		}

		return string(runes), nil
	},
	"trim": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("trim keyword expects a string value")
		}

		// Parse input as a comma-separated list of indices
		start, end, err := parseRange(input, len(strVal))
		if err != nil {
			return nil, fmt.Errorf("trim keyword: %v", err)
		}

		// Trim whitespace in the specified range
		trimmed := strings.TrimSpace(strVal[start:end])
		return strVal[:start] + trimmed + strVal[end:], nil
	},
	"truncate": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("truncate keyword expects a string value")
		}

		length, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("truncate keyword: invalid length")
		}

		if length < 0 || length > len(strVal) {
			return nil, fmt.Errorf("truncate keyword: length out of bounds")
		}

		return strVal[:length], nil
	},
	"uppercase": func(value interface{}, context map[string]interface{}, input string) (interface{}, error) {
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("uppercase keyword expects a string value")
		}

		// Parse input as a comma-separated list of indices
		start, end, err := parseRange(input, len(strVal))
		if err != nil {
			return nil, fmt.Errorf("uppercase keyword: %v", err)
		}

		// Apply uppercase transformation to the specified range
		runes := []rune(strVal)
		for i := start; i < end; i++ {
			runes[i] = rune(strings.ToUpper(string(runes[i]))[0])
		}

		return string(runes), nil
	},
}

// parseRange parses the input string as a comma-separated range of integers and returns start and end indices.
func parseRange(input string, maxLen int) (int, int, error) {
	if input == "" {
		return 0, maxLen, nil // Default to the entire string if no range is specified
	}

	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid input range, expected two comma-separated values")
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start index: %v", err)
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end index: %v", err)
	}

	if start < 0 {
		start = 0
	}

	if end > maxLen {
		end = maxLen
	}

	return start, end, nil
}

// Regular expression to detect and process templated strings with flexible keywords
var templateRegex = regexp.MustCompile(`\${([^|}]+)(\|[^}]*)?}`)

// RegisterKeyword allows developers to register keyword handlers, overriding stock handlers if necessary.
func RegisterKeyword(keyword string, handler KeywordHandler) {
	keywordHandlers[keyword] = handler
}

func TransformData(schema map[string]interface{}, data map[string]interface{}) (map[string]interface{}, error) {
	transformed := make(map[string]interface{})
	context := map[string]interface{}{}

	for key, value := range schema {
		transformedKey, err := applyKeywordsToString(key, data, context)
		if err != nil {
			return nil, fmt.Errorf("failed to apply keywords to key '%s': %v", key, err)
		}

		strKey, ok := transformedKey.(string)
		if !ok {
			return nil, fmt.Errorf("key '%v' did not resolve to a string", transformedKey)
		}

		switch v := value.(type) {
		case string:
			transformedValue, err := applyKeywordsToString(v, data, context)
			if err != nil {
				return nil, fmt.Errorf("failed to apply keywords to value of key '%s': %v", strKey, err)
			}
			transformed[strKey] = transformedValue
		case map[string]interface{}:
			nestedTransformed, err := TransformData(v, data)
			if err != nil {
				return nil, fmt.Errorf("failed to transform nested schema for key '%s': %v", strKey, err)
			}
			transformed[strKey] = nestedTransformed
		case []interface{}:
			transformedArray := []interface{}{}
			for _, item := range v {
				switch itemVal := item.(type) {
				case string:
					transformedItem, err := applyKeywordsToString(itemVal, data, context)
					if err != nil {
						return nil, fmt.Errorf("failed to apply keywords to array item under key '%s': %v", strKey, err)
					}
					transformedArray = append(transformedArray, transformedItem)
				case map[string]interface{}:
					nestedTransformed, err := TransformData(itemVal, data)
					if err != nil {
						return nil, fmt.Errorf("failed to transform nested array schema for key '%s': %v", strKey, err)
					}
					transformedArray = append(transformedArray, nestedTransformed)
				default:
					transformedArray = append(transformedArray, itemVal)
				}
			}
			transformed[strKey] = transformedArray
		default:
			transformed[strKey] = v
		}
	}
	return transformed, nil
}

func applyKeywordsToString(value string, data map[string]interface{}, context map[string]interface{}) (interface{}, error) {
	if templateRegex.MatchString(value) {
		resolvedValue := templateRegex.ReplaceAllStringFunc(value, func(match string) string {
			matches := templateRegex.FindStringSubmatch(match)
			fieldName := matches[1]
			keywords := matches[2]

			context["currentKey"] = fieldName

			fieldValue := ResolveField(fieldName, data)

			if keywords != "" {
				keywordList := strings.Split(keywords, "|")[1:]
				for _, keyword := range keywordList {
					keyword = strings.TrimSpace(keyword)

					var keywordName string
					input := ""
					if strings.Contains(keyword, "(") && strings.HasSuffix(keyword, ")") {
						keywordName = keyword[:strings.Index(keyword, "(")]
						input = keyword[strings.Index(keyword, "(")+1 : len(keyword)-1]
					} else {
						keywordName = keyword
					}

					handler, exists := keywordHandlers[keywordName]
					if !exists {
						handler, exists = stockKeywordHandlers[keywordName]
					}

					if exists {
						var err error
						fieldValue, err = handler(fieldValue, context, input)
						if err != nil {
							return fmt.Sprintf("Error: %v", err)
						}
					}
				}
			}

			if fieldValue == nil {
				return ""
			}

			return fmt.Sprintf("%v", fieldValue)
		})
		return resolvedValue, nil
	}

	if strings.Contains(value, "|") {
		parts := strings.SplitN(value, "|", 2)
		baseValue := parts[0]
		keywords := parts[1]

		var resolvedValue interface{}
		resolvedValue = baseValue

		if keywords != "" {
			keywordList := strings.Split(keywords, "|")
			for _, keyword := range keywordList {
				keyword = strings.TrimSpace(keyword)

				var keywordName string
				input := ""
				if strings.Contains(keyword, "(") && strings.HasSuffix(keyword, ")") {
					keywordName = keyword[:strings.Index(keyword, "(")]
					input = keyword[strings.Index(keyword, "(")+1 : len(keyword)-1]
				} else {
					keywordName = keyword
				}

				context["currentKey"] = baseValue

				handler, exists := keywordHandlers[keywordName]
				if !exists {
					handler, exists = stockKeywordHandlers[keywordName]
				}

				if exists {
					var err error
					resolvedValue, err = handler(resolvedValue, context, input)
					if err != nil {
						return nil, fmt.Errorf("error applying keyword '%s': %v", keywordName, err)
					}
				}
			}
		}

		return resolvedValue, nil
	}

	return value, nil
}

// ResolveField retrieves nested fields from the original data.
func ResolveField(fieldName string, data map[string]interface{}) interface{} {
	if strings.Contains(fieldName, "->") {
		keys := strings.Split(fieldName, "->")
		return resolveRecursive(keys, data)
	}
	return data[fieldName]
}

func resolveRecursive(keys []string, data interface{}) interface{} {
	if len(keys) == 0 {
		return data
	}
	currentKey := keys[0]
	remainingKeys := keys[1:]

	// Handle array indexing
	if strings.Contains(currentKey, "[") && strings.HasSuffix(currentKey, "]") {
		bracketIndex := strings.Index(currentKey, "[")
		arrayKey := currentKey[:bracketIndex]
		indexPart := currentKey[bracketIndex+1 : len(currentKey)-1]

		// Resolve the index or condition
		if nestedMap, ok := data.(map[string]interface{}); ok {
			if nestedData, ok := nestedMap[arrayKey]; ok {
				switch nestedArray := nestedData.(type) {
				case []interface{}:
					// Try to parse as an integer index
					if index, err := strconv.Atoi(indexPart); err == nil {
						if index >= 0 && index < len(nestedArray) {
							return resolveRecursive(remainingKeys, nestedArray[index])
						}
					} else {
						// Handle conditional selection
						conditionParts := strings.SplitN(indexPart, "=", 2)
						if len(conditionParts) == 2 {
							conditionField := conditionParts[0]
							conditionValue := conditionParts[1]

							for _, item := range nestedArray {
								if itemMap, ok := item.(map[string]interface{}); ok {
									if itemValue, exists := itemMap[conditionField]; exists {
										// Convert both to strings for type-agnostic comparison
										itemValueStr := fmt.Sprintf("%v", itemValue)
										conditionValueStr := fmt.Sprintf("%v", conditionValue)

										if itemValueStr == conditionValueStr {
											return resolveRecursive(remainingKeys, itemMap)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		// Handle regular map keys
		if nestedMap, ok := data.(map[string]interface{}); ok {
			if nestedData, ok := nestedMap[currentKey]; ok {
				return resolveRecursive(remainingKeys, nestedData)
			}
		}
	}

	return nil
}
