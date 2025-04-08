package main

import (
	"fmt"

	rubberring "github.com/Skrip42/rubberRing"
)

func main() {
	// Let's create a ring buffer of numbers, 4 elements in size, with 2 chunks (2*2=4),
	// and a buffer of free chunks equal to 2.
	// When the buffer is full, it will increase by 4 elements (2 chunks with 2 elements)
	rr := rubberring.NewRubberRing[int](
		rubberring.WithStartChankCount(2),
		rubberring.WithStartChankSize(2),
		rubberring.WithGrowStrategy(func(_ int) (int, int) { return 2, 2 }),
		rubberring.WithPassiveChankBufferSize(2),
	)

	result := []int{}

	// Size:0                 Capacity:4
	// ActiveChanks:2         ActiveCapacity:4
	// PassiveChanks:0        PassiveCapacity:0
	// ActiveChanksSize:[2 2] EndChankNo:0
	// StartPosition:0        EndPosition:0
	fmt.Printf("%+v\n", rr.Stat())

	// This will put 7 messages in the buffer.
	// Since this is more than the initial capacity, 2 more chunks of 2 will be created
	// (according to the growStrategy(2), memory for 4 elements should be allocated
	// and distributed SplitFactor two chunks)
	for i := 0; i < 7; i++ {
		rr.Push(i)
	}

	// Size:7                     Capacity:8
	// ActiveChanks:4             ActiveCapacity:8
	// PassiveChanks:0            PassiveCapacity:0
	// ActiveChanksSize:[2 2 2 2] EndChankNo:3
	// StartPosition:0            EndPosition:7
	fmt.Printf("%+v\n", rr.Stat())

	//When the capacity is filled, memory is immediately allocated for the next element.
	rr.Push(8)

	// Size:8                         Capacity:12
	// ActiveChanks:6                 ActiveCapacity:12
	// PassiveChanks:0                PassiveCapacity:0
	// ActiveChanksSize:[2 2 2 2 2 2] EndChankNo:4
	// StartPosition:0                EndPosition:0
	fmt.Printf("%+v\n", rr.Stat())

	// Now we subtract 6 elements from the buffer.
	// This will clear the first 3 chunks, and they will try to go to passive chunks,
	// but since we have a buffer of passive chunks = 2, the third chunk is utilized.
	// This will free up excess memory
	for i := 0; i < 6; i++ {
		v, err := rr.Pull()
		if err != nil {
			panic("unexpected error")
		}
		result = append(result, v)
	}

	// [0 1 2 3 4 5]
	fmt.Println(result)

	// Size:2                   Capacity:10
	// ActiveChanks:3           ActiveCapacity:6
	// PassiveChanks:2          PassiveCapacity:4
	// ActiveChanksSize:[2 2 2] EndChankNo:1
	// StartPosition:0          EndPosition:2
	fmt.Printf("%+v\n", rr.Stat())

	// If we put more elements into the buffer,
	// the passive capacity will become active. No new allocation will occur.
	for i := 0; i < 7; i++ {
		rr.Push(i + 8)
	}

	// Size:9                       Capacity:10
	// ActiveChanks:5               ActiveCapacity:10
	// PassiveChanks:0              PassiveCapacity:0
	// ActiveChanksSize:[2 2 2 2 2] EndChankNo:4
	// StartPosition:0              EndPosition:9
	fmt.Printf("%+v\n", rr.Stat())

	for i := 0; i < 9; i++ {
		v, err := rr.Pull()
		if err != nil {
			panic("unexpected error")
		}
		result = append(result, v)
	}

	// [0 1 2 3 4 5 6 8 8 9 10 11 12 13 14]
	fmt.Println(result)

	// When trying to get an element from an empty buffer, we get an error io.EOF
	_, err := rr.Pull()

	// EOF
	fmt.Println(err)

	// Size:0                       Capacity:6
	// ActiveChanks:1               ActiveCapacity:2
	// PassiveChanks:2              PassiveCapacity:4
	// ActiveChanksSize:[2 2 2 2 2] EndChankNo:0
	// StartPosition:1              EndPosition:1
	fmt.Printf("%+v\n", rr.Stat())
}
