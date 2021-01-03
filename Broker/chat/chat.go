package chat

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

var ipActual int = 1

type Server struct {
}

func (s *Server) RecibirDeAdmin(ctx context.Context, in *Message) (*Message, error) { //cuando un admin envia una peticion
	log.Printf("Administrador solicita servidor DNS. Petición: %s", in.Mensaje)
	if ipActual == 1 {
		ipActual++
		return &Message{Mensaje: "10.10.28.155:9001"}, nil
	} else if ipActual == 2 {
		ipActual++
		return &Message{Mensaje: "10.10.28.156:9002"}, nil
	} else {
		ipActual = 1
		return &Message{Mensaje: "10.10.28.157:9003"}, nil //10.10.28.157
	}
	return &Message{Mensaje: "10.10.28.157:9003"}, nil //No debería llegar aqui
}

func sendToDNS(comando string, ip string) string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se logró conectar: %s", err)
	}
	defer conn.Close()

	c := NewChatServiceClient(conn)

	response, err := c.RecibirDeBroker(context.Background(), &Message{Mensaje: comando})
	if err != nil {
		log.Fatalf("Error al tratar de conectar con el Broker: %s", err)
	}
	return response.Mensaje
}

func (s *Server) RecibirDeCliente(ctx context.Context, in *Message) (*Message, error) { //cuando un admin envia una peticion
	log.Printf("Cliente solicita servidor DNS. Petición: %s", in.Mensaje)
	if ipActual == 1 {
		ipActual++
		fmt.Println("BROKER TO DNS1")
		fmt.Println(in.Mensaje)
		IpResponse := sendToDNS(in.Mensaje, "10.10.28.155:9001")
		return &Message{Mensaje: IpResponse}, nil //10.10.28.155
	} else if ipActual == 2 {
		ipActual++
		fmt.Println("BROKER TO DNS2")
		fmt.Println(in.Mensaje)
		IpResponse := sendToDNS(in.Mensaje, "10.10.28.156:9002")
		return &Message{Mensaje: IpResponse}, nil //10.10.28.156
	} else {
		ipActual = 1
		fmt.Println("BROKER TO DNS3")
		fmt.Println(in.Mensaje)
		IpResponse := sendToDNS(in.Mensaje, "10.10.28.157:9003")
		return &Message{Mensaje: IpResponse}, nil //10.10.28.157
	}
	return &Message{Mensaje: "10.10.28.157:9003"}, nil //No debería llegar aqui

}

func (s *Server) RecibirDeBroker(ctx context.Context, in *Message) (*Message, error) {
	return &Message{Mensaje: "aki no llega"}, nil
}
