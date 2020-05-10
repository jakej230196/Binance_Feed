package main
import (
"fmt"
"sync"
"time"
"strings"
"path"
"runtime"
"os"
log "github.com/sirupsen/logrus"
"github.com/jakej230196/prometheus"
"github.com/jakej230196/Kafka/Producer"
)

func GetLogFormatter() *log.JSONFormatter {
	Formatter := &log.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return f.Function, fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
		FieldMap: log.FieldMap{
			log.FieldKeyFile:  "@file",
			log.FieldKeyTime:  "@timestamp",
			log.FieldKeyLevel: "@level",
			log.FieldKeyMsg:   "@message",
			log.FieldKeyFunc:  "@caller",
		},
	}
	Formatter.TimestampFormat = "2006-01-02T15:04:05.999999999Z"
	return Formatter
}

func LogLevel(lvl string) log.Level {
	level := strings.ToLower(lvl)
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	default:
		panic(fmt.Sprintf("Log level (%s) is not supported", lvl))
	}
}
func init() {
	log.SetReportCaller(true)
	Formatter := GetLogFormatter()
	log.SetFormatter(Formatter)
	log.SetLevel(LogLevel(os.Getenv("LOG_LEVEL")))
	log.Info("LOG_LEVEL: ", os.Getenv("LOG_LEVEL"))
}

func main(){
	ProducerFeed := KafkaProducer.Start(os.Getenv("KAFKA_BROKERLIST"), os.Getenv("KAFKA_OUTPUT_TOPIC"), Prometheus.InfoGauge, Prometheus.TotalErrorsReadGauge, Prometheus.ReadDurationHistogram)
	var WaitGroup sync.WaitGroup
	WaitGroup.Add(1)
	go WhatsTheTime(&WaitGroup, ProducerFeed)
	WaitGroup.Wait()
}


func WhatsTheTime(WaitGroup *sync.WaitGroup, ProducerFeed chan []byte){
	defer WaitGroup.Done()
	for{
		fmt.Println("The Time is:", time.Now())
		ProducerFeed <- []byte(time.Now().String())
		time.Sleep(10 * time.Second)
	}
}	
