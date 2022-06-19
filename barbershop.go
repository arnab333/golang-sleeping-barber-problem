package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (bs *BarberShop) addBarber(barber string) {
	bs.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients.", barber)

		for {
			// if there are no clients, the barber goes to sleep
			if len(bs.ClientsChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap.", barber)
				isSleeping = true
			}

			client, ok := <-bs.ClientsChan

			// if "ok" that means we are reciving value from the client channel. so, there are clients still waiting. so the shop is open
			if ok {
				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barber)
					isSleeping = false
				}

				// cut hair
				bs.cutHair(barber, client)
			} else {
				// shop is closed, so send the barber home and close this goroutine
				bs.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (bs *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair.", barber, client)
	time.Sleep(bs.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)
}

func (bs *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	bs.BarbersDoneChan <- true
}

func (bs *BarberShop) closeShopForDay() {
	color.Cyan("Closing shop for the day.")

	close(bs.ClientsChan)
	bs.Open = false

	for a := 1; a <= bs.NumberOfBarbers; a++ {
		<-bs.BarbersDoneChan
	}

	close(bs.BarbersDoneChan)

	color.Green("----------------------------------------------------------------------")
	color.Green("The barber shop is now closed for the day, and everyone has gone home.")
}

func (bs *BarberShop) addClient(client string) {
	color.Green("*** %s arrives!", client)
	if bs.Open {
		select {
		case bs.ClientsChan <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default:
			color.Red("The waiting room is full, so %s leaves.", client)
		}
	} else {
		color.Red("The shop is already closed, so %s leaves!", client)
	}
}
