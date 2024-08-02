package htmx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func IsHTMXRequest(w http.ResponseWriter, r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func TriggerEvent(w http.ResponseWriter, r *http.Request, events ...string) {
	w.Header().Set("HX-Trigger", strings.Join(events, ", "))
}

type TriggerPayload map[string]any

func Trigger(w http.ResponseWriter, r *http.Request, payload TriggerPayload) {
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode(payload)

	w.Header().Set("HX-Trigger", b.String())
}
