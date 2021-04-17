package handler

import (
	"context"
	"encoding/json"
	"errors"

	"authkey/cmd/apiserver/internal/config"
	"authkey/cmd/apiserver/internal/middleware"
	"authkey/cmd/apiserver/internal/svc"
	"authkey/cmd/apiserver/internal/types"
	"authkey/pkg/lib"
	mvalidator "authkey/pkg/lib/validator"
	"authkey/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var (
	tracer trace.Tracer
)

func RootCtx(c *gin.Context) context.Context {
	return metadata.AppendToOutgoingContext(c.Request.Context(), "lang", middleware.GetLang(c))
}

func Bind(ctx context.Context, c *gin.Context, req interface{}) error {
	ctx, span := tracer.Start(ctx, "bind")
	defer span.End()

	trans, _ := mvalidator.Uni.GetTranslator(middleware.GetLang(c))

	err := c.ShouldBind(req)
	if err != nil {
		switch ev := err.(type) {
		case validator.ValidationErrors:
			err = errors.New(mvalidator.Translate(ev, trans))
		case *json.UnmarshalTypeError:
		//	unmarshalTypeError := err.(*json.UnmarshalTypeError)
		//	errStr = fmt.Errorf("%s 类型错误，期望类型 %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
		default:
		}

		lib.SpanLog(ctx, span, zap.ErrorLevel, "bind", attribute.String("error", err.Error()))

		c.JSON(400, gin.H{
			"error": util.I18nError(c, "arg.parse", err.Error()),
		})

		return err
	}

	return nil
}

func InitHandler() *gin.Engine {
	mvalidator.InitGin()

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("middleware"))
	r.Use(middleware.I18n())
	if config.GlobalConfig.Common.IsDev() {
		r.Use(middleware.DebugReq())
	}

	// apiv1 := r.Group("/api/v1")
	// apiv1.Use(middleware.JWT())

	// {
	// 	apiv1.POST("/users2", _CreateUser)
	// }

	// r.POST("/users", _CreateUser)

	adminv1 := r.Group("/api/v1/admin")
	adminv1.Use(middleware.JWT())
	adminv1.Use(middleware.IsAdmin())

	{
		// client
		clients := adminv1.Group("/clients")
		{
			groups := clients.Group("/groups")
			{
				groups.GET("", _ListClientGroup)
				groups.POST("", _CreateClientGroup)
				groups.PATCH("/{id}", _UpdateClientGroup)
				groups.DELETE("", _DeleteClientGroup)
			}
		}
	}

	tracer = otel.Tracer("handler")

	return r
}

func _CreateUser(c *gin.Context) {
	ctx, span := tracer.Start(RootCtx(c), "_CreateUser")
	defer span.End()

	r := &types.CreateUserReq{}
	if err := Bind(ctx, c, r); err != nil {
		return
	}

	users, err := svc.CreateUser(ctx, r)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}
