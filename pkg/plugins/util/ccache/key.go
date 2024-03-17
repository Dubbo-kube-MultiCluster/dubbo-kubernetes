package ccache

import "strings"

func GenerateDCacheKey(name string, mesh string) string {
	if mesh == "" {
		return name
	}
	return name + "/" + mesh
}

func DeDCacheKey(key string) (name string, mesh string) {
	parts := strings.Split(key, "/")
	if len(parts) == 1 {
		return parts[0], ""
	} else if len(parts) == 2 {
		return parts[0], parts[1]
	} else {
		return "", ""
	}
}
