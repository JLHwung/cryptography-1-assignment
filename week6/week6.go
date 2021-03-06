package week6

import (
	"crypto/rsa"
	"errors"
	"math/big"
)

func isSquareNumber(n *big.Int) (bool, *big.Int) {
	sqrtFloor := new(big.Int).Sqrt(n)
	squareSqrtFloor := new(big.Int).Mul(sqrtFloor, sqrtFloor)
	return squareSqrtFloor.Cmp(n) == 0, sqrtFloor
}

// FactorCloselyFactorSemiPrime finds p, q such that
// N = p*q when |p - q| < 2N^{1/4} and p <= q
func FactorCloselyFactorSemiPrime(N *big.Int) (*big.Int, *big.Int, error) {
	return FactorNearlyFactorSemiPrime(N, 0)
}

// FactorNearlyFactorSemiPrime finds p, q such that
// N = p*q when |p - q| < 2^{magnitude + 1} N^{1/4} and p <= q
// Notice that when magnitude = 1, it reduces to `FactorCloselyFactorSemiPrime`
func FactorNearlyFactorSemiPrime(N *big.Int, magnitude uint) (*big.Int, *big.Int, error) {
	return FactorProportionalFactorSemiPrime(N, magnitude, big.NewRat(1, 1))
}

// FactorProportionalFactorSemiPrime finds p, q such that
// N = p*q when |ap - bq| < 2^{magnitude + 1} N^{1/4}, λ = a/b is a rational number
// Notice that when a = b = 1, it reduces to `FactorNearlyFactorSemiPrime`
func FactorProportionalFactorSemiPrime(N *big.Int, magnitude uint, proportion *big.Rat) (*big.Int, *big.Int, error) {
	// extract λ = a/b s.t. a, b is even
	num := proportion.Num()
	num.Add(num, num)
	denom := proportion.Denom()
	denom.Add(denom, denom)

	// calculate a*b*N
	numDenomN := new(big.Int).Set(N)
	numDenomN.Mul(numDenomN, num)
	numDenomN.Mul(numDenomN, denom)

	numDenomNSqrt := new(big.Int).Sqrt(numDenomN)

	AvgGuess := new(big.Int).Set(numDenomNSqrt)
	AvgGuessSqaure := new(big.Int).Mul(AvgGuess, AvgGuess)
	Difference := new(big.Int).Sub(AvgGuessSqaure, numDenomN)
	one := big.NewInt(1)

	for i := 0; i < (1 << (magnitude + magnitude)); i++ {
		// increment Difference by 2*A_guess + 1
		Difference.Add(Difference, AvgGuess)
		Difference.Add(Difference, AvgGuess)
		Difference.Add(Difference, one)

		// increment A_guess by 1
		AvgGuess.Add(AvgGuess, one)

		// check if x^2 is square number
		if isSquare, SqrtDifference := isSquareNumber(Difference); isSquare {
			pMultiple := new(big.Int).Sub(AvgGuess, SqrtDifference)
			qMultiple := new(big.Int).Add(AvgGuess, SqrtDifference)

			quo := new(big.Int)
			if quo.Mod(pMultiple, num); quo.BitLen() == 0 {
				pMultiple.Div(pMultiple, num)
				qMultiple.Div(qMultiple, denom)
				return pMultiple, qMultiple, nil
			} else if quo.Mod(pMultiple, denom); quo.BitLen() == 0 {
				pMultiple.Div(pMultiple, denom)
				qMultiple.Div(qMultiple, num)
				return pMultiple, qMultiple, nil
			} else {
				continue
			}
		}
	}

	return nil, nil, errors.New("The factor is not closely enough for efficient factoring")
}

// DecryptRSAPKCSv15WithCloselyFactor will try to decrypt cipherText given RSA public key assuming the SemiPrimi can be
// factored into two close prime
func DecryptRSAPKCSv15WithCloselyFactor(pubKey *rsa.PublicKey, cipherText []byte) ([]byte, error) {
	p, q, err := FactorCloselyFactorSemiPrime(pubKey.N)
	one := big.NewInt(1)
	if err != nil {
		panic(err)
	}

	Primes := make([]*big.Int, 2)
	Primes[0] = p
	Primes[1] = q

	// Compute Euler Function φ(N) = (p - 1)(q - 1)
	pSub1 := new(big.Int).Sub(p, one)
	qSub1 := new(big.Int).Sub(q, one)
	phi := new(big.Int).Mul(pSub1, qSub1)
	// Compute Private Component from E, φ(N)
	D := new(big.Int)
	new(big.Int).GCD(D, new(big.Int), big.NewInt(int64(pubKey.E)), phi)

	privateKey := new(rsa.PrivateKey)
	privateKey.PublicKey = *pubKey
	privateKey.D = D
	privateKey.Primes = Primes

	return rsa.DecryptPKCS1v15(nil, privateKey, cipherText)
}
