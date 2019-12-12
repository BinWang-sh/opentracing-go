package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	traceconfig "binTest/jaegerTest/CSJaeger/tracelib"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	URL        = "http://localhost:8080"
	LIST_API   = "/getList"
	RESULT_API = "/getResult"
)

var (
	flag = make(chan bool)
)

func saveResponse(response []byte) error {
	err := ioutil.WriteFile("response.txt", response, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func sendRequest(req *http.Request, ctx context.Context) {
	reqPrepareSpan, _ := opentracing.StartSpanFromContext(ctx, "Client_sendRequest")
	defer reqPrepareSpan.Finish()

	go func(req *http.Request) {
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Printf("Do send requst failed(%s)\n", err)
			return
		}

		respSpan, _ := opentracing.StartSpanFromContext(ctx, "Client_response")
		defer respSpan.Finish()

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ReadAll error(%s)\n", err)
			return
		}

		if resp.StatusCode != 200 {
			return
		}

		fmt.Printf("Response:%s\n", string(body))

		respSpan.LogFields(
			log.String("event", "getResponse"),
			log.String("value", string(body)),
		)

		saveResponse(body)

		flag <- true
	}(req)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Argument error(getlist or getresult number) ")
		os.Exit(1)
	}

	tracer, closer := traceconfig.TraceInit("CS-tracing", "const", 1)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span := tracer.StartSpan(fmt.Sprintf("%s trace", os.Args[1]))
	span.SetTag("trace to", os.Args[1])
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	api := ""
	var err error

	if os.Args[1] == "getlist" {
		api = LIST_API
	} else if os.Args[1] == "getresult" {
		api = RESULT_API
		num, err := strconv.Atoi(os.Args[2])

		if err != nil || num <= 0 {
			fmt.Println("getresult input parameter error!")
			os.Exit(1)
		}
	}

	reqSpan, _ := opentracing.StartSpanFromContext(ctx, "Client_"+api+" request")
	defer reqSpan.Finish()

	reqURL := URL + api
	req, err := http.NewRequest("GET", reqURL, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ext.SpanKindRPCClient.Set(reqSpan)
	ext.HTTPUrl.Set(reqSpan, reqURL)
	ext.HTTPMethod.Set(reqSpan, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	if os.Args[1] == "getresult" {
		q := req.URL.Query()
		q.Add("num", os.Args[2])
		req.URL.RawQuery = q.Encode()
	}

	fmt.Println(req.URL.String())
	reqSpan.LogFields(
		log.String("event", api),
		log.String("value", api),
	)

	sendRequest(req, ctx)

	<-flag
}
