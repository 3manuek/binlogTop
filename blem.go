package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	//https://github.com/percona/go-mysql
)

var (
	// ./blem -port=22695 -password="msandbox" -user="msandbox" -interval=5
	serverid      = flag.Uint("serverid", 9999, "Server Id (must be unique)")
	flavor        = flag.String("flavor", "mysql", "Flavor: mysql or mariadb")
	user          = flag.String("user", "root", "MySQL user, must have replication privilege")
	password      = flag.String("password", "", "MySQL password.")
	port          = flag.Uint("port", 3306, "MySQL port.")
	host          = flag.String("host", "127.0.0.1", "MySQL host.")
	interval      = flag.Int("interval", 5, "Interval in seconds.")
	binfile       = flag.String("binfile", "", "Binlog File name")
	binpos        = flag.Int("binpos", 0, "Binglog File Pos")
	keyParsing    = flag.String("keyParsing", "", "Which keys do you want to analyze?")
	tableParsing  = flag.String("tableParsing", "", "Which table do you want to trace?")
	columnParsing = flag.String("columnParsing", "", "Which column do you want to examinate?")
	dbpath        = flag.String("dbpath", "./binlogtop.db", "Generally you won't set this unless you have issues with space")
	removeStatsDB = flag.Bool("removeStatsDB", true, "Do I remove the stats DB or do you want to keep it?")
)

func main() {

	flag.Parse()

	var hostport = fmt.Sprintf("%s:%d", *host, *port)
	var convTime = (time.Duration)(*interval * 1000)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(2)
	}()

	statsdb := InitDB(*dbpath)
	defer statsdb.Close()
	InitTable(statsdb)

	//MapTable := make(map[uint64]replication.TableMapEvent)
	//MapCounters := make(map[TypeKeyEvent]uint64)
	TableMap := make(map[uint64]TypeTableName)
	MapStats := make(map[TypeKeyEvent]TypeDataEvent)

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

	conn, err := client.Connect(hostport, cfg.User, cfg.Password, "test")
	if err != nil {
		fmt.Printf("Cannot connect to host: %v \n", errors.ErrorStack(err))
		os.Exit(3)
	}

	currPos := mysql.Position{*binfile, uint32(*binpos)}
	if *binfile == "" || *binpos == 0 {
		currPos = getCoordinates(conn) //currPos is mysql.Position
	}

	conn.Close()

	Syncer := replication.NewBinlogSyncer(&cfg)
	defer Syncer.Close()

	streamer, _ := Syncer.StartSync(currPos)

	go feedingThread(streamer, TableMap, MapStats) //better to add a time context and send info by channel?

	//basic timer method.
	for {
		for ix, val := range MapStats {
			//fmt.Println("Ix", ix)
			fmt.Printf("Table: %s.%s | Event: %s |  Accum Size: %d  | Counted Events: %d \n", TableMap[ix.TableID].Schema, TableMap[ix.TableID].Table, ix.Event, val.AccumSize, val.Counted)
			//fmt.Println(MapStats)

		}
		//fmt.Printf("--\n")
		time.Sleep(convTime * time.Millisecond) //Calculate time by getting the Event timestamp
	}

	if *removeStatsDB {
		statsdb.Close()
		os.Remove(*dbpath)
	}
}
