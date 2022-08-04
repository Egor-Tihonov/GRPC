package main

import (
	pb "awesomeProjectGRPC/proto"
	"bufio"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("dont connect")
	}
	defer conn.Close()
	client := pb.NewCRUDClient(conn)
	fmt.Println("\nWelcome to a simple gRPC/PotrgreSQL based app that performs CRUD",
		" operations!")
	fmt.Println("Enter the one of the following choices below:")
	fmt.Print("1 Registration; 2 Show information about you; 3 Show information about all users; 4 to remove: ")

	choice := bufio.NewReader(os.Stdin)
	text, _ := choice.ReadString('\n')

	switch text {
	case "1\n":
		// Registration operation
		// Read the name
		err = Registration(client)
		if err != nil {
			log.Fatalf("error: ", err)
		}
	case "2\n":
		// GetAllUsers operation
		err = GetAllUsers(client)
		if err != nil {
			log.Fatalf("error: ", err)
		}

	case "3\n":

	case "4\n":

	default:
		fmt.Println("\nWrong option!")
	}
}
func Registration(client pb.CRUDClient) error {
	var choice string
	person := &pb.RegistrationRequest{}
	_, err := fmt.Scan(&person.Person.Name)
	if err != nil {
		return err
	}
	fmt.Scanf("Введите возраст: %v", person.Person.Name)
	fmt.Scanf("Вы работаете?: %v", choice)
	switch choice {
	case "нет":
		person.Person.Works = true
	case "да":
		person.Person.Works = false
	}
	fmt.Scanf("Введите пароль: %v", person.Person.Password)
	id, err := client.Registration(context.Background(), person)
	if err != nil {
		return err
	}
	fmt.Printf("Your id: %v", id)
	return nil
}
func GetAllUsers(client pb.CRUDClient) error {
	q := &pb.GetAllUsersRequest{}
	stream, err := client.GetAllUsers(context.Background(), q)
	if err != nil {
		return err
	}
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error ", err)
		}
		fmt.Printf("person: %v\n", row)
	}
	return nil
}
