package main

import (
	"errors"
	"fmt"
)

func main() {
	po := new(PurchaseOrder)
	po.Value = 30.85

	SavePO(po, false).Then(func (obj interface{}) error  {
		po := obj.(*PurchaseOrder)
		fmt.Printf("Purchase Number Saved With iD: %d\n", po.Number)

		return nil
	}, func (err error)  {
		fmt.Printf("Purchase Number Failed With Error: " + err.Error() + "\n")
	}).Then(func (obj interface{}) error {
		fmt.Printf("Second Promise Success")
		return nil
	}, func (err error)  {
		fmt.Printf("Second Promise Failed: " + err.Error() + "\n")
	})

	fmt.Scanln()

}

type PurchaseOrder struct{
	Number int
	Value float64
}

func SavePO(po *PurchaseOrder, shouldFail bool) *Promise{
	result := new(Promise)

	result.successChannel = make(chan interface{}, 1)
	result.failureChannel = make(chan error, 1)

	go func() {
		if shouldFail {
			result.failureChannel <- errors.New("Failed To Save Purchase Order")
		}else{
				po.Number = 1234
				result.successChannel <- po
		}
	}()

	return result
}


type Promise struct {
	successChannel chan interface{}
	failureChannel chan error
}

func (this *Promise) Then(success func(interface{}) error, failure func(error)) *Promise {
	result := new(Promise)

	result.successChannel = make(chan interface{}, 1)
	result.failureChannel = make(chan error, 1)

	go func() {
		select{
		case obj :=  <- this.successChannel:
			newErr := success(obj)
			if newErr == nil {
				result.successChannel <- obj
			}else{
				result.failureChannel <- newErr
			}
		case err := <- this.failureChannel:
			failure(err)
			result.failureChannel <- err
		}
	}()

	return result
}
