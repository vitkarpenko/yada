package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/vitkarpenko/yada/internal/bot"
	"github.com/vitkarpenko/yada/internal/config"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env not found")
	}

	var cfg config.Config
	envconfig.MustProcess("YADA", &cfg)

	yada := bot.NewYada(cfg)

	swearFile, err := os.Open("data/swear.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer swearFile.Close()

	r, err := gzip.NewReader(swearFile)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	uncompressed := bytes.NewBuffer(nil)
	io.Copy(uncompressed, r)

	yada.Run()
}
