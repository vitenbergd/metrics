package metrics

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shaj13/go-guardian/v2/auth/strategies/token"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"log"
	"net/http"
)

var (
	myCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "my_own_private_metrics_example_counter",
		Help: "My own private metrics example counter",
	})
	mySummaryVec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "my_own_private_metrics_example_summary",
			Help:       "My own private metrics example summary",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service"},
	)
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(mySummaryVec)
}

func recordCounter() {
	go func() {
		for {
			myCounter.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func recordSummary() {
	go func() {
		for {
			r := rand.Float64()
			mySummaryVec.WithLabelValues("random_summary_value").Observe(r)
			time.Sleep(2 * time.Millisecond)
		}
	}()
}

func genStrategy(tokenFile string) (union.Union, error) {
	tokenStrategy, err := token.NewStaticFromFile(tokenFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't open token file '%s': %s", tokenFile, err))
	}
	// Union strategy provides useful AuthenticateRequest method
	// which returns auth.Info including user
	return union.New(tokenStrategy), nil
}
func middleware(strategy union.Union, handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, user, err := strategy.AuthenticateRequest(r)
		remoteAddr := r.RemoteAddr
		if err != nil {
			log.Printf("Authentication failed for '%s': %s", remoteAddr, err)
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		log.Printf("Authenticaied: '%s', from '%s'", user.GetUserName(), remoteAddr)
		handler.ServeHTTP(w, r)
	}
}

func Metrics(httpAddr string, metricsPath string, tokenFile string) error {
	recordCounter()
	recordSummary()
	strategy, err := genStrategy(tokenFile)
	if err != nil {
		return err
	}
	http.Handle(metricsPath, middleware(strategy, promhttp.Handler()))
	log.Printf("Starting: addr '%s', metrics path '%s', token file '%s'", httpAddr, metricsPath, tokenFile)
	return http.ListenAndServe(httpAddr, nil)
}
