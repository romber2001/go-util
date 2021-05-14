package prometheus

import (
	"fmt"
	"testing"
	"time"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
)

const (
	defaultAddr = "192.168.10.210:80/prometheus"
	defaultUser = "admin"
	defaultPass = "admin"
)

var conn = initConn()

func initConn() *Conn {
	config := NewConfigWithBasicAuth(defaultAddr, defaultUser, defaultPass)
	c, err := NewConnWithConfig(config)
	if err != nil {
		log.Error(fmt.Sprintf("initAppRepo() failed.\n%s", err.Error()))
		return nil
	}

	return c
}

func TestConn_Execute(t *testing.T) {
	asst := assert.New(t)

	var (
		err    error
		query  string
		result *Result
	)

	start := time.Now().Add(-time.Hour)
	end := time.Now()
	step := time.Minute
	// r := apiv1.Range{
	// 	Start: start,
	// 	End:   end,
	// 	Step:  step,
	// }

	// query := "1"
	query = `rate(mysql_global_status_queries)[1m]`
	// query := `mysql_global_status_queries`
	result, err = conn.Execute(query, start, end, step)
	asst.Nil(err, "test Execute() failed")
	s, err := result.GetString(constant.ZeroInt, constant.ZeroInt)
	asst.Nil(err, "test Execute() failed")
	ts, err := result.GetString(constant.ZeroInt, 1)
	asst.Nil(err, "test Execute() failed")
	t.Log(s, ts)

	query = `sum(avg by (node_name,mode) (clamp_max(((avg by (mode,node_name) ( (clamp_max(rate(node_cpu_seconds_total{node_name=~"192-168-137-11",mode!="idle"}[5s]),1)) or (clamp_max(irate(node_cpu_seconds_total{node_name=~"192-168-137-11",mode!="idle"}[5m]),1)) ))*100 or (avg_over_time(node_cpu_average{node_name=~"192-168-137-11", mode!="total", mode!="idle"}[5s]) or avg_over_time(node_cpu_average{node_name=~"192-168-137-11", mode!="total", mode!="idle"}[5m]))),100)))`
	result, err = conn.Execute(query, start, end, step)
	asst.Nil(err, "test Execute() failed")
	t.Log(result)

	query = `topk(10, avg by (service_name,schema,table) (sum(mysql_info_schema_table_rows{service_name=~"192-168-10-210:3306"}) by (service_name, schema, table))) > 0`
	result, err = conn.Execute(query, start, end, step)
	asst.Nil(err, "test Execute() failed")
	t.Log(result)
}
