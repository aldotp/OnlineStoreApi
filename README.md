# Online Store API Documentation

Welcome to the Online Store API Documentation! This document provides information on how to interact with the Online Store API, which offers various functionalities for managing an online store.

## Requirements

To use the Online Store API, ensure you have the following:
- MySQL: Required for storing product information, user data, and transaction details.
- Redis: Used for caching and improving API performance.
- Docker: Required for containerizing the API application.
- Golang: Required for running the API server.
- Postman: Required for API testing.

## Features

1. **User Authentication**
   - **Login:** `/login` (POST)
     - Description: Endpoint for user login.
   - **Registration:** `/register` (POST)
     - Description: Endpoint for user registration.

2. **Product Management**
   - **Get Products by Category:** `/products/category/{id}` (GET)
     - Description: Retrieves products based on the specified category.
   - **Update Product:** `/product/{id}` (PUT)
     - Description: Updates a product with the specified ID.
   - **Get Product by ID:** `/product/{id}` (GET)
     - Description: Retrieves product details by ID.
   - **Store Product:** `/product` (POST)
     - Description: Stores a new product.
   - **Delete Product:** `/product/{id}` (DELETE)
     - Description: Deletes a product with the specified ID.
   - **Get All Products:** `/products` (GET)
     - Description: Retrieves all products.

3. **Category Management**
   - **Get Category by ID:** `/category/{id}` (GET)
     - Description: Retrieves category details by ID.
   - **Get All Categories:** `/categories` (GET)
     - Description: Retrieves all categories.
   - **Delete Category:** `/category/{id}` (DELETE)
     - Description: Deletes a category with the specified ID.
   - **Update Category:** `/category/{id}` (PUT)
     - Description: Updates a category with the specified ID.
   - **Store Category:** `/category` (POST)
     - Description: Stores a new category.

4. **Shopping Cart Management**
   - **View Shopping Cart:** `/cart` (GET)
     - Description: Retrieves the contents of the user's shopping cart.
   - **Add Product to Cart:** `/cart` (POST)
     - Description: Adds a product to the user's shopping cart.
   - **Delete Product from Cart:** `/cart/product/{id}` (DELETE)
     - Description: Removes a product from the user's shopping cart.
   - **Empty Cart:** `/cart` (DELETE)
     - Description: Empties the user's shopping cart.
   - **Modify Cart:** `/cart/product/{id}` (PUT)
     - Description: Modifies the quantity of a product in the user's shopping cart.

5. **Checkout**
   - **Checkout and Make Payment:** `/checkout` (POST)
     - Description: Allows the user to complete the purchase and make payment transactions.
   - **View Checkout History:** `/checkout/history` (GET)
     - Description: Retrieves the user's checkout history.
     
## How to Use

### Using Docker Compose

1. Ensure Docker is installed on your system.
2. Navigate to the project directory containing the `docker-compose.yaml` file.
3. Run the command `docker-compose up` to start the server.
4. The server will be running on port 8080 locally.

### Using Local Environment

1. Set up MySQL and Redis on your local environment.
2. Create a `.env` file using the provided `example.env`.
3. Run the command `make run` to start the API.
4. The server will be running on port 8080 locally.

## Testing

- Export the collection and environment files `.json` located in the `collection` folder to Postman Apps for testing.

## Deployed App

The API is deployed and accessible at [https://onlinestore-golang-api-4aeba344fcf6.herokuapp.com](https://onlinestore-golang-api-4aeba344fcf6.herokuapp.com).

Feel free to explore and utilize the Online Store API for managing your online store efficiently!
