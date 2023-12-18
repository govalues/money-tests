package money_test

import (
	"encoding/binary"
	"io"
	"os"
	"strconv"
	"testing"

	rh "github.com/Rhymond/go-money"
	bo "github.com/bojanz/currency"
	gd "github.com/govalues/decimal"
	gm "github.com/govalues/money"
)

var (
	resultErr      error
	resultString   string
	resultRhymond  *rh.Money
	resultGovalues gm.Amount
	resultBojanz   bo.Amount
)

func BenchmarkAmount_Add(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, _ := gd.New(200, 2)
			e, _ := gd.New(300, 2)
			a, _ := gm.NewAmountFromDecimal(gm.USD, d)
			c, _ := gm.NewAmountFromDecimal(gm.USD, e)
			resultGovalues, resultErr = a.Add(c)
		}
	})

	b.Run("mod=rhymond", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := rh.New(200, rh.USD)
			c := rh.New(300, rh.USD)
			resultRhymond, resultErr = a.Add(c)
		}
	})

	b.Run("mod=bojanz", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a, _ := bo.NewAmountFromInt64(200, "USD")
			c, _ := bo.NewAmountFromInt64(300, "USD")
			resultBojanz, resultErr = a.Add(c)
		}
	})
}

func BenchmarkAmount_Mul(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, _ := gd.New(200, 2)
			e, _ := gd.New(3, 0)
			a, _ := gm.NewAmountFromDecimal(gm.USD, d)
			resultGovalues, resultErr = a.Mul(e)
		}
	})

	b.Run("mod=rhymond", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := rh.New(200, rh.USD)
			resultRhymond = a.Multiply(3)
		}
	})

	b.Run("mod=bojanz", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a, _ := bo.NewAmountFromInt64(200, "USD")
			resultBojanz, resultErr = a.Mul("3")
		}
	})
}

func BenchmarkAmount_QuoFinite(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, _ := gd.New(200, 2)
			e, _ := gd.New(4, 0)
			a, _ := gm.NewAmountFromDecimal(gm.USD, d)
			resultGovalues, resultErr = a.Quo(e)
		}
	})

	b.Run("mod=rhymond", func(b *testing.B) {
		b.Skip("rhymond does not support division")
	})

	b.Run("mod=bojanz", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a, _ := bo.NewAmountFromInt64(200, "USD")
			resultBojanz, resultErr = a.Div("4")
		}
	})
}

func BenchmarkAmount_QuoInfinite(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, _ := gd.New(200, 2)
			e, _ := gd.New(3, 0)
			a, _ := gm.NewAmountFromDecimal(gm.USD, d)
			resultGovalues, resultErr = a.Quo(e)
		}
	})

	b.Run("mod=rhymond", func(b *testing.B) {
		b.Skip("rhymond does not support division")
	})

	b.Run("mod=bojanz", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a, _ := bo.NewAmountFromInt64(200, "USD")
			resultBojanz, resultErr = a.Div("3")
		}
	})
}

func BenchmarkAmount_Split(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, _ := gd.New(200, 2)
			a, _ := gm.NewAmountFromDecimal(gm.USD, d)
			_, resultErr = a.Split(10)
		}
	})

	b.Run("mod=rhymond", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := rh.New(200, rh.USD)
			_, resultErr = a.Split(10)
		}
	})

	b.Run("mod=bojanz", func(b *testing.B) {
		b.Skip("bojanz does not support split")
	})
}

func BenchmarkAmount_Conv(b *testing.B) {
	b.Run("mod=govalues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, _ := gd.New(200, 2)
			e, _ := gd.New(8000, 4)
			a, _ := gm.NewAmountFromDecimal(gm.USD, d)
			r, _ := gm.NewExchRateFromDecimal(gm.USD, gm.EUR, e)
			resultGovalues, resultErr = r.Conv(a)
		}
	})

	b.Run("mod=rhymond", func(b *testing.B) {
		b.Skip("rhymond does not support currency conversion")
	})

	b.Run("mod=bojanz", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a, _ := bo.NewAmountFromInt64(200, "USD")
			resultBojanz, resultErr = a.Convert("EUR", "0.8000")
		}
	})
}

func BenchmarkParseAmount(b *testing.B) {
	tests := []string{
		"1",
		"123.456",
		"123456789.1234567890",
	}

	for _, s := range tests {
		b.Run(s, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					resultGovalues, resultErr = gm.ParseAmount("USD", s)
				}
			})

			b.Run("mod=rhymond", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					f, _ := strconv.ParseFloat(s, 64)
					resultRhymond = rh.NewFromFloat(f, "USD")
				}
			})

			b.Run("mod=bojanz", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					resultBojanz, resultErr = bo.NewAmount(s, "USD")
				}
			})
		})
	}
}

func BenchmarkAmount_String(b *testing.B) {
	tests := []string{
		"1",
		"123.456",
		"123456789.1234567890",
	}

	for _, s := range tests {
		b.Run(s, func(b *testing.B) {
			b.Run("mod=govalues", func(b *testing.B) {
				a, err := gm.ParseAmount("USD", s)
				if err != nil {
					b.Fatal(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultString = a.String()
				}
			})

			b.Run("mod=rhymond", func(b *testing.B) {
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					b.Fatal(err)
				}
				a := rh.NewFromFloat(f, rh.USD)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultString = a.Display()
				}
			})

			b.Run("mod=bojanz", func(b *testing.B) {
				a, err := bo.NewAmount(s, "USD")
				if err != nil {
					b.Fatal(err)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					resultString = a.String()
				}
			})
		})
	}
}

func readTelcoTests() ([]int64, error) {
	file, err := os.Open("expon180.1e6b")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := make([]int64, 0, 1000000)
	buf := make([]byte, 8)
	for {
		_, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		num := binary.BigEndian.Uint64(buf)
		data = append(data, int64(num))
	}
	return data, nil
}

// BenchmarkAmount_Telco implements computational part of "[Telco benchmark]"
// by Mike Cowlishaw.
// I/O part is not implemented.
//
// [Telco benchmark]: https://speleotrove.com/decimal/telco.html
func BenchmarkAmount_Telco(b *testing.B) {
	tests, err := readTelcoTests()
	if err != nil {
		b.Fatal(err)
	}

	b.Run("mod=govalues", func(b *testing.B) {
		totalFinalPrice := gm.MustParseAmount("USD", "0.00")
		totalBaseTax := gm.MustParseAmount("USD", "0.00")
		totalDistTax := gm.MustParseAmount("USD", "0.00")
		baseRate := gm.MustParseAmount("USD", "0.0013")
		distRate := gm.MustParseAmount("USD", "0.00894")
		baseTaxRate := gd.MustParse("0.0675")
		distTaxRate := gd.MustParse("0.0341")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			tt := tests[i%len(tests)]
			callType := tt & 0x01

			// Duration, Seconds
			duration := gd.MustNew(tt, 0)

			// Price
			var price gm.Amount
			if callType == 0 {
				price, err = baseRate.Mul(duration)
			} else {
				price, err = distRate.Mul(duration)
			}
			if err != nil {
				b.Fatal(err)
			}
			price = price.RoundToCurr()

			// Base Tax
			baseTax, err := price.Mul(baseTaxRate)
			if err != nil {
				b.Fatal(err)
			}
			baseTax = baseTax.TruncToCurr()
			totalBaseTax, err = totalBaseTax.Add(baseTax)
			if err != nil {
				b.Fatal(err)
			}
			finalPrice, err := price.Add(baseTax)
			if err != nil {
				b.Fatal(err)
			}

			// Distance Tax
			if callType != 0 {
				distTax, err := price.Mul(distTaxRate)
				if err != nil {
					b.Fatal(err)
				}
				distTax = distTax.TruncToCurr()
				totalDistTax, err = totalDistTax.Add(distTax)
				if err != nil {
					b.Fatal(err)
				}
				finalPrice, err = finalPrice.Add(distTax)
				if err != nil {
					b.Fatal(err)
				}
			}

			// Final Price
			totalFinalPrice, err = totalFinalPrice.Add(finalPrice)
			if err != nil {
				b.Fatal(err)
			}
			resultString = finalPrice.String()
		}
	})

	b.Run("mod=rhymond", func(b *testing.B) {
		b.Skip("rhymond does not support multiplication by fraction")
	})

	b.Run("mod=bojanz", func(b *testing.B) {
		totalFinalPrice, _ := bo.NewAmount("0.00", "USD")
		totalBaseTax, _ := bo.NewAmount("USD", "0.00")
		totalDistTax, _ := bo.NewAmount("USD", "0.00")
		baseRate, _ := bo.NewAmount("USD", "0.0013")
		distRate, _ := bo.NewAmount("USD", "0.00894")
		baseTaxRate := "0.0675"
		distTaxRate := "0.0341"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			tt := tests[i%len(tests)]
			callType := tt & 0x01

			// Duration, seconds
			duration := strconv.FormatInt(tt, 10)

			// Price
			var price bo.Amount
			if callType == 0 {
				price, err = baseRate.Mul(duration)
			} else {
				price, err = distRate.Mul(duration)
			}
			if err != nil {
				b.Fatal(err)
			}
			price = price.RoundTo(bo.DefaultDigits, bo.RoundHalfEven)

			// Base Tax
			baseTax, err := price.Mul(baseTaxRate)
			if err != nil {
				b.Fatal(err)
			}
			baseTax = baseTax.RoundTo(bo.DefaultDigits, bo.RoundDown)
			totalBaseTax, err = totalBaseTax.Add(baseTax)
			if err != nil {
				b.Fatal(err)
			}
			finalPrice, err := price.Add(baseTax)
			if err != nil {
				b.Fatal(err)
			}

			// Distance Tax
			if callType != 0 {
				distTax, err := price.Mul(distTaxRate)
				if err != nil {
					b.Fatal(err)
				}
				distTax = distTax.RoundTo(bo.DefaultDigits, bo.RoundDown)
				totalDistTax, err = totalDistTax.Add(distTax)
				if err != nil {
					b.Fatal(err)
				}
				finalPrice, err = finalPrice.Add(distTax)
				if err != nil {
					b.Fatal(err)
				}
			}

			// Final Price
			totalFinalPrice, err = totalFinalPrice.Add(finalPrice)
			if err != nil {
				b.Fatal(err)
			}
			resultString = finalPrice.String()
		}
	})
}
