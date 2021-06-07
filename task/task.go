package task

import (
	"fmt"
	"gbot/commerce/banking"
	"time"
)

func PayInterest(interval time.Duration, bank *banking.Bank) {
	fmt.Println("Started 'PayInterest' task")
	for range time.Tick(interval) {
		err := bank.PayInterest()
		fmt.Println("Paid Interest")
		if err != nil {
			panic(err)
		}
	}
}
