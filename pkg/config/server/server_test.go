package server

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	stnrv1 "github.com/l7mp/stunner/pkg/apis/v1"
	"github.com/l7mp/stunner/pkg/config/client"
	"github.com/l7mp/stunner/pkg/logger"
)

var testerLogLevel = zapcore.Level(-4)

//var testerLogLevel = zapcore.DebugLevel
//var testerLogLevel = zapcore.ErrorLevel

const stunnerLogLevel = "all:TRACE"

func init() {
	// setup a fast pinger so that we get a timely error notification
	client.PingPeriod = 500 * time.Millisecond
	client.PongWait = 800 * time.Millisecond
	client.WriteWait = 200 * time.Millisecond
	client.RetryPeriod = 250 * time.Millisecond

}

func TestServerLoad(t *testing.T) {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(testerLogLevel)
	z, err := zc.Build()
	assert.NoError(t, err, "logger created")
	zlogger := zapr.NewLogger(z)
	log := zlogger.WithName("tester")

	logger := logger.NewLoggerFactory(stunnerLogLevel)
	testLog := logger.NewLogger("test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testLog.Debug("create server")
	server := New(stnrv1.DefaultConfigDiscoveryAddress, log)
	assert.NotNil(t, server, "server")
	err = server.Start(ctx)
	assert.NoError(t, err, "start")

	time.Sleep(20 * time.Millisecond)

	testLog.Debug("create client")
	client1, err := client.New("127.0.0.1:13478", "ns1/gw1", logger)
	assert.NoError(t, err, "client 1")
	client2, err := client.New("127.0.0.1:13478", "ns1/gw2", logger)
	assert.NoError(t, err, "client 2")
	// nonexistent
	client3, err := client.New("127.0.0.1:13478", "ns1/gw3", logger)
	assert.NoError(t, err, "client 3")

	testLog.Debug("load: error")
	c, err := client1.Load()
	assert.Error(t, err, "load")
	assert.Nil(t, c, "conf")
	c, err = client2.Load()
	assert.Error(t, err, "load")
	assert.Nil(t, c, "conf")
	c, err = client3.Load()
	assert.Error(t, err, "load")
	assert.Nil(t, c, "conf")

	c1 := testConfig("ns1/gw1", "realm1")
	c2 := testConfig("ns1/gw2", "realm1")
	err = server.UpdateConfig([]Config{c1, c2})
	assert.NoError(t, err, "update")

	cs := server.configs.Snapshot()
	assert.Len(t, cs, 2, "snapshot len")
	sc1 := server.configs.Get("ns1/gw1")
	assert.NotNil(t, sc1, "get 1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq")
	sc2 := server.configs.Get("ns1/gw2")
	assert.NotNil(t, sc2, "get 2")
	assert.True(t, c2.Config.DeepEqual(sc2), "deepeq")
	sc3 := server.configs.Get("ns1/gw3")
	assert.Nil(t, sc3, "get 3")

	testLog.Debug("load: config ok")
	c, err = client1.Load()
	assert.NoError(t, err, "load")
	assert.True(t, c.DeepEqual(sc1), "deepeq")
	c, err = client2.Load()
	assert.NoError(t, err, "load")
	assert.True(t, c.DeepEqual(sc2), "deepeq")
	c, err = client3.Load()
	assert.Error(t, err, "load")
	assert.Nil(t, c, "conf")

	testLog.Debug("remove 2 configs")
	err = server.UpdateConfig([]Config{})
	assert.NoError(t, err, "update")

	cs = server.configs.Snapshot()
	assert.Len(t, cs, 0, "snapshot len")

	testLog.Debug("load: no result")
	_, err = client1.Load()
	assert.Error(t, err, "load")
	_, err = client2.Load()
	assert.Error(t, err, "load")
	_, err = client3.Load()
	assert.Error(t, err, "load")
	assert.Nil(t, c, "conf")
}

func TestServerPoll(t *testing.T) {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(testerLogLevel)
	z, err := zc.Build()
	assert.NoError(t, err, "logger created")
	zlogger := zapr.NewLogger(z)
	log := zlogger.WithName("tester")

	logger := logger.NewLoggerFactory(stunnerLogLevel)
	testLog := logger.NewLogger("test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testLog.Debug("create server")
	server := New(stnrv1.DefaultConfigDiscoveryAddress, log)
	assert.NotNil(t, server, "server")
	err = server.Start(ctx)
	assert.NoError(t, err, "start")

	time.Sleep(20 * time.Millisecond)

	testLog.Debug("create client")
	client1, err := client.New("127.0.0.1:13478", "ns1/gw1", logger)
	assert.NoError(t, err, "client 1")
	client2, err := client.New("127.0.0.1:13478", "ns1/gw2", logger)
	assert.NoError(t, err, "client 2")
	client3, err := client.New("127.0.0.1:13478", "ns1/gw3", logger)
	assert.NoError(t, err, "client 3")

	testLog.Debug("poll: no result")
	ch1 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch1)
	ch2 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch2)
	ch3 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch3)

	go func() {
		err = client1.Poll(ctx, ch1)
		assert.NoError(t, err, "client 1 cancelled")
	}()
	go func() {
		err = client2.Poll(ctx, ch2)
		assert.NoError(t, err, "client 2 cancelled")
	}()
	go func() {
		err = client3.Poll(ctx, ch2)
		assert.NoError(t, err, "client 3 cancelled")
	}()

	s := watchConfig(ch1, 10*time.Millisecond)
	assert.Nil(t, s, "config 1")
	s = watchConfig(ch2, 10*time.Millisecond)
	assert.Nil(t, s, "config 2")
	s = watchConfig(ch3, 10*time.Millisecond)
	assert.Nil(t, s, "config 3")

	testLog.Debug("poll: one result")
	c1 := testConfig("ns1/gw1", "realm1")
	c2 := testConfig("ns1/gw2", "realm1")
	err = server.UpdateConfig([]Config{c1, c2})
	assert.NoError(t, err, "update")

	cs := server.configs.Snapshot()
	assert.Len(t, cs, 2, "snapshot len")
	sc1 := server.configs.Get("ns1/gw1")
	assert.NotNil(t, sc1, "get 1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq")
	sc2 := server.configs.Get("ns1/gw2")
	assert.NotNil(t, sc2, "get 2")
	assert.True(t, c2.Config.DeepEqual(sc2), "deepeq")
	sc3 := server.configs.Get("ns1/gw3")
	assert.Nil(t, sc3, "get 3")

	// poll should have fed the configs to the channels
	s = watchConfig(ch1, 500*time.Millisecond)
	assert.NotNil(t, s, "config 1")
	assert.True(t, s.DeepEqual(sc1), "deepeq 1")
	s = watchConfig(ch2, 500*time.Millisecond)
	assert.NotNil(t, s, "config 2")
	assert.True(t, s.DeepEqual(sc2), "deepeq 2")
	s = watchConfig(ch3, 500*time.Millisecond)
	assert.Nil(t, s, "config 3")

	testLog.Debug("remove 2 configs")
	err = server.UpdateConfig([]Config{})
	assert.NoError(t, err, "update")

	cs = server.configs.Snapshot()
	assert.Len(t, cs, 0, "snapshot len")

	testLog.Debug("poll: zeroconfig")
	s = watchConfig(ch1, 10*time.Millisecond)
	assert.Nil(t, s, "config")
	s = watchConfig(ch2, 10*time.Millisecond)
	assert.Nil(t, s, "config")
	s = watchConfig(ch3, 10*time.Millisecond)
	assert.Nil(t, s, "config")
}

func TestServerWatch(t *testing.T) {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(testerLogLevel)
	z, err := zc.Build()
	assert.NoError(t, err, "logger created")
	zlogger := zapr.NewLogger(z)
	log := zlogger.WithName("tester")

	logger := logger.NewLoggerFactory(stunnerLogLevel)
	testLog := logger.NewLogger("test")

	serverCtx, serverCancel := context.WithCancel(context.Background())

	testLog.Debug("create server")
	server := New(stnrv1.DefaultConfigDiscoveryAddress, log)
	assert.NotNil(t, server, "server")
	err = server.Start(serverCtx)
	assert.NoError(t, err, "start")

	testLog.Debug("create client")
	client1, err := client.New("127.0.0.1:13478", "ns1/gw1", logger)
	assert.NoError(t, err, "client 1")
	client2, err := client.New("127.0.0.1:13478", "ns1/gw2", logger)
	assert.NoError(t, err, "client 2")
	client3, err := client.New("127.0.0.1:13478", "ns1/gw3", logger)
	assert.NoError(t, err, "client 3")

	testLog.Debug("watch: no result")
	ch1 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch1)
	ch2 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch2)
	ch3 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch3)

	clientCtx, clientCancel := context.WithCancel(context.Background())
	defer clientCancel()
	err = client1.Watch(clientCtx, ch1)
	assert.NoError(t, err, "client 1 watch")
	err = client2.Watch(clientCtx, ch2)
	assert.NoError(t, err, "client 2 watch")
	err = client3.Watch(clientCtx, ch3)
	assert.NoError(t, err, "client 3 watch")

	s := watchConfig(ch1, 150*time.Millisecond)
	assert.Nil(t, s, "config 1")
	s = watchConfig(ch2, 150*time.Millisecond)
	assert.Nil(t, s, "config 2")
	s = watchConfig(ch3, 150*time.Millisecond)
	assert.Nil(t, s, "config 3")

	testLog.Debug("poll: one result")
	c1 := testConfig("ns1/gw1", "realm1")
	c2 := testConfig("ns1/gw2", "realm1")
	err = server.UpdateConfig([]Config{c1, c2})
	assert.NoError(t, err, "update")

	cs := server.configs.Snapshot()
	assert.Len(t, cs, 2, "snapshot len")
	sc1 := server.configs.Get("ns1/gw1")
	assert.NotNil(t, sc1, "get 1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq")
	sc2 := server.configs.Get("ns1/gw2")
	assert.NotNil(t, sc2, "get 2")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq")
	sc3 := server.configs.Get("ns1/gw3")
	assert.Nil(t, sc3, "get 3")

	// poll should have fed the configs to the channels
	s = watchConfig(ch1, 500*time.Millisecond)
	assert.NotNil(t, s, "config 1")
	assert.True(t, s.DeepEqual(sc1), "deepeq 1")
	s = watchConfig(ch2, 500*time.Millisecond)
	assert.NotNil(t, s, "config 2")
	assert.True(t, s.DeepEqual(sc2), "deepeq 2")
	s = watchConfig(ch3, 500*time.Millisecond)
	assert.Nil(t, s, "config 3")

	testLog.Debug("update: conf 1 and conf 3")
	c1 = testConfig("ns1/gw1", "realm-new")
	c3 := testConfig("ns1/gw3", "realm3")
	err = server.UpdateConfig([]Config{c1, c2, c3})
	assert.NoError(t, err, "update")

	cs = server.configs.Snapshot()
	assert.Len(t, cs, 3, "snapshot len")
	sc1 = server.configs.Get("ns1/gw1")
	assert.NotNil(t, sc1, "get 1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq 1")
	sc2 = server.configs.Get("ns1/gw2")
	assert.NotNil(t, sc2, "get 2")
	assert.True(t, c2.Config.DeepEqual(sc2), "deepeq 2")
	sc3 = server.configs.Get("ns1/gw3")
	assert.NotNil(t, sc3, "get 3")
	assert.True(t, c3.Config.DeepEqual(sc3), "deepeq 3")

	// poll should have fed the configs to the channels
	s = watchConfig(ch1, 500*time.Millisecond)
	assert.NotNil(t, s, "config 1")
	assert.True(t, s.DeepEqual(sc1), "deepeq 1")
	s = watchConfig(ch2, 500*time.Millisecond)
	assert.Nil(t, s, "config 2")
	s = watchConfig(ch3, 500*time.Millisecond)
	assert.NotNil(t, s, "config 3")
	assert.True(t, s.DeepEqual(sc3), "deepeq 3")

	testLog.Debug("restarting server")
	serverCancel()
	// let the server shut down and restart
	time.Sleep(50 * time.Millisecond)
	serverCtx, serverCancel = context.WithCancel(context.Background())
	defer serverCancel()
	server = New(stnrv1.DefaultConfigDiscoveryAddress, log)
	assert.NotNil(t, server, "server")
	err = server.Start(serverCtx)
	assert.NoError(t, err, "start")
	err = server.UpdateConfig([]Config{c1, c2, c3})
	assert.NoError(t, err, "update")

	// obtain the initial configs: this may take a while
	s = watchConfig(ch1, 5000*time.Millisecond)
	assert.NotNil(t, s, "config 1")
	assert.True(t, s.DeepEqual(sc1), "deepeq 1")
	s = watchConfig(ch2, 500*time.Millisecond)
	assert.NotNil(t, s, "config 2")
	assert.True(t, s.DeepEqual(sc2), "deepeq 2")
	s = watchConfig(ch3, 500*time.Millisecond)
	assert.NotNil(t, s, "config 3")
	assert.True(t, s.DeepEqual(sc3), "deepeq 3")

	testLog.Debug("remove 1 config (the 2nd)")
	err = server.UpdateConfig([]Config{c1, c3})
	assert.NoError(t, err, "update")

	cs = server.configs.Snapshot()
	assert.Len(t, cs, 2, "snapshot len")
	sc1 = server.configs.Get("ns1/gw1")
	assert.NotNil(t, sc1, "get 1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq 1")
	sc2 = server.configs.Get("ns1/gw2")
	assert.Nil(t, sc2, "get 2")
	sc3 = server.configs.Get("ns1/gw3")
	assert.NotNil(t, sc3, "get 3")
	assert.True(t, c3.Config.DeepEqual(sc3), "deepeq 3")

	s = watchConfig(ch1, 50*time.Millisecond)
	assert.Nil(t, s, "config 1")
	s = watchConfig(ch2, 50*time.Millisecond)
	assert.Nil(t, s, "config 2")
	s = watchConfig(ch3, 50*time.Millisecond)
	assert.Nil(t, s, "config 3")

	testLog.Debug("remove remaining 2 configs")
	err = server.UpdateConfig([]Config{})
	assert.NoError(t, err, "update")

	cs = server.configs.Snapshot()
	assert.Len(t, cs, 0, "snapshot len")

	testLog.Debug("poll: no config")
	s = watchConfig(ch1, 10*time.Millisecond)
	assert.Nil(t, s, "config")
	s = watchConfig(ch2, 10*time.Millisecond)
	assert.Nil(t, s, "config")
	s = watchConfig(ch3, 10*time.Millisecond)
	assert.Nil(t, s, "config")
}

// test APIs
func TestServerAPI(t *testing.T) {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(testerLogLevel)
	z, err := zc.Build()
	assert.NoError(t, err, "logger created")
	zlogger := zapr.NewLogger(z)
	log := zlogger.WithName("tester")

	logger := logger.NewLoggerFactory(stunnerLogLevel)
	testLog := logger.NewLogger("test")

	serverCtx, serverCancel := context.WithCancel(context.Background())

	testLog.Debug("create server")
	server := New(stnrv1.DefaultConfigDiscoveryAddress, log)
	assert.NotNil(t, server, "server")
	err = server.Start(serverCtx)
	assert.NoError(t, err, "start")

	testLog.Debug("create client")
	client1, err := client.NewAllConfigsAPI("127.0.0.1:13478", logger.NewLogger("all-config-client"))
	assert.NoError(t, err, "client 1")
	client2, err := client.NewConfigsNamespaceAPI("127.0.0.1:13478", "ns1", logger.NewLogger("ns-config-client-ns1"))
	assert.NoError(t, err, "client 2")
	client3, err := client.NewConfigsNamespaceAPI("127.0.0.1:13478", "ns2", logger.NewLogger("ns-config-client-ns2"))
	assert.NoError(t, err, "client 3")
	client4, err := client.NewConfigNamespaceNameAPI("127.0.0.1:13478", "ns1", "gw1", logger.NewLogger("gw-config-client"))
	assert.NoError(t, err, "client 4")

	testLog.Debug("watch: no result")
	ch1 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch1)
	ch2 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch2)
	ch3 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch3)
	ch4 := make(chan stnrv1.StunnerConfig, 8)
	defer close(ch4)

	clientCtx, clientCancel := context.WithCancel(context.Background())
	defer clientCancel()
	err = client1.Watch(clientCtx, ch1)
	assert.NoError(t, err, "client 1 watch")
	err = client2.Watch(clientCtx, ch2)
	assert.NoError(t, err, "client 2 watch")
	err = client3.Watch(clientCtx, ch3)
	assert.NoError(t, err, "client 3 watch")
	err = client4.Watch(clientCtx, ch4)
	assert.NoError(t, err, "client 4 watch")

	s := watchConfig(ch1, 50*time.Millisecond)
	assert.Nil(t, s, "config 1")
	s = watchConfig(ch2, 50*time.Millisecond)
	assert.Nil(t, s, "config 2")
	s = watchConfig(ch3, 50*time.Millisecond)
	assert.Nil(t, s, "config 3")
	s = watchConfig(ch4, 50*time.Millisecond)
	assert.Nil(t, s, "config 4")

	testLog.Debug("--------------------------------")
	testLog.Debug("Update1: ns1/gw1 + ns2/gw1      ")
	testLog.Debug("--------------------------------")
	testLog.Debug("poll: one result")
	c1 := testConfig("ns1/gw1", "realm1")
	c2 := testConfig("ns2/gw1", "realm1")
	err = server.UpdateConfig([]Config{c1, c2})
	assert.NoError(t, err, "update")

	cs := server.configs.Snapshot()
	assert.Len(t, cs, 2, "snapshot len")
	sc1 := server.configs.Get("ns1/gw1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq 1")
	assert.NotNil(t, sc1, "get 1")
	sc2 := server.configs.Get("ns2/gw1")
	assert.NotNil(t, sc2, "get 2")
	assert.True(t, c2.Config.DeepEqual(sc2), "deepeq 2")

	testLog.Debug("load")

	// all-configs should result sc1 and sc2
	scs, err := client1.Get(clientCtx)
	assert.NoError(t, err, "load 1")
	assert.Len(t, scs, 2, "load 1")
	co := findConfById(scs, "ns1/gw1")
	assert.NotNil(t, co, "c1")
	assert.True(t, co.DeepEqual(sc1), "deepeq")
	co = findConfById(scs, "ns2/gw1")
	assert.NotNil(t, co, "c2")
	assert.True(t, co.DeepEqual(sc2), "deepeq")

	// ns1 client should yield 1 config
	scs, err = client2.Get(clientCtx)
	assert.NoError(t, err, "load 2")
	assert.Len(t, scs, 1, "load 2")
	assert.True(t, scs[0].DeepEqual(sc1), "deepeq")

	// ns2 client should yield 1 config
	scs, err = client3.Get(clientCtx)
	assert.NoError(t, err, "load 3")
	assert.Len(t, scs, 1, "load 3")
	assert.True(t, scs[0].DeepEqual(sc2), "deepeq")

	// ns1/gw1 client should yield 1 config
	scs, err = client4.Get(clientCtx)
	assert.NoError(t, err, "load 4")
	assert.Len(t, scs, 1, "load 4")
	assert.True(t, scs[0].DeepEqual(sc1), "deepeq")

	// two configs from client1 watch
	s1 := watchConfig(ch1, 50*time.Millisecond)
	assert.NotNil(t, s1)
	s2 := watchConfig(ch1, 50*time.Millisecond)
	assert.NotNil(t, s2)
	s3 := watchConfig(ch1, 50*time.Millisecond)
	assert.Nil(t, s3)
	lst := []*stnrv1.StunnerConfig{s1, s2}
	assert.NotNil(t, findConfById(lst, "ns1/gw1"))
	assert.True(t, findConfById(lst, "ns1/gw1").DeepEqual(sc1), "deepeq 1")
	assert.NotNil(t, findConfById(lst, "ns2/gw1"))
	assert.True(t, findConfById(lst, "ns2/gw1").DeepEqual(sc2), "deepeq 1")

	// 1 config from client2 watch
	s = watchConfig(ch2, 50*time.Millisecond)
	assert.NotNil(t, s)
	assert.True(t, s.DeepEqual(sc1))
	s = watchConfig(ch2, 50*time.Millisecond)
	assert.Nil(t, s)

	// 1 config from client3 watch
	s = watchConfig(ch3, 50*time.Millisecond)
	assert.NotNil(t, s, "config 3")
	assert.True(t, s.DeepEqual(sc2))
	s = watchConfig(ch3, 50*time.Millisecond)
	assert.Nil(t, s)

	// 1 config from client4 watch
	s = watchConfig(ch4, 50*time.Millisecond)
	assert.NotNil(t, s)
	assert.True(t, s.DeepEqual(sc1))
	s = watchConfig(ch4, 50*time.Millisecond)
	assert.Nil(t, s)

	testLog.Debug("--------------------------------")
	testLog.Debug("Update1: ns1/gw1 + ns1/gw2      ")
	testLog.Debug("--------------------------------")
	testLog.Debug("update: conf 1 and conf 3")
	c1 = testConfig("ns1/gw1", "realm-new")
	c3 := testConfig("ns1/gw2", "realm3")
	err = server.UpdateConfig([]Config{c1, c2, c3})
	assert.NoError(t, err, "update")

	cs = server.configs.Snapshot()
	assert.Len(t, cs, 3, "snapshot len")
	sc1 = server.configs.Get("ns1/gw1")
	assert.NotNil(t, sc1, "get 1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq")
	sc2 = server.configs.Get("ns2/gw1")
	assert.NotNil(t, sc2, "get 2")
	assert.True(t, c2.Config.DeepEqual(sc2), "deepeq")
	sc3 := server.configs.Get("ns1/gw2")
	assert.NotNil(t, sc3, "get 3")
	assert.True(t, c3.Config.DeepEqual(sc3), "deepeq")

	// all-configs should result sc1 and sc2 and sc3
	scs, err = client1.Get(clientCtx)
	assert.NoError(t, err, "load 1")
	assert.Len(t, scs, 3, "load 1")
	co = findConfById(scs, "ns1/gw1")
	assert.NotNil(t, co, "c1")
	assert.True(t, co.DeepEqual(sc1), "deepeq")
	co = findConfById(scs, "ns2/gw1")
	assert.NotNil(t, co, "c2")
	assert.True(t, co.DeepEqual(sc2), "deepeq")
	co = findConfById(scs, "ns1/gw2")
	assert.NotNil(t, co, "c3")
	assert.True(t, co.DeepEqual(sc3), "deepeq")

	// ns1 client should yield 2 configs
	scs, err = client2.Get(clientCtx)
	assert.NoError(t, err, "load 2")
	assert.Len(t, scs, 2, "load 2")
	assert.NotNil(t, findConfById(scs, "ns1/gw1"))
	assert.True(t, findConfById(scs, "ns1/gw1").DeepEqual(sc1), "deepeq")
	assert.NotNil(t, findConfById(scs, "ns1/gw2"))
	assert.True(t, findConfById(scs, "ns1/gw2").DeepEqual(sc3), "deepeq")

	// ns2 client should yield 1 config
	scs, err = client3.Get(clientCtx)
	assert.NoError(t, err, "load 3")
	assert.Len(t, scs, 1, "load 3")
	assert.True(t, scs[0].DeepEqual(sc2), "deepeq")

	// ns1/gw1 client should yield 1 config
	scs, err = client4.Get(clientCtx)
	assert.NoError(t, err, "load 4")
	assert.Len(t, scs, 1, "load 4")
	assert.True(t, scs[0].DeepEqual(sc1), "deepeq")

	// 2 configs from client1 watch
	s1 = watchConfig(ch1, 1500*time.Millisecond)
	assert.NotNil(t, s1)
	s2 = watchConfig(ch1, 150*time.Millisecond)
	assert.NotNil(t, s2)
	s3 = watchConfig(ch1, 150*time.Millisecond)
	assert.Nil(t, s3)
	lst = []*stnrv1.StunnerConfig{s1, s2}
	assert.NotNil(t, findConfById(lst, "ns1/gw1"))
	assert.True(t, findConfById(lst, "ns1/gw1").DeepEqual(sc1), "deepeq")
	assert.NotNil(t, findConfById(lst, "ns1/gw2"))
	assert.True(t, findConfById(lst, "ns1/gw2").DeepEqual(sc3), "deepeq")

	// 2 configs from client2 watch
	s1 = watchConfig(ch2, 1500*time.Millisecond)
	assert.NotNil(t, s1)
	s2 = watchConfig(ch2, 150*time.Millisecond)
	assert.NotNil(t, s2)
	s3 = watchConfig(ch2, 50*time.Millisecond)
	assert.Nil(t, s3)
	lst = []*stnrv1.StunnerConfig{s1, s2}
	assert.NotNil(t, findConfById(lst, "ns1/gw1"))
	assert.True(t, findConfById(lst, "ns1/gw1").DeepEqual(sc1), "deepeq")
	assert.NotNil(t, findConfById(lst, "ns1/gw2"))
	assert.True(t, findConfById(lst, "ns1/gw2").DeepEqual(sc3), "deepeq")

	// 0 config from client3 watch
	s = watchConfig(ch3, 50*time.Millisecond)
	assert.Nil(t, s, "config 3")

	// 1 config from client4 watch
	s = watchConfig(ch4, 50*time.Millisecond)
	assert.NotNil(t, s)
	assert.True(t, s.DeepEqual(sc1), "deepeq")

	testLog.Debug("--------------------------------")
	testLog.Debug("Restart + Update1: ns1/gw1 + ns2/gw1 + ns1/gw2")
	testLog.Debug("--------------------------------")
	testLog.Debug("restarting server")
	serverCancel()
	// let the server shut down and restart
	time.Sleep(50 * time.Millisecond)
	serverCtx, serverCancel = context.WithCancel(context.Background())
	defer serverCancel()
	server = New(stnrv1.DefaultConfigDiscoveryAddress, log)
	assert.NotNil(t, server, "server")
	err = server.Start(serverCtx)
	assert.NoError(t, err, "start")
	err = server.UpdateConfig([]Config{c1, c2, c3})
	assert.NoError(t, err, "update")

	cs = server.configs.Snapshot()
	assert.Len(t, cs, 3, "snapshot len")
	sc1 = server.configs.Get("ns1/gw1")
	assert.NotNil(t, sc1, "get 1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq")
	sc2 = server.configs.Get("ns2/gw1")
	assert.NotNil(t, sc2, "get 2")
	assert.True(t, c2.Config.DeepEqual(sc2), "deepeq")
	sc3 = server.configs.Get("ns1/gw2")
	assert.NotNil(t, sc3, "get 3")
	assert.True(t, c3.Config.DeepEqual(sc3), "deepeq")

	// all-configs should result sc1 and sc2 and sc3
	scs, err = client1.Get(clientCtx)
	assert.NoError(t, err, "load 1")
	assert.Len(t, scs, 3, "load 1")
	co = findConfById(scs, "ns1/gw1")
	assert.NotNil(t, co, "c1")
	assert.True(t, co.DeepEqual(sc1), "deepeq")
	co = findConfById(scs, "ns2/gw1")
	assert.NotNil(t, co, "c2")
	assert.True(t, co.DeepEqual(sc2), "deepeq")
	co = findConfById(scs, "ns1/gw2")
	assert.NotNil(t, co, "c3")
	assert.True(t, co.DeepEqual(sc3), "deepeq")

	// ns1 client should yield 2 configs
	scs, err = client2.Get(clientCtx)
	assert.NoError(t, err, "load 2")
	assert.Len(t, scs, 2, "load 2")
	assert.NotNil(t, findConfById(scs, "ns1/gw1"))
	assert.True(t, findConfById(scs, "ns1/gw1").DeepEqual(sc1), "deepeq")
	assert.NotNil(t, findConfById(scs, "ns1/gw2"))
	assert.True(t, findConfById(scs, "ns1/gw2").DeepEqual(sc3), "deepeq")

	// ns2 client should yield 1 config
	scs, err = client3.Get(clientCtx)
	assert.NoError(t, err, "load 3")
	assert.Len(t, scs, 1, "load 3")
	assert.True(t, scs[0].DeepEqual(sc2), "deepeq")

	// ns1/gw1 client should yield 1 config
	scs, err = client4.Get(clientCtx)
	assert.NoError(t, err, "load 4")
	assert.Len(t, scs, 1, "load 4")
	assert.True(t, scs[0].DeepEqual(sc1), "deepeq")

	// 3 configs from client1 watch
	s1 = watchConfig(ch1, 5000*time.Millisecond)
	assert.NotNil(t, s1)
	s2 = watchConfig(ch1, 100*time.Millisecond)
	assert.NotNil(t, s2)
	s3 = watchConfig(ch1, 100*time.Millisecond)
	assert.NotNil(t, s2)
	s4 := watchConfig(ch1, 100*time.Millisecond)
	assert.Nil(t, s4)
	lst = []*stnrv1.StunnerConfig{s1, s2, s3}
	assert.NotNil(t, findConfById(lst, "ns1/gw1"))
	assert.True(t, findConfById(lst, "ns1/gw1").DeepEqual(sc1), "deepeq")
	assert.NotNil(t, findConfById(lst, "ns1/gw2"))
	assert.True(t, findConfById(lst, "ns2/gw1").DeepEqual(sc2), "deepeq")
	assert.NotNil(t, findConfById(lst, "ns2/gw1"))
	assert.True(t, findConfById(lst, "ns1/gw2").DeepEqual(sc3), "deepeq")

	// 2 configs from client2 watch
	s1 = watchConfig(ch2, 50*time.Millisecond)
	assert.NotNil(t, s1)
	s2 = watchConfig(ch2, 50*time.Millisecond)
	assert.NotNil(t, s2)
	s3 = watchConfig(ch2, 50*time.Millisecond)
	assert.Nil(t, s3)
	lst = []*stnrv1.StunnerConfig{s1, s2}
	assert.NotNil(t, findConfById(lst, "ns1/gw1"))
	assert.True(t, findConfById(lst, "ns1/gw1").DeepEqual(sc1), "deepeq")
	assert.NotNil(t, findConfById(lst, "ns1/gw2"))
	assert.True(t, findConfById(lst, "ns1/gw2").DeepEqual(sc3), "deepeq")

	// 1 config from client3 watch
	s = watchConfig(ch3, 50*time.Millisecond)
	assert.NotNil(t, s, "config 3")
	assert.True(t, s.DeepEqual(sc2))
	s = watchConfig(ch3, 50*time.Millisecond)
	assert.Nil(t, s)

	// 1 config from client4 watch
	s = watchConfig(ch4, 50*time.Millisecond)
	assert.NotNil(t, s)
	assert.True(t, s.DeepEqual(sc1))
	s = watchConfig(ch4, 50*time.Millisecond)
	assert.Nil(t, s)

	testLog.Debug("--------------------------------")
	testLog.Debug("Update1: ns1/gw1 + ns3/gw1      ")
	testLog.Debug("--------------------------------")
	testLog.Debug("update: conf 1, remove conf 3, and add conf 4")
	c1 = testConfig("ns1/gw1", "realm-newer")
	c4 := testConfig("ns3/gw1", "realm4")
	err = server.UpdateConfig([]Config{c1, c2, c4})
	assert.NoError(t, err, "update")

	cs = server.configs.Snapshot()
	assert.Len(t, cs, 3, "snapshot len")
	sc1 = server.configs.Get("ns1/gw1")
	assert.NotNil(t, sc1, "get 1")
	assert.True(t, c1.Config.DeepEqual(sc1), "deepeq")
	sc2 = server.configs.Get("ns2/gw1")
	assert.NotNil(t, sc2, "get 2")
	assert.True(t, c2.Config.DeepEqual(sc2), "deepeq")
	sc4 := server.configs.Get("ns3/gw1")
	assert.NotNil(t, sc3, "get 3")
	assert.True(t, c4.Config.DeepEqual(sc4), "deepeq")

	// all-configs should result sc1 and sc2 and sc4
	scs, err = client1.Get(clientCtx)
	assert.NoError(t, err, "load 1")
	assert.Len(t, scs, 3, "load 1")
	co = findConfById(scs, "ns1/gw1")
	assert.NotNil(t, co, "c1")
	assert.True(t, co.DeepEqual(sc1), "deepeq")
	co = findConfById(scs, "ns2/gw1")
	assert.NotNil(t, co, "c2")
	assert.True(t, co.DeepEqual(sc2), "deepeq")
	co = findConfById(scs, "ns3/gw1")
	assert.NotNil(t, co, "c4")
	assert.True(t, co.DeepEqual(sc4), "deepeq")

	// ns1 client should yield 1 config
	scs, err = client2.Get(clientCtx)
	assert.NoError(t, err, "load 2")
	assert.Len(t, scs, 1, "load 2")
	assert.True(t, scs[0].DeepEqual(sc1), "deepeq")

	// ns2 client should yield 1 config
	scs, err = client3.Get(clientCtx)
	assert.NoError(t, err, "load 3")
	assert.Len(t, scs, 1, "load 3")
	assert.True(t, scs[0].DeepEqual(sc2), "deepeq")

	// ns1/gw1 client should yield 1 config
	scs, err = client4.Get(clientCtx)
	assert.NoError(t, err, "load 4")
	assert.Len(t, scs, 1, "load 4")
	assert.True(t, scs[0].DeepEqual(sc1), "deepeq")

	// 2 configs from client1 watch
	s1 = watchConfig(ch1, 5000*time.Millisecond)
	assert.NotNil(t, s1)
	s2 = watchConfig(ch1, 500*time.Millisecond)
	assert.NotNil(t, s2)
	s3 = watchConfig(ch1, 500*time.Millisecond)
	assert.Nil(t, s3)
	lst = []*stnrv1.StunnerConfig{s1, s2}
	assert.NotNil(t, findConfById(lst, "ns1/gw1"))
	assert.True(t, findConfById(lst, "ns1/gw1").DeepEqual(sc1), "deepeq")
	assert.NotNil(t, findConfById(lst, "ns3/gw1"))
	assert.True(t, findConfById(lst, "ns3/gw1").DeepEqual(sc4), "deepeq")

	// 1 config from client2 watch (removed config never updated)
	s1 = watchConfig(ch2, 50*time.Millisecond)
	assert.NotNil(t, s1)
	s2 = watchConfig(ch2, 50*time.Millisecond)
	assert.Nil(t, s2)
	assert.True(t, s1.DeepEqual(sc1), "deepeq")

	// no config from client3 watch
	s = watchConfig(ch3, 50*time.Millisecond)
	assert.Nil(t, s, "config 3")

	// 1 config from client4 watch
	s = watchConfig(ch4, 50*time.Millisecond)
	assert.NotNil(t, s)
	assert.True(t, s.DeepEqual(sc1), "deepeq")
}

// only differ in id and realm
func testConfig(id, realm string) Config {
	c := client.ZeroConfig(id)
	c.Auth.Realm = realm

	return Config{id, c}
}

// wait for some configurable time for a watch element
func watchConfig(ch chan stnrv1.StunnerConfig, d time.Duration) *stnrv1.StunnerConfig {
	select {
	case c := <-ch:
		// fmt.Println("++++++++++++ got config ++++++++++++: ", c.String())
		return &c
	case <-time.After(d):
		// fmt.Println("++++++++++++ timeout ++++++++++++")
		return nil
	}
}

func findConfById(cs []*stnrv1.StunnerConfig, id string) *stnrv1.StunnerConfig {
	for _, c := range cs {
		if c != nil && c.Admin.Name == id {
			return c
		}

	}

	return nil
}
