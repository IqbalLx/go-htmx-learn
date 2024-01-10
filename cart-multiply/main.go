package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type Product struct {
	ID int
	Name string
	MinQty int
	QtyMultiple int
	Price int
}

type Cart struct {
	ID int
	Product Product
	Quantity int
}

func formatPrice(price int) string {
	currencyStr := strconv.Itoa(price)
	length := len(currencyStr)

	var formatted string

	if length <= 3 {
		formatted = "Rp. " + currencyStr
	} else {
		formatted = "Rp. " + currencyStr[:length-3] + "." + currencyStr[length-3:]
	}

	return formatted + ",00"
}

func formatQuantity(quantity int) string {
	if quantity % 1000 == 0 {
		return strconv.Itoa(quantity / 1000)
	}

	return fmt.Sprintf("%.2f", float64(quantity) / 1000)
}

func calcTotal(quantity int, price int) int {
	total := quantity * price

	return int(math.Ceil(float64(total) / 1000))
}

func updateQuantity(carts *[]Cart, cartId int, quantity int) Cart {
	for i, cart := range *carts {
		if (cart.ID == cartId) {
			remainder := quantity % cart.Product.QtyMultiple
			if ( remainder == 0) {
				(*carts)[i].Quantity = quantity
			} else {
				(*carts)[i].Quantity = int(math.Max(float64(quantity - remainder), 
					float64(cart.Product.QtyMultiple)))
			}

			return (*carts)[i]
		}
	}

	return Cart{}
}

func main() {
	carts := []Cart{
		{
			ID: 1,
			Product: Product{
				ID: 1,
				Name: "Kentang Goreng",
				MinQty: 1300,
				QtyMultiple: 1300,
				Price: 10000,
			},
			Quantity: 2600,
		},
		{
			ID: 2,
			Product: Product{
				ID: 3,
				Name: "Susu Diamond",
				MinQty: 1000,
				QtyMultiple: 1000,
				Price: 11700,
			},
			Quantity: 1500,
		},
		{
			ID: 3,
			Product: Product{
				ID: 3,
				Name: "Beras Kencur",
				MinQty: 2030,
				QtyMultiple: 2030,
				Price: 12345,
			},
			Quantity: 500,
		},
	}

	initialCartTotal := 0
	for _, cart := range carts {
		initialCartTotal += calcTotal(cart.Quantity, cart.Product.Price)
	}

	http.Handle("/", templ.Handler(root(cartList(carts, initialCartTotal))))

	http.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		cartId, err := strconv.Atoi(r.FormValue("cartId")); if err != nil {
			showError(err.Error(), "500").Render(r.Context(), w)
			return
		}
		quantity, err := strconv.Atoi(r.FormValue("quantity")); if err != nil {
			showError(err.Error(), "500").Render(r.Context(), w)
			return
		}

		updatedCart := updateQuantity(&carts, cartId, quantity)

		w.Header().Add("HX-Trigger", "totalUpdate")

		cartCounter(updatedCart).Render(r.Context(), w)
	})

	http.HandleFunc("/cart/total", func(w http.ResponseWriter, r *http.Request) {
		total := 0
		for _, cart := range carts {
			total += calcTotal(cart.Quantity, cart.Product.Price)
		}

		cartTotal(total).Render(r.Context(), w)
	})

	http.ListenAndServe(":3000", nil)
}