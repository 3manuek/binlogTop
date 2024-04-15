# binlogTop (archived)

Binlog Event (sort of) top. Kind deprecated for latest MySQL versions.

## Notice

This tool is old, and it was used to debug the Binary Log events in MySQL stream
when a transaction is split in between files. This tool helped to figure out the 
transactions that were having the issue.

## Howto run

```
./blem -port=22695 -password="msandbox" -user="msandbox" -interval=5
```

Help: `./blem --help`

## Considering advanced monitoring

See [metric tuning for parallel replication](http://jfg-mysql.blogspot.com.ar/2017/02/metric-for-tuning-parallel-replication-mysql-5-7.html)




## Acknowledgments

Based on the great tool developed by Siddontang ("github.com/siddontang/go-mysql/").

## TODO

Add https://github.com/SentimensRG/sigctx for graceful shutdown (not extremely necessary).

