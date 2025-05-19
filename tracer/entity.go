package tracer

type TracerOption struct {
	CollectorURL string
	InsecureMode bool
	ServiceName  string
	Environment  string
	SignozToken  string
}
