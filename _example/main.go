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
		logger.Diagnostics().OnEvent(logger.EventRequest, req)
		rw := logger.NewResponseWriter(res)
		handler(rw, req)
		logger.Diagnostics().OnEvent(logger.EventRequestComplete, req, rw.StatusCode(), rw.ContentLength(), time.Now().Sub(start))
	}
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"ok!"}`))
}

func fatalErrorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	logger.Diagnostics().Fatal(exception.New("this is a fatal exception"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func errorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	logger.Diagnostics().Error(exception.New("this is an exception"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func warningHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	logger.Diagnostics().Warning(exception.New("this is warning"))
	res.Write([]byte(`{"status":"not ok."}`))
}

func postHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	b := pool.Get()
	defer req.Body.Close()
	defer pool.Put(b)
	b.ReadFrom(req.Body)
	res.Write([]byte(fmt.Sprintf(`{"status":"ok!","received_bytes":%d}`, b.Len())))
	logger.Diagnostics().OnEvent(logger.EventPostBody, b.Bytes())
}

func port() string {
	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		return envPort
	}
	return "8888"
}

func main() {
	logger.SetDiagnostics(logger.NewDiagnosticsAgentFromEnvironment())
	logger.Diagnostics().EventQueue().UseSynchronousDispatch() //events fire in order, but will hang if queue is full
	logger.Diagnostics().EventQueue().SetMaxWorkItems(1 << 20) //make the queue size enormous (~1mm items)
	logger.Diagnostics().AddEventListener(logger.EventRequest,
		logger.NewRequestHandler(func(writer logger.Logger, ts logger.TimeSource, req *http.Request) {
			logger.WriteRequest(writer, ts, req)
		}))
	logger.Diagnostics().AddEventListener(logger.EventRequestComplete,
		logger.NewRequestCompleteHandler(func(writer logger.Logger, ts logger.TimeSource, req *http.Request, statusCode, contentLengthBytes int, elapsed time.Duration) {
			logger.WriteRequestComplete(writer, ts, req, statusCode, contentLengthBytes, elapsed)
		}))
	logger.Diagnostics().AddEventListener(logger.EventPostBody,
		logger.NewRequestBodyHandler(func(writer logger.Logger, ts logger.TimeSource, body []byte) {
			logger.WriteRequestBody(writer, ts, body)
		}))
	logger.Diagnostics().AddEventListener(logger.EventError, func(wr logger.Logger, ts logger.TimeSource, e uint64, args ...interface{}) {
		//ping an external service?
		//log something to the db?
		//this action will be handled by a separate go-routine
	})

	http.HandleFunc("/", logged(indexHandler))
	http.HandleFunc("/fatalerror", logged(fatalErrorHandler))
	http.HandleFunc("/error", logged(errorHandler))
	http.HandleFunc("/warning", logged(warningHandler))
	http.HandleFunc("/post", logged(postHandler))
	logger.Diagnostics().Infof("Listening on :%s", port())
	logger.Diagnostics().Infof("Diagnostics %s", logger.ExpandEventNames(logger.Diagnostics().Verbosity()))
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}
