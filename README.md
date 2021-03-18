# terraform-stateserver-go

- Updated the code from this project https://github.com/MerlinDMC/go-terraform-stateserver
- Use any of the self-sgined howto's. I used this https://github.com/denji/golang-tls
- Created subfolder states to match URL in terraform config

Run http server example
---
````bash
  $ go run server.go -data_path=./ -listen_address=192.168.1.235:8080
````

Run https server example
---
````bash
  $ sudo go run server.go -certfile="server.crt" -keyfile="server.key" -data_path=./ -listen_address=192.168.1.235:443
````

Example terraform syntax for http
---
````bash
➜  terraform-poc tail -5 main.tf 
  terraform {
    backend "http" {
      address = "http://192.168.1.235:8080/states/terraform.tfstate"
    }
  }
````
Example terraform syntax for https
---
````bash
  ➜  terraform-poc tail -6 main.tf 
  terraform {
    backend "http" {
      address = "https://192.168.1.235/states/terraform.tfstate"
      skip_cert_verification = true
    }
  }
````
TODO
---
1. Add terraform lock configuration using configuration like this example used in a different setup using python and flask:
````bash
➜  terraform-poc cat main.tf 
terraform {
  backend "http" {
    address = "http://192.168.1.235:5000/terraform_state/4cdd0c76-d78b-11e9-9bea-db9cd8374f3a"
    lock_address = "http://192.168.1.235:5000/terraform_lock/4cdd0c76-d78b-11e9-9bea-db9cd8374f3a"
    lock_method = "PUT"
    unlock_address = "http://192.168.1.235:5000/terraform_lock/4cdd0c76-d78b-11e9-9bea-db9cd8374f3a"
    unlock_method = "DELETE"
  }
}
````
2. Build locking into server so each GET or POST will include lock and unlock automatically ie no terraform configuration
