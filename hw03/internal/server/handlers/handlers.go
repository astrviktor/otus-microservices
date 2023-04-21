package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"otus-microservices/hw03/internal/config"
	"otus-microservices/hw03/internal/storage"
	storagegorm "otus-microservices/hw03/internal/storage/gorm"
	"strconv"
	"time"
)

type Handler struct {
	storage storage.Storage
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
		log:     log,
	}, nil
}

type ResponseHealth struct {
	Status string `json:"status"`
}

type ResponseUserID struct {
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

func (h *Handler) Metrics(ctx *fasthttp.RequestCtx) {
	metrics.WritePrometheus(ctx.Response.BodyWriter(), true)
}

func (h *Handler) HandleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseHealth{Status: "OK"})

	return
}

func (h *Handler) HandleTesting(ctx *fasthttp.RequestCtx) {
	delay := rand.Intn(800) + 200
	h.log.Info(fmt.Sprintf("[Testing] sleeping %d millisecond", delay))
	time.Sleep(time.Duration(delay) * time.Millisecond)

	code := [...]int{
		fasthttp.StatusOK,
		fasthttp.StatusOK,
		fasthttp.StatusOK,
		fasthttp.StatusOK,
		fasthttp.StatusOK,
		fasthttp.StatusOK,
		fasthttp.StatusBadRequest,
		fasthttp.StatusNotFound,
		fasthttp.StatusNotFound,
		fasthttp.StatusInternalServerError,
	}
	idx := rand.Intn(10)

	ctx.SetStatusCode(code[idx])
	return
}

func (h *Handler) CreateUser(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		_, _ = ctx.WriteString("body is empty")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()
	user := storage.User{}

	if err := json.Unmarshal(body, &user); err != nil {
		_, _ = ctx.WriteString("error in body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id, err := h.storage.CreateUser(user)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseUserID{
		ID: id,
	})

	return
}

func (h *Handler) ReadUser(ctx *fasthttp.RequestCtx) {
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

	user, err := h.storage.ReadUser(id)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &user)
	return
}

func (h *Handler) UpdateUser(ctx *fasthttp.RequestCtx) {
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

	user, err := h.storage.ReadUser(id)
	if err != nil {
		_, _ = ctx.WriteString("user not found")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()

	if err := json.Unmarshal(body, &user); err != nil {
		_, _ = ctx.WriteString("error in body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err = h.storage.UpdateUser(id, user)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseUserID{
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

	err = h.storage.DeleteUser(id)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseUserID{
		ID: id,
	})
	return
}
