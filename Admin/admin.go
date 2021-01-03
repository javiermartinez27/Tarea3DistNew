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

var consistencia [][]string

func sendToDNS(puerto string, accion string, busca string) string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se logró conectar: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	if busca == "normal" { //el comando se manda a un puerto en el que se sabe que está el registro
		response, err := c.RecibirDeAdmin(context.Background(), &chat.Message{Mensaje: accion})
		if err != nil {
			log.Fatalf("Error al conectar con el DNS: %s", err)
		}
		return response.Mensaje
	} else { //Hay que buscar el registro en algún DNS
		response, err := c.BuscaRegistro(context.Background(), &chat.Message{Mensaje: accion})
		if err != nil {
			log.Fatalf("Error al conectar con el DNS: %s", err)
		}
		return response.Mensaje
	}
	return "aqui no llego hehe"
}

func sendAccion(accion string) string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se logró conectar: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.RecibirDeAdmin(context.Background(), &chat.Message{Mensaje: accion})
	if err != nil {
		log.Fatalf("Error al tratar de conectar con el Broker: %s", err)
	}
	return response.Mensaje
}

func main() {
	fmt.Println("Bienvenido al administrador; ingrese el comando que neecsita, o ingrese 'exit' para salir:\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		comando := scanner.Text()
		if comando == "exit" {
			fmt.Println("Saliendo...")
			break
		}

		check := strings.Split(comando, " ")
		if check[0] == "create" {
			puerto := sendAccion(comando)
			reloj := sendToDNS(puerto, comando, "normal")
			fmt.Println(reloj)
			paraAppend := []string{check[1], puerto, reloj}
			consistencia = append(consistencia, paraAppend)
		} else if check[0] == "update" || check[0] == "delete" {
			encontrado := false
			for _, s := range consistencia { //Se busca si ya creo el registro
				if s[0] == check[1] {
					encontrado = true
					puerto := s[1]
					reloj := sendToDNS(puerto, comando, "normal")
					fmt.Println(reloj)
					break
				}
			}
			if encontrado == false { //se busca en los DNS
				reloj2 := "no encontrado"
				reloj2 = sendToDNS(":9001", comando, "busca")
				if reloj2 == "no encontrado" {
					reloj2 = sendToDNS(":9002", comando, "busca")
					if reloj2 == "no encontrado" {
						reloj2 = sendToDNS(":9003", comando, "busca")
						if reloj2 == "no encontrado" {
							fmt.Println("El registro no existe, puede crearlo usando 'create nombre.dominio IP'")
						} else {
							fmt.Println(reloj2)
						}
					} else {
						fmt.Println(reloj2)
					}
				} else {
					fmt.Println(reloj2)
				}
			}
		} else {
			fmt.Println("Por favor, ingrese un comando válido")
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Se encontró un error:", err)
	}
}
