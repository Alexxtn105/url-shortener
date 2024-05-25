package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/Alexxtn105/protos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

// New конструктор клиента для SSO/Auth
func New(
	ctx context.Context, //
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	// опции для ретраев
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded), // указываем, какие коды нужно ретраить
		grpcretry.WithMax(uint(retriesCount)),                                      //количество ретраев
		grpcretry.WithPerRetryTimeout(timeout),                                     //таймаут ретрая
	}

	// Опции для интерцептора gprclog (логирование запросов и ответов)
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// Создаём соединение с gRPC-сервером SSO для клиента
	// по-хорошему здесь нужно создавать защищенное соединение, здесь будет insecure
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor( //обертка для двух следующих интерсепторов (создаем цепочку интерцепторов, чтобы все интерцепторы вызывались по очереди)
			// этот интерцептор будет тело каждого запроса и ответа,
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			// этот интерцептор будет делать ретраи в случае неудачных запросов
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))

	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// создаем gRPC-клиент SSO/Auth
	//  Она создаёт gRPC клиент для нужного сервиса,
	// в нашем случае — это SSO, а конкретно — Auth.
	// Напомню, что мы внутри proto-файла сделали группировку методов
	// на отдельные сервисы: Auth и в будущем Permissions, UserInfo и т.п.
	grpcClient := ssov1.NewAuthClient(cc)

	return &Client{
		api: grpcClient,
	}, nil

}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
// Копипаст из gRPC middleware, для логирования запросов и ответов.
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "grpc.IsAdmin"

	// основная часть
	resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: userID,
	})

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.IsAdmin, nil
}
