package labels

import "strings"

type Labels struct {
	Enabled bool     `json:"goisolator"`
	Traefik bool     `json:"goisolator.traefik"`
	Ignore  bool     `json:"goisolator.ignore"`
	LinkTo  []string `json:"goisolator.linkto"`
}

func MapToLabels(labels map[string]string) Labels {
	l := Labels{}
	if v, ok := labels["goisolator.linkto"]; ok {
		l.LinkTo = strings.Split(v, ",")
	}
	if _, ok := labels["goisolator.ignore"]; ok {
		l.Ignore = true
	}
	if _, ok := labels["goisolator.traefik"]; ok {
		l.Traefik = true
	}
	if _, ok := labels["goisolator"]; ok {
		l.Enabled = true
	}
	return l
}
