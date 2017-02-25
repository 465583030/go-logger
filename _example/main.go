package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	exception "github.com/blendlabs/go-exception"
	logger "github.com/blendlabs/go-logger"
)

var pool = logger.NewBufferPool(16)

func logged(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		logger.Default().OnEvent(logger.EventWebRequestStart, req)
		rw := logger.NewResponseWriter(res)
		handler(rw, req)
		logger.Default().OnEvent(logger.EventWebRequest, req, rw.StatusCode(), rw.ContentLength(), time.Now().Sub(start))
	}
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"ok!"}`))
}

func fatalErrorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	logger.Default().Fatal(exception.New("this is a fatal exception"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func errorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	logger.Default().Error(exception.New("this is an exception"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func warningHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	logger.Default().Warning(exception.New("this is warning"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func postHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	b := pool.Get()
	defer req.Body.Close()
	defer pool.Put(b)
	b.ReadFrom(req.Body)
	res.Write([]byte(fmt.Sprintf(`{"status":"ok!","received_bytes":%d}`, b.Len())))
	logger.Default().OnEvent(logger.EventWebRequestPostBody, b.Bytes())
}

func port() string {
	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		return envPort
	}
	return "8888"
}

func main() {
	logger.SetDefault(logger.NewFromEnvironment())

	logger.Default().AddEventListener(logger.EventWebRequestStart,
		logger.NewRequestStartListener(func(writer logger.Logger, ts logger.TimeSource, req *http.Request) {
			logger.WriteRequestStart(writer, ts, req)
		}))
	logger.Default().AddEventListener(logger.EventWebRequest,
		logger.NewRequestListener(func(writer logger.Logger, ts logger.TimeSource, req *http.Request, statusCode, contentLengthBytes int, elapsed time.Duration) {
			logger.WriteRequest(writer, ts, req, statusCode, contentLengthBytes, elapsed)
		}))
	logger.Default().AddEventListener(logger.EventWebRequestPostBody,
		logger.NewRequestBodyListener(func(writer logger.Logger, ts logger.TimeSource, body []byte) {
			logger.WriteRequestBody(writer, ts, body)
		}))
	logger.Default().AddEventListener(logger.EventError, func(wr logger.Logger, ts logger.TimeSource, e logger.EventFlag, args ...interface{}) {
		//ping an external service?
		//log something to the db?
		//this action will be handled by a separate go-routine
	})

	http.HandleFunc("/", logged(indexHandler))
	http.HandleFunc("/fatalerror", logged(fatalErrorHandler))
	http.HandleFunc("/error", logged(errorHandler))
	http.HandleFunc("/warning", logged(warningHandler))
	http.HandleFunc("/post", logged(postHandler))
	logger.Default().Infof("Listening on :%s", port())
	logger.Default().Infof("Events %s", logger.Default().Events().String())
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}
