package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	traceconfig "binTest/jaegerTest/CSJaeger/tracelib"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var (
	tracer opentracing.Tracer
)

func GetListProc(w http.ResponseWriter, req *http.Request) {

	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	span := tracer.StartSpan("GetListProc", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	serialNum := span.BaggageItem("serialno")

	respList := []string{"l1", "l2", "l3", "l4", "l5"}
	respString := serialNum + ":"

	for _, v := range respList {
		respString += v + ","
	}

	fmt.Println(respString)
	io.WriteString(w, respString)
}

func GetResultProc(w http.ResponseWriter, req *http.Request) {

	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	span := tracer.StartSpan("GetResultProc", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	serialNum := span.BaggageItem("serialno")

	keys, ok := req.URL.Query()["num"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("No request parameter 'num' error! ")
		return
	}

	num, err := strconv.Atoi(keys[0])
	if err != nil {
		fmt.Println("num invalidate")
		return
	}

	result := 0

	for i := 0; i < num; i++ {
		result += i
	}

	respString := fmt.Sprintf("Result(%s):%d", serialNum, result)

	fmt.Println(respString)
	io.WriteString(w, respString)
}

func main() {
	var closer io.Closer
	tracer, closer = traceconfig.TraceInit("Trace-Server", "const", 1)
	defer closer.Close()

	http.HandleFunc("/getList", GetListProc)
	http.HandleFunc("/getResult", GetResultProc)

	http.ListenAndServe(":8080", nil)
}
