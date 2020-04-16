package prometheus

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func getGatewayURL() string {
	url := os.Getenv("PUSHGATEWAY_URL")
	if len(url) == 0 {
		url = "http://localhost:9091"
	}
	return url
}

func TestPushGatewayExporter_Push(t *testing.T) {
	// no parallel

	// create exporter
	serviceName := "push_gateway_test"
	exporter, err := NewPushGatewayExporter(serviceName, getGatewayURL())
	require.Nil(t, err)

	// create view
	m := stats.Int64("tests/foo", "foo", stats.UnitDimensionless)
	v := &view.View{
		Name:        m.Name(),
		Description: m.Description(),
		Measure:     m,
		Aggregation: view.Count(),
	}
	err = view.Register(v)
	require.Nil(t, err)
	defer view.Unregister(v)

	// record and push some data
	stats.Record(context.Background(), m.M(1))
	stats.Record(context.Background(), m.M(1))
	time.Sleep(30 * time.Millisecond)
	exporter.push()

	expectedMetric := `
# HELP push_gateway_test_tests_foo foo
# TYPE push_gateway_test_tests_foo counter
push_gateway_test_tests_foo{instance="",job="push_gateway_test"} 2
`
	verifyData(t, expectedMetric)

	// record and push some data
	stats.Record(context.Background(), m.M(1))
	stats.Record(context.Background(), m.M(1))
	time.Sleep(30 * time.Millisecond)
	exporter.push()

	expectedMetric = `
# HELP push_gateway_test_tests_foo foo
# TYPE push_gateway_test_tests_foo counter
push_gateway_test_tests_foo{instance="",job="push_gateway_test"} 4
`
	verifyData(t, expectedMetric)
}

func verifyData(t *testing.T, expected string) {
	resp, err := http.Get(fmt.Sprintf("%s/metrics", getGatewayURL()))
	require.Nil(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	err = resp.Body.Close()
	require.Nil(t, err)

	output := string(body)

	// simple check some conditions
	if strings.Contains(output, "collected before with the same name and label values") {
		t.Fatal("metric name and labels being duplicated but must be unique")
	}

	if strings.Contains(output, "error(s) occurred") {
		t.Fatal("error reported by prometheus registry")
	}

	if !strings.Contains(output, expected) {
		t.Errorf("output does not contain correct opencensus counter. Output: %s want: %s", output, expected)
	}
}
