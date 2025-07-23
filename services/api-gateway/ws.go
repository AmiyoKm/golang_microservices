package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebsocket(w http.ResponseWriter, r *http.Request) {
	conn , err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v",err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("No user ID provided")
		return
	}
	for {
		_ , message , err :=	conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v",err)
			break
		}
		log.Printf("recieved message: %v",message)
	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn , err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v",err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("No userID provided")
		return
	}
	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Printf("No package slug provided")
		return
	}
	type Driver struct {
		Id string `json:"id"`
		Name string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate string `json:"carPlate"`
		PackageSlug string `json:"packageSlug"`
	}

	msg := contracts.WSMessage{
		Type : "driver.cmd.register",
		Data: Driver{
			Id: userID,
			Name : "Amiyo",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate: "ABC123",
			PackageSlug: packageSlug,
		},
	}
	if err := conn.WriteJSON(msg) ; err != nil {
		log.Printf("error sending message: %v",err)
		return
	}

	for {
		_ , message , err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v",err)
			return
		}

		log.Printf("recieved message: %s",message)
	}

}


