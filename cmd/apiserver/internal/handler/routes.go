package handler

import (
	"context"
	"encoding/json"
	"errors"

	"authkey/cmd/apiserver/internal/config"
	"authkey/cmd/apiserver/internal/middleware"
	"authkey/pkg/lib"
	mvalidator "authkey/pkg/lib/validator"
	"authkey/pkg/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/meilihao/water"
	"github.com/meilihao/water-contrib/otelwater"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var (
	tracer trace.Tracer
)

func RootCtx(c *water.Context) context.Context {
	return metadata.AppendToOutgoingContext(c.Request.Context(), "lang", middleware.GetLang(c))
}

func Bind(ctx context.Context, c *water.Context, req interface{}) error {
	ctx, span := tracer.Start(ctx, "bind")
	defer span.End()

	trans, _ := mvalidator.Uni.GetTranslator(middleware.GetLang(c))

	b := binding.Default(c.Request.Method, c.ContentType())
	err := b.Bind(c.Request, req)
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

		c.JSON(400, water.H{
			"error": util.I18nError(c, "arg.parse", err.Error()),
		})

		return err
	}

	return nil
}

func InitHandler2() *water.Engine {
	mvalidator.InitGin()

	r := water.NewRouter()

	r.Classic()
	r.Use(otelwater.Middleware("middleware"))
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

	r.Group("/api/v1/admin", func(r *water.Router) {
		r.Use(middleware.JWT())
		r.Use(middleware.IsAdmin())

		// client
		r.Group("/clients", func(r *water.Router) {
			r.Group("/groups", func(r *water.Router) {
				r.GET("", _ListClientGroup)
				r.POST("", _CreateClientGroup)
				r.PATCH("/{id}", _UpdateClientGroup)
				r.DELETE("", _DeleteClientGroup)
			})
		})
	})

	tracer = otel.Tracer("handler")

	return r.Handler()
}

// func _CreateUser(c *gin.Context) {
// 	ctx, span := tracer.Start(RootCtx(c), "_CreateUser")
// 	defer span.End()

// 	r := &types.CreateUserReq{}
// 	if err := Bind(ctx, c, r); err != nil {
// 		return
// 	}

// 	users, err := svc.CreateUser(ctx, r)
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"error": err.Error(),
// 		})

// 		return
// 	}

// 	c.JSON(200, gin.H{
// 		"users": users,
// 	})
// }
