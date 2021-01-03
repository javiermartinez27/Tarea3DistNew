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

func sendToBroker(accion string) string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.10.28.154:9000", grpc.WithInsecure())
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
			fmt.Println("Informacion de los dominios que ha solicitado:")
			fmt.Printf("%+q", consistencia)
			fmt.Println("")
			fmt.Println("Saliendo...")
			break
		}
		check := strings.Split(comando, " ")
		if check[0] == "get" { //Funcion encargada de mandar los get al broker para despues recibir los mensajes desde este.
			mensaje := sendToBroker(comando)
			if len(consistencia) == 0 {
				paraAppend := []string{check[1] + " " + mensaje}
				consistencia = append(consistencia, paraAppend)
				fmt.Println("Recibido:")
				fmt.Printf("%+q", paraAppend)
				fmt.Println("")
			} else {
				actualizar := 0 // nos dice si debemos actualizar alguna informacion en memoria, de los dominios que se han solicitado.
				aux := ""       //sera usado para guardar el slice que luego debe coincidir en el otro for, asi actuliza correctamente el slice correspondiente.
				for _, s := range consistencia {
					sComoString := strings.Join(s, " ")
					if strings.Contains(sComoString, check[1]) { //funcion encargada de actualizar dominio
						actualizar = 1
						aux = sComoString
						break
					}

				}
				for i, s := range consistencia {
					sComoString := strings.Join(s, " ")
					if actualizar == 1 && aux == sComoString { //si se debe actualizar pasa por aqui, sComoString debe coincidir con el sComoString del for pasado para que este actualice
						paraAppend := []string{check[1] + " " + mensaje}
						consistencia[i] = paraAppend
						fmt.Printf("%+q", paraAppend)
						fmt.Println("")
						break
					}
					if actualizar != 1 { //si no se cumple lo anterior simplemente se agrega.
						// fmt.Println("lo agregue :(")
						paraAppend := []string{check[1] + " " + mensaje}
						consistencia = append(consistencia, paraAppend)
						// fmt.Println("Recibido en el else:")
						fmt.Printf("%+q", paraAppend)
						fmt.Println("")
						break
					}

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
