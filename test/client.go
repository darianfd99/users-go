package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/darianfd99/users-go/pkg/proto"
	"google.golang.org/grpc"
)

const (
	address  = "localhost:8080"
	timeout  = 100
	filename = "data.json"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := proto.NewUsersServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var usersList []*proto.RequestUser

	err = json.Unmarshal(file, &usersList)
	if err != nil {
		log.Fatal(err)
	}

	errorsChannel := make(chan error)
	for _, user := range usersList {

		go func(user *proto.RequestUser) {
			_, err = c.CreateUser(ctx, user)
			errorsChannel <- err
		}(user)

	}

	for i := 0; i < len(usersList); i++ {
		err = <-errorsChannel
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Get all:")

	protoUsersList, err := c.GetAllUsers(context.Background(), &proto.Null{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(protoUsersList.GetUserList())

	protoUsersList, err = c.GetAllUsers(context.Background(), &proto.Null{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(protoUsersList.GetUserList())

	time.Sleep(63 * time.Second)

	protoUsersList, err = c.GetAllUsers(context.Background(), &proto.Null{})
	if err != nil {
		log.Fatal(err)
	}

	protoUsersList, err = c.GetAllUsers(context.Background(), &proto.Null{})
	if err != nil {
		log.Fatal(err)
	}

	protoUsersList, err = c.GetAllUsers(context.Background(), &proto.Null{})
	if err != nil {
		log.Fatal(err)
	}

}
