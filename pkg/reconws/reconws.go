/*
   reconws is websocket client that automatically reconnects


   modified from github.com/timdrysdale/relay/pkg/reconws to remove
   additional elements relating to usage with hub/agg etc in relay

   Copyright (C) 2019 Timothy Drysdale <timothy.d.drysdale@gmail.com>
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as
   published by the Free Software Foundation, either version 3 of the
   License, or (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.
   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package reconws

import (
	"context"
	"errors"
	"math/rand"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

type WsMessage struct {
	Data []byte
	Type int
}

// connects (retrying/reconnecting if necessary) to websocket server at url

type ReconWs struct {
	ForwardIncoming bool
	In              chan WsMessage
	Out             chan WsMessage
	Retry           RetryConfig
	Url             string
	ID              string
}

type RetryConfig struct {
	Factor  float64
	Jitter  bool
	Min     time.Duration
	Max     time.Duration
	Timeout time.Duration
}

func New() *ReconWs {
	r := &ReconWs{
		In:              make(chan WsMessage),
		Out:             make(chan WsMessage),
		ForwardIncoming: true,
		Retry: RetryConfig{Factor: 2,
			Min:     1 * time.Second,
			Max:     10 * time.Second,
			Timeout: 1 * time.Second,
			Jitter:  false},
		ID: uuid.New().String()[0:6],
	}
	return r
}

// run this in a separate goroutine so that the connection can be
// ended from where it was initialised, by close((* ReconWs).Stop)
func (r *ReconWs) Reconnect(ctx context.Context, url string) {

	boff := &backoff.Backoff{
		Min:    r.Retry.Min,
		Max:    r.Retry.Max,
		Factor: r.Retry.Factor,
		Jitter: r.Retry.Jitter,
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// try dialling ....

	for {

		select {
		case <-ctx.Done():
			return
		default:

			dialCtx, cancel := context.WithCancel(ctx)

			err := r.Dial(dialCtx, url)
			cancel()

			log.WithField("error", err).Debug("Dial finished")
			if err == nil {
				boff.Reset()
			} else {
				time.Sleep(boff.Duration())
			}
			//TODO immediate return if cancelled....
		}
	}
}

// Dial the websocket server once.
// If dial fails then return immediately
// If dial succeeds then handle message traffic until
// the context is cancelled
func (r *ReconWs) Dial(ctx context.Context, urlStr string) error {

	id := "reconws.Dial(" + r.ID + ")"

	var err error

	if urlStr == "" {
		log.Errorf("%s: Can't dial an empty Url", id)
		return errors.New("Can't dial an empty Url")
	}

	// parse to check, dial with original string
	u, err := url.Parse(urlStr)

	if err != nil {
		log.Errorf("%s: error with url because %s:", id, err.Error())
		return err
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		log.Errorf("%s: Url needs to start with ws or wss", id)
		return errors.New("Url needs to start with ws or wss")
	}

	if u.User != nil {
		log.Errorf("%s: Url can't contain user name and password", id)
		return errors.New("Url can't contain user name and password")
	}

	// start dialing ....

	log.WithField("To", u).Tracef("%s: connecting to %s", id, u)

	//assume our context has been given a deadline if needed
	c, _, err := websocket.DefaultDialer.DialContext(ctx, urlStr, nil)
	//	defer c.Close()

	if err != nil {
		log.WithField("error", err).Errorf("%s: dialing error because %s", id, err.Error())
		return err
	}

	log.WithField("To", u).Tracef("%s: connected to %s", id, u)
	// handle our reading tasks

	readClosed := make(chan struct{})

	go func() {
	LOOP:
		for {
			select {
			case <-ctx.Done():
			default:
			}
			//assume this will produce non-nil err on context.Done
			mt, data, err := c.ReadMessage()

			// Check for errors, e.g. caused by writing task closing conn
			// because we've been instructed to exit
			// log as info since we expect an error here on a normal exit
			if err != nil {
				log.WithField("error", err).Infof("%s: error reading from conn; closing", id)
				close(readClosed)
				break LOOP
			}

			// optionally forward messages
			if r.ForwardIncoming {
				r.In <- WsMessage{Data: data, Type: mt}
				log.Debugf("%s: received %d-byte message", id, len(data))
			} else {
				log.Debugf("%s: ignored %d-byte message", id, len(data))
			}
		}
	}()

	// handle our writing tasks
LOOPWRITING:
	for {
		select {
		case <-readClosed:
			err = nil // nil error resets the backoff
			break LOOPWRITING
		case msg := <-r.Out:

			err := c.WriteMessage(msg.Type, msg.Data)
			if err != nil {
				log.WithField("error", err).Infof("%s: error writing to conn; closing", id)
				break LOOPWRITING
			}
			log.Debugf("%s: sent %d-byte message", id, len(msg.Data))

		case <-ctx.Done(): // context has finished, either timeout or cancel
			//TODO - do we need to do this?
			// Cleanly close the connection by sending a close message
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.WithField("error", err).Infof("%s: error sending close message; closing", id)
			} else {
				log.Infof("%s: connection closed", id)
			}
			c.Close()
			break LOOPWRITING
		}
	}
	log.Debugf("%s: done", id)
	return err

}
