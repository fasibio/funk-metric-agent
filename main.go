package main

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fasibio/funk-metric-agent/logger"
	"github.com/fasibio/funk-metric-agent/tracker"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

const (
	// ClikeyInsecureSkipVerify see description in main methode
	ClikeyInsecureSkipVerify string = "insecureSkipVerify"
	// ClikeyFunkserver see description in main methode
	ClikeyFunkserver string = "funkserver"
	// ClikeyConnectionkey see description in main methode
	ClikeyConnectionkey string = "connectionkey"
	// ClikeyLoglevel see description in main methode
	ClikeyLoglevel string = "loglevel"
	// StatsIntervall set the second where statsinfo will be send
	StatsIntervall string = "statsintervall"
	// SearchIndex set the elasicsearch destination
	SearchIndex string = "searchindex"
	// StaticContent set extra content as json
	StaticContent string = "staticcontent"
)

func main() {

	app := cli.NewApp()
	app.Name = "Funk Agent"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   ClikeyInsecureSkipVerify,
			EnvVar: "INSECURE_SKIP_VERIFY",
			Usage:  "Allow insecure serverconnections",
		},
		cli.StringFlag{
			Name:   ClikeyFunkserver,
			EnvVar: "FUNK_SERVER",
			Value:  "ws://localhost:3000",
			Usage:  "the url of the funk_server",
		},
		cli.StringFlag{
			Name:   ClikeyConnectionkey,
			EnvVar: "CONNECTION_KEY",
			Value:  "changeMe04cf242924f6b5f96",
			Usage:  "The connectionkey given to the funk-server to connect",
		},
		cli.StringFlag{
			Name:   ClikeyLoglevel,
			EnvVar: "LOG_LEVEL",
			Value:  "info",
			Usage:  "debug, info, warn, error ",
		},
		cli.StringFlag{
			Name:   StatsIntervall,
			EnvVar: "STATSINTERVALL",
			Usage:  "set the second where statsinfo will be send to server",
			Value:  "15",
		},
		cli.StringFlag{
			Name:   SearchIndex,
			EnvVar: "SEARCHINDEX",
			Usage:  "set the elasticsearch destination index",
			Value:  "default",
		},

		cli.StringFlag{
			Name:   StaticContent,
			EnvVar: "STATICCONTENT",
			Usage:  "set extra content as json",
			Value:  "{}",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.Get().Fatalw("Global error: " + err.Error())
	}
}

type Holder struct {
	streamCon       *websocket.Conn
	Props           Props
	itSelfNamedHost string
	writeToServer   Serverwriter
}

// Props hold all cli given information
type Props struct {
	funkServerURL      string
	InsecureSkipVerify bool
	Connectionkey      string
	SearchIndex        string
	StaticContent      string
}

func run(c *cli.Context) error {
	logger.Initialize(c.String(ClikeyLoglevel))
	searchindex := c.String(SearchIndex)
	hostname, err := os.Hostname()
	if err != nil {
		logger.Get().Warnw("Could not read hostname set to searchindex" + err.Error())
		hostname = searchindex
	}

	holder := Holder{
		Props: Props{
			funkServerURL:      c.String(ClikeyFunkserver),
			InsecureSkipVerify: c.Bool(ClikeyInsecureSkipVerify),
			Connectionkey:      c.String(ClikeyConnectionkey),
			SearchIndex:        searchindex,
			StaticContent:      c.String(StaticContent),
		},
		itSelfNamedHost: hostname,
		writeToServer:   WriteToServer,
	}
	err = holder.openSocketConn(false)
	for err != nil {
		err = holder.openSocketConn(false)
		logger.Get().Errorw("No connection to Server... Wait 5s and try again later")
		time.Sleep(5 * time.Second)
	}
	statsSecond, err := strconv.ParseInt(c.String(StatsIntervall), 10, 64)
	if err != nil {
		return err
	}
	statsTicker := time.NewTicker(time.Duration(statsSecond) * time.Second)
	holder.uploadMetricInformation(statsTicker)

	return nil
}

func (w *Holder) uploadMetricInformation(intervall *time.Ticker) {
	for {
		for range intervall.C {
			w.SaveMetrics()
		}
	}
}

func (w *Holder) SaveMetrics() {

	metrics, err := tracker.GetSystemMetrics()

	if err != nil {
		logger.Get().Error(err)
		return
	}
	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		logger.Get().Error(err)
		return
	}
	var data []string
	data = append(data, string(jsonMetrics))
	msg := []Message{
		Message{
			Type:          MessageTypeStats,
			Data:          data,
			Time:          time.Now(),
			SearchIndex:   w.Props.SearchIndex + "_metrics_cumlated",
			StaticContent: w.Props.StaticContent,
			Attributes: Attributes{
				Host: w.itSelfNamedHost,
			},
		},
	}

	disk, err := tracker.GetDisksMetrics()
	if err == nil {
		var diskdata []string
		for _, one := range disk {
			jsonDiskMetrics, err := json.Marshal(one)
			if err != nil {
				logger.Get().Error("Error by parsing disk info: " + err.Error())
			}
			diskdata = append(diskdata, string(jsonDiskMetrics))

		}
		msg = append(msg,
			Message{
				Type:          MessageTypeDisk,
				Data:          diskdata,
				Time:          time.Now(),
				SearchIndex:   w.Props.SearchIndex + "_metrics_cumlated",
				StaticContent: w.Props.StaticContent,
				Attributes: Attributes{
					Host: w.itSelfNamedHost,
				},
			})
	} else {
		logger.Get().Error(err)
	}

	if len(msg) != 0 {
		err := w.writeToServer(w.streamCon, msg)
		if err != nil {
			logger.Get().Warnw("Error by write Data to Server" + err.Error() + " try to reconnect")

			err := w.openSocketConn(true)
			if err != nil {
				logger.Get().Warnw("Can not connect try again later: " + err.Error())
			} else {
				logger.Get().Infow("Connected to Funk-Server")
			}
		}
	}

}

func openSocketConnection(url string, connectionString string) (*websocket.Conn, error) {
	d := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpHeader := make(http.Header)
	httpHeader.Add("funk.connection", connectionString)

	c, _, err := d.Dial(url, httpHeader)
	if err != nil {
		return nil, err
	}
	return c, nil

}

func (w *Holder) openSocketConn(force bool) error {
	if w.streamCon == nil || force {
		d, err := openSocketConnection(w.Props.funkServerURL+"/data/subscribe", w.Props.Connectionkey)
		if err != nil {
			return err
		}
		w.streamCon = d
		// go h.handleInterrupt(&done)
		// go h.checkConnAndPoll(&conn, &done)
	}
	return nil
}
