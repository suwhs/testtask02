package rest

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	swaggerui "github.com/esurdam/go-swagger-ui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"whs.su/rusprofile/src/rpc"
)

//go:embed assets/rpc/rusprofile.swagger.json
var docs []byte


var DefaultAssetFn = func(s string) ([]byte, error) { return docs, nil }

func RunRestServer(ctx context.Context, grpc_addr string, http_port int) {
	ctx, cancel := context.WithCancel(ctx)

	rmux := runtime.NewServeMux(runtime.WithErrorHandler(func(ctx context.Context, sm *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("err: %v",err)
	}))

	mux := swaggerui.NewServeMux(DefaultAssetFn, "swagger.json")
	mux.Handle("/", rmux)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := rpc.RegisterRusprofileHandlerFromEndpoint(ctx, rmux, grpc_addr, opts); err != nil {
		panic(fmt.Sprintf("could not start rest servce"))
	}
	rest_addr := fmt.Sprintf("localhost:%d",http_port)
	log.Printf("http server listening at %s",rest_addr)
	srv := http.Server{Addr: rest_addr, Handler: mux}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http server shutdown: %s", err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("context closed")
				cancel()
				srv.Shutdown(ctx)
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

}
