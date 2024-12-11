# eCommerce-Store

This is a full-fledged eCommerce platform built using Go and the Gin web framework. It allows users to register, log in, manage their cart, apply for discounts and place orders. Platform admin can view sales analytics.

## Features

- **User Regisration**: Register and log in using a username.
- **Product Management**: Sellers can add, update, and view products.
- **Cart Management**: Users can add products to their cart and view it.
- **Order Management**: Users can place orders, apply discount coupons, and view order details.
- **Admin Analytics**: Admins can view analytics such as total items sold, total purchase amount, and discount coupons applied.

## Technologies Used

- **Go (Golang)**: Backend programming language.
- **Gin**: Web framework for building the API.
- **Postman**: API testing tool.

## Installation

### Prerequisites

1. Install **Go 1.18+** from the [official Go website](https://golang.org/dl/).
2. Install **Git** from the [official Git website](https://git-scm.com/).

### Steps to Run Locally

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/tejash-9/eCommerce-store.git
   cd eCommerce-store
   ```
2. **Install Dependencies**:
    ```bash
    go mod tidy
    ```
3. **Set Up Environment Variables**:
    ```bash
    DISCOUNT_INTERVAL=5
    PORT=8080
    GIN_MODE=release
    ```
4. **Run the Application**:
    ```bash
    go run main.go
    ```

## Contributing

We welcome contributions! If you would like to contribute to this project, please follow these steps:

1. **Fork the repository** by clicking on the "Fork" button in the top right corner of the repository page.
2. **Clone your fork** to your local machine:
    ```bash
    git clone https://github.com/your-username/eCommerce-store.git
    ```
3. **Create a new branch** for your changes
   ```bash
    git checkout -b your-branch-name
    ```
4. **Make your changes and Commit**
    ```bash
    git add .
    git commit -m "Description of your changes"
    ```
5. **Push your changes** to your forked repository
    ```bash
    git push origin your-branch-name
    ```
6. **Open a pull request** on the original repository, describing the changes you've made.



