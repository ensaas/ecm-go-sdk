# ECM sdk for go
## Usage

1. cd gopath

2. git clone https://gitlab.wise-paas.com/WISE-PaaS-4.0-Ops/ecm-sdk-go.git

3. cd gopath/your_project

4. modify the go.mod of your project, add the  following two lines:

   require "ecm-sdk-go" v0.0.0
   replace "ecm-sdk-go" => " gopath/ecm-sdk-go"

   