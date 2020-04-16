package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/push"
	"go.opencensus.io/stats/view"
	"log"
)

// PushGatewayExporter exporter for the Prometheus's push gateway
type PushGatewayExporter struct {
	service    string
	gatewayURL string
	prometheus *Exporter
	pusher     *push.Pusher
}

// NewPushGatewayExporter returns the new PushGatewayExporter instance
func NewPushGatewayExporter(service string, gatewayURL string) (*PushGatewayExporter, error) {
	prometheus, err := NewExporter(Options{
		Namespace: service,
	})
	if err != nil {
		return nil, err
	}

	pusher := push.New(gatewayURL, service).Gatherer(prometheus.g)
	return &PushGatewayExporter{
		service:    service,
		gatewayURL: gatewayURL,
		prometheus: prometheus,
		pusher:     pusher,
	}, nil
}

// ExportView implements the views interface
// Deprecated: don't need to do anything. prometheus uses the metricexport.Reader interface.
// which is implemented in: collector > metricExporter > go.opencensus.io/metric/metricexport.ReadAndExport
func (p *PushGatewayExporter) ExportView(viewData *view.Data) {
}

func (p *PushGatewayExporter) push() {
	err := p.pusher.Add()
	if err != nil {
		log.Println(p.service, "Could not push to the Pushgateway", p.gatewayURL, err)
	}
}
