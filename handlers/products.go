package handlers

import (
	"log"
	"net/http"

	"github.com/ayushpaharia/microservices-with-go/data"
)

type ProductLogger struct {
	l *log.Logger
}


func NewProducts(l *log.Logger) *ProductLogger {
	return &ProductLogger{l}
}



func (p *ProductLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// type Products data.Products
// func (p *Products) toJSON(w io.Writer) error {
// 	e := json.NewEncoder(w)
// 	return e.Encode(p)
// }
func (p *ProductLogger) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// lp_new := Products(lp)
	// d, err := json.Marshal(sampleList)
	
	err := data.JSONify(&lp, rw)
	// err := lp_new.toJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	// rw.Write(d)
}