An example of using opentracing.

This show how to trace over the process boundaries and RPC calls.

Client side

<1>import

    import (
        "github.com/opentracing/opentracing-go/ext"
    )
    
<2>inject

ext.SpanKindRPCClient.Set(reqSpan)
ext.HTTPUrl.Set(reqSpan, reqURL)
ext.HTTPMethod.Set(reqSpan, "GET")
span.Tracer().Inject(
     span.Context(),
     opentracing.HTTPHeaders,
     opentracing.HTTPHeadersCarrier(req.Header),
)


Server side

<1>import

import (
     opentracing "github.com/opentracing/opentracing-go"
     "github.com/opentracing/opentracing-go/ext"
     otlog "github.com/opentracing/opentracing-go/log"
     "github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

<2>Extrace span context from incoming http request
spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

<3>Creates a ChildOf reference to the passed spanCtx as well as sets a span.kind=server tag on the new span by using a special option RPCServerOption
    span := tracer.StartSpan("format", ext.RPCServerOption(spanCtx))
    defer span.Finish()
