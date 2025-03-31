package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	JoinedAt  time.Time `json:"joined_at"`
	IsAdmin   bool      `json:"is_admin"`
}

// Product represents a product in the catalog
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	Category    string    `json:"category"`
	InStock     bool      `json:"in_stock"`
	CreatedAt   time.Time `json:"created_at"`
	Tags        []string  `json:"tags"`
}

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ProductID string  `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
}

// Cart represents a user's shopping cart
type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
	Total  float64    `json:"total"`
}

// MockData provides sample data for the demo
var MockData = struct {
	CurrentUser User
	Products    []Product
	UserCart    Cart
}{
	CurrentUser: User{
		ID:        "user123",
		Name:      "Jane Smith",
		Email:     "jane@example.com",
		AvatarURL: "https://via.placeholder.com/40x40?text=User",
		JoinedAt:  time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		IsAdmin:   true,
	},
	Products: []Product{
		{
			ID:          "prod1",
			Name:        "Mechanical Keyboard",
			Description: "High-quality mechanical keyboard with customizable switches",
			Price:       129.99,
			ImageURL:    "https://via.placeholder.com/300x200?text=Keyboard",
			Category:    "electronics",
			InStock:     true,
			CreatedAt:   time.Now().AddDate(0, -2, 0),
			Tags:        []string{"keyboard", "gaming", "mechanical"},
		},
		{
			ID:          "prod2",
			Name:        "Ergonomic Mouse",
			Description: "Comfortable ergonomic mouse designed for long usage",
			Price:       59.99,
			ImageURL:    "https://via.placeholder.com/300x200?text=Mouse",
			Category:    "electronics",
			InStock:     true,
			CreatedAt:   time.Now().AddDate(0, -1, -15),
			Tags:        []string{"mouse", "ergonomic", "office"},
		},
		{
			ID:          "prod3",
			Name:        "Ultrawide Monitor",
			Description: "34-inch curved ultrawide monitor for immersive experience",
			Price:       349.99,
			ImageURL:    "https://via.placeholder.com/300x200?text=Monitor",
			Category:    "electronics",
			InStock:     false,
			CreatedAt:   time.Now().AddDate(0, -3, -5),
			Tags:        []string{"monitor", "ultrawide", "screen"},
		},
		{
			ID:          "prod4",
			Name:        "Desk Lamp",
			Description: "Adjustable desk lamp with multiple brightness settings",
			Price:       39.99,
			ImageURL:    "https://via.placeholder.com/300x200?text=Lamp",
			Category:    "home",
			InStock:     true,
			CreatedAt:   time.Now().AddDate(0, 0, -20),
			Tags:        []string{"lamp", "desk", "lighting"},
		},
	},
}

func init() {
	// Initialize cart with a few items
	MockData.UserCart = Cart{
		UserID: MockData.CurrentUser.ID,
		Items: []CartItem{
			{
				ProductID: "prod1",
				Product:   MockData.Products[0],
				Quantity:  1,
			},
			{
				ProductID: "prod2",
				Product:   MockData.Products[1],
				Quantity:  2,
			},
		},
		Total: MockData.Products[0].Price + (MockData.Products[1].Price * 2),
	}
}