package prometheus

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	client "github.com/prometheus/client_golang/api"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/romberli/go-util/constant"
)

const (
	defaultHTTPPrefix          = "http://"
	defaultHTTPSPrefix         = "https://"
	defaultDialTimeout         = 30 * time.Second
	defaultKeepAlive           = 30 * time.Second
	defaultTLSHandshakeTimeout = 10 * time.Second
)

type Config struct {
	client.Config
}

// DefaultRoundTripper is used if no RoundTripper is set in Config,
var DefaultRoundTripper http.RoundTripper = &http.Transport{
	Proxy:               http.ProxyFromEnvironment,
	DialContext:         (&net.Dialer{Timeout: defaultDialTimeout, KeepAlive: defaultKeepAlive}).DialContext,
	TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
}

// NewConfig returns a new client.Config with given address and round tripper
func NewConfig(addr string, rt http.RoundTripper) Config {
	address := strings.ToLower(addr)
	if !strings.HasPrefix(address, defaultHTTPPrefix) && !strings.HasPrefix(address, defaultHTTPSPrefix) {
		addr = defaultHTTPPrefix + addr
	}

	if rt == nil {
		rt = DefaultRoundTripper
	}

	return Config{
		client.Config{
			Address:      addr,
			RoundTripper: rt,
		},
	}
}

// NewConfigWithDefaultRoundTripper returns a new client.Config with given address and default round tripper
func NewConfigWithDefaultRoundTripper(addr string) Config {
	address := strings.ToLower(addr)
	if !strings.HasPrefix(address, defaultHTTPPrefix) && !strings.HasPrefix(address, defaultHTTPSPrefix) {
		addr = defaultHTTPSPrefix + addr
	}

	return Config{
		client.Config{
			Address:      addr,
			RoundTripper: DefaultRoundTripper,
		},
	}
}

// NewConfigWithBasicAuth returns a new client.Config with given address, user and password
func NewConfigWithBasicAuth(addr, user, pass string) Config {
	address := strings.ToLower(addr)
	if !strings.HasPrefix(address, defaultHTTPPrefix) && !strings.HasPrefix(address, defaultHTTPSPrefix) {
		addr = defaultHTTPPrefix + addr
	}

	return Config{
		client.Config{
			Address:      addr,
			RoundTripper: config.NewBasicAuthRoundTripper(user, config.Secret(pass), constant.EmptyString, DefaultRoundTripper),
		},
	}
}

type Conn struct {
	apiv1.API
}

// NewConn returns a new *Conn with given address and round tripper
func NewConn(addr string, rt http.RoundTripper) (*Conn, error) {
	return NewConnWithConfig(NewConfig(addr, rt))
}

// NewConnWithConfig returns a new *Conn with given config
func NewConnWithConfig(config Config) (*Conn, error) {
	cli, err := client.NewClient(config.Config)
	if err != nil {
		return nil, err
	}

	return &Conn{apiv1.NewAPI(cli)}, nil
}

func (conn *Conn) CheckInstanceStatus() bool {
	query := "1"
	result, err := conn.Execute(query)
	if err != nil {
		return false
	}

	status, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	if err != nil {
		return false
	}

	return status == 1
}

// Execute executes given command with arguments and returns a result
func (conn *Conn) Execute(command string, args ...interface{}) (*Result, error) {
	return conn.executeContext(context.Background(), command, args...)
}

// ExecuteContext executes given command with arguments and returns a result
func (conn *Conn) ExecuteContext(ctx context.Context, command string, args ...interface{}) (*Result, error) {
	return conn.executeContext(ctx, command, args...)
}

// executeContext executes given command with arguments and returns a result.
// if args length is 0:
// 		it uses time.Now() as the time series
// if args length is 1:
// 		argument type must be either time.Time or TimeRange
// if args length is 2:
// 		argument types must be time.Time and time.Time, represent start time and end time, it uses 1 minute as step
// if args length is 3:
//		argument types muse be in order of time.Time, time.Time and time.Duration, represent start time, end time and step
// if args length is larger than 3:
// 		it returns error
func (conn *Conn) executeContext(ctx context.Context, command string, args ...interface{}) (*Result, error) {
	var (
		arg      interface{}
		value    model.Value
		warnings apiv1.Warnings
		err      error
	)

	switch len(args) {
	case 0:
		arg = time.Now()
	case 1:
		arg = args[constant.ZeroInt]
	case 2:
		start, startOK := args[0].(time.Time)
		end, endOK := args[1].(time.Time)
		if !(startOK && endOK) {
			return nil, errors.New(
				"args length is 2, should be in order of time.Time, time.Time, represent start time, end time")
		}

		arg = NewTimeRange(start, end, DefaultStep)
	case 3:
		start, startOK := args[0].(time.Time)
		end, endOK := args[1].(time.Time)
		step, stepOK := args[2].(time.Duration)
		if !(startOK && endOK && stepOK) {
			return nil, errors.New(
				"args length is 3, should be in order of time.Time, time.Time and time.Duration, represent start time, end time and step")
		}

		arg = NewTimeRange(start, end, step)
	default:
		return nil, errors.New(fmt.Sprintf("args length shoud be less or equal to 3, %d is not valid", len(args)))
	}

	switch arg.(type) {
	case time.Time:
		value, warnings, err = conn.Query(ctx, command, arg.(time.Time))
		if err != nil {
			return nil, err
		}
	case TimeRange:
		value, warnings, err = conn.QueryRange(ctx, command, arg.(TimeRange).GetRange())
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("unsupported argument type: %T", arg))
	}

	return NewResult(value, warnings), nil
}