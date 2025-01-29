package metrics

import (
	"context"
	"sync/atomic"

	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Keys
var (
	KeyMessageType = "message_type"
	KeyPeerID      = "peer_id"
	// KeyInstanceID identifies a dht instance by the pointer address.
	// Useful for differentiating between different dhts that have the same peer id.
	KeyInstanceID = "instance_id"
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
