package logger

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nbari/violetear"
)

func JSONLogger(w *violetear.ResponseWriter, r *http.Request) {
    j := map[string]interface{}{
        "Time":        time.Now().UTC().Format(time.RFC3339),
        "RemoteAddr":  r.RemoteAddr,
        "URL":         r.URL.String(),
        "Status":      w.Status(),
        "Size":        w.Size(),
        "RequestTime": w.RequestTime(),
        "RequestID":   w.RequestID(),
    }
    if err := json.NewEncoder(os.Stdout).Encode(j); err != nil {
        log.Println(err)
    }
}

