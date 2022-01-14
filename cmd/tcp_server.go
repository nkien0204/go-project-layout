package cmd

import (
	"github.com/nkien0204/projectTemplate/configs"
	"github.com/nkien0204/projectTemplate/internal/log"
	"github.com/nkien0204/projectTemplate/internal/network/tcp_handler/tcp_server"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var runServerCmd = &cobra.Command{
	Use:   "server",
	Short: "start tcp server",
	Run:   runServer,
}

func init() {
	serveCmd.AddCommand(runServerCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	configs.Config = configs.InitConfigs()
	ServerManager := tcp_server.GetServer()
	go ServerManager.Listen()
	go tcp_server.RunTcpTimer()

	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	log.Logger().Warn("shutdown app")
}
