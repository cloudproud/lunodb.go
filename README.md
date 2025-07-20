# LunoDB Go SDK

LunoDB is a federated query engine that allows you to connect to external data sources and query them using standard SQL. This Go SDK provides everything you need to implement your own custom connector, making your data accessible within LunoDB’s unified query environment.

A Connector bridges LunoDB and your external data source-be it an API, database, or anything else you can query. You define what data is available and how to fetch it, and LunoDB takes care of the rest.

This SDK lets you implement:

- Schema introspection (Fetch)
- Streaming data retrieval (Scan)
- Health checks (Ping)

## 📁 Examples

Check out the [example directory](https://github.com/cloudproud/lunodb.go/tree/main/examples) for working connector implementations that showcase integrations with real APIs and data sources.

```go
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

func (c *Connector) Scan(ctx context.Context, plan *plan.Literal, writer lunodb.Writer) error {
	city := "Amsterdam"
	endpoint := fmt.Sprintf("https://wttr.in/%s?format=j1", city)
	res, err := http.Get(endpoint)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	weather := WeatherResponse{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&weather)
	if err != nil {
		return err
	}

	current := weather.CurrentCondition[0]
	return writer.Write(ctx, []any{city, current.TempC, current.Humidity})
}
```
