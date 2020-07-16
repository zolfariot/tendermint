package proxy

import (
	"sync"

	"github.com/pkg/errors"

	abcicli "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/example/counter"
	"github.com/tendermint/tendermint/abci/example/kvstore"
	"github.com/tendermint/tendermint/abci/types"
)

// NewABCIClient returns newly connected client
type ClientCreator interface {
	NewABCIClient() (abcicli.Client, error)
	WithMetrics(*abcicli.Metrics) ClientCreator
}

//----------------------------------------------------
// local proxy uses a mutex on an in-proc app

type localClientCreator struct {
	mtx     *sync.Mutex
	app     types.Application
	metrics *abcicli.Metrics
}

func NewLocalClientCreator(app types.Application) ClientCreator {
	return &localClientCreator{
		mtx: new(sync.Mutex),
		app: app,
	}
}

func (l *localClientCreator) WithMetrics(metrics *abcicli.Metrics) ClientCreator {
	l.metrics = metrics
	return l
}

func (l *localClientCreator) NewABCIClient() (abcicli.Client, error) {
	client := abcicli.NewLocalClient(l.mtx, l.app)
	if l.metrics != nil {
		client = client.WithMetrics(l.metrics)
	}
	return client, nil
}

//---------------------------------------------------------------
// remote proxy opens new connections to an external app process

type remoteClientCreator struct {
	addr        string
	transport   string
	mustConnect bool
}

func NewRemoteClientCreator(addr, transport string, mustConnect bool) ClientCreator {
	return &remoteClientCreator{
		addr:        addr,
		transport:   transport,
		mustConnect: mustConnect,
	}
}

func (r *remoteClientCreator) WithMetrics(_ *abcicli.Metrics) ClientCreator {
	return r
}

func (r *remoteClientCreator) NewABCIClient() (abcicli.Client, error) {
	remoteApp, err := abcicli.NewClient(r.addr, r.transport, r.mustConnect)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to proxy")
	}
	return remoteApp, nil
}

//-----------------------------------------------------------------
// default

func DefaultClientCreator(addr, transport, dbDir string) ClientCreator {
	switch addr {
	case "counter":
		return NewLocalClientCreator(counter.NewApplication(false))
	case "counter_serial":
		return NewLocalClientCreator(counter.NewApplication(true))
	case "kvstore":
		return NewLocalClientCreator(kvstore.NewApplication())
	case "persistent_kvstore":
		return NewLocalClientCreator(kvstore.NewPersistentKVStoreApplication(dbDir))
	case "noop":
		return NewLocalClientCreator(types.NewBaseApplication())
	default:
		mustConnect := false // loop retrying
		return NewRemoteClientCreator(addr, transport, mustConnect)
	}
}
