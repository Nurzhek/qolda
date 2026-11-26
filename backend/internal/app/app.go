package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"qolda/internal/config"
	"qolda/internal/constant"
	domainListingRepoDbP "qolda/internal/domain/listing/repo/db"
	domainListingServiceP "qolda/internal/domain/listing/service"
	domainSessionServiceP "qolda/internal/domain/session/service"
	domainUsrRepoDbP "qolda/internal/domain/usr/repo/db"
	domainUsrServiceP "qolda/internal/domain/usr/service"
	handlerGrpcP "qolda/internal/handler/grpc"
	handlerHttpP "qolda/internal/handler/http"
	serviceFilestoreServiceP "qolda/internal/service/filestore/service"
	serviceGoogleServiceP "qolda/internal/service/google/service"
	usecaseAuthP "qolda/internal/usecase/auth"
	usecaseFileP "qolda/internal/usecase/file"
	usecaseListingP "qolda/internal/usecase/listing"
	usecaseUsrP "qolda/internal/usecase/usr"
	"qolda/pkg/proto/qolda_v1"
)

type App struct {
	pgpool *pgxpool.Pool

	grpcServer       *GrpcServer
	httpServer       *http.Server
	systemHttpServer *http.Server

	ctx       context.Context
	ctxCancel context.CancelFunc

	exitCode int
}

func (a *App) Init() {
	var err error

	a.ctx, a.ctxCancel = context.WithCancel(context.Background())

	// logger
	initLogger(config.Conf.Debug, config.Conf.LogLevel)

	if config.Conf.JwtSecret == "dev-secret-change-me" && !config.Conf.Debug {
		errCheck(errors.New("JWT_SECRET is not set"), "config")
	}

	// pgpool
	a.pgpool, err = initPgPool(config.Conf.PgDsn)
	errCheck(err, "pgpool init")

	// migrations
	{
		runMigrations()
		slog.Info("PG-migrations have been successfully applied")
	}

	// session
	sessionService := domainSessionServiceP.New(config.Conf.JwtSecret, config.Conf.AccessTtl, config.Conf.RefreshTtl)

	// google
	hostedDomain := ""
	if len(config.Conf.AllowedDomains) == 1 {
		hostedDomain = config.Conf.AllowedDomains[0]
	}
	googleService := serviceGoogleServiceP.New(
		config.Conf.GoogleClientId,
		config.Conf.GoogleClientSecret,
		strings.TrimRight(config.Conf.PublicUrl, "/")+"/api/auth/google/callback",
		hostedDomain,
	)

	// filestore
	filestoreService, err := serviceFilestoreServiceP.New(config.Conf.UploadDir)
	errCheck(err, "filestore init")

	// usr
	usrRepo := domainUsrRepoDbP.New(a.pgpool)
	usrService := domainUsrServiceP.New(usrRepo)
	usrUsecase := usecaseUsrP.New(usrService, sessionService)
	usrHandler := handlerGrpcP.NewUsr(usrUsecase)

	// auth
	authUsecase := usecaseAuthP.New(usrService, sessionService, googleService)
	authHandler := handlerGrpcP.NewAuth(authUsecase)
	authHttpHandler := handlerHttpP.NewAuth(authUsecase, config.Conf.PublicUrl)

	// listing
	listingRepo := domainListingRepoDbP.New(a.pgpool)
	listingService := domainListingServiceP.New(listingRepo)
	listingUsecase := usecaseListingP.New(listingService, usrService, sessionService)
	listingHandler := handlerGrpcP.NewListing(listingUsecase)

	// file
	fileUsecase := usecaseFileP.New(filestoreService, sessionService)
	fileHttpHandler := handlerHttpP.NewFile(fileUsecase, sessionService)

	// dict
	dictHandler := handlerGrpcP.NewDict(googleService.Enabled())

	// grpc server
	{
		a.grpcServer = NewGrpcServer("main", sessionService, func(server *grpc.Server) {
			qolda_v1.RegisterDictServiceServer(server, dictHandler)
			qolda_v1.RegisterAuthServiceServer(server, authHandler)
			qolda_v1.RegisterUsrServiceServer(server, usrHandler)
			qolda_v1.RegisterListingServiceServer(server, listingHandler)
		})
	}

	// http-gw server
	{
		var handler http.Handler

		handler, err = GrpcGatewayCreateHandler(func(mux *runtime.ServeMux) error {
			opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

			var conn *grpc.ClientConn
			conn, err = grpc.NewClient("localhost:"+config.Conf.GrpcPort, opts...)
			errCheck(err, "grpc.Dial")

			// register grpc handlers
			handlers := []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
				qolda_v1.RegisterDictServiceHandler,
				qolda_v1.RegisterAuthServiceHandler,
				qolda_v1.RegisterUsrServiceHandler,
				qolda_v1.RegisterListingServiceHandler,
			}
			for _, h := range handlers {
				err = h(context.Background(), mux, conn)
				if err != nil {
					return fmt.Errorf("grpc-gateway: register grpc-handler: %w", err)
				}
			}

			// custom http handlers
			httpHandlers := []struct {
				method  string
				path    string
				handler runtime.HandlerFunc
			}{
				{"GET", "/auth/google", authHttpHandler.GoogleStart},
				{"GET", "/auth/google/callback", authHttpHandler.GoogleCallback},
				{"POST", "/file", fileHttpHandler.Upload},
			}
			for _, h := range httpHandlers {
				err = mux.HandlePath(h.method, h.path, h.handler)
				if err != nil {
					return fmt.Errorf("grpc-gateway: register http-handler: %w", err)
				}
			}

			return nil
		})
		errCheck(err, "grpcGatewayCreateHandler")

		// server: API, uploaded photos and SPA on one port
		a.httpServer = &http.Server{
			Addr:              ":" + config.Conf.HttpPort,
			Handler:           rootHandler(handler, config.Conf.UploadDir, config.Conf.WebDir),
			ReadHeaderTimeout: 2 * time.Second,
			ReadTimeout:       time.Minute,
			MaxHeaderBytes:    300 * 1024,
		}
	}

	// system http server (healthcheck, docs)
	{
		a.systemHttpServer = SystemHttpServerCreate()
	}
}

func (a *App) PreStartHook() {
	slog.Info("PreStartHook")
}

func (a *App) Start() {
	slog.Info("Starting", "version", constant.Version, "commit", constant.Commit, "built_at", constant.BuiltAt)

	// grpc server
	{
		err := a.grpcServer.Start()
		errCheck(err, "grpcServer.Start")
	}

	// http-gw server
	{
		go func() {
			err := a.httpServer.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCheck(err, "http-server stopped")
			}
		}()
		slog.Info("http-server started " + a.httpServer.Addr)
	}

	// system http server
	{
		go func() {
			err := a.systemHttpServer.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCheck(err, "system-http-server stopped")
			}
		}()
		slog.Info("system-http-server started " + a.systemHttpServer.Addr)
	}
}

func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCtxCancel()

	// wait signal
	<-signalCtx.Done()
}

func (a *App) Stop() {
	slog.Info("Shutting down...")

	// stop context
	a.ctxCancel()

	// http-gw server
	{
		ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer ctxCancel()

		if err := a.httpServer.Shutdown(ctx); err != nil {
			slog.Error("http-server shutdown error", "error", err)
			a.exitCode = 1
		}
	}

	// system http server
	{
		ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer ctxCancel()

		if err := a.systemHttpServer.Shutdown(ctx); err != nil {
			slog.Error("system-http-server shutdown error", "error", err)
			a.exitCode = 1
		}
	}

	// grpc server
	a.grpcServer.Stop()
}

func (a *App) WaitJobs() {
	slog.Info("waiting jobs")
}

func (a *App) Exit() {
	slog.Info("Exit")

	a.pgpool.Close()

	os.Exit(a.exitCode)
}

func errCheck(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}
