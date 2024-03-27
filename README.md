# **Go Fiber Fish E-Commerce**

This project is a backend API for a Fish Store themed e-commerce platform designed to facilitate product browsing, user management & authorization, shopping cart operations, and payment processing.

## Getting Started

### Prerequisites

Before setting up the project, please ensure you have the following prerequisites:

#### Software Requirements

- **Go**: Version 1.21.4 or later.
- **PostgreSQL**: Required for the database.
- **Midtrans Account**: Needed for payment processing functionalities.

#### Key Dependencies

The project uses `go mod` for dependency management. Key dependencies include:

| Dependency             | Package                             |
|------------------------|-------------------------------------|
| Web framework          | `github.com/gofiber/fiber/v2`       |
| ORM                    | `github.com/jinzhu/gorm`            |
| Cryptography           | `golang.org/x/crypto`               |
| JWT handling           | `github.com/golang-jwt/jwt/v4`      |
| PostgreSQL integration | `github.com/lib/pq`                 |
| Payment processing     | `github.com/midtrans/midtrans-go`   |

Run `go get` to install dependencies *independently*. For example, to install `Fiber`:

```bash
go get -u github.com/gofiber/fiber/v2
```

### Installation

1. **Clone the repository**

Start by cloning the repository to your local machine:

```bash
git clone https://github.com/ricky-kiva/be-go-fiber-ecommerce.git
cd be-go-fiber-ecommerce
```

2. **Set Up the PostgreSQL Database**

Make sure PostgreSQL is installed and running. Create a new database for the project:

```bash
CREATE DATABASE your_db_name;
```

Modify the DATABASE_URL in your `.env` file to reflect your PostgreSQL connection details:

```bash
DATABASE_URL="postgresql://username:password@localhost/dbname?sslmode=disable"
```

3. **Register for a Midtrans Account**
If you haven't already, sign up for a [Midtrans](https://www.midtrans.com/) account. Navigate to the settings section to obtain your `MIDTRANS_SERVER_KEY`.

4. **Configure Environmental Variables**
Within the root directory of your project, create a `.env` file and fill it with your configuration settings:

```bash
DATABASE_URL="postgresql://username:password@localhost/dbname?sslmode=disable"
JWT_SECRET="ANY_CUSTOM_KEY"
MIDTRANS_SERVER_KEY="YOUR_MIDTRANS_SERVER_KEY"
ENV="DEV"
```

5. **Install Dependencies**

Execute the following command to install the project dependencies:

```bash
go mod tidy
```

This will download and install the necessary Go modules as specified in your go.mod file.

6. **Launch the Application**

With the setup complete and the database operational, you can start the application by running:

```bash
go run main.go
```

7. **Interacting with the Backend**

Your server should now be up and running, accessible via `http://localhost:3000` or another port if configured differently.

### Entity-Relationship Diagram (ERD)

![ERD for Go Fiber Fish E-Commerce Project](https://github.com/ricky-kiva/be-go-fiber-ecommerce/blob/main/assets/rickyf_ERD.png?raw=true)

## Usage

This section provides guidance on how to interact with the API, including constructing requests and understanding responses.

### Base Url

All API requests are made to the base URL, which is structured as follows for version 1 of the API:

```plaintext
http://localhost:3000/v1
```

### Making requests

- **GET**
Example, to fetch products by category ID:

```plaintext
GET /v1/products/:categoryId
```

Replace `:categoryId` with the actual ID of the category you wish to query.

- **POST**
Example, to register a new user:

```plaintext
POST /v1/users/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

Ensure the email is valid and the password meets the specified criteria.

- **PUT**
Example, to update a product's information:

```plaintext
PUT /v1/products/:productId
Content-Type: application/json
Authorization: Bearer YOUR_JWT_TOKEN

{
  "name": "Updated Product Name",
  "price": 20.99
}
```

Replace :productId with the product ID you wish to update. The request body should contain the updated product data. This endpoint requires authentication.

- **DELETE**
Example, to delete an item from the shopping cart:

```plaintext
DELETE /v1/cart/:productId
Authorization: Bearer YOUR_JWT_TOKEN
```

Replace `:productId` with the ID of the product you wish to remove. For this endpoint, the request must be authenticated with a valid JWT token.

### Response Format

Responses from the API can generally be categorized into three types: success, fail, and error.

- **Success**:
```json
{
  "status": "success",
  "message": "Cart updated successfully"
}
```

- **Fail**:
```json
{
  "status": "fail",
  "message": "Incorrect password"
}
```

- **Error**
```json
{
  "status": "error",
  "message": "Internal Server Error"
}
```

### Authentication

Authentication is required for certain endpoints. Here's how to register a user and log in to receive an authentication token:

- **Register a User**:
```plaintext
POST /v1/users/register
```

Provide an email and password in the request body.

- **Login**:
```plaintext
POST /v1/users/login
```

Upon successful login, you'll receive a JWT token in the response. Use this token in the `Authorization` header for subsequent requests that require authentication.

### Endpoints

Below is a list of all available endpoints with brief descriptions.

| Method | Endpoint                        | Description                                                    |
|--------|---------------------------------|----------------------------------------------------------------|
| GET    | `/`                             | Displays information about the project.                        |
| POST   | `/v1/register`                  | Registers a new user.                                          |
| POST   | `/v1/login`                     | Authenticates a user and returns a token.                      |
| GET    | `/v1/cart`                      | Retrieves the current user's shopping cart.                    |
| POST   | `/v1/cart`                      | Adds a new item to the shopping cart.                          |
| DELETE | `/v1/cart/items/:productID`     | Removes an item from the shopping cart.                        |
| GET    | `/v1/cart/checkout`             | Processes the checkout for the current user's cart.            |
| GET    | `/v1/cart/pay`                  | Initiates a payment process for the items in the shopping cart.|
| GET    | `/v1/products`                  | Retrieves a list of all products.                              |
| GET    | `/v1/products/info`             | Retrieves details of a specific product by its ID.             |
| GET    | `/v1/products/categories/:categoryID` | Retrieves a list of products filtered by a specific category ID.|
| GET    | `/v1/products/categories`       | Retrieves a list of all product categories.                            |

For more detailed information, including request parameters, headers, response objects, and examples, please refer to our: [**Postman API Documentation**](https://www.postman.com/rickyslash/workspace/syn-fish-ecommerce-go-fiber/documentation/26442482-22509c1c-ad14-4182-9aa2-5ebc91b17cd6)

## Example

This section outlines a typical user interaction flow with our API, from signing up to making a payment.

1. **Register**

This step creates a new user account. Make sure to replace the email and password with your desired credentials.

- **Request**
```plaintext
POST /v1/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

- **Response**
```plaintext
{
  "status": "success",
  "message": "User registered successfully."
}
```

2. **Login (get `Authorization` token)**

The response after login includes a JWT token that you must use in subsequent requests to authenticate.

- **Request**
```plaintext
POST /v1/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

- **Response**
```plaintext
{
  "token": "YOUR_JWT_TOKEN"
}
```

3. **Use the Authorization Token**

For any request that requires authentication, add the obtained JWT token in the `Authorization` field of the request headers.

4. **Get All Items**

This retrieves a list of all products. You can choose which products to add to your cart.

- **Request**
```plaintext
GET /v1/products
Authorization: Bearer YOUR_JWT_TOKEN
```

- **Response**
```plaintext
[
  {
    "id": 1,
    "name": "Product Name",
    "description": "Product Description",
    "price": 100
  }
]
```

5. **Add to Cart**

Repeat this step for each product you wish to add to your cart, modifying the `productId` and `quantity` as necessary.

- **Request**
```plaintext
POST /v1/cart
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "productId": 1,
  "quantity": 2
}
```

- **Response**
```plaintext
{
  "status": "success",
  "message": "Product added to cart successfully."
}
```

6. **Changing Mind? Delete an Item from the Cart**

If you change your mind, you can remove an item from your cart. Replace `:productId` with the ID of the product you wish to remove from your cart.

- **Request**
```plaintext
DELETE /v1/cart/items/:productId
Authorization: Bearer YOUR_JWT_TOKEN
```

- **Response**
```plaintext
{
  "status": "success",
  "message": "Item removed from cart successfully."
}
```

7. **Checkout**

Review your cart and the total amount to be paid before proceeding with payment.

- **Request**
```plaintext
GET /v1/cart/checkout
Authorization: Bearer YOUR_JWT_TOKEN
```

- **Response**
```plaintext
{
  "Total": 200,
  "Items": [
    {
      "Name": "Product Name",
      "Quantity": 2,
      "Price": 100,
      "Total": 200
    }
  ]
}
```

8. **Make a Payment**

step generates a transaction token and provides a URL to proceed with the payment through Midtrans.

- **Request**
```plaintext
GET /v1/cart/pay
Authorization: Bearer YOUR_JWT_TOKEN
```

- **Response**
```plaintext
{
  "status": "success",
  "data": {
    "Token": "payment-token",
    "RedirectUrl": "https://payment-gateway-url.com/transaction/payment-token"
  }
}
```

9. **Proceed with Payment**

Clicking on the provided URL takes you to Midtrans secure payment gateway, where you can complete the transaction.

## Support

For support, please create an issue in the [Issues tab](https://github.com/ricky-kiva/be-go-fiber-ecommerce/issues). Include a clear title and detailed description for questions or issues. We appreciate your feedback!

## Contributing

Your contributions are welcome! Here’s how to get started:

1. **Search Issues**: Check if an [issue](https://github.com/ricky-kiva/be-go-fiber-ecommerce/issues) already exists for your concern or idea.
2. **Report or Suggest**: For new ideas or bugs, please create a new issue with a descriptive title and details.
3. **Submit a Pull Request**: Feel free to submit PRs against the `main` branch for code contributions. Ensure your changes are clear and concise.

Thank you for helping to improve this project!

## License

[MIT](https://choosealicense.com/licenses/mit/)
