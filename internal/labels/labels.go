package labels

import "strings"

type Labels struct {
	LinkTo []string `json:"goisolator.linkto"`
}

func MapToLabels(labels map[string]string) Labels {
	l := Labels{}
	if v, ok := labels["goisolator.linkto"]; ok {
		l.LinkTo = strings.Split(v, ",")
	}
	return l
}
