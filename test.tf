provider "wait" {

}

data "wait_tcp" "test" {
  host = "pasztor.at"
  port = 443
}

output "test" {
  value = data.wait_tcp.test.host
}
