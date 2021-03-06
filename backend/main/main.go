package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dmtr/mail_me_all/backend/app"
	"github.com/dmtr/mail_me_all/backend/mail"
	"github.com/dmtr/mail_me_all/backend/twapi"
	"github.com/dmtr/mail_me_all/backend/twproxy"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	runAPI           string = "api"
	runTwProxy       string = "run-tw-proxy"
	check            string = "check-new-subscriptions"
	prepare          string = "prepare-subscriptions"
	send             string = "send-subscriptions"
	testEmail        string = "test-email"
	sendConfirmation string = "send-confirmation"
	removeTweets     string = "remove-old-tweets"
)

func handleSignals(server *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT:
				log.Infof("Received shutdown signal: %s", s)
				err := server.Close()
				if err != nil {
					log.Errorf("Web server closed with error: %s", err)
				}

			default:
				log.Warningf("Unknown signal!?, %s", s)
			}
		}
	}()
}

func startAPIServer(app *app.App) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", app.Conf.Host, app.Conf.Port),
		Handler: app.Router,
	}

	go handleSignals(server)

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Info("Web server shutdown complete")
		} else {
			log.Errorf("Web server closed unexpect: %s", err)
		}
	}
	log.Info("Exiting")
}

func startTwProxy(app *app.App) {
	log.Info("Starting Twitter API proxy server")
	lsnr, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", app.Conf.TwProxyPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	var opts []grpc.ServerOption
	if app.Conf.PemFile != "" && app.Conf.KeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(app.Conf.PemFile, app.Conf.KeyFile)
		if err != nil {
			log.Fatalf("Cant read credentials file: %s", err)
		}
		log.Info("Load credentials")

		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	t := twapi.NewTwitter(app.Conf.TwConsumerKey, app.Conf.TwConsumerSecret)
	s := twproxy.NewServiceServer(t)
	pb.RegisterTwProxyServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	grpcServer.Serve(lsnr)
}

func main() {
	flag.String("tw-consumer-key", "", "twitter consumer key")
	flag.String("tw-consumer-secret", "", "twitter consumer secret")

	flag.String("auth-key", "", "auth key")
	flag.String("encrypt-key", "", "encryption key")

	var subscriptionIDs *string = flag.String("subscription-ids", "", "subscription IDs")
	var subject *string = flag.String("subject", "", "email subject")
	var to *string = flag.String("to", "", "email to")
	var body *string = flag.String("body", "", "email body")

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	var IDs []uuid.UUID
	if *subscriptionIDs != "" {
		for _, s := range strings.Split(*subscriptionIDs, ",") {
			id, err := uuid.Parse(s)
			if err == nil {
				IDs = append(IDs, id)
			}
		}
	}

	cmd := flag.Arg(0)
	if cmd == "" {
		cmd = runAPI
	}

	var a *app.App
	defer func() { a.Close() }()

	if cmd == runAPI {
		a = app.GetApp(true)
		startAPIServer(a)
	} else if cmd == runTwProxy {
		a = app.GetApp(false)
		startTwProxy(a)
	} else if cmd == check {
		a = app.GetApp(false, true, true, true)
		checkNewSubscriptions(a, IDs...)
	} else if cmd == prepare {
		a = app.GetApp(false, true, true, true)
		prepareSubscriptions(a, IDs...)
	} else if cmd == send {
		a = app.GetApp(false, true, true, true)
		sendSubscriptions(a, IDs...)
	} else if cmd == testEmail {
		a = app.GetApp(false, false, false, false)
		sender := mail.NewEmailSender(a.Conf)
		sender.Send(a.Conf.From, *to, *subject, *body)
	} else if cmd == sendConfirmation {
		a = app.GetApp(false, true, true, true)
		sendConfirmationEmail(a)
	} else if cmd == removeTweets {
		a = app.GetApp(false, true, false, true)
		removeOldTweets(a)
	} else {
		fmt.Printf("Unknown command %s", cmd)
		os.Exit(1)
	}
}
