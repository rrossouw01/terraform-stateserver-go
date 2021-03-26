# terraform-stateserver-go

POC to test terraform http backend.  

A few different proof of concepts I tried can be described as follow
---
1. python implmentation storing state file as a file in the local folder.  project/code here: https://github.com/rrossouw01/terraform-stateserver-py
2. golang implementation storing state file as a file in the local folder. code in this project in file server.go
3. golang implementation storing state in mongodb collection. code is in this project in file server-mongodb.go 

References
---
- Initial source/idea from this project https://github.com/MerlinDMC/go-terraform-stateserver
- Use any of the self-sgined howto's. I used this https://github.com/denji/golang-tls
- Created subfolder states to match URL in terraform config

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

http server example for go and state stored in local folder/file
---
````bash
$ go run server.go -data_path=./ -listen_address=192.168.1.235:8080
````

https server example for go and state stored in local folder/file
---
````bash
$ sudo go run server.go -certfile="server.crt" -keyfile="server.key" -data_path=./ -listen_address=192.168.1.235:443
````

http server example for go and state stored in mongodb document
---
````bash
$ go run server-mongodb.go -listen_address=192.168.1.235:443
````

https server example for go and state stored in mongodb document
---
NOTE: this is TBD not yet done
````bash
$ go run server-mongodb.go -certfile="server.crt" -keyfile="server.key" -listen_address=192.168.1.235:443
````

TODO
---
1. Add terraform lock configuration using configuration like this example:
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
