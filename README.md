# Online Store API Documentation

Welcome to the Online Store API Documentation! This document provides information on how to interact with the Online Store API, which offers various functionalities for managing an online store.

## Requirement

To use the Online Store API, ensure you have the following:
-- MySQL: Required for storing product information, user data, and transaction details.
-- Redis: Used for caching and improving API performance.
-- Docker: Required for containerizing the API application.
-- Postman: Required for API testing.

## Features

1. User Authentication (Login and Registration)
Endpoint: /login, /register
Method: POST
Description: Provides endpoints for user login and registration.

2. View Product List by Product Category
Endpoint: /products?category=
Method: GET
Description: Retrieves a list of products based on the specified category.

3. Get All Category
Endpoint: /categories
Method: GET
Description: Retrives a list of categories

4. Add Product to Shopping Cart
Endpoint: /cart
Method: POST
Description: Adds a product to the user's shopping cart.

5. View Shopping Cart
Endpoint: /cart
Method: GET
Description: Retrieves the list of products in the user's shopping cart.

6. Delete Product from Shopping Cart
Endpoint: /cart/{productID}
Method: DELETE
Description: Removes a product from the user's shopping cart.

7. Checkout and Make Payment
Endpoint: /checkout
Method: POST
Description: Allows the user to complete the purchase and make payment transactions.

## How to Use

- Using Docker-compose.yaml, run: docker-compose up
- Server running on port 8080 (local) 
or
- Using local environment (setup and MySQL, Redis first)
- Create .env using example.env 
- make run (to start the api) 


## TestCase Testing  

- Export file collection and environment .json in folder collection to Postman Apps

## Deployed Apps

URL: https://onlinestore-golang-api-4aeba344fcf6.herokuapp.com


