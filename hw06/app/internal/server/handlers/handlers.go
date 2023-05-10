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
	"math/rand"
	"otus-microservices/hw06/internal/config"
	"otus-microservices/hw06/internal/storage"
	"otus-microservices/hw06/internal/storage/memory"
	"strconv"
	"time"
)

type Handler struct {
	storage storage.Storage
	log     *zap.Logger
}

func New(cfg config.Config, log *zap.Logger) (*Handler, error) {

	storage := storagememory.New(cfg.Storage)
	//storage := storagegorm.New(cfg.Storage)

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

type ResponseProfileID struct {
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

func (h *Handler) CreateProfile(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		_, _ = ctx.WriteString("body is empty")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()
	Profile := storage.Profile{}

	if err := json.Unmarshal(body, &Profile); err != nil {
		_, _ = ctx.WriteString("error in body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id, err := h.storage.CreateProfile(Profile)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseProfileID{
		ID: id,
	})

	return
}

func (h *Handler) ReadProfile(ctx *fasthttp.RequestCtx) {
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

	Profile, err := h.storage.ReadProfile(id)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &Profile)
	return
}

func (h *Handler) UpdateProfile(ctx *fasthttp.RequestCtx) {
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

	Profile, err := h.storage.ReadProfile(id)
	if err != nil {
		_, _ = ctx.WriteString("Profile not found")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()

	if err := json.Unmarshal(body, &Profile); err != nil {
		_, _ = ctx.WriteString("error in body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err = h.storage.UpdateProfile(id, Profile)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseProfileID{
		ID: id,
	})

	return
}

func (h *Handler) DeleteProfile(ctx *fasthttp.RequestCtx) {
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

	err = h.storage.DeleteProfile(id)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseProfileID{
		ID: id,
	})
	return
}

func (h *Handler) Auth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		sessionID := ctx.Request.Header.Cookie("session_id")
		if len(sessionID) == 0 {
			_, _ = ctx.WriteString("Need to login")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		profile, err := h.storage.GetProfileForSession(string(sessionID))
		if err != nil {
			_, _ = ctx.WriteString("Profile not found")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		pathID, ok := ctx.UserValue("id").(string)
		if !ok {
			_, _ = ctx.WriteString("id is wrong in request path")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(pathID, 10, 64)
		if err != nil {
			_, _ = ctx.WriteString("id is not int type")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		if profile.ProfileId != id {
			_, _ = ctx.WriteString("Wrong credentions")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		next(ctx)
	}
}

func (h *Handler) Login(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		_, _ = ctx.WriteString("body is empty")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	profile := storage.Profile{}

	body := ctx.Request.Body()

	if err := json.Unmarshal(body, &profile); err != nil {
		_, _ = ctx.WriteString("error in body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	sessionID, err := h.storage.SetSessionForProfile(profile.Username)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	sessionCookie := fasthttp.Cookie{}
	sessionCookie.SetKey("session_id")
	sessionCookie.SetValue(sessionID)
	sessionCookie.SetMaxAge(3600)
	sessionCookie.SetDomain("arch.homework")
	//sessionCookie.SetPath(("/"))

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetCookie(&sessionCookie)

	return
}

func (h *Handler) Logout(ctx *fasthttp.RequestCtx) {
	sessionID := ctx.Request.Header.Cookie("session_id")
	if len(sessionID) == 0 {
		_, _ = ctx.WriteString("Need to login")
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	profile, err := h.storage.GetProfileForSession(string(sessionID))
	if err != nil {
		_, _ = ctx.WriteString("Profile not found")
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	err = h.storage.ClearSessionForProfileId(profile.ProfileId)
	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	sessionCookie := fasthttp.Cookie{}
	sessionCookie.SetKey("session_id")
	sessionCookie.SetMaxAge(-1)
	sessionCookie.SetDomain("arch.homework")
	//sessionCookie.SetPath(("/"))

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetCookie(&sessionCookie)

	return
}
