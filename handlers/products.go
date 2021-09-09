package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ayushpaharia/microservices-with-go/data"
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
		p.l.Println("PUT", r.URL.Path)
		// expect the id in the URI
		rs := regexp.MustCompile(`/([0-9]+)`)
		g := rs.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.l.Println(id)
		p.updateProduct(id, rw, r)
		return
	}

	// PUT: Update Prodcuts

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

	// err := product.fromJSON(r.Body)

	// lp := data.AddProducts(product)
	// d, err := ioutil.ReadAll(r.Body)
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