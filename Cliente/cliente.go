package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
	"google.golang.org/grpc"
)

func sendToBroker(accion string) string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se logró conectar: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.RecibirDeCliente(context.Background(), &chat.Message{Mensaje: accion})
	if err != nil {
		log.Fatalf("Error al tratar de conectar con el Broker: %s", err)
	}
	return response.Mensaje
}

func main() {
	fmt.Println("Bienvenido al Cliente; ingrese el comando  get nombre.dominio, o ingrese 'exit' para salir:\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		comando := scanner.Text()
		if comando == "exit" {
			fmt.Println("Saliendo...")
			break
		}
		check := strings.Split(comando, " ")
		if check[0] == "get" {
			fmt.Println(sendToBroker(comando))
		} else {
			fmt.Println("Por favor, ingrese un comando válido")
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Se encontró un error:", err)
	}
}
