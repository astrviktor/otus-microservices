package main

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
)

const (
	Host = ""
	Port = 8001
)

var balance int64 = 0

func main() {
	r := router.New()

	r.GET("/health/", HandleHealth)

	r.GET("/payment/balance", HandleGetBalance)
	r.POST("/payment/credit", HandlePaymentCredit)
	r.POST("/payment/debit", HandlePaymentDebit)

	srv := &fasthttp.Server{
		Handler: r.Handler,
	}

	addr := fmt.Sprintf("%s:%d", Host, Port)
	log.Println("http server starting on address: " + addr)

	if err := srv.ListenAndServe(addr); err != nil {
		log.Fatal("error ListenAndServe()")
	}
}

func WriteResponse(ctx *fasthttp.RequestCtx, resp interface{}) {
	respBuf, err := json.Marshal(resp)
	if err != nil {
		log.Println(fmt.Sprintf("response marshal error: %s", err))
	}

	respBuf = append(respBuf, []byte("\n")...)
	ctx.Response.SetBody(respBuf)

	ctx.SetContentType("application/json; charset=utf-8")
}

type Response struct {
	Status string `json:"status"`
}

type RequestAmount struct {
	Amount int64 `json:"amount"`
}

type ResponseBalance struct {
	Balance int64 `json:"balance"`
}

func HandleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &Response{Status: "OK"})

	return
}

func HandleGetBalance(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseBalance{Balance: balance})

	return
}

func HandlePaymentCredit(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		_, _ = ctx.WriteString("body is empty\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()
	amount := RequestAmount{}

	if err := json.Unmarshal(body, &amount); err != nil {
		_, _ = ctx.WriteString("error in body\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	balance += amount.Amount

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseBalance{Balance: balance})

	return
}

func HandlePaymentDebit(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.Body()) == 0 {
		_, _ = ctx.WriteString("body is empty\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body := ctx.Request.Body()
	amount := RequestAmount{}

	if err := json.Unmarshal(body, &amount); err != nil {
		_, _ = ctx.WriteString("error in body\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if balance < amount.Amount {
		_, _ = ctx.WriteString("balance < amount\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	balance -= amount.Amount

	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &ResponseBalance{Balance: balance})

	return
}
