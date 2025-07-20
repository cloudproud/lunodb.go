package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudproud/lunodb.api/proto/plan"
	lunodb "github.com/cloudproud/lunodb.go"
	"github.com/cloudproud/lunodb.go/types"
	"go.uber.org/zap"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	connector, err := lunodb.NewConnector(
		lunodb.WithSource(os.Getenv("LUNODB_SOURCE")),
		lunodb.WithToken(os.Getenv("LUNODB_TOKEN")),
		lunodb.WithLogger(zap.Must(zap.NewDevelopment())),
	)
	if err != nil {
		return err
	}

	return connector.Serve(context.Background(), &Connector{})
}

type Connector struct{}

func (c *Connector) Ping(ctx context.Context) error {
	return nil
}

func (c *Connector) Fetch(ctx context.Context) (lunodb.Tables, error) {
	weather := lunodb.Table{
		Name:    "weather",
		Schema:  "public",
		Catalog: "meteorology",
		Columns: []lunodb.Column{
			{
				Name: "city",
				Type: types.BasicString,
				Operators: []lunodb.Operator{
					{
						Statement: lunodb.StatementEqual,
						ComparisonTypes: []lunodb.ComparisonType{
							lunodb.VariableConstant,
						},
						Required: true,
					},
				},
			},
			{
				Name: "temperature",
				Type: types.BasicString,
			},
			{
				Name: "humidity",
				Type: types.BasicString,
			},
		},
	}

	return []lunodb.Table{weather}, nil
}

type CurrentCondition struct {
	Humidity string `json:"humidity"`
	TempC    string `json:"temp_C"`
}

type WeatherResponse struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
}

func (c *Connector) Scan(ctx context.Context, plan *plan.Literal, writer lunodb.Writer) error {
	city := "Amsterdam"
	endpoint := fmt.Sprintf("https://wttr.in/%s?format=j1", city)
	res, err := http.Get(endpoint)
	if err != nil {
		return err
	}

	defer res.Body.Close() //nolint:errcheck

	weather := WeatherResponse{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&weather)
	if err != nil {
		return err
	}

	current := weather.CurrentCondition[0]
	return writer.Write(ctx, []any{city, current.TempC, current.Humidity})
}
