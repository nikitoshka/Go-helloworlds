package main

import (
	"log"
	"the_go_progr_lang/9/1/bank"
)

func main() {
	done := make(chan struct{})

	go func() {
		log.Printf("[first] balance: %d\n", bank.Balance())
		log.Printf("[first] depositing 300\n")
		bank.Deposit(300)
		log.Printf("[first] balance: %d\n", bank.Balance())

		done <- struct{}{}
	}()

	go func() {
		log.Printf("[second] balance: %d\n", bank.Balance())
		log.Printf("[second] depositing 500\n")
		bank.Deposit(500)
		log.Printf("[second] balance: %d\n", bank.Balance())

		done <- struct{}{}
	}()

	go func() {
		log.Printf("[third] balance: %d\n", bank.Balance())
		log.Printf("[third] withdrawing 900\n")
		if ok := bank.Withdraw(900); ok {
			log.Printf("[third] withdrawing succeeded\n")
		} else {
			log.Printf("[third] failed to withdraw\n")
		}

		done <- struct{}{}
	}()

	<-done
	<-done
	<-done

	log.Printf("[main] balance: %d\n", bank.Balance())
}
