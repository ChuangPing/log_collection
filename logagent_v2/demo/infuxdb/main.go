package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
	"time"
)

// connect
func connInflux() (client influxdb2.Client, err error) {
	client = influxdb2.NewClient("http://127.0.0.1:8086", "EmLDPZP2sRY3vZLaLAgXDcL1nc7X7v4XVB3TYzGu7RAQ9AKB23eQEwohiT5OFHqbvtXLti0cZDlsV7bRHzOSuw==")
	res, err := client.Ping(context.Background())
	if err != nil {
		fmt.Println("client ping server failed,err:", err)
		return
	}
	if !res {
		fmt.Println("connect influxdb failed")
		return
	}
	logrus.Info("connect influxdb success")
	return
}

// insert
func insert(client influxdb2.Client) {
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking("test", "test")
	// Create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45.0},
		time.Now())
	// write point immediately
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Println("writAPI:write point failed,err:", err)
		return
	}
	// Create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45.0).
		SetTime(time.Now())
	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Println("write point failed,err:", err)
		return
	}
	// Or write directly line protocol
	line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	err = writeAPI.WriteRecord(context.Background(), line)
	if err != nil {
		fmt.Println("WriteRecord failed,err:", err)
		return
	}
}

// get query
func query(client influxdb2.Client) {
	queryAPI := client.QueryAPI("test")
	// get parser flux query result
	result, err := queryAPI.Query(context.Background(), `from(bucket:"test")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}
	// Ensures background processes finishes
	client.Close()
}

func main() {
	client, _ := connInflux()
	//insert(client)
	logrus.Infof("insert success\n")
	query(client)
}
