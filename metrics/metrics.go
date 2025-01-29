package metrics

import (
	"context"
	"sync/atomic"

	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	defaultBytesDistribution        = view.Distribution(1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296)
	defaultMillisecondsDistribution = view.Distribution(0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000)
)

// Keys
var (
	KeyMessageType = "message_type"
	KeyPeerID, _   = tag.NewKey("peer_id")
	// KeyInstanceID identifies a dht instance by the pointer address.
	// Useful for differentiating between different dhts that have the same peer id.
	KeyInstanceID, _ = tag.NewKey("instance_id")
)

// UpsertMessageType is a convenience upserts the message type
// of a pb.Message into the KeyMessageType.
func UpsertMessageType(m *pb.Message) metric.MeasurementOption {
	return metric.WithAttributes(attribute.String(KeyMessageType, m.Type.String()))
}

// Measures
var (
	meter = otel.Meter("libp2p.io/dht/kad")

	// Define OpenTelemetry counters
	ReceivedMessages, _ = meter.Int64Counter(
		"libp2p.io/dht/kad/received_messages",
		metric.WithDescription("Total number of messages received per RPC"),
	)

	ReceivedMessageErrors, _ = meter.Int64Counter(
		"libp2p.io/dht/kad/received_message_errors",
		metric.WithDescription("Total number of errors for messages received per RPC"),
	)

	ReceivedBytes, _ = meter.Int64Counter(
		"libp2p.io/dht/kad/received_bytes",
		metric.WithDescription("Total received bytes per RPC"),
		metric.WithUnit("By"),
	)

	InboundRequestLatency, _ = meter.Float64Histogram(
		"libp2p.io/dht/kad/inbound_request_latency",
		metric.WithDescription("Latency per RPC"),
		metric.WithUnit("ms"),
	)

	OutboundRequestLatency, _ = meter.Float64Histogram(
		"libp2p.io/dht/kad/outbound_request_latency",
		metric.WithDescription("Latency per RPC"),
		metric.WithUnit("ms"),
	)

	SentMessages, _ = meter.Int64Counter(
		"libp2p.io/dht/kad/sent_messages",
		metric.WithDescription("Total number of messages sent per RPC"),
	)

	SentMessageErrors, _ = meter.Int64Counter(
		"libp2p.io/dht/kad/sent_message_errors",
		metric.WithDescription("Total number of errors for messages sent per RPC"),
	)

	SentRequests, _ = meter.Int64Counter(
		"libp2p.io/dht/kad/sent_requests",
		metric.WithDescription("Total number of requests sent per RPC"),
	)

	SentRequestErrors, _ = meter.Int64Counter(
		"libp2p.io/dht/kad/sent_request_errors",
		metric.WithDescription("Total number of errors for requests sent per RPC"),
	)

	SentBytes, _ = meter.Int64Counter(
		"libp2p.io/dht/kad/sent_bytes",
		metric.WithDescription("Total sent bytes per RPC"),
		metric.WithUnit("By"),
	)

	networkSize int64
)

func init() {
	// Register an observable gauge
	meter.Int64ObservableGauge(
		"libp2p.io/dht/kad/network_size",
		metric.WithDescription("Network size estimation"),
		metric.WithInt64Callback(networkSizeCallback),
	)
}

func networkSizeCallback(ctx context.Context, observer metric.Int64Observer) error {
	observer.Observe(atomic.LoadInt64(&networkSize), metric.WithAttributes(attribute.String("instance", "default")))
	return nil
}

// Function to update the network size estimation
func SetNetworkSize(size int64) {
	atomic.StoreInt64(&networkSize, size)
}
