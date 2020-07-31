package main

import "my-filestore-server/service/apigw/route"

func main() {
	route.SetUp().Run(":8080")
}
