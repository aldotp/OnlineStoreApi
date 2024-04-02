# Online Store API Documentation

Welcome to the Online Store API Documentation! This document provides information on how to interact with the Online Store API, which offers various functionalities for managing an online store.

## Requirements

To use the Online Store API, ensure you have the following:
- MySQL: Required for storing product information, user data, and transaction details.
- Redis: Used for caching and improving API performance.
- Docker: Required for containerizing the API application.
- Postman: Required for API testing.

## Features

1. **User Authentication (Login and Registration)**
   - **Endpoint:** `/login`, `/register`
   - **Method:** POST
   - **Description:** Provides endpoints for user login and registration.

2. **View Product List by Product Category**
   - **Endpoint:** `/products?category=`
   - **Method:** GET
   - **Description:** Retrieves a list of products based on the specified category.

3. **Get All Categories**
   - **Endpoint:** `/categories`
   - **Method:** GET
   - **Description:** Retrieves a list of all product categories.

4. **Add Product to Shopping Cart**
   - **Endpoint:** `/cart`
   - **Method:** POST
   - **Description:** Adds a product to the user's shopping cart.

5. **View Shopping Cart**
   - **Endpoint:** `/cart`
   - **Method:** GET
   - **Description:** Retrieves the list of products in the user's shopping cart.

6. **Delete Product from Shopping Cart**
   - **Endpoint:** `/cart/{productID}`
   - **Method:** DELETE
   - **Description:** Removes a product from the user's shopping cart.

7. **Checkout and Make Payment**
   - **Endpoint:** `/checkout`
   - **Method:** POST
   - **Description:** Allows the user to complete the purchase and make payment transactions.

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
