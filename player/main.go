package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"aluance.io/wordleplayer/internal/config"
	"aluance.io/wordleplayer/internal/playerbot"
)

func main() {
	numGames := 15000
	numBots := 5
	ch := make(chan string)

	bots := make([]playerbot.Playerbot, numBots)
	for i := 0; i < numBots; i++ {
		b, err := playerbot.CreateBot(playerbot.ONEBOT_NAME)
		if err != nil {
			log.Fatal(err)
		}
		bots[i] = b
	}

	count := 0
	for i := 0; i < numGames; i++ {
		go bots[i%numBots].PlayGame(&ch)
		count++
		time.Sleep(config.CONFIG_BOT_THROTTLE * time.Millisecond)
	}
	fmt.Println("Games launched:", count)

	for i := count; i > 0; i-- {
		select {
		case s := <-ch:
			fmt.Println(s)
		case <-time.After(60 * time.Second):
			log.Fatal("Timed out before bot finished playing.")
		}
	}
	close(ch)
}

func initLogger() {

}
