package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ayushpaharia/rest-golang/data"
)

type Products struct {
	l *log.Logger
}


func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}



func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// GET: Get all Products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	
	if r.Method == http.MethodPut {
		id, err := findId(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} 
		p.updateProduct(id, rw, r)
		return
	}

	if r.Method == http.MethodDelete {
		id, err := findId(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} 
		p.deleteProduct(id, rw, r)
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
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET:getProducts")
	lp := data.GetProducts()
	// lp_new := Products(lp)
	// d, err := json.Marshal(sampleList)
	
	err := data.ToJSONFunc(&lp, rw)
	// err := lp.toJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	// rw.Write(d)
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST:addProduct")
	np := &data.Product{}

	err := data.FromJSONFunc(np, r.Body)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}

	p.l.Printf("Product: %#v", np)

	data.AddProduct(np)
}

func (p *Products) deleteProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle DELETE:deleteProduct")

	pl, err := data.DeleteProduct(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	err = data.ToJSONFunc(&pl, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT:updateProduct")
	
	np := &data.Product{}
	err := data.FromJSONFunc(np, r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, np)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

func findId(r *http.Request) (int, error) {
	rs := regexp.MustCompile(`/([0-9]+)`)
		g := rs.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1  {
			return -1, fmt.Errorf("Invalid URI Code:%v", http.StatusBadRequest)
		}
		if  len(g[0]) != 2  {
			return -1, fmt.Errorf("Invalid URI Code:%v", http.StatusNotFound)
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			return -1, fmt.Errorf("Invalid URI unable to convert to number Code:%v", http.StatusBadRequest)
		}
		return id, nil
}