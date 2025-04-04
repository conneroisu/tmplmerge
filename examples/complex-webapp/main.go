// Package main is a simple example of how to use twerge.
// It demonstrates how to use ClassMapStr to populate ClassMap.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/conneroisu/twerge/examples/complex-webapp/components"
	"github.com/conneroisu/twerge/examples/complex-webapp/models"
	"github.com/conneroisu/twerge/examples/complex-webapp/pages"

	"github.com/a-h/templ"
)

func main() {
	// No initialization needed

	// Set up the server
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/products", handleProducts)
	mux.HandleFunc("/cart", handleCart)

	// API endpoints for HTMX
	mux.HandleFunc("/api/cart/add/", handleCartAdd)
	mux.HandleFunc("/api/cart/remove/", handleCartRemove)
	mux.HandleFunc("/api/cart/increase/", handleCartIncrease)
	mux.HandleFunc("/api/cart/decrease/", handleCartDecrease)

	// Start the server
	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// For a real application, you would fetch the featured products from a database
	component := pages.HomePage(models.MockData.CurrentUser, models.MockData.Products)
	templ.Handler(component).ServeHTTP(w, r)
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	// Filter products by category if specified
	filteredProducts := models.MockData.Products
	if category != "" {
		filtered := []models.Product{}
		for _, product := range models.MockData.Products {
			if product.Category == category {
				filtered = append(filtered, product)
			}
		}
		filteredProducts = filtered
	}

	component := pages.ProductsPage(models.MockData.CurrentUser, filteredProducts, category)
	templ.Handler(component).ServeHTTP(w, r)
}

func handleCart(w http.ResponseWriter, r *http.Request) {
	component := pages.CartPage(models.MockData.CurrentUser, models.MockData.UserCart)
	templ.Handler(component).ServeHTTP(w, r)
}

// API handlers for HTMX interactions

func handleCartAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract product ID from URL
	productID := strings.TrimPrefix(r.URL.Path, "/api/cart/add/")

	// Find the product
	var product models.Product
	found := false
	for _, p := range models.MockData.Products {
		if p.ID == productID {
			product = p
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Check if product is already in cart
	for i, item := range models.MockData.UserCart.Items {
		if item.ProductID == productID {
			// Update quantity
			models.MockData.UserCart.Items[i].Quantity++
			models.MockData.UserCart.Total += product.Price

			// In a real app, you'd return an updated cart count or other data
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, "%d", len(models.MockData.UserCart.Items))
			return
		}
	}

	// Add new item to cart
	cartItem := models.CartItem{
		ProductID: productID,
		Product:   product,
		Quantity:  1,
	}

	models.MockData.UserCart.Items = append(models.MockData.UserCart.Items, cartItem)
	models.MockData.UserCart.Total += product.Price

	// In a real app, you'd return an updated cart count or other data
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%d", len(models.MockData.UserCart.Items))
}

func handleCartRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract product ID from URL
	productID := strings.TrimPrefix(r.URL.Path, "/api/cart/remove/")

	// Remove item from cart
	newItems := []models.CartItem{}
	for _, item := range models.MockData.UserCart.Items {
		if item.ProductID != productID {
			newItems = append(newItems, item)
		} else {
			models.MockData.UserCart.Total -= item.Product.Price * float64(item.Quantity)
		}
	}

	models.MockData.UserCart.Items = newItems

	// If this is an HTMX request, return empty content to indicate the element should be removed
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Otherwise redirect back to cart page
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func handleCartIncrease(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract product ID from URL
	productID := strings.TrimPrefix(r.URL.Path, "/api/cart/increase/")

	// Find and update item
	for i, item := range models.MockData.UserCart.Items {
		if item.ProductID == productID {
			models.MockData.UserCart.Items[i].Quantity++
			models.MockData.UserCart.Total += item.Product.Price

			// Return updated cart item
			if r.Header.Get("HX-Request") == "true" {
				component := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					div := fmt.Sprintf(`<div data-cart-item id="cart-item-%s">`, productID)
					_, err := w.Write([]byte(div))
					if err != nil {
						return err
					}

					err = components.CartItem(models.MockData.UserCart.Items[i]).Render(ctx, w)
					if err != nil {
						return err
					}

					_, err = w.Write([]byte("</div>"))
					return err
				})

				templ.Handler(component).ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, "/cart", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Product not found in cart", http.StatusNotFound)
}

func handleCartDecrease(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract product ID from URL
	productID := strings.TrimPrefix(r.URL.Path, "/api/cart/decrease/")

	// Find and update item
	for i, item := range models.MockData.UserCart.Items {
		if item.ProductID == productID {
			if models.MockData.UserCart.Items[i].Quantity > 1 {
				models.MockData.UserCart.Items[i].Quantity--
				models.MockData.UserCart.Total -= item.Product.Price

				// Return updated cart item
				if r.Header.Get("HX-Request") == "true" {
					component := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
						div := fmt.Sprintf(`<div data-cart-item id="cart-item-%s">`, productID)
						_, err := w.Write([]byte(div))
						if err != nil {
							return err
						}

						err = components.CartItem(models.MockData.UserCart.Items[i]).Render(ctx, w)
						if err != nil {
							return err
						}

						_, err = w.Write([]byte("</div>"))
						return err
					})

					templ.Handler(component).ServeHTTP(w, r)
					return
				}
			}

			http.Redirect(w, r, "/cart", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Product not found in cart", http.StatusNotFound)
}
