package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
	"nhooyr.io/websocket"
)

type Receive struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Username string `json:"username"`
	Amount   string `json:"amount"`
	Coin     string `json:"coin"`
	Duration string `json:"duration"`
	Fadein   string `json:"fadein"`
	Fadeout  string `json:"fadeout"`
	IMG      string `json:"img"`
}

type donation struct {
	app.Compo

	image   string
	message string
	donator string
	amount  string
	coin    string
	opacity float64
}

const (
	Address = "ws://mokky.iptime.org:12345/echo"
)

func (d *donation) ShowDonation(img, duration, fadein, fadeout int, msg, donator, amount, coin string) {
	switch img {
	case 0:
		d.image = "/web/image/money.png"
	case 1:
		d.image = "/web/image/money.gif"
	default:
		d.image = "/web/image/money.png"
	}

	d.message = msg
	d.donator = donator
	d.amount = amount
	d.coin = coin
	d.FadeIn(fadein)

	time.Sleep(time.Duration(duration) * time.Second)

	d.FadeOut(fadeout)
}

func (d *donation) FadeIn(fadein int) {
	for i := 0; i <= 100; i++ {
		d.opacity = float64(i) / 100
		d.Update()

		time.Sleep(time.Duration(fadein) * time.Millisecond)
	}

	d.Update()
}

func (d *donation) FadeOut(fadeout int) {
	for i := 100; i >= 0; i-- {
		d.opacity = float64(i) / 100
		d.Update()

		time.Sleep(time.Duration(fadeout) * time.Millisecond)
	}

	d.Update()
}

func (d *donation) OnMount(ctx app.Context) {
	d.opacity = 0
	d.Update()

	urlPaths := strings.Split(app.Window().URL().Path, "/")

	if len(urlPaths) < 3 {
		return
	}

	alertID := urlPaths[2]

	c, _, err := websocket.Dial(ctx, Address+"/"+alertID, nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "StatusInternalError")

	go func() {
		for {
			_, message, err := c.Read(ctx)
			if err != nil {
				log.Println("read:", err)

				return
			}

			var receive Receive
			json.Unmarshal(message, &receive)

			var duration, fadein, fadeout, img int

			duration, err = strconv.Atoi(receive.Duration)
			if err != nil {
				duration = 5
			}

			fadein, err = strconv.Atoi(receive.Fadein)
			if err != nil {
				fadein = 30
			}

			fadeout, err = strconv.Atoi(receive.Fadeout)
			if err != nil {
				fadeout = 15
			}

			img, err = strconv.Atoi(receive.IMG)
			if err != nil {
				img = 0
			}

			d.ShowDonation(img, duration, fadein, fadeout, receive.Message, receive.Username, receive.Amount, receive.Coin)
		}
	}()
}

func (d *donation) Render() app.UI {
	return app.Div().Body(
		app.Img().
			Class("image").
			Src(d.image),
		app.Div().Body(
			app.Span().Body(
				app.Text(d.donator),
			),
			app.Text("님이 "),
			app.Span().Body(
				app.Text(d.amount),
			),
			app.Text(d.coin+"을 주셨어요!"),
		).
			Class("summary"),
		app.Div().Body(
			app.Text(d.message),
		).
			Class("message"),
	).
		Class("content").
		Style("opacity", fmt.Sprintf("%.2f", d.opacity))
}

func main() {
	app.RouteWithRegexp("/alert/[0-9a-z]+", &donation{})
	app.Run()
}
