package tracelog

import (
	"context"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var globalLogger zerolog.Logger

func SetGlobalLogger(logger zerolog.Logger) {
	globalLogger = logger
}

type TraceLogger struct {
	span   trace.Span
	logger zerolog.Logger
}

func New(ctx context.Context, name string) (tl *TraceLogger, ctxWithTraceId context.Context) {
	tl = &TraceLogger{logger: globalLogger}
	ctx, tl.span = otel.Tracer("").Start(ctx, name)
	return tl, ctx
}

func (t *TraceLogger) SetAttributes(attrs ...KeyValue) {
	for _, attr := range attrs {
		if !attr.IsValid() {
			continue
		}

		switch attr.vtype {
		case vTypeString:
			t.span.SetAttributes(attribute.String(attr.key, attr.value.(string)))
			t.logger = t.logger.With().Str(attr.key, attr.value.(string)).Logger()
		case vTypeInt:
			t.span.SetAttributes(attribute.Int(attr.key, attr.value.(int)))
			t.logger = t.logger.With().Int(attr.key, attr.value.(int)).Logger()
		case vTypeBool:
			t.span.SetAttributes(attribute.Bool(attr.key, attr.value.(bool)))
			t.logger = t.logger.With().Bool(attr.key, attr.value.(bool)).Logger()
		}
	}
}

func (t *TraceLogger) AddEvent(eventName string, logLevel zerolog.Level, attrs ...KeyValue) {
	var eventAttrs []trace.EventOption
	var tmpLogger = &t.logger

	for _, attr := range attrs {
		if !attr.IsValid() {
			continue
		}

		switch attr.vtype {
		case vTypeString:
			eventAttrs = append(eventAttrs, trace.WithAttributes(attribute.String(attr.key, attr.value.(string))))
			*tmpLogger = tmpLogger.With().Str(attr.key, attr.value.(string)).Logger()
		case vTypeInt:
			eventAttrs = append(eventAttrs, trace.WithAttributes(attribute.Int(attr.key, attr.value.(int))))
			*tmpLogger = tmpLogger.With().Int(attr.key, attr.value.(int)).Logger()
		case vTypeBool:
			eventAttrs = append(eventAttrs, trace.WithAttributes(attribute.Bool(attr.key, attr.value.(bool))))
			*tmpLogger = tmpLogger.With().Bool(attr.key, attr.value.(bool)).Logger()
		}
	}

	t.span.AddEvent(eventName, eventAttrs...)
	tmpLogger.WithLevel(logLevel).Msg(eventName)
}

func (t *TraceLogger) EndSpanWithRecordError(err error) {
	if t.span == nil {
		return
	}

	defer t.span.End()

	if err != nil {
		t.span.RecordError(err)
		t.span.SetStatus(codes.Error, err.Error())
	} else {
		t.span.SetStatus(codes.Ok, "")
	}
}

func (t *TraceLogger) EndSpan() {
	if t.span == nil {
		return
	}

	t.span.End()
}

func (t *TraceLogger) TraceId() (id string) {
	if t.span != nil {
		return t.span.SpanContext().TraceID().String()
	}
	return
}
