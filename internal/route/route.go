package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aldotp/OnlineStore/internal/config"
	"github.com/aldotp/OnlineStore/internal/handler"
	"github.com/aldotp/OnlineStore/internal/middleware"
	"github.com/aldotp/OnlineStore/internal/repositories"
	"github.com/aldotp/OnlineStore/internal/services"
	"github.com/gorilla/mux"
)

type Route struct {
	config *config.BootstrapConfig
}

func NewRouter(conf *config.BootstrapConfig) *Route {
	return &Route{
		config: conf,
	}
}

func (route *Route) Router() *mux.Router {

	// instance
	redisInstance := config.NewRedisClient(route.config.Viper)

	// repositories
	userRepo := repositories.NewUserRepository(route.config.DB)
	productRepo := repositories.NewProductRepository(route.config.DB)
	categoryRepo := repositories.NewCategoryRepository(route.config.DB)
	cartRepo := repositories.NewCartRepository(route.config.DB)
	cartItemsRepo := repositories.NewCartItemsRepository(route.config.DB)
	orderRepo := repositories.NewOrderRepository(route.config.DB)
	orderDetailRepo := repositories.NewOrderDetailRepository(route.config.DB)

	// services
	userService := services.NewUser(userRepo, route.config)
	productService := services.NewProduct(productRepo, categoryRepo, redisInstance)
	paymentService := services.NewPayment()
	cartService := services.NewCart(cartRepo, cartItemsRepo, productRepo)
	checkoutService := services.NewCheckout(orderRepo, cartRepo, orderDetailRepo, paymentService)
	categoryService := services.NewCategory(redisInstance, categoryRepo)

	// handlers
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService, categoryRepo, redisInstance)
	cartHandler := handler.NewCartHandler(cartService, cartRepo, cartItemsRepo, productRepo)
	checkoutHandler := handler.NewCheckoutHandler(checkoutService)

	// router
	r := mux.NewRouter()

	v1 := r.PathPrefix("/v1").Subrouter()
	api := v1.PathPrefix("/api").Subrouter()

	public := api.PathPrefix("/public").Subrouter()
	protected := api.PathPrefix("/protected").Subrouter()

	public.HandleFunc("/login", userHandler.LoginUser).Methods("POST")
	public.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")

	jwt := middleware.NewJWT(route.config)
	protected.Use(jwt.AuthMiddleware)

	protected.HandleFunc("/products/category/{id}", productHandler.GetProductsByCategory).Methods("GET")
	protected.HandleFunc("/product/{id}", productHandler.UpdateProduct).Methods("PUT")
	protected.HandleFunc("/product/{id}", productHandler.GetProductByID).Methods("GET")
	protected.HandleFunc("/product", productHandler.StoreProducts).Methods("POST")
	protected.HandleFunc("/product/{id}", productHandler.DeleteProduct).Methods("DELETE")
	protected.HandleFunc("/products", productHandler.GetProducts).Methods("GET")

	protected.HandleFunc("/category/{id}", categoryHandler.GetCategoryByID).Methods("GET")
	protected.HandleFunc("/categories", categoryHandler.GetCategories).Methods("GET")
	protected.HandleFunc("/category/{id}", categoryHandler.DeleteCategory).Methods("DELETE")
	protected.HandleFunc("/category/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	protected.HandleFunc("/category", categoryHandler.StoreCategory).Methods("POST")

	protected.HandleFunc("/cart", cartHandler.Cart).Methods("GET")
	protected.HandleFunc("/cart", cartHandler.AddToCart).Methods("POST")
	protected.HandleFunc("/cart/product/{id}", cartHandler.DeleteProductFromCart).Methods("DELETE")
	protected.HandleFunc("/cart", cartHandler.EmptyCart).Methods("DELETE")
	protected.HandleFunc("/cart/product/{id}", cartHandler.ModifyCart).Methods("PUT")

	protected.HandleFunc("/checkout", checkoutHandler.CheckoutHandler).Methods("POST")
	protected.HandleFunc("/checkout/history", checkoutHandler.CheckoutHistory).Methods("GET")

	return r
}

func (route *Route) Run() {
	router := route.Router()
	log.Printf("Server is running on %s:%s", route.config.Host, route.config.WebPort)
	http.ListenAndServe(fmt.Sprintf(":%s", route.config.WebPort), router)
}
