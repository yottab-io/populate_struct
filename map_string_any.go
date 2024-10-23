package populate_struct

// Gol sets VAL in DATA at the path specified by PATH
func AddToMapStringAny(data map[string]any, path []string, val any) {
	lastIndex := len(path) - 1

	// Traverse the map and create intermediate maps if needed
	current := data
	for _, key := range path[:lastIndex] {
		// If key doesn't exist or isn't a map, create a new map
		if next, ok := current[key].(map[string]any); ok {
			current = next
		} else {
			newMap := make(map[string]any)
			current[key] = newMap
			current = newMap
		}
	}

	// Set the value at the final key in the path
	finalKey := path[lastIndex]
	current[finalKey] = val
}

func GetFromMapStringAny(m map[string]any, path []string) (any, error) {
	var current any = m

	for _, field := range path {
		switch v := current.(type) {
		case map[string]any:
			var ok bool
			current, ok = v[field]
			if !ok {
				return nil, FieldNotFound
			}
		default:
			return nil, InvalidPath
		}
	}

	return current, nil
}
