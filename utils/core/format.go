package utils

import "strings"

// Format — {key} şeklindeki placeholderları verilen map ile replace eder
//
// Örnek:
//   Format("Merhaba {user}, {guild} sunucusundasın!", map[string]string{
//       "user":  "VoidBot",
//       "guild": "Void",
//   })
//   → "Merhaba VoidBot, Void sunucusundasın!"
func FormatKeys(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}
