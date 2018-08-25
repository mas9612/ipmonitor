package ipmonitor

// ErrorResponse represents error message
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// HostsResponse represents response of /hosts endpoint
type HostsResponse struct {
	Count int    `json:"count"`
	Hosts []Host `json:"hosts"`
}
