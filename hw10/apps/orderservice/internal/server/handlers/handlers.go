package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.uber.org/zap"
	"log"
	"otus-microservices/orderservice/internal/broker"
	"otus-microservices/orderservice/internal/broker/rabbitmq"
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
	Broker  broker.InterfaceBroker
}

func New(cfg config.Config, log *zap.Logger) (*Handler, error) {

	//storage := storagememory.New(cfg.Storage)
	storage := storagegorm.New(cfg.Storage)

	broker := rabbitmq.New(cfg.Rabbitmq)

	err := storage.Connect()
	if err != nil {
		return nil, err
	}

	return &Handler{
		storage: storage,
		cfg:     cfg,
		log:     log,
		Broker:  broker,
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

type RequestBilling struct {
	ClientID int64 `json:"client_id,omitempty"`
	Balance  int64 `json:"balance,omitempty"`
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

	// billing
	billingErrorChan := make(chan error)

	go func() {
		billingClient := new(fasthttp.Client)

		billingRequest := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(billingRequest)

		uri := fmt.Sprintf("http://%s:%d/billing/%d/debit", h.cfg.Billing.Host, h.cfg.Billing.Port, order.ClientID)
		h.log.Info("uri", zap.String("uri", uri))

		billingRequest.SetRequestURI(uri)
		billingRequest.Header.SetMethod("POST")

		billingRequest.Header.Add("content-type", "application/json")
		data, err := json.Marshal(&RequestBilling{
			ClientID: order.ClientID,
			Balance:  order.Total,
		})

		if err != nil {
			billingErrorChan <- err
			return
		}
		billingRequest.SetBody(data)

		billingResponse := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(billingResponse)

		err = billingClient.Do(billingRequest, billingResponse)
		if err != nil {
			billingErrorChan <- err
			return
		}

		if billingResponse.StatusCode() != fasthttp.StatusOK {
			billingErrorChan <- fmt.Errorf("billing status not OK: %v\n", billingResponse.StatusCode())
			return
		}

		billingErrorChan <- nil
	}()

	billingErr := <-billingErrorChan
	message := broker.Message{
		ClientID: order.ClientID,
	}

	if billingErr != nil {
		message.OrderID = -1
		message.Theme = "Bad"
		message.Message = "Bad"

		err := h.Broker.Send(&message)
		if err != nil {
			h.log.Error("error sending RMQ message")
		}

		_, _ = ctx.WriteString(billingErr.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id, err := h.storage.CreateOrder(order)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	message.OrderID = id
	message.Theme = "Good"
	message.Message = "Good"

	err = h.Broker.Send(&message)
	if err != nil {
		h.log.Error("error sending RMQ message")
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
