package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.uber.org/zap"
	"log"
	"otus-microservices/orderservice/internal/config"
	"otus-microservices/orderservice/internal/storage"
	storagegorm "otus-microservices/orderservice/internal/storage/gorm"
	"strconv"
	"time"
)

type Handler struct {
	storage storage.Storage
	cfg     config.Config
	log     *zap.Logger
}

func New(cfg config.Config, log *zap.Logger) (*Handler, error) {

	//storage := storagememory.New(cfg.Storage)
	storage := storagegorm.New(cfg.Storage)

	err := storage.Connect()
	if err != nil {
		return nil, err
	}

	return &Handler{
		storage: storage,
		cfg:     cfg,
		log:     log,
	}, nil
}

type ResponseHealth struct {
	Status string `json:"status"`
}

type ResponseOrderID struct {
	ID int64 `json:"id"`
}

func WriteResponse(ctx *fasthttp.RequestCtx, resp interface{}) {
	respBuf, err := json.Marshal(resp)
	if err != nil {
		log.Println(fmt.Sprintf("response marshal error: %s", err))
	}

	respBuf = append(respBuf, []byte("\n")...)
	ctx.Response.SetBody(respBuf)

	//if err != nil {
	//	log.Println(fmt.Sprintf("response write error: %s", err))
	//}

	ctx.SetContentType("application/json; charset=utf-8")
}

func (h *Handler) PrometheusHandler() fasthttp.RequestHandler {
	return fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
}

func (h *Handler) Metrics(ctx *fasthttp.RequestCtx) {
	metrics.WritePrometheus(ctx.Response.BodyWriter(), true)
}

func (h *Handler) HandleHealth(ctx *fasthttp.RequestCtx) {
	time.Sleep(time.Duration(100) * time.Millisecond)

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseHealth{Status: "OK"})

	return
}

func (h *Handler) CreateOrder(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		_, _ = ctx.WriteString("body is empty")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()
	order := storage.Order{}

	if err := json.Unmarshal(body, &order); err != nil {
		_, _ = ctx.WriteString("error in body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// Idempotency
	const headerRequestId = "X-Request-ID"
	requestID := string(ctx.Request.Header.Peek(headerRequestId))
	if requestID == "" {
		requestID = uuid.New().String()
		ctx.Response.Header.Add(headerRequestId, requestID)
	}

	order.RequestId = requestID

	id, err := h.storage.CreateOrder(order)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseOrderID{
		ID: id,
	})

	return
}

type RequestAmount struct {
	Amount int64 `json:"amount"`
}

func (h *Handler) CreateOrderSaga(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		_, _ = ctx.WriteString("body is empty")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()
	order := storage.Order{}

	if err := json.Unmarshal(body, &order); err != nil {
		_, _ = ctx.WriteString("error in body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// payment
	paymentErrorChan := make(chan error)

	go func() {
		paymentClient := new(fasthttp.Client)

		paymentRequest := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(paymentRequest)

		uri := fmt.Sprintf("http://%s:%d/%s", h.cfg.Payment.Host, h.cfg.Payment.Port, "payment/debit")
		h.log.Info("uri", zap.String("uri", uri))

		paymentRequest.SetRequestURI(uri)
		paymentRequest.Header.SetMethod("POST")

		paymentRequest.Header.Add("content-type", "application/json")
		data, err := json.Marshal(&RequestAmount{Amount: order.Total})
		if err != nil {
			paymentErrorChan <- err
			return
		}
		paymentRequest.SetBody(data)

		paymentResponse := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(paymentResponse)

		err = paymentClient.Do(paymentRequest, paymentResponse)
		if err != nil {
			paymentErrorChan <- err
			return
		}

		if paymentResponse.StatusCode() != fasthttp.StatusOK {
			paymentErrorChan <- fmt.Errorf("payment status not OK: %v\n", paymentResponse.StatusCode())
			return
		}

		paymentErrorChan <- nil
	}()

	// warehouse
	go func() {
		h.log.Info("warehouse process")
	}()

	// delivery
	go func() {
		h.log.Info("delivery process")
	}()

	paymentErr := <-paymentErrorChan
	if paymentErr != nil {
		_, _ = ctx.WriteString(paymentErr.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// Idempotency
	const headerRequestId = "X-Request-ID"
	requestID := string(ctx.Request.Header.Peek(headerRequestId))
	if requestID == "" {
		requestID = uuid.New().String()
		ctx.Response.Header.Add(headerRequestId, requestID)
	}

	order.RequestId = requestID

	id, err := h.storage.CreateOrder(order)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseOrderID{
		ID: id,
	})

	return
}

func (h *Handler) ReadOrder(ctx *fasthttp.RequestCtx) {
	stringID, ok := ctx.UserValue("id").(string)
	if !ok {
		_, _ = ctx.WriteString("id is wrong in request path")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		_, _ = ctx.WriteString("id is not int type")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	user, err := h.storage.ReadOrder(id)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &user)
	return
}

func (h *Handler) UpdateOrder(ctx *fasthttp.RequestCtx) {
	stringID, ok := ctx.UserValue("id").(string)
	if !ok {
		_, _ = ctx.WriteString("id is wrong in request path")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		_, _ = ctx.WriteString("id is not int type")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if len(ctx.Request.Body()) == 0 {
		_, _ = ctx.WriteString("body is empty")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	order, err := h.storage.ReadOrder(id)
	if err != nil {
		_, _ = ctx.WriteString("order not found")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()

	if err := json.Unmarshal(body, &order); err != nil {
		_, _ = ctx.WriteString("error in body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err = h.storage.UpdateOrder(id, order)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseOrderID{
		ID: id,
	})

	return
}

func (h *Handler) DeleteUser(ctx *fasthttp.RequestCtx) {
	stringID, ok := ctx.UserValue("id").(string)
	if !ok {
		_, _ = ctx.WriteString("id is wrong in request path")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		_, _ = ctx.WriteString("id is not int type")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err = h.storage.DeleteOrder(id)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseOrderID{
		ID: id,
	})
	return
}
