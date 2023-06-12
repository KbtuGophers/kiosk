package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/config"
	"github.com/KbtuGophers/kiosk/account/internal/handler"
	"github.com/KbtuGophers/kiosk/account/internal/repository"
	"github.com/KbtuGophers/kiosk/account/internal/service/smsc"
	//service2 "github.com/KbtuGophers/kiosk/account/internal/service"
	"github.com/KbtuGophers/kiosk/account/internal/service/account"
	"github.com/KbtuGophers/kiosk/account/internal/service/otp"
	"github.com/KbtuGophers/kiosk/account/pkg/log"
	"github.com/KbtuGophers/kiosk/account/pkg/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	schema      = "account"
	version     = "1.0.0"
	description = "account-service"
)

func Run() {
	logger := log.New(version, description)
	cfg, err := config.New()

	if err != nil {
		logger.Error("ERR_INIT_CONFIG", zap.Error(err))
		return
	}

	repo, err := repository.New(repository.WithPostgresStore(schema, cfg.POSTGRES.DSN))
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORY", zap.Error(err))
		return
	}
	defer repo.Close()

	accountService, err := account.New(account.WithAccountRepository(repo.Account))
	if err != nil {
		logger.Error("ERR_INIT_SERVICE", zap.Error(err))
		return
	}

	//twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
	//	Username: cfg.TWILIO.Username,
	//	Password: cfg.TWILIO.Password,
	//})

	//smscClient := smsc.Client{
	//	Login:    cfg.SMSC.Login,
	//	Password: cfg.SMSC.Password,
	//	Sender:   cfg.SMSC.Sender,
	//	Tinyurl:  "1",
	//}

	smscClient, err := smsc.New(cfg.SMSC.Login, cfg.SMSC.Password, cfg.SMSC.Sender)
	if err != nil {
		logger.Error("ERR_NEW_SMSC_CLIENT")
	}

	attempts, err := strconv.ParseInt(cfg.OTP.Attempts, 10, 64)
	if err != nil {
		logger.Error("ERR_PARSE_ATTEMPTS_TO_INT")
	}
	interval, err := strconv.ParseInt(cfg.OTP.Interval, 10, 64)
	if err != nil {
		logger.Error("ERR_PARSE_INTERVAL_TO_INT")
	}

	otpService, err := otp.NewOtpService(
		smscClient,
		int(attempts),
		int(interval),
		otp.WithOtpRepository(repo.Otp),
	)
	if err != nil {
		logger.Error("ERR_NEW_OTP_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(handler.Dependencies{
		AccountService: accountService, OtpService: otpService},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_SERVICE", zap.Error(err))
		return
	}

	servers, err := server.New(server.WithHTTPServer(
		handlers.HTTP, cfg.HTTP.Port))
	if err != nil {
		logger.Error("ERR_INIT_SERVER", zap.Error(err))
		return
	}

	if err = servers.Run(logger); err != nil {
		logger.Error("ERR_RUN_SERVER", zap.Error(err))
		return
	}

	// Graceful Shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1) // create channel to signify a signal being sent

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-quit                                             // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err = servers.Stop(ctx); err != nil {
		panic(err) // failure/timeout shutting down the httpServer gracefully
	}

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here

	fmt.Println("Server was successful shutdown.")

}
