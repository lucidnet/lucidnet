package main

import (
	"github.com/lucidnet/lucidnet/internal/app/vertex"
	"github.com/lucidnet/lucidnet/internal/pkg/env"
)

func main() {
	vertex.NewServer(&vertex.ServerConfig{
		Host:     env.GetOrDefault("HOST", "0.0.0.0"),
		Port:     env.GetOrDefault("PORT", "3000"),
		DataPath: env.GetOrDefault("DATA_PATH", "data"),
	}).Start()
}
