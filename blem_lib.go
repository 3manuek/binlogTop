package main

import (
	"golang.org/x/net/context"

	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/replication"
)


type TypeKeyEvent struct {
	TableID	uint64
	//Table 	string // Slices can't be on the index key []byte
	//Schema	string //[]byte 
	Event   replication.EventType
}

type TypeTableName struct{
	Table string
	Schema string
}

type TypeDataEvent struct {
	AccumSize uint64
	Counted   uint64
}

func getCoordinates(conn *client.Conn) (p mysql.Position) {
	r, _ := conn.Execute("SHOW MASTER STATUS")

	binFile, _ := r.GetString(0, 0)
	binPos, _ := r.GetInt(0, 1)
	p.Name = binFile
	p.Pos = uint32(binPos)
	return p
}


/*
  Standard feedingThread using Maps. Will be deprecated as I'm planning to use SQLite
  for a more extendeable aggregations. 
*/
func feedingThread (streamer *replication.BinlogStreamer, TableMap map[uint64]TypeTableName, MapStats map[TypeKeyEvent]TypeDataEvent ) {

	for {

		ev, _ := streamer.GetEvent(context.Background())

		switch ev.Header.EventType { 
			case replication.WRITE_ROWS_EVENTv0,
				replication.UPDATE_ROWS_EVENTv0,
				replication.DELETE_ROWS_EVENTv0,
				replication.WRITE_ROWS_EVENTv1,
				replication.DELETE_ROWS_EVENTv1,
				replication.UPDATE_ROWS_EVENTv1,
				replication.WRITE_ROWS_EVENTv2,
				replication.UPDATE_ROWS_EVENTv2,
				replication.DELETE_ROWS_EVENTv2:

				tableMap_   := (*replication.TableMapEvent)(ev.Event.(*replication.RowsEvent).Table)
				key_ := TypeKeyEvent{tableMap_.TableID, ev.Header.EventType}

				TableMap[key_.TableID] = TypeTableName{string(tableMap_.Schema),string(tableMap_.Table)}
    			//MapTable := ev.Event.(*replication.RowEvent).tables

    			bufferStats := MapStats[key_]
    			bufferStats.Counted++
    			bufferStats.AccumSize = bufferStats.AccumSize + (uint64)(ev.Header.EventSize)
    			MapStats[key_] = TypeDataEvent{bufferStats.AccumSize,bufferStats.Counted}

		} 
	}
}


/*
  Standard feedingThread using Maps. Will be deprecated as I'm planning to use SQLite
  for a more extendeable aggregations. 
  DB  is  *driver.Conn type
*/

/*
func feedSQLiteThread (streamer *replication.BinlogStreamer, DB *sql.DB ) {

	for {

		ev, _ := streamer.GetEvent(context.Background())

		switch ev.Header.EventType { 
			case replication.WRITE_ROWS_EVENTv0,
				replication.UPDATE_ROWS_EVENTv0,
				replication.DELETE_ROWS_EVENTv0,
				replication.WRITE_ROWS_EVENTv1,
				replication.DELETE_ROWS_EVENTv1,
				replication.UPDATE_ROWS_EVENTv1,
				replication.WRITE_ROWS_EVENTv2,
				replication.UPDATE_ROWS_EVENTv2,
				replication.DELETE_ROWS_EVENTv2:

				tableMap_   := (*replication.TableMapEvent)(ev.Event.(*replication.RowsEvent).Table)
				key_ := TypeKeyEvent{tableMap_.TableID, ev.Header.EventType}

				TableMap[key_.TableID] = TypeTableName{string(tableMap_.Schema),string(tableMap_.Table)}
    			//MapTable := ev.Event.(*replication.RowEvent).tables

    			bufferStats := MapStats[key_]
    			bufferStats.Counted++
    			bufferStats.AccumSize = bufferStats.AccumSize + (uint64)(ev.Header.EventSize)
    			MapStats[key_] = TypeDataEvent{bufferStats.AccumSize,bufferStats.Counted}

		} 
	}
}
*/


// Stealing Code section :p . Taken from siddontang's go-mysql package.

/*

I want to support almost all the events in a pluggable way.

			case QUERY_EVENT:
				e = &QueryEvent{}
			case XID_EVENT:
				e = &XIDEvent{}
			case TABLE_MAP_EVENT:
				te := &TableMapEvent{}
				if p.format.EventTypeHeaderLengths[TABLE_MAP_EVENT-1] == 6 {
					te.tableIDSize = 4
				} else {
					te.tableIDSize = 6
				}
				e = te
			case WRITE_ROWS_EVENTv0,
				UPDATE_ROWS_EVENTv0,
				DELETE_ROWS_EVENTv0,
				WRITE_ROWS_EVENTv1,
				DELETE_ROWS_EVENTv1,
				UPDATE_ROWS_EVENTv1,
				WRITE_ROWS_EVENTv2,
				UPDATE_ROWS_EVENTv2,
				DELETE_ROWS_EVENTv2:
				e = p.newRowsEvent(h)
			case ROWS_QUERY_EVENT:
				e = &RowsQueryEvent{}
			case GTID_EVENT:
				e = &GTIDEvent{}
			case BEGIN_LOAD_QUERY_EVENT:
				e = &BeginLoadQueryEvent{}
			case EXECUTE_LOAD_QUERY_EVENT:
				e = &ExecuteLoadQueryEvent{}
			case MARIADB_ANNOTATE_ROWS_EVENT:
				e = &MariadbAnnotaeRowsEvent{}
			case MARIADB_BINLOG_CHECKPOINT_EVENT:
				e = &MariadbBinlogCheckPointEvent{}
			case MARIADB_GTID_LIST_EVENT:
				e = &MariadbGTIDListEvent{}
			case MARIADB_GTID_EVENT:
				ee := &MariadbGTIDEvent{}
				ee.GTID.ServerID = h.ServerID
				e = ee
			default:
				e = &GenericEvent{}

*/