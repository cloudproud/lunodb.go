package lunodbgo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	lunopb "github.com/cloudproud/lunodb.api/proto"
	"github.com/cloudproud/lunodb.go/value"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// ConnectorOption defines a functional option for configuring a Connector. It
// allows for flexible and composable Connector construction via NewConnector.
type ConnectorOption func(*Connector) error

// WithLogger sets a custom zap.Logger for the Connector.
// If not specified, a no-op logger is used by default.
func WithLogger(logger *zap.Logger) ConnectorOption {
	return func(connector *Connector) error {
		connector.logger = logger
		return nil
	}
}

// WithStargateAddress sets the address of the Stargate gRPC endpoint.
func WithStargateAddress(address string) ConnectorOption {
	return func(connector *Connector) error {
		connector.StargateAddress = address
		return nil
	}
}

// WithInsecure configures the Connector to use insecure transport credentials (non-TLS).
func WithInsecure(insecure bool) ConnectorOption {
	return func(connector *Connector) error {
		connector.Insecure = insecure
		return nil
	}
}

// WithSource sets the source identifier used to register this Connector with Stargate.
func WithSource(source string) ConnectorOption {
	return func(connector *Connector) (err error) {
		connector.Source, err = strconv.ParseUint(source, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid source uid: %w", err)
		}

		return nil
	}
}

// WithToken sets the secret token used for authenticating this Connector with Stargate.
func WithToken(token string) ConnectorOption {
	return func(connector *Connector) error {
		connector.Token = token
		return nil
	}
}

// DefaultConnectorHeartbeat represents the default heartbeat interval in which a given
// source controller will ping the configured source.
const DefaultConnectorHeartbeat = 5 * time.Second

// NewConnector constructs a new Connector instance with optional configuration overrides.
func NewConnector(options ...ConnectorOption) (*Connector, error) {
	connector := Connector{
		logger:          zap.NewNop(),
		StargateAddress: os.Getenv("LUNODB_STARGATE_ADDRESS"),
		Insecure:        os.Getenv("LUNODB_INSECURE") == "true",
	}

	if connector.StargateAddress == "" {
		connector.StargateAddress = DefaultStargateAddress
	}

	for _, option := range options {
		err := option(&connector)
		if err != nil {
			return nil, err
		}
	}

	return &connector, nil
}

// Connector is a gRPC client that connects to the LunoDB Stargate server. It
// maintains connection settings such as address, TLS preferences, and source
// identity.
type Connector struct {
	logger          *zap.Logger
	StargateAddress string
	Insecure        bool
	Source          uint64
	Token           string
	mu              sync.Mutex
	healthy         bool
}

// Healthy returns true if the Connector is currently healthy and able to
// process Stargate requests.
func (connector *Connector) Healthy() bool {
	connector.mu.Lock()
	defer connector.mu.Unlock()

	return connector.healthy
}

// Serve establishes a gRPC connection to the configured Stargate server and
// starts the message receive loop using the provided handler. It blocks until
// the stream ends or an error occurs.
func (connector *Connector) Serve(ctx context.Context, handler Handler) error {
	options := []grpc.DialOption{}

	if connector.Insecure {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(connector.StargateAddress, options...)
	if err != nil {
		return err
	}

	client := lunopb.NewStargateClient(conn)
	return connector.serveLoop(ctx, client, handler)
}

func (connector *Connector) defaultOutgoingContext(ctx context.Context) context.Context {
	md := map[string][]string{
		"authorization": {"Bearer " + connector.Token},
		"source":        {fmt.Sprintf("%d", connector.Source)},
	}

	return metadata.NewOutgoingContext(ctx, md)
}

func (connector *Connector) serveLoop(ctx context.Context, client lunopb.StargateClient, handler Handler) error {
	heartbeat := time.NewTicker(DefaultConnectorHeartbeat)
	defer heartbeat.Stop()

	for {
		logger := connector.logger.With(zap.String("address", connector.StargateAddress))
		logger.Info("attempting to connect to Stargate")

		connector.serveTick(ctx, client, handler)

		logger.Info("connection closed, attempting to reconnect", zap.Duration("heartbeat", DefaultConnectorHeartbeat))

		select {
		case <-ctx.Done():
			connector.logger.Info("context cancelled, stopping connector")
			return nil
		case <-heartbeat.C:
			logger.Info("attempting to reconnect to Stargate...")
		}
	}
}

func (connector *Connector) serveTick(ctx context.Context, client lunopb.StargateClient, handler Handler) {
	ctx = connector.defaultOutgoingContext(ctx)
	stream, err := client.Connector(ctx)
	if err != nil {
		connector.logger.Error("failed to connect to Stargate", zap.Error(err))
		return
	}

	connector.logger.Info("connected to Stargate")
	err = connector.recvLoop(ctx, stream, handler)
	if err != nil {
		connector.logger.Error("unexpected error in receive loop", zap.Error(err))
		return
	}
}

func (connector *Connector) recvLoop(ctx context.Context, stream grpc.BidiStreamingClient[lunopb.ConnectorResponse, lunopb.ConnectorRequest], handler Handler) error {
	connector.health(true)
	defer connector.health(false)

	logger := connector.logger.With(zap.String("address", connector.StargateAddress))
	logger.Info("starting message receive loop")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		msg, err := stream.Recv()
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil
		}

		if err != nil {
			return err
		}

		logger.Debug("received message", zap.Uint32("id", msg.Id))

		switch state := msg.State.(type) {
		case *lunopb.ConnectorRequest_Ping:
			go func() {
				err := connector.ping(ctx, msg.Id, stream, handler)
				if err != nil {
					logger.Error("failed to ping", zap.Error(err))
					cancel()
				}
			}()
		case *lunopb.ConnectorRequest_Fetch:
			go func() {
				err := connector.fetch(ctx, msg.Id, stream, handler)
				if err != nil {
					logger.Error("failed to fetch tables", zap.Error(err))
					cancel()
				}
			}()
		case *lunopb.ConnectorRequest_ExecuteStatement:
			go func() {
				err := connector.execute(ctx, msg.Id, state.ExecuteStatement, stream, handler)
				if err != nil {
					logger.Error("failed to execute statement", zap.Error(err))
					cancel()
				}
			}()
		}
	}
}

func (connector *Connector) health(status bool) {
	connector.mu.Lock()
	defer connector.mu.Unlock()

	connector.healthy = status
}

func (connector *Connector) ping(ctx context.Context, id uint32, stream grpc.BidiStreamingClient[lunopb.ConnectorResponse, lunopb.ConnectorRequest], handler Handler) error {
	logger := connector.logger.With(zap.Uint32("id", id))
	logger.Debug("ping connector")

	pong := &lunopb.PingResponse{}
	err := handler.Ping(ctx)
	if err != nil {
		logger.Error("unexpected error while pinging", zap.Error(err))
		pong.Error = &lunopb.Error{
			Message: err.Error(),
		}
	}

	logger.Debug("ping complete")
	return stream.Send(&lunopb.ConnectorResponse{
		Id: id,
		State: &lunopb.ConnectorResponse_Ping{
			Ping: pong,
		},
	})
}

func (connector *Connector) fetch(ctx context.Context, id uint32, stream grpc.BidiStreamingClient[lunopb.ConnectorResponse, lunopb.ConnectorRequest], handler Handler) error {
	logger := connector.logger.With(zap.Uint32("id", id))
	logger.Debug("fetching tables")

	fetch := &lunopb.FetchResponse{}
	tables, err := handler.Fetch(ctx)
	if err != nil {
		logger.Error("unexpected error while fetching tables", zap.Error(err))
		fetch.Error = &lunopb.Error{
			Message: err.Error(),
		}
	}

	if tables != nil {
		fetch.Tables = tables.Proto()
	}

	logger.Debug("tables fetched", zap.Int("count", len(fetch.Tables)))
	return stream.Send(&lunopb.ConnectorResponse{
		Id: id,
		State: &lunopb.ConnectorResponse_Fetch{
			Fetch: fetch,
		},
	})
}

func (connector *Connector) execute(ctx context.Context, id uint32, state *lunopb.ExecuteStatementRequest, stream grpc.BidiStreamingClient[lunopb.ConnectorResponse, lunopb.ConnectorRequest], handler Handler) error {
	plan := state.Plan

	logger := connector.logger.With(zap.Uint32("id", id))
	logger.Debug("executing statement")

	writer := WriterFunc(func(ctx context.Context, values []any) (err error) {
		row := make([][]byte, len(values))
		for index, col := range values {
			_, row[index], err = value.Encode(col, nil)
			if err != nil {
				return err
			}
		}

		logger.Debug("writing row")
		return stream.Send(&lunopb.ConnectorResponse{
			Id: id,
			State: &lunopb.ConnectorResponse_ExecuteStatement{
				ExecuteStatement: &lunopb.ExecuteStatementResponse{
					Result: &lunopb.ExecuteStatementResponse_Data{
						Data: &lunopb.Row{
							Values: row,
						},
					},
				},
			},
		})
	})

	err := handler.Scan(ctx, plan, writer)
	if err != nil {
		logger.Error("unexpected error while scanning", zap.Error(err))
		return stream.Send(&lunopb.ConnectorResponse{
			Id: id,
			State: &lunopb.ConnectorResponse_ExecuteStatement{
				ExecuteStatement: &lunopb.ExecuteStatementResponse{
					Result: &lunopb.ExecuteStatementResponse_Error{
						Error: &lunopb.Error{
							Message: err.Error(),
						},
					},
				},
			},
		})
	}

	logger.Debug("statement executed successfully")
	return stream.Send(&lunopb.ConnectorResponse{
		Id: id,
		State: &lunopb.ConnectorResponse_ExecuteStatement{
			ExecuteStatement: &lunopb.ExecuteStatementResponse{
				Result: &lunopb.ExecuteStatementResponse_EOE{},
			},
		},
	})
}
