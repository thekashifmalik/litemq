package main

import (
	"github.com/thekashifmalik/litemq/internal"
)

func main() {
	server := internal.NewServer()
	server.Serve()
}

// mux := http.NewServeMux()
// mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
// 	slog.Info(fmt.Sprintf("HEALTH"))
// 	w.WriteHeader(http.StatusOK)
// })

// mux.HandleFunc("GET /{queue}", func(w http.ResponseWriter, r *http.Request) {
// })

// mux.HandleFunc("POST /{queue}/purge", func(w http.ResponseWriter, r *http.Request) {
// 	queueName := r.PathValue("queue")
// 	slog.Info(fmt.Sprintf("DEL %v", queueName))
// 	// queue, _ := queues[queueName]
// 	delete(server.queues, queueName)
// 	w.WriteHeader(http.StatusOK)
// })

// mux.HandleFunc("POST /{queue}", func(w http.ResponseWriter, r *http.Request) {
// 	queueName := r.PathValue("queue")
// 	data, _ := io.ReadAll(r.Body)
// 	request := &gen.EnqueueRequest{
// 		Queue: queueName,
// 		Data:  data,
// 	}
// 	_, err := server.Enqueue(r.Context(), request)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// })

// mux.HandleFunc("DELETE /{queue}", func(w http.ResponseWriter, r *http.Request) {
// 	queueName := r.PathValue("queue")
// 	queue := server.queues[queueName]
// 	slog.Info(fmt.Sprintf("DEQ %v", queueName))
// 	for {
// 		if len(queue.data) > 0 {
// 			break
// 		}
// 		select {
// 		case <-time.After(time.Second / 10):
// 			queue = server.queues[queueName]
// 		case <-r.Context().Done():
// 			return
// 		}
// 	}
// 	data := queue.data[0]
// 	queue.data = queue.data[1:]
// 	w.Write(data)
// })

// go http.ListenAndServe(":42080", mux)
