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
	Port = 8003
)

func main() {
	r := router.New()

	r.GET("/health/", HandleHealth)
	r.POST("/delivery/create", HandleCreateDelivery)
	r.POST("/delivery/delete", HandleDeleteDelivery)

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

func HandleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &Response{Status: "OK"})

	return
}

func HandleCreateDelivery(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &Response{Status: "OK"})

	return
}

func HandleDeleteDelivery(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	WriteResponse(ctx, &Response{Status: "OK"})

	return
}
