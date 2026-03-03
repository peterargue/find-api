package command

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// FormatResult formats a result for printing based on format and filter flags.
func FormatResult(result Result, filter, format string) (string, error) {
	if filter != "" {
		return filterField(result, filter)
	}
	switch strings.ToLower(format) {
	case "json":
		b, err := json.Marshal(result.JSON())
		if err != nil {
			return "", fmt.Errorf("JSON marshal: %w", err)
		}
		return string(b), nil
	case "inline":
		return result.Oneliner(), nil
	default:
		return result.String(), nil
	}
}

// filterField extracts a single field from the JSON representation.
func filterField(result Result, field string) (string, error) {
	m, ok := result.JSON().(map[string]any)
	if !ok {
		return "", fmt.Errorf("cannot filter: result is not a JSON object")
	}
	v := m[field]
	if v == nil {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return "", fmt.Errorf("field %q not found; available: %s", field, strings.Join(keys, ", "))
	}
	return fmt.Sprintf("%v", v), nil
}

// printResult writes the formatted result to stdout (or a file if --save is set).
func printResult(result string, format, filter, save string) {
	if save != "" {
		if err := os.WriteFile(save, []byte(result), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving file: %s\n", err)
		}
		fmt.Printf("Result saved to: %s\n", save)
		return
	}
	if format == "inline" || filter != "" {
		fmt.Print(result)
	} else {
		fmt.Printf("\n%s\n\n", result)
	}
}
