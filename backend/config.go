// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package main

type ConsulConfig struct {
	Address    string `yaml:"address"`
	DataCenter string `yaml:"datacenter"`
	Token      string `yaml:"admin_token"`
}

type Config struct {
	Consul   ConsulConfig `yaml:"consul"`
	LogLevel string       `yaml:"log_level"`
	LogFile  string       `yaml:"log_file"`
	Port     int          `yaml:"port"`
}

var config Config = Config{
	ConsulConfig{"http://127.0.0.1:8500", "dc1", ""},
	"info",
	"",
	3668,
}
