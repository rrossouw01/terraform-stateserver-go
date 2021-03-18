# terraform-stateserver-go

Updated the code from this project https://github.com/MerlinDMC/go-terraform-stateserver
Use any of the self-sgined howto's. I used this https://github.com/denji/golang-tls

Run http server example:
$ go run tf-stateserver.go -data_path=./ -listen_address=192.168.1.235:8080

Run https server example:
$ sudo go run tf-stateserver.go -certfile="server.crt" -keyfile="server.key" -data_path=./ -listen_address=192.168.1.235:443

Example terraform syntax for http:
➜  terraform-poc tail -8 main.tf 
#}
## POC to test a go server
terraform {
  backend "http" {
    address = "http://192.168.1.235:8080/states/terraform.tfstate"
  }
}

Example terraform syntax for https:
➜  terraform-poc tail -8 main.tf 
#}
## POC to test a go server
terraform {
  backend "http" {
    address = "https://192.168.1.235/states/terraform.tfstate"
    skip_cert_verification = true
  }
}
