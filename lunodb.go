package lunodbgo

import (
	"context"

	"github.com/cloudproud/lunodb.api/proto/plan"
)

const DefaultStargateAddress = "stargate.lunodb.io"

// Handler defines the interface that must be implemented by a connector.
// It handles the core lifecycle methods for data introspection and execution.
type Handler interface {
	// Ping is called by the server to check if the connector is alive and healthy.
	// It should return an error if the connector is not ready or unhealthy.
	Ping(ctx context.Context) error

	// Fetch returns a slice of tables that the connector supports, including their schema definitions.
	// This is typically called during initialization or introspection.
	Fetch(ctx context.Context) (Tables, error)

	// Scan executes a literal query plan and writes the resulting rows using the provided Writer.
	// The implementation is responsible for pushing all matching rows to the writer.
	Scan(ctx context.Context, plan *plan.Literal, writer Writer) error
}

// Writer defines the interface for streaming or collecting rows during a scan operation.
// It is typically called once per matching result row.
type Writer interface {
	// Write sends a single row of values, typically from a Scan implementation.
	// The values should match the output schema defined in the query plan.
	Write(ctx context.Context, values []any) error
}

// WriterFunc is a helper type to allow using ordinary functions as Writer implementations.
// It enables writing rows using a simple function.
type WriterFunc func(ctx context.Context, values []any) error

func (fn WriterFunc) Write(ctx context.Context, values []any) error {
	return fn(ctx, values)
}
