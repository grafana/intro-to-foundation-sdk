# Introduction to Grafana Foundation SDK

This repository contains supplementary content for a series of videos introducing the **Grafana Foundation SDK**. It includes several examples to help you get started with building and managing Grafana dashboards programmatically.

## Prerequisites
- GoLang installed on your system.
- Docker (for the complex example stack).
- Basic understanding of Grafana, Prometheus, and Loki.

## Examples

### 1. Basic Example: Generating Dashboard JSON in GoLang
This example demonstrates how to use the Grafana Foundation SDK to programmatically generate dashboard JSON using GoLang. It serves as a simple starting point for understanding the SDK's capabilities.

### 2. Advanced Example: Full Stack with Metrics and Logs
This example sets up a complete stack, including:
- **Grafana**
- **Prometheus**
- **Loki**

It also includes a dummy web application that:
- Uses **Alloy** to pull metrics and logs.
- Has the capability to generate and deploy its own dashboard.

This example showcases how to integrate multiple tools in the Grafana ecosystem and automate dashboard creation.

## Getting Started
To explore the examples, clone this repository and follow the instructions provided in each example's directory.

```bash
git clone https://github.com/grafana/intro-to-foundation-sdk.git
cd intro-to-foundation-sdk
```

## Instructions for Running the Examples

### 1. Basic Example
To generate a Grafana dashboard JSON file, follow these steps:
1. Navigate to the `basic-example` directory:
    ```bash
    cd basic-example
    ```
2. Run the following command to generate the dashboard JSON in the terminal output:
    ```bash
    go run main.go
    ```
3. To start Grafana using Docker Compose, run:
    ```bash
    docker compose up
    ```
4. Once Grafana is running, access it at [http://localhost:3000](http://localhost:3000).
5. Import the generated dashboard JSON into Grafana manually.

This process allows you to visualize the generated dashboard in a local Grafana instance.

### 2. Generate and Deploy Example
To spin up the entire stack, follow these steps:
1. Navigate to the `generate-and-deploy-example` directory:
    ```bash
    cd generate-and-deploy-example
    ```
2. Use Docker Compose to build and start the stack:
    ```bash
    docker compose up --build
    ```
3. Once the stack is running, you can access the following services:
    - **Grafana**: [http://localhost:3000](http://localhost:3000)
    - **Web Application**: [http://localhost:5001](http://localhost:5001)
    - **Prometheus**: [http://localhost:9090](http://localhost:9090)
    - **Alloy**: [http://localhost:12345](http://localhost:12345)

The k6 script will automatically generate web traffic, and the web service dashboard will be provisioned in Grafana.

## Contributing
Contributions are welcome! Feel free to open issues or submit pull requests to improve the examples or add new ones.

---
Happy coding!