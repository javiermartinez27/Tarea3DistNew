package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
	"google.golang.org/grpc"
)

func escuchar() { //Funcion que permite al DNS recibir los comandos
	puerto := "10.10.28.155:9001"
	fmt.Println("DNS escuchando en el puerto " + puerto)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", puerto))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := chat.Server{}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

//De aquí para abajo, todo es para hacer el merge de los DNS

func connToDNS(puerto string, mensaje string) []string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se logró conectar: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.Consistencia(context.Background(), &chat.Message{Mensaje: mensaje})
	if err != nil {
		log.Fatalf("Error al tratar de conectar con con los DNS: %s", err)
	}
	return response.ConsistenciaList
}

func crearRegistro(registro string, ip string) {
	registroSeparado := strings.Split(registro, ".")
	nombre := "registros_zf/registro_" + registroSeparado[1] + ".txt"
	f, err := os.OpenFile(nombre, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := f.WriteString(registro + " IN A " + ip + "\n")
	if err2 != nil {
		log.Fatal(err)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func updateReloj(registro string, dns int) {
	registroSeparado := strings.Split(registro, ".")
	nombre := "relojes/reloj_" + registroSeparado[1] + ".txt"
	if _, err := os.Stat(nombre); err == nil { //actualiza el reloj
		reloj, err := readLines(nombre)
		if err != nil {
			log.Fatal(err)
		}
		separar := strings.Split(reloj[0], ",")
		if dns == 2 {
			i, err := strconv.Atoi(separar[1])
			i++
			s := strconv.Itoa(i)
			newReloj := separar[0] + "," + s + "," + separar[2]
			err = ioutil.WriteFile(nombre, []byte(newReloj), 0644)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			i, err := strconv.Atoi(separar[2])
			i++
			s := strconv.Itoa(i)
			newReloj := separar[0] + "," + separar[1] + "," + s
			err = ioutil.WriteFile(nombre, []byte(newReloj), 0644)
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else { //primera vez que se añade algo a este registro
		f, err := os.OpenFile(nombre, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if dns == 2 {
			_, err2 := f.WriteString("0,1,0")
			if err2 != nil {
				log.Fatal(err)
			}
		} else {
			_, err2 := f.WriteString("0,0,1")
			if err2 != nil {
				log.Fatal(err)
			}
		}
	}
}

func updateRegistro(registro string, cambio string) {
	registroSeparado := strings.Split(registro, ".")
	nombre := "registros_zf/registro_" + registroSeparado[1] + ".txt"
	input, err := ioutil.ReadFile(nombre)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		registroSep := strings.Split(registro, ".")
		if strings.Contains(line, registro) {
			lineaSep := strings.Split(line, " ")
			if strings.Contains(cambio, registroSep[1]) {
				lines[i] = cambio + " IN A " + lineaSep[3]
			} else {
				lines[i] = lineaSep[0] + " IN A " + cambio
			}
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(nombre, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func borrarRegistro(registro string) {
	registroSeparado := strings.Split(registro, ".")
	nombre := "registros_zf/registro_" + registroSeparado[1] + ".txt"
	input, err := ioutil.ReadFile(nombre)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, registro) {
			lines[i] = ""
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(nombre, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func actualizaLog(archivos []string, dns int) {
	flagContains := false              //flag para ver si el DNS tiene info de algun nombre.dominio
	for _, archivo := range archivos { //itera sobre el log recibido
		logSeparado := strings.Split(archivo, "+")
		paraNombre1 := strings.Split(logSeparado[0], " ")
		paraNombre2 := strings.Split(paraNombre1[1], ".")
		nombre := "logs/log_" + paraNombre2[1] + ".txt"
		var aux []string

		if _, err := os.Stat(nombre); err == nil { //si el archivo ya existe

			input, err := ioutil.ReadFile(nombre)
			if err != nil {
				log.Fatalln(err)
			}

			lines := strings.Split(string(input), "\n")

			for _, line := range lines {
				if line != "" {
					aux = append(aux, line)
				}
			}
			for _, accion := range logSeparado {
				if accion == "" {
					continue
				}
				commandSeparado := strings.Split(accion, " ")
				for _, line := range lines {
					if strings.Contains(line, commandSeparado[1]) { //si el nombre.dominio está en el archivo, se deja solo lo que está en el nodo dominante, es decir este
						flagContains = true
					}
				}
				if flagContains == false {
					aux = append(aux, accion)
					if commandSeparado[0] == "create" {
						crearRegistro(commandSeparado[1], commandSeparado[2])
						if dns == 2 {
							updateReloj(commandSeparado[1], 2)
						} else {
							updateReloj(commandSeparado[1], 3)
						}
					} else if commandSeparado[0] == "update" {
						updateRegistro(commandSeparado[1], commandSeparado[2])
						if dns == 2 {
							updateReloj(commandSeparado[1], 2)
						} else {
							updateReloj(commandSeparado[1], 3)
						}
					} else {
						borrarRegistro(commandSeparado[1])
						if dns == 2 {
							updateReloj(commandSeparado[1], 2)
						} else {
							updateReloj(commandSeparado[1], 3)
						}
					}
				}
			}

			output := strings.Join(aux, "\n")
			err = ioutil.WriteFile(nombre, []byte(output), 0644)
			if err != nil {
				log.Fatalln(err)
			}
		} else { //el archivo no existe
			f, err := os.OpenFile(nombre, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			for _, accion := range logSeparado {
				if accion != "" {
					_, err2 := f.WriteString(accion)
					if err2 != nil {
						log.Fatal(err)
					}
					commandSeparado := strings.Split(accion, " ")
					if commandSeparado[0] == "create" {
						crearRegistro(commandSeparado[1], commandSeparado[2])
						if dns == 2 {
							updateReloj(commandSeparado[1], 2)
						} else {
							updateReloj(commandSeparado[1], 3)
						}
					} else if commandSeparado[0] == "update" {
						updateRegistro(commandSeparado[1], commandSeparado[2])
						if dns == 2 {
							updateReloj(commandSeparado[1], 2)
						} else {
							updateReloj(commandSeparado[1], 3)
						}
					} else {
						borrarRegistro(commandSeparado[1])
						if dns == 2 {
							updateReloj(commandSeparado[1], 2)
						} else {
							updateReloj(commandSeparado[1], 3)
						}
					}
				}
			}
		}
	}
}

func connToDNS2(puerto string, nombre1 string, registro1 []byte, log1 []byte, reloj1 []byte) string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se logró conectar: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.VueltaArchivos(context.Background(), &chat.Archivos{Registro: registro1, Otroarchivo: log1, Reloj: reloj1, Nombre: nombre1})
	if err != nil {
		log.Fatalf("Error al tratar de conectar con con los DNS: %s", err)
	}
	return response.Mensaje
}

func consistencia() {
	for {
		time.Sleep(300 * time.Second)
		var log1 []string
		var log2 []string
		log1 = connToDNS("10.10.28.156:9002", "inicio") //Recibe el log del DNS2
		log2 = connToDNS("10.10.28.157:9003", "inicio") //Recibe el log del DNS3
		actualizaLog(log1, 2)                           //Actualiza el log de este DNS con la info del DNS2
		actualizaLog(log2, 3)                           //Actualiza el log de este DNS con la info del DNS2
		//aqui se busca cada archivo
		files, err := ioutil.ReadDir("logs")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files { //por cada archivo en la carpeta 'logs'
			nombre := file.Name()
			log1, err := ioutil.ReadFile("logs/" + nombre) // se deja el log en 'log1'
			if err != nil {
				log.Fatalln(err)
			}
			nombreSep1 := strings.Split(nombre, ".")
			nombreSep2 := strings.Split(nombreSep1[0], "_")
			registro, err2 := ioutil.ReadFile("registros_zf/registro_" + nombreSep2[1] + ".txt") //se deja el regustro en 'registro'
			if err2 != nil {
				log.Fatalln(err)
			}
			reloj, err3 := ioutil.ReadFile("relojes/reloj_" + nombreSep2[1] + ".txt") //se deja el reloj en 'reloj'
			if err3 != nil {
				log.Fatalln(err)
			}
			fmt.Println(connToDNS2("10.10.28.156:9002", nombreSep2[1], registro, log1, reloj)) //envia los archivos como bytes al DNS2
			fmt.Println(connToDNS2("10.10.28.157:9003", nombreSep2[1], registro, log1, reloj)) //envia los archivos como bytes al DNS3
		}
	}
}

func main() {
	go consistencia() //Hace el merge cada 5 minutos
	escuchar()
}
