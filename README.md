# Wait provider for Terraform

This is a provider that helps you wait for certain events to happen in Terraform.

## Download

You can download the latest release [from the releases section](https://github.com/janoszen/terraform-provider-wait/releases).

## Installation

Unpack the archive for your platform and place it into one of the following directories:

- Directly into your Terraform project directory
- On Linux and MacOS into `~/.terraform/plugins`
- On Windows into `%APPDATA%\terraform.d\plugins`

## Wait for TCP server to become available

```hcl-terraform
data "wait_tcp" "myserver" {
  host = "127.0.0.1",
  port = 8080
}

resource "your_resource" "test" {
  depends_on = [data.wait_tcp.myserver]
}
```

## Wait for HTTP server to become available

```hcl-terraform
data "wait_http" "myserver" {
  url = "http://localhost:8080"
}

resource "your_resource" "test" {
  depends_on = [data.wait_http.myserver]
}
```

Optional parameters:

- `ca_certificate_pem`: CA certificate of the server in PEM format.
- `skip_certificate_verify`: Completely skip certificate validation.
- `expect_status`: Only pass if the HTTP status matches this. Defaults to 200. Set to 0 to skip check.
- `expect_content`: Only pass if the response contains the expected content.

Response parameters:

- `response_status`: contains the HTTP status code in the response
- `response_body`: contains the HTTP response body as string

**Note:** certificate validation does not work without `ca_certificate_pem` on Windows due to 
a bug in Golang ( golang/go#16736 ).