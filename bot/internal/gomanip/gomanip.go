package gomanip

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/trollLemon/DiscordBot/internal/util"
)

var tracer = otel.Tracer("discord-bot/gomanip")

func RandomFilter(gomanipClient *GoManip, image []byte, contentType string, kernelSize, lower, higher int64, normalize bool) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.RandomFilter", trace.WithAttributes(
		attribute.String("gomanip.operation", "randomFilter"),
		attribute.Int64("gomanip.kernel_size", kernelSize),
		attribute.Int64("gomanip.lower_bound", lower),
		attribute.Int64("gomanip.upper_bound", higher),
		attribute.Bool("gomanip.normalize", normalize),
	))
	defer span.End()
	queries := util.RandomFilterQuery(kernelSize, lower, higher, normalize)
	bytes, err := gomanipClient.Do(ctx, image, contentType, "randomFilter", queries)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}

func InvertImage(gomanipClient *GoManip, image []byte, contentType string) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.InvertImage", trace.WithAttributes(
		attribute.String("gomanip.operation", "invert"),
	))
	defer span.End()
	bytes, err := gomanipClient.Do(ctx, image, contentType, "invert", "")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}

func SaturateImage(gomanipClient *GoManip, image []byte, contentType string, saturation int64) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.SaturateImage", trace.WithAttributes(
		attribute.String("gomanip.operation", "saturate"),
		attribute.Int64("gomanip.saturation", saturation),
	))
	defer span.End()
	saturationNorm := float32(saturation) / 100.0
	queries := util.SaturateQuery(saturationNorm)
	bytes, err := gomanipClient.Do(ctx, image, contentType, "saturate", queries)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}

func EdgeDetect(gomanipClient *GoManip, image []byte, contentType string, lower, higher int64) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.EdgeDetect", trace.WithAttributes(
		attribute.String("gomanip.operation", "edgeDetection"),
		attribute.Int64("gomanip.lower_bound", lower),
		attribute.Int64("gomanip.upper_bound", higher),
	))
	defer span.End()
	queries := util.EdgeDetectQuery(lower, higher)
	bytes, err := gomanipClient.Do(ctx, image, contentType, "edgeDetection", queries)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}

func DilateImage(gomanipClient *GoManip, image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.DilateImage", trace.WithAttributes(
		attribute.String("gomanip.operation", "dilate"),
		attribute.Int64("gomanip.kernel_size", kernelSize),
		attribute.Int64("gomanip.iterations", iterations),
	))
	defer span.End()
	queries := util.DilateQuery(kernelSize, iterations)
	bytes, err := gomanipClient.Do(ctx, image, contentType, "morphology", queries)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}

func ErodeImage(gomanipClient *GoManip, image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.ErodeImage", trace.WithAttributes(
		attribute.String("gomanip.operation", "erode"),
		attribute.Int64("gomanip.kernel_size", kernelSize),
		attribute.Int64("gomanip.iterations", iterations),
	))
	defer span.End()
	queries := util.ErodeQuery(kernelSize, iterations)
	bytes, err := gomanipClient.Do(ctx, image, contentType, "morphology", queries)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}

func AddText(gomanipClient *GoManip, image []byte, contentType string, text string, fontScale float32, xPercentage float32, yPercentage float32) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.AddText", trace.WithAttributes(
		attribute.String("gomanip.operation", "text"),
		attribute.Float64("gomanip.font_scale", float64(fontScale)),
		attribute.Float64("gomanip.x_percentage", float64(xPercentage)),
		attribute.Float64("gomanip.y_percentage", float64(yPercentage)),
	))
	defer span.End()
	queries := util.AddTextQuery(text, fontScale, xPercentage, yPercentage)
	bytes, err := gomanipClient.Do(ctx, image, contentType, "text", queries)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}

func Reduced(gomanipClient *GoManip, image []byte, contentType string, quality float32) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.Reduced", trace.WithAttributes(
		attribute.String("gomanip.operation", "reduced"),
		attribute.Float64("gomanip.quality", float64(quality)),
	))
	defer span.End()
	queries := util.ReduceQuery(quality)
	bytes, err := gomanipClient.Do(ctx, image, contentType, "reduced", queries)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}

func Shuffle(gomanipClient *GoManip, image []byte, contentType string, partitions int64) ([]byte, error) {
	ctx, span := tracer.Start(context.Background(), "gomanip.Shuffle", trace.WithAttributes(
		attribute.String("gomanip.operation", "shuffle"),
		attribute.Int64("gomanip.partitions", partitions),
	))
	defer span.End()
	queries := util.ShuffleQuery(partitions)
	bytes, err := gomanipClient.Do(ctx, image, contentType, "shuffle", queries)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return bytes, err
}
