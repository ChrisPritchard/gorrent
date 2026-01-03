package bencode

import "fmt"

func Get[T any](m map[string]any, key string) (T, error) {
	var nilT T
	val, exists := m[key]
	if !exists {
		return nilT, fmt.Errorf("key %s was not in map", key)
	}
	res, ok := val.(T)
	if !ok {
		return nilT, fmt.Errorf("key %s's value was an invalid type: %v", key, val)
	}
	return res, nil
}

func GetStrings(m map[string]any, key string) ([]string, error) {
	list, err := Get[[]any](m, key)
	if err != nil {
		return nil, err
	}
	results := []string{}
	for _, v := range list {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("a non-string value was in the list: %v", v)
		}
		results = append(results, string(s))
	}
	return results, nil
}
