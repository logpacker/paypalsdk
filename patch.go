package paypal

type Patch struct {
	Operation string                 `json:"op"`
	Path      string                 `json:"path"`
	Values    map[string]interface{} `json:"value"`
}
