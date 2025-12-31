// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

//go:build !premium

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	httpadapter "github.com/FlyingOnion/consee/backend/adapter/http"
	"github.com/FlyingOnion/consee/backend/consul"
	"github.com/FlyingOnion/consee/backend/infra"
	"github.com/FlyingOnion/consee/backend/service"
	"github.com/spf13/pflag"

	"github.com/goccy/go-yaml"
)

var (
	configfile string
	verbose    int
	t          string
	port       int
)

func parseCmd() {
	pflag.StringVarP(&configfile, "config", "c", "config/config.yaml", "config file")
	pflag.StringVarP(&t, "token", "t", "", "consee admin token")
	pflag.IntVarP(&verbose, "verbose", "v", 0, "show more output")
	pflag.IntVarP(&port, "port", "p", 3668, "http server port")
	pflag.Parse()
}

func parseConfig() {
	if configfile == "" {
		configfile = "config/config.yaml"
	}
	f, err := os.Open(configfile)
	if err != nil {
		slog.Warn("failed to open config file, try to parse from flags", "file", configfile)
		goto parseFromFlags
	}
	yaml.NewDecoder(f).Decode(&config)
	f.Close()
parseFromFlags:
	if t != "" {
		config.Consul.Token = t
	}
	if verbose > 0 {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		config.LogLevel = "debug"
	}
	if port > 0 {
		config.Port = port
	}
	slog.Debug("consee configuration", "config", config)
}

func main() {
	parseCmd()
	parseConfig()

	client := consul.NewClient()
	qAdmin, wAdmin := &consul.QueryOptions{
		Datacenter: config.Consul.DataCenter,
		Token:      config.Consul.Token,
	}, &consul.WriteOptions{
		Datacenter: config.Consul.DataCenter,
		Token:      config.Consul.Token,
	}

	adminRepo := infra.NewAdmin(client, qAdmin, wAdmin)
	kvRepo := infra.NewKV(client)
	aclRepo := infra.NewACL(client)

	adminService := service.NewAdminService(adminRepo)
	kvService := service.NewKVService(kvRepo, adminService)
	aclService := service.NewACLService(aclRepo, adminService)
	a2 := service.NewA2(kvService, aclService, adminService)

	ctx, cancel := context.WithCancel(context.Background())
	initCtx := consul.ContextWithQueryOptions(ctx, qAdmin)
	initCtx = consul.ContextWithWriteOptions(initCtx, wAdmin)

	if err := a2.Initialize(initCtx); err != nil {
		slog.Error("failed to initialize", "error", err)
		cancel()
		os.Exit(1)
	}

	httpAdapter := httpadapter.NewAdapter(a2, kvService, aclService, adminService)
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: httpAdapter.Handler(),

		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("starting HTTP server", "address", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server terminated", "err", err)
		}
	}()

	slog.Info("starting consee")
	select {
	case <-ctx.Done():
	case sig := <-sigC:
		slog.Info("received signal", "signal", sig.String())
		slog.Info("shutting down HTTP server")
		ctx1, cancel1 := context.WithTimeout(context.Background(), 5)
		httpServer.Shutdown(ctx1)
		cancel1()
	}
	cancel()
	signal.Stop(sigC)
	slog.Info("stopped")
}
