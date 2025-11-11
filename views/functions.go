package views

import "avileads-web/http/routers"

func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func HasAccess(route string, currentRules []string) bool {
	return routers.HasAccessToPath(route, currentRules)
}
