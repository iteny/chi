package middleware

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"github.com/pressly/chi"
)

// TODO: keep profiler here or move to goware/?

// Profiler is a convenient subrouter used for mounting net/http/pprof. ie.
//
// func MyService() http.Handler {
//   r := chi.NewRouter()
//   // ..middlewares
//   r.Mount("/debug", profiler.Router())
//   // ..routes
//   return r
// }
func Profiler() http.Handler {
	r := chi.NewRouter()
	r.Use(NoCache)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", 301)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", 301)
	})

	r.HandleFunc("/pprof/", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.HandleFunc("/vars", expVars)

	return r
}

// Replicated from expvar.go as not public.
func expVars(w http.ResponseWriter, r *http.Request) {
	first := true
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
