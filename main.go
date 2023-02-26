package main

import (
	"flag"

	"github.com/ronanzindev/go-cacher-example/cache"
)

func main() {
	// conn, err := net.Dial("tcp", ":3000")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// conn.Write([]byte("SET Foo bar 2500000000000"))

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// select {}
	// return
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the serve")

		leaderAddr = flag.String("leaderaddr", "", "listen addres of the leader")
	)
	flag.Parse()

	opts := ServerOpts{
		ListenAddr:  *listenAddr,
		IsLeader:    len(*leaderAddr) == 0,
		LeaderAdder: *leaderAddr,
	}

	// go func() {
	// 	time.Sleep(time.Second * 2)

	// 	conn, err := net.Dial("tcp", ":3000")

	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	conn.Write([]byte("SET Foo bar 2500000000000"))

	// 	time.Sleep(time.Second * 2)

	// 	conn.Write([]byte("GET Foo"))

	// 	buf := make([]byte, 1000)
	// 	n, _ := conn.Read(buf)
	// 	fmt.Println(string(buf[:n]))
	// }()

	server := NewServer(opts, cache.New())
	server.Start()
}
