package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int 	`json:"id"`
	Name        string  `json:"name"`
	Description string	`json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product)  {
//  /* Convert ID to string before use */
// 	p.ID = strings.Replace(uuid.New().String(),"-","",-1)
	
p.ID = getNextID()
productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p

	return nil
}

func filter(pl Products, test func(int) bool) (Products) {
	var npl = Products{}
    for i, p := range pl {
        if test(i) {
            npl = append(npl, p)
        }
    }
	productList = npl
	return npl
}

func DeleteProduct(id int) (Products, error) {
	_, _, err := findProduct(id)
	if err != nil {
		return productList, err
	}

	mytest := func(i int) bool { return i == id }
	npl := filter(productList, mytest)
	
	return npl, nil
}


var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound 
}

func getNextID() int {
	lp := productList[len(productList) - 1]
	return lp.ID + 1
}

// func (p *Products) toJSON(w io.Writer) error {
// 	e := json.NewEncoder(w)
// 	return e.Encode(p)
// }

// func (p *Product) fromJSON(r io.Reader) error {
// 	e:= json.NewDecoder(r)
// 	return e.Decode(p)
// }

func FromJSONFunc(p *Product, r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
} 

func ToJSONFunc(p *Products, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

var productList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "random-slug-1",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "random-slug-2",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}