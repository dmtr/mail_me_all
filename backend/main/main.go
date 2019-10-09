package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/dmtr/mail_me_all/backend/app"
	"github.com/dmtr/mail_me_all/backend/fbproxy"
	"github.com/dmtr/mail_me_all/backend/fbwrapper"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	runAPI          string = "api"
	verifyFbLogin   string = "verify-fb-login"
	generateFbToken string = "generate-fb-token"
	runFBProxy      string = "run-fb-proxy"
)

func startAPIServer(app *app.App) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", app.Conf.Host, app.Conf.Port),
		Handler: app.Router,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Info("Recieve interrupt signal")
		err := server.Close()
		if err != nil {
			log.Errorf("Web server closed : %v", err)
		}

	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Info("Web server shutdown complete")
		} else {
			log.Errorf("Web server closed unexpect: %s", err)
		}
	}
	log.Info("Exiting")
}

func startFBProxy(app *app.App) {
	log.Info("Starting FB proxy server")
	lsnr, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", app.Conf.FBProxyPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	f := fbwrapper.NewFacebook(app.Conf.FbAppID, app.Conf.AppSecret, app.Conf.FbRedirectURI)
	s := fbproxy.NewServiceServer(f)
	pb.RegisterFbProxyServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	grpcServer.Serve(lsnr)
}

func main() {
	flag.String("app-secret", "", "app secret")
	flag.String("auth-key", "", "auth key")
	flag.String("encrypt-key", "", "encryption key")
	var accessToken *string = flag.String("access-token", "", "access token")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	cmd := flag.Arg(0)
	if cmd == "" {
		cmd = runAPI
	}

	var a *app.App
	defer func() { fmt.Print("Shutting down"); a.Close() }()

	if cmd == runAPI {
		a = app.GetApp(true)
		startAPIServer(a)
	} else if cmd == verifyFbLogin {
		a = app.GetApp(false)
		f := fbwrapper.NewFacebook(a.Conf.FbAppID, a.Conf.AppSecret, a.Conf.FbRedirectURI)
		VerifyFbLogin(*accessToken, f)
	} else if cmd == generateFbToken {
		a = app.GetApp(false)
		f := fbwrapper.NewFacebook(a.Conf.FbAppID, a.Conf.AppSecret, a.Conf.FbRedirectURI)
		GenerateFbToken(*accessToken, f)
	} else if cmd == runFBProxy {
		a = app.GetApp(false)
		startFBProxy(a)
	} else {
		fmt.Printf("Unknown command %s", cmd)
		os.Exit(1)
	}
}
