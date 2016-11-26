package bank

type withdrawDescr struct {
	sum    int
	result chan bool
}

var deposits = make(chan int)
var balances = make(chan int)
var withdraws = make(chan withdrawDescr)

func init() {
	go monitor()
}

// Deposit serves to put a sum to the bank
func Deposit(sum int) {
	deposits <- sum
}

// Balance serves to get the current balance
func Balance() int {
	return <-balances
}

// Withdraw serves to get a sum from the bank
func Withdraw(sum int) bool {
	result := make(chan bool)
	withdraw := withdrawDescr{sum, result}
	withdraws <- withdraw

	return <-result
}

func monitor() {
	var balance int

	for {
		select {
		case sum := <-deposits:
			balance += sum

		case balances <- balance:

		case withdraw := <-withdraws:
			if withdraw.sum > balance {
				withdraw.result <- false
			} else {
				balance -= withdraw.sum
				withdraw.result <- true
			}
		}
	}
}
