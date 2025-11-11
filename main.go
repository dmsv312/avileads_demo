package main

import (
	http_filters "avileads-web/http/filters"
	_ "avileads-web/http/routers"
	"avileads-web/services/background"
	"avileads-web/services/booking"
	"avileads-web/utils"
	"avileads-web/views"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func main() {
	beego.InsertFilter(`/*`, beego.BeforeRouter, http_filters.AuthFilter())
	beego.InsertFilter(`/*`, beego.BeforeRouter, http_filters.GuardFilter())

	initBeegoOrmWebFunctions()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := Bootstrap(ctx); err != nil {
		log.Fatal(err)
	}

	beego.Run()
}

func init() {
	initDbConnection()
	initOrmModels()
}

func initDbConnection() {
	dbNameDefault := utils.GetEnv("PG_NAME", "avileads")
	dbHostDefault := utils.GetEnv("PG_HOST", "127.0.0.1")
	dbPortDefault := utils.GetEnv("PG_PORT", "4444")
	dbUserDefault := utils.GetEnv("PG_USER", "test")
	dbPasswordDefault := utils.GetEnv("PG_PASSWORD", "test")

	dbNameChatbot := utils.GetEnv("PG_NAME_CHATBOT", "avileads-full")
	dbHostChatbot := utils.GetEnv("PG_HOST_CHATBOT", "127.0.0.1")
	dbPortChatbot := utils.GetEnv("PG_PORT_CHATBOT", "5555")
	dbUserChatbot := utils.GetEnv("PG_USER_CHATBOT", "test")
	dbPasswordChatbot := utils.GetEnv("PG_PASSWORD_CHATBOT", "test")

	timeZone := utils.GetEnv("PG_TIMEZONE", "UTC")
	pgTimeZone := ""
	if len(timeZone) > 0 {
		pgTimeZone = "&timezone=" + timeZone
	}

	dbUrlDefault := fmt.Sprintf(
		"host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable TimeZone=%s search_path=%s",
		dbHostDefault, dbPortDefault, dbUserDefault, dbPasswordDefault,
		dbNameDefault, timeZone, "public,avileads_geo,avileads_config",
	)
	dbUrlChatBot := "postgres://" + dbUserChatbot + ":" + dbPasswordChatbot + "@" + dbHostChatbot + ":" + dbPortChatbot + "/" + dbNameChatbot + "?sslmode=disable" + pgTimeZone

	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", dbUrlDefault)
	orm.RegisterDataBase("chatbot", "postgres", dbUrlChatBot)
}

func initBeegoOrmWebFunctions() {
	beego.AddFuncMap("contains", views.Contains)
	beego.AddFuncMap("hasAccess", views.HasAccess)
	beego.AddFuncMap("cond", func(ok bool, a, b interface{}) interface{} {
		if ok {
			return a
		}
		return b
	})
	beego.AddFuncMap("mul", func(a, b int) int { return a * b })
	beego.AddFuncMap("seq", func(start, end int64) []int64 {
		var s []int64
		for i := start; i <= end; i++ {
			s = append(s, i)
		}
		return s
	})
	beego.AddFuncMap("add", func(a, b int) int { return a + b })
	beego.AddFuncMap("sub", func(a, b int) int { return a - b })
}

func initDadataService() (*booking.DadataService, error) {
	token := utils.GetEnv("DADATA_TOKEN", "")
	if token == "" {
		return nil, fmt.Errorf("DADATA_TOKEN is required")
	}

	secret := utils.GetEnv("DADATA_SECRET", "")
	if secret == "" {
		return nil, fmt.Errorf("DADATA_SECRET is required")
	}

	httpTimeoutSec := utils.MustGetInt("DADATA_HTTP_TIMEOUT_SEC", 10)
	retries := utils.MustGetInt("DADATA_RETRIES", 5)
	backoffMs := utils.MustGetInt("DADATA_BACKOFF_BASE_MS", 500)

	dadata := booking.NewDadataService(
		token,
		secret,
		time.Duration(httpTimeoutSec)*time.Second,
		retries,
		time.Duration(backoffMs)*time.Millisecond,
	)

	return dadata, nil
}

func startBackgroundServices(ctx context.Context, dadata *booking.DadataService) error {
	newOrm := func() orm.Ormer {
		o := orm.NewOrm()
		_ = o.Using("default")
		return o
	}

	background.Start(ctx, newOrm, dadata)

	log.Println("[bootstrap] background started: rps=17, workers=12 (global)")
	return nil
}

func Bootstrap(ctx context.Context) error {
	dadata, err := initDadataService()
	if err != nil {
		return err
	}

	if err := startBackgroundServices(ctx, dadata); err != nil {
		return err
	}

	return nil
}
