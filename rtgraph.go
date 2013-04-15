package main

import (
	"log"
	"flag"
 	"time"
	"math/rand"
	"net/http"
	"code.google.com/p/go.net/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	
	http.HandleFunc("/", index_handler)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/websocket", websocket.Handler(websocket_handler))
	
	log.Println("starting web server at", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalln("http.ListenAndServe:", err)
	}
}


func index_handler(w http.ResponseWriter, req *http.Request) {
	log.Println("request for \"index.html\" received")
	http.ServeFile(w, req, "./index.html")
}

func websocket_handler(c *websocket.Conn) {
	
	log.Println("request for \"websocket\" received")
	
	// watch  for close frames in a seperate go routine
	go func() {
		
		defer c.Close();
		
		var message string;
		for {
			err := websocket.Message.Receive(c, &message);
			
			if err != nil {
				log.Println("websocket.Message.Receive:", err);
				// the following close will also stop the sender
				return;
			} else {
				log.Println("received", message);
			}
		}
		
	}()
	
	
	var val int32 = 50;
	var timestamp int32 = 0;
	
	for {
		val += rand.Int31n(11) - 5
		if val < 0 {
			val = 0
		} else if val > 100 {
			val = 100
		}
		
		err := websocket.JSON.Send(c,[2]int32{timestamp, val}) 
		if err != nil {
			log.Println("encoder.Encode", err)
			return
		}
		
		log.Println("sent [", timestamp, ",", val, "]")
		
		timestamp++;
		time.Sleep(50 * time.Millisecond)
	}
	
}
