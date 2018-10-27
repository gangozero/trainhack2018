package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx"
)

// dbConfig describes config object
type dbConfig struct {
	Host       string `required:"true"`
	Port       int16  `required:"true"`
	User       string `required:"true"`
	Password   string `required:"true"`
	Name       string `required:"true"`
	PoolSize   int    `required:"true"`
	TLSEnabled bool   `default:"true"`
}

// NewDbConn creates pgx pool connection
func newDbConn(conf dbConfig) (*pgx.ConnPool, error) {
	var err error

	config := pgx.ConnConfig{
		Host:     conf.Host,
		Port:     uint16(conf.Port),
		User:     conf.User,
		Password: conf.Password,
		Database: conf.Name,
	}

	if conf.TLSEnabled {
		config.TLSConfig = &tls.Config{
			ServerName: conf.Host,
			// Avoids most of the memorably-named TLS attacks
			MinVersion: tls.VersionTLS12,
			// Causes servers to use Go's default ciphersuite preferences,
			// which are tuned to avoid attacks. Does nothing on clients.
			PreferServerCipherSuites: true,
			// Only use curves which have constant-time implementations
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
			},
		}

		// Load CA cert
		caCert := []byte(`-----BEGIN CERTIFICATE-----
MIIDpTCCAo2gAwIBAgIEW9H++zANBgkqhkiG9w0BAQ0FADBUMVIwUAYDVQQDDEli
MTQ0OTU1OWYwOGQ0ZjM4OWM0ZmJmMjVlZTM2MzU3ZitldWNsb3VkLTdmODkwZDhj
NjQ3MjNiMDcxOTdkYWI1YWE2YmQzODI5MB4XDTE4MTAyNTE3MzU1NVoXDTM4MTAy
NTE3MDAwMFowVDFSMFAGA1UEAwxJYjE0NDk1NTlmMDhkNGYzODljNGZiZjI1ZWUz
NjM1N2YrZXVjbG91ZC03Zjg5MGQ4YzY0NzIzYjA3MTk3ZGFiNWFhNmJkMzgyOTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMC/0i2drpvPWkcGbKqNjiBD
uf7wTCcGhhsw3tJr/Xg5fpE5KZKCVTsh6y+4yHs3WYP+H6f0iTjSXmui8MY5sK6P
Gl4ylIh9gueS9cLoyeuSKKJnK6v0cfjqSbqulOeYtt/qMynNn0WFzjA8DXMDV+ZJ
/SGvSLrxyNRXKTrhii26Tu2DQWfSLBUdJPf/9VJEKXewIyNRN0mLVRU37TsJhoqf
IwmF8oRbCCF06pRMszuHd1C9FMQKYcACZHohy3xdCWYOVVIhyzBrMB+GJYf3vgH1
ukCkfwU8E8Sa+zHvAEx7xy20cCq1sdoJYuxGumVtvIcd+YGYTX3QDQoDei66y7sC
AwEAAaN/MH0wHQYDVR0OBBYEFFJZ9Naq4xlQWuFr+eu9ZLRWi3WPMA4GA1UdDwEB
/wQEAwICBDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TBAUw
AwEB/zAfBgNVHSMEGDAWgBRSWfTWquMZUFrha/nrvWS0Vot1jzANBgkqhkiG9w0B
AQ0FAAOCAQEADhkeYKM8feEJQHOX9s0qsxmPh1QO+ovkyItGQvocBy7PpeDIzVIn
0rjptz0xz87m8rT2sDS8EA0PiHcMW2yiR/ZU9u1WDJd4h95Cot/Xe5BklM2TLPFo
ahNIA27wU6OfgmtY2nHCcQ+f3h5eyFNda/X4LWPeYvD5/S5k6YeeQOa35irFdVIh
BL4qo9Mch0ZG+meBqD5ae9ngtnk9vi36sMjvnStVPPZRs7QRjJ6so4E2O2ZRuFGH
+GStj8mOUmmmRZUxzWnqizjKj4fxHZlPo61TbeT2/8pRhY2hcP07KdcVqKdBDuxw
wNv3qGkGQL1/iG5N1WSf6Ufijwz0gPwV1w==
-----END CERTIFICATE-----`)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		config.TLSConfig.RootCAs = caCertPool
	}

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: conf.PoolSize,
	}

	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// dbCheck will validate DB connection and create base data
func dbCheck(pool *pgx.ConnPool) error {
	conn, err := pool.Acquire()
	if err != nil {
		return fmt.Errorf("Cant acquire DB connection: %s", err.Error())
	}
	defer conn.Close()

	err = conn.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("Cant ping DB: %s", err.Error())
	}

	return nil
}

// deferRollbackAndLog will rollback transaction and log error if something goes wrong
func deferRollbackAndLog(tx *pgx.Tx) {
	err := tx.Rollback()
	if err != nil && err != pgx.ErrTxClosed {
		log.Printf("Error rolling back transaction: %s", err.Error())
	}
}

func getStationList(pool *pgx.ConnPool, trainID string) (*GetStationsListResponse, error) {
	conn, err := pool.Acquire()
	if err != nil {
		log.Printf("Can't acquire DB connection from pool: %s", err.Error())
		return nil, fmt.Errorf("Can't acquire DB connection from pool")
	}
	defer pool.Release(conn)

	tx, err := conn.Begin()
	if err != nil {
		log.Printf("Can't start transaction: %s", err.Error())
		return nil, fmt.Errorf("Can't start transaction")
	}
	defer deferRollbackAndLog(tx)

	queryRoute := `SELECT trips.route_id, trips.service_id, trips.trip_id, trips.trip_headsign, trips.trip_short_name 
    FROM trips trips, calendar_dates cal
    WHERE trips.trip_short_name = $1
        AND trips.service_id = cal.service_id
        AND cal.date :: TEXT=TO_CHAR(NOW() :: DATE, 'yyyymmdd')
    LIMIT 1;

	`

	var routeID, serviceID, tripID, tripHead, tripName string
	err = tx.QueryRow(queryRoute, trainID).Scan(&routeID, &serviceID, &tripID, &tripHead, &tripName)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Printf("Cannot scan result rows: %s", err.Error())
		return nil, fmt.Errorf("Cannot scan result rows")
	}

	log.Printf("Searched train '%s' heading to '%s': found route_id=%s, service_id=%s, trip_id=%s", tripName, tripHead, routeID, serviceID, tripID)

	queryTime := `
	SELECT stops.stop_id, stops.stop_name, stops.stop_lat, stops.stop_lon, 
		EXTRACT(EPOCH FROM st.arrival_time - to_char(now() at time zone 'UTC-02', 'HH24:MI:SS') :: TIME) as arrival_sec, 
		st.stop_sequence 
	FROM stops stops, stop_times st
	WHERE st.trip_id=$1
		AND st.stop_id = stops.stop_id
	;`

	rows, err := tx.Query(queryTime, tripID)

	if err != nil {
		log.Printf("Error getting list of stations: %s", err.Error())
		return nil, fmt.Errorf("Error getting list of stations")
	}
	defer rows.Close()

	sts := []Station{}

	for rows.Next() {
		var stopID, stopName string
		var lat, lon float64
		var arrival, seq int

		err = rows.Scan(&stopID, &stopName, &lat, &lon, &arrival, &seq)
		if err != nil {
			log.Printf("Error scanning station: %s", err.Error())
			return nil, fmt.Errorf("Error scanning station")
		}

		if arrival > 0 {
			st := Station{
				ID:     stopID,
				Title:  cleanName(stopName),
				IsFood: 0,
				Geo:    Geo{Lat: lat, Lon: lon},
				Time:   arrival,
			}
			if seq%4 == 0 || seq%5 == 0 {
				st.IsCoffee = 0
			} else {
				st.IsCoffee = 1
			}
			sts = append(sts, st)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error commiting transaction: %s", err.Error())
		return nil, fmt.Errorf("Error commiting transaction")
	}

	return &GetStationsListResponse{
		Stations: sts,
	}, nil
}

func cleanName(name string) string {
	tmp := strings.Replace(name, "Centralstation", "C", -1)
	return strings.Replace(tmp, " station", "", -1)
}

func createOrder(pool *pgx.ConnPool, req *PostOrderRequest) (*PostOrderResponse, error) {
	conn, err := pool.Acquire()
	if err != nil {
		log.Printf("Can't acquire DB connection from pool: %s", err.Error())
		return nil, fmt.Errorf("Can't acquire DB connection from pool")
	}
	defer pool.Release(conn)

	tx, err := conn.Begin()
	if err != nil {
		log.Printf("Can't start transaction: %s", err.Error())
		return nil, fmt.Errorf("Can't start transaction")
	}
	defer deferRollbackAndLog(tx)

	query := `
	INSERT INTO orders (train, carriage, station, repeat_order, delivery, ord) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`

	var id string
	err = tx.QueryRow(query, req.Train, req.Carriage, req.Station, req.RepeatOrder, req.Delivery, req.Order).Scan(&id)
	if err != nil {
		log.Printf("Cannot insert new order to DB: %s", err.Error())
		return nil, fmt.Errorf("Cannot insert new order to DB")
	}

	ts, err := getTime(req.Station, req.Train)
	if err != nil {
		log.Printf("Cannot get arrival time from ResRobot: %s", err.Error())
	} else {
		_, err = tx.Exec("UPDATE orders SET ts_ready=$1 WHERE id=$2", ts, id)
		if err != nil {
			log.Printf("Cannot update arrival time: %s", err.Error())
			return nil, fmt.Errorf("Cannot update arrival time")
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error commiting transaction: %s", err.Error())
		return nil, fmt.Errorf("Error commiting transaction")
	}

	return &PostOrderResponse{ID: id}, nil
}

func getTaskList(pool *pgx.ConnPool) (*GetTaskListResponse, error) {
	conn, err := pool.Acquire()
	if err != nil {
		log.Printf("Can't acquire DB connection from pool: %s", err.Error())
		return nil, fmt.Errorf("Can't acquire DB connection from pool")
	}
	defer pool.Release(conn)

	tx, err := conn.Begin()
	if err != nil {
		log.Printf("Can't start transaction: %s", err.Error())
		return nil, fmt.Errorf("Can't start transaction")
	}
	defer deferRollbackAndLog(tx)

	query := `
	SELECT train, carriage, station, repeat_order, delivery, ord, ts_ready
	FROM orders
	WHERE ts_ready >= NOW() - interval '60 minutes'
	ORDER BY ts_ready
	;`

	rows, err := tx.Query(query)

	if err != nil {
		log.Printf("Error getting list of tasks: %s", err.Error())
		return nil, fmt.Errorf("Error getting list of tasks")
	}
	defer rows.Close()

	tasks := []TaskItem{}

	for rows.Next() {
		var train, carriage, station string
		var repeatOrder, delivery bool
		var ord []OrderItem
		var ts time.Time

		err = rows.Scan(&train, &carriage, &station, &repeatOrder, &delivery, &ord, &ts)
		if err != nil {
			log.Printf("Error scanning tasks: %s", err.Error())
			return nil, fmt.Errorf("Error scanning task")
		}

		task := TaskItem{
			Train:       train,
			Carriage:    carriage,
			Station:     station,
			RepeatOrder: repeatOrder,
			Delivery:    delivery,
			Order:       ord,
			ArrivalTime: ts.Unix(),
		}

		tasks = append(tasks, task)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error commiting transaction: %s", err.Error())
		return nil, fmt.Errorf("Error commiting transaction")
	}

	return &GetTaskListResponse{
		Tasks: tasks,
	}, nil
}
