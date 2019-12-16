An example of using opentracing.

This example shows how to trace over the process boundaries and RPC calls.

Jaeger start

    docker run \
    -p 5775:5775/udp \
    -p 16686:16686 \
    -p 6831:6831/udp \
    -p 6832:6832/udp \
    -p 5778:5778 \
    -p 14268:14268 \
    jaegertracing/all-in-one:latest
    
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
