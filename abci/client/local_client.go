package abcicli

import (
	"sync"
	"time"

	types "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/service"
)

var _ Client = (*localClient)(nil)

type LocalClient interface {
	Client
	WithMetrics(metrics *Metrics) LocalClient
}

// NOTE: use defer to unlock mutex because Application might panic (e.g., in
// case of malicious tx or query). It only makes sense for publicly exposed
// methods like CheckTx (/broadcast_tx_* RPC endpoint) or Query (/abci_query
// RPC endpoint), but defers are used everywhere for the sake of consistency.
type localClient struct {
	service.BaseService

	mtx *sync.Mutex
	types.Application
	Callback
	metrics *Metrics
}

func NewLocalClient(mtx *sync.Mutex, app types.Application) LocalClient {
	if mtx == nil {
		mtx = new(sync.Mutex)
	}
	cli := &localClient{
		mtx:         mtx,
		Application: app,
		metrics:     NopMetrics(),
	}
	cli.BaseService = *service.NewBaseService(nil, "localClient", cli)
	return cli
}

func (app *localClient) WithMetrics(metrics *Metrics) LocalClient {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	app.metrics = metrics
	return app
}

func (app *localClient) SetResponseCallback(cb Callback) {
	defer app.leave(app.enter("SetResponseCallback"))
	defer app.mtx.Unlock()
	app.Callback = cb
}

// TODO: change types.Application to include Error()?
func (app *localClient) Error() error {
	return nil
}

func (app *localClient) FlushAsync() *ReqRes {
	// Do nothing
	return newLocalReqRes(types.ToRequestFlush(), nil)
}

func (app *localClient) EchoAsync(msg string) *ReqRes {
	defer app.leave(app.enter("EchoAsync"))
	defer app.mtx.Unlock()

	return app.callback(
		types.ToRequestEcho(msg),
		types.ToResponseEcho(msg),
	)
}

func (app *localClient) InfoAsync(req types.RequestInfo) *ReqRes {
	defer app.leave(app.enter("InfoAsync"))
	defer app.mtx.Unlock()

	res := app.Application.Info(req)
	return app.callback(
		types.ToRequestInfo(req),
		types.ToResponseInfo(res),
	)
}

func (app *localClient) SetOptionAsync(req types.RequestSetOption) *ReqRes {
	defer app.leave(app.enter("SetOptionsAsync"))
	defer app.mtx.Unlock()

	res := app.Application.SetOption(req)
	return app.callback(
		types.ToRequestSetOption(req),
		types.ToResponseSetOption(res),
	)
}

func (app *localClient) DeliverTxAsync(params types.RequestDeliverTx) *ReqRes {
	defer app.leave(app.enter("DeliverTxAsync"))
	defer app.mtx.Unlock()

	res := app.Application.DeliverTx(params)
	return app.callback(
		types.ToRequestDeliverTx(params),
		types.ToResponseDeliverTx(res),
	)
}

func (app *localClient) CheckTxAsync(req types.RequestCheckTx) *ReqRes {
	defer app.leave(app.enter("CheckTxAsync"))
	defer app.mtx.Unlock()

	res := app.Application.CheckTx(req)
	return app.callback(
		types.ToRequestCheckTx(req),
		types.ToResponseCheckTx(res),
	)
}

func (app *localClient) QueryAsync(req types.RequestQuery) *ReqRes {
	defer app.leave(app.enter("QueryAsync"))
	defer app.mtx.Unlock()

	res := app.Application.Query(req)
	return app.callback(
		types.ToRequestQuery(req),
		types.ToResponseQuery(res),
	)
}

func (app *localClient) CommitAsync() *ReqRes {
	defer app.leave(app.enter("CommitAsync"))
	defer app.mtx.Unlock()

	res := app.Application.Commit()
	return app.callback(
		types.ToRequestCommit(),
		types.ToResponseCommit(res),
	)
}

func (app *localClient) InitChainAsync(req types.RequestInitChain) *ReqRes {
	defer app.leave(app.enter("InitChainAsync"))
	defer app.mtx.Unlock()

	res := app.Application.InitChain(req)
	return app.callback(
		types.ToRequestInitChain(req),
		types.ToResponseInitChain(res),
	)
}

func (app *localClient) BeginBlockAsync(req types.RequestBeginBlock) *ReqRes {
	defer app.leave(app.enter("BeginBlockAsync"))
	defer app.mtx.Unlock()

	res := app.Application.BeginBlock(req)
	return app.callback(
		types.ToRequestBeginBlock(req),
		types.ToResponseBeginBlock(res),
	)
}

func (app *localClient) EndBlockAsync(req types.RequestEndBlock) *ReqRes {
	defer app.leave(app.enter("EndBlockAsync"))
	defer app.mtx.Unlock()

	res := app.Application.EndBlock(req)
	return app.callback(
		types.ToRequestEndBlock(req),
		types.ToResponseEndBlock(res),
	)
}

//-------------------------------------------------------

func (app *localClient) FlushSync() error {
	return nil
}

func (app *localClient) EchoSync(msg string) (*types.ResponseEcho, error) {
	return &types.ResponseEcho{Message: msg}, nil
}

func (app *localClient) InfoSync(req types.RequestInfo) (*types.ResponseInfo, error) {
	defer app.leave(app.enter("InfoSync"))
	defer app.mtx.Unlock()

	res := app.Application.Info(req)
	return &res, nil
}

func (app *localClient) SetOptionSync(req types.RequestSetOption) (*types.ResponseSetOption, error) {
	defer app.leave(app.enter("SetOptionSync"))
	defer app.mtx.Unlock()

	res := app.Application.SetOption(req)
	return &res, nil
}

func (app *localClient) DeliverTxSync(req types.RequestDeliverTx) (*types.ResponseDeliverTx, error) {
	defer app.leave(app.enter("DeliverTxSync"))
	defer app.mtx.Unlock()

	res := app.Application.DeliverTx(req)
	return &res, nil
}

func (app *localClient) CheckTxSync(req types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	defer app.leave(app.enter("CheckTxSync"))
	defer app.mtx.Unlock()

	res := app.Application.CheckTx(req)
	return &res, nil
}

func (app *localClient) QuerySync(req types.RequestQuery) (*types.ResponseQuery, error) {
	defer app.leave(app.enter("QuerySync"))
	defer app.mtx.Unlock()

	res := app.Application.Query(req)
	return &res, nil
}

func (app *localClient) CommitSync() (*types.ResponseCommit, error) {
	defer app.leave(app.enter("CommitSync"))
	defer app.mtx.Unlock()

	res := app.Application.Commit()
	return &res, nil
}

func (app *localClient) InitChainSync(req types.RequestInitChain) (*types.ResponseInitChain, error) {
	defer app.leave(app.enter("InitChainSync"))
	defer app.mtx.Unlock()

	res := app.Application.InitChain(req)
	return &res, nil
}

func (app *localClient) BeginBlockSync(req types.RequestBeginBlock) (*types.ResponseBeginBlock, error) {
	defer app.leave(app.enter("BeginBlockSync"))
	defer app.mtx.Unlock()

	res := app.Application.BeginBlock(req)
	return &res, nil
}

func (app *localClient) EndBlockSync(req types.RequestEndBlock) (*types.ResponseEndBlock, error) {
	defer app.leave(app.enter("EndBlockSync"))
	defer app.mtx.Unlock()

	res := app.Application.EndBlock(req)
	return &res, nil
}

//-------------------------------------------------------

func (app *localClient) callback(req *types.Request, res *types.Response) *ReqRes {
	app.Callback(req, res)
	return newLocalReqRes(req, res)
}

func (app *localClient) enter(method string) (string, time.Time, time.Time) {
	start := time.Now()
	app.mtx.Lock()
	lockedAt := time.Now()
	delta := lockedAt.Sub(start) / time.Millisecond
	app.metrics.LockWaitDuration.
		WithLabelValues(metricsMethodKey, method).
		Observe(float64(delta))
	return method, start, lockedAt
}

func (app *localClient) leave(method string, start time.Time, lockedAt time.Time) {
	now := time.Now()

	delta := lockedAt.Sub(now) / time.Millisecond
	app.metrics.UnlockedDuration.
		WithLabelValues(metricsMethodKey, method).
		Observe(float64(delta))

	delta = start.Sub(now) / time.Millisecond
	app.metrics.TotalDuration.
		WithLabelValues(metricsMethodKey, method).
		Observe(float64(delta))
}

func newLocalReqRes(req *types.Request, res *types.Response) *ReqRes {
	reqRes := NewReqRes(req)
	reqRes.Response = res
	reqRes.SetDone()
	return reqRes
}
