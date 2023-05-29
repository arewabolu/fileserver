package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
)

var (
	dbConn *pg.DB
	svc    *s3.Client
)

func init() {
	dbConn = pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Password: os.Getenv("PG_PWD"),
		Database: "fileserver",
	})
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Could not load config")
	}

	svc = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = region
		o.UseAccelerate = true
	})

}
