# Wait provider for Terraform

This is a provider that helps you wait for certain events to happen in Terraform.

## Installation

Currently this provider is only available in source form. To install it you have to have a working Go 
build environment.

To build:

```
go build
```

To run:

Copy the `terraform-provider-wait` or `terraform-provider-wait.exe` file into your project directory or in
the directory next to Terraform.

## Wait for TCP server to become available

```hcl-terraform
data "wait_tcp" "myserver" {
  host: "127.0.0.1",
  port: 8080
}

resource "your_resource" "test" {
  depends_on = [data.wait_tcp.myserver]
}
```

## Wait for HTTP server to become available

```hcl-terraform
data "wait_http" "myserver" {
  url: "http://localhost:8080"
}

resource "your_resource" "test" {
  depends_on = [data.wait_http.myserver]
}
```

Optional parameters:

- `ca_certificate_pem`: CA certificate of the server in PEM format.
- `skip_certificate_verify`: Completely skip certificate validation.
- `expect_content`: Only pass if the response contains the expected content.