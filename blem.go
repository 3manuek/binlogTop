package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"flag"

	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/replication"
	"github.com/juju/errors"
)

var (
	// Default configurations
	/*
	serverid	  = uint32(9999)
	flavor		  = "mysql"	
	host          = "127.0.0.1"
	port          = uint16(22695)
	user          = "msandbox"
	password      = "msandbox"
	database      = "test" */

	serverid	  = flag.Uint("serverid", 9999, "Server Id (must be unique)")
	flavor 		  = flag.String("flavor", "mysql", "Flavor: mysql or mariadb")
	user 		  = flag.String("user", "root", "MySQL user, must have replication privilege")
	password 	  = flag.String("password", "", "MySQL password")
	port 		  = flag.Uint("port", 3306, "MySQL port")
	host 		  = flag.String("host", "127.0.0.1", "MySQL host")


)


func main() {
	flag.Parse()

	// What to do on kill (verify conn is dead.)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(2)
	}()

	//MapTable := make(map[uint64]replication.TableMapEvent)
	//MapCounters := make(map[TypeKeyEvent]uint64)
	TableMap    := make(map[uint64]TypeTableName)
	MapStats    := make(map[TypeKeyEvent]TypeDataEvent)

	//
	// Connection 
	//
	cfg := replication.BinlogSyncerConfig{
		ServerID:        uint32(*serverid),
		Flavor:          *flavor,
		Host:            *host,
		Port:            uint16(*port),
		User:            *user,
		Password:        *password,
		SemiSyncEnabled: false,
		//RawModeEnabled: false,
	}
	var hostport = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	conn, err := client.Connect(hostport, cfg.User, cfg.Password, "test")
	if err != nil {
		fmt.Printf("Cannot connect to host: %v \n", errors.ErrorStack(err))
		os.Exit(3)
	}
	currPos := getCoordinates(conn) //currPos is mysql.Position

	conn.Close()

	Syncer := replication.NewBinlogSyncer(&cfg)
	defer Syncer.Close()

	streamer, _ := Syncer.StartSync(currPos)

	go feedingThread(streamer, TableMap, MapStats)

	//basic timer method.
	for {
		for ix, val := range MapStats {
			//fmt.Println("Ix", ix)
			fmt.Printf("Table: %s.%s | Event: %s |  Accum: %d  | Counted: %d \n", TableMap[ix.TableID].Schema,TableMap[ix.TableID].Table , ix.Event, val.AccumSize, val.Counted)
			//fmt.Println(MapStats)

		}
		fmt.Printf("--\n")
		time.Sleep(5000 * time.Millisecond)
	}
}
