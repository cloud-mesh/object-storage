package main

import (
	"github.com/cloud-mesh/object-storage-sdk/impl/huaweicloud_obs"
	"github.com/cloud-mesh/object-storage-sdk/impl/huaweicloud_obs/obs"
	"github.com/cloud-mesh/object-storage/handler/http"
	"github.com/cloud-mesh/object-storage/model/repository/storage"
	"github.com/cloud-mesh/object-storage/model/usecase"
	"github.com/cloud-mesh/object-storage/utils"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

const (
	storageVendorHuawei = "huawei"
)

var (
	// 基本配置
	storageVendor = utils.Env("STORAGE_VENDOR", storageVendorHuawei)
	httpPort      = utils.EnvInt("HTTP_PORT", 80)
	dataPath      = utils.Env("DATA_PATH", "/data/")
	logLevel      = utils.Env("LOG_LEVEL", "info")

	// 华为对象存储
	huaweiEndpoint = utils.Env("HUAWEI_OBS_ENDPOINT", "")
	huaweiLocation = utils.Env("HUAWEI_OBS_LOCATION", "")
	huaweiAK       = utils.Env("HUAWEI_OBS_AK", "")
	huaweiSK       = utils.Env("HUAWEI_OBS_SK", "")
)

func main() {
	mainLogFile := openFile(filepath.Join(dataPath, "main.log"))
	defer mainLogFile.Close()
	accessLogFile := openFile(filepath.Join(dataPath, "access.log"))
	defer accessLogFile.Close()

	storageRepo := getStorage()
	ucase := usecase.NewUseCase(storageVendor, storageRepo)

	configLog(mainLogFile, logLevel)
	http.ServerHTTP(httpPort, accessLogFile, func(e *echo.Echo) {
		http.NewHandler(ucase).Route(e)
	})
}

func getStorage() storage.Repo {
	switch storageVendor {
	case storageVendorHuawei:
		// 华为云对象存储
		huaweiObsClient, err := obs.New(huaweiAK, huaweiSK, huaweiEndpoint)
		if err != nil {
			log.WithError(err).Fatalf("new huawei obs client")
		}
		huaweiObsClientAdapter := huaweicloud_obs.NewClient(huaweiLocation, huaweiObsClient)
		return storage.New(huaweiObsClientAdapter)
	default:
		panic("unsupported storage vendor")
	}
}

func openFile(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.WithError(err).Fatalf("failed to open file: %s", filePath)
	}
	return file
}

func configLog(file *os.File, level string) {
	if l, err := log.ParseLevel(level); err == nil {
		log.SetLevel(l)
	} else {
		log.WithError(err).Errorf("parse level failed: level=%s", level)
	}
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	log.SetReportCaller(true)
}
