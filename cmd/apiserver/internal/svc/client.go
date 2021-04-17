package svc

import (
	"context"

	"authkey/cmd/apiserver/internal/types"
	"authkey/pkg/lib"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func ListClientGroup(ctx context.Context) ([]*types.ClientGroup, error) {
	ctx, span := tracer.Start(ctx, "ListClientGroup")
	defer span.End()

	var ls []*types.ClientGroup
	_, err := l.NewFindSession().Find(&ls)
	if err != nil {
		span.SetStatus(codes.Error, "error")
		lib.SpanLog(ctx, span, zap.ErrorLevel, "grcp",
			attribute.String("reason", err.Error()),
		)
	}

	lib.SpanLog(ctx, span, zap.DebugLevel, "client groups", attribute.Array("groups", ls))

	return ls, nil
}
