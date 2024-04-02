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
	productService := services.NewProduct(productRepo, categoryRepo)
	paymentService := services.NewPayment()
	cartService := services.NewCart(cartRepo, cartItemsRepo)
	checkoutService := services.NewCheckout(orderRepo, cartRepo, orderDetailRepo, paymentService)

	// handlers
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryRepo, redisInstance)
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
	protected.HandleFunc("/products/{category}", productHandler.GetProductsByCategory).Methods("GET")
	protected.HandleFunc("/product", productHandler.StoreProducts).Methods("POST")

	protected.HandleFunc("/category", categoryHandler.StoreCategory).Methods("POST")
	protected.HandleFunc("/categories", categoryHandler.GetCategory).Methods("GET")

	protected.HandleFunc("/cart", cartHandler.Cart).Methods("GET")
	protected.HandleFunc("/cart", cartHandler.AddToCart).Methods("POST")
	protected.HandleFunc("/cart/{productID}", cartHandler.DeleteProductFromCart).Methods("DELETE")

	protected.HandleFunc("/checkout", checkoutHandler.CheckoutHandler).Methods("POST")

	return r
}

func (route *Route) Run() {
	router := route.Router()
	log.Printf("Server OnlineStoreAPI is running on %s:%s", route.config.Host, route.config.WebPort)
	http.ListenAndServe(fmt.Sprintf(":%s", route.config.WebPort), router)
}
