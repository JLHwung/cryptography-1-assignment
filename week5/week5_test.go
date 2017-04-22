package week5

import (
	"fmt"
	"math/big"
)

func ExampleDLog() {
	x := big.NewInt(0)
	y := big.NewInt(0)
	m := big.NewInt(0)
	x.SetString("11717829880366207009516117596335367088558084999998952205599979459063929499736583746670572176471460312928594829675428279466566527115212748467589894601965568", 10)
	y.SetString("3239475104050450443565264378728065788649097520952449527834792452971981976143292558073856937958553180532878928001494706097394108577585732452307673444020333", 10)
	m.SetString("13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006084171", 10)
	result, _ := DLog(x, y, m, 40)
	fmt.Printf("%d\n", result)
	// Output: 375374217830
}
