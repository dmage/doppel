package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dmage/doppel/pkg/database"
	"github.com/dmage/doppel/pkg/hash2"
)

func popcount(x uint32) int {
	i := 0
	for x > 0 {
		if x&1 == 1 {
			i++
		}
		x >>= 1
	}
	return i
}

func hdist(a, b uint32) int {
	return popcount(a ^ b)
}

func jaccard(a, b map[string]struct{}) float64 {
	union := 0
	intersect := 0
	for i := range a {
		union++
		if _, ok := b[i]; ok {
			intersect++
		}
	}
	for i := range b {
		if _, ok := a[i]; !ok {
			union++
		}
	}
	return float64(intersect) / float64(union)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: <progname> <id>")
	}

	id := os.Args[1]

	db, err := database.NewDefault()
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	sig, err := db.GetSignature(id)
	if err != nil {
		log.Fatal(err)
	}
	m1, err := hash2.Trigrams(strings.NewReader(sig.Signature))
	if err != nil {
		log.Fatal(err)
	}

	sigsiter, err := db.ListSignatures()
	if err != nil {
		log.Fatal(err)
	}
	defer sigsiter.Close()
	//total, fp1, fp2, fp, tp, fn, tn := 0, 0, 0, 0, 0, 0, 0
	for {
		s, ok := sigsiter.Next()
		if !ok {
			break
		}

		m2, err := hash2.Trigrams(strings.NewReader(s.Signature))
		if err != nil {
			log.Fatal(err)
		}
		hd := hdist(sig.Hash2, s.Hash2)

		var h3d float64
		if sig.Hash3 == 0 {
			if s.Hash3 != 0 {
				h3d = 1
			} else {
				h3d = 0
			}
		} else {
			h3d = float64(s.Hash3) / float64(sig.Hash3)
			if h3d > 1 {
				h3d = 1 / h3d
			}
			h3d = 1 - h3d
		}

		if hd <= 7 && h3d < 0.2 {
			jac := jaccard(m1, m2)
			if jac > 0.8 {
				fmt.Println(s.Hash1, hd, jac)
			}
		}

		/*
			jac := jaccard(m1, m2)
			total++
			if hd <= 7 && h3d < 0.2 {
				if jac > 0.8 {
					fmt.Println(id, hd, jac)
					tp++
				} else {
					fp++
				}
			} else {
				if jac > 0.8 {
					fn++
				} else {
					tn++
				}
			}
			if jac <= 0.8 {
				if hd <= 7 {
					fp1++
				}
				if h3d < 0.2 {
					fp2++
				}
			}
		*/
	}

	/*
		log.Printf("stats: tp=%-5d, tn=%-5d, fp1=%-5d (%7.4f), fp2=%-5d (%7.4f), fp=%-5d (%7.4f), fn=%-5d",
			tp, tn,
			fp1, float64(fp1*100)/float64(total),
			fp2, float64(fp2*100)/float64(total),
			fp, float64(fp*100)/float64(total),
			fn)
		if fp+fn+tp+tn != total {
			panic(fmt.Errorf("%d %d", fp+fn+tp+tn, total))
		}
	*/
}
