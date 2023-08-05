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
	"otus-microservices/notification/internal/broker"
	"otus-microservices/notification/internal/config"
	"otus-microservices/notification/internal/storage"
	storagegorm "otus-microservices/notification/internal/storage/gorm"
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

func (h *Handler) ReadNotification(ctx *fasthttp.RequestCtx) {
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

	notification, err := h.storage.ReadNotification(id)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &notification)
	return
}

func (h *Handler) ProcessMessages(broker broker.InterfaceBroker) {
	for {
		message, err := broker.Receive()
		if err != nil {
			h.log.Error(fmt.Sprintf("error get message from brocker: %s", err.Error()))
			continue
		}

		if message != nil {
			notification := storage.Notification{
				ClientID: message.ClientID,
				OrderID:  message.OrderID,
				Theme:    message.Theme,
				Message:  message.Message,
			}

			err = h.storage.CreateNotification(notification)
			if err != nil {
				h.log.Error(fmt.Sprintf("error create notification in storage: %s", err.Error()))
			}
		}
	}
}
