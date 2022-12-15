package main

import (
	"testing"
)

func Substr1(str string, lo, hi int) string {
	return string([]byte(str[lo:hi]))
}

func Substr2(str string, lo, hi int) string {
	sub := str[lo:hi]
	return (sub + " ")[:len(sub)]
}

func Substr3(str string, lo, hi int) string {
	sub := str[lo:hi]
	if len(sub) == 0 {
		return ""
	}
	return sub[0:1] + sub[1:]
}

var substrFunctions = []struct {
	name     string
	function func(str string, lo, hi int) string
}{
	{"Substr1", Substr1},
	{"Substr2", Substr2},
	{"Substr3", Substr3},
}

var substrBenchmarks = []struct {
	name                 string
	strLen, subLo, subHi int
}{
	{"Zero  ", 1, 1, 1},
	{"Small ", 4, 1, 4},
	{"Medium", 256, 1, 256},
	{"Large ", 4096, 1, 4096},
}

func benchmarkSubstrSize() {
	//client, err := dbclient.NewBenchClient()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer func(cli *dbclient.Client) {
	//	err := cli.BenchClientClose(cli)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//}(client)
	//fmt.Println("BenchmarkSubstrSize:")
	//for _, benchmark := range substrBenchmarks {
	//	str := strings.Repeat("abc", benchmark.strLen)
	//	for _, function := range substrFunctions {
	//		benchmarkFunc := func(b *testing.B) {
	//			b.ResetTimer()
	//			for i := 0; i < b.N; i++ {
	//				function.function(str, benchmark.subLo, benchmark.subHi)
	//			}
	//			b.StopTimer()
	//		}
	//		results := testing.Benchmark(benchmarkFunc)
	//		out := fmt.Sprintf(
	//			"%10s %10s %5d ns/op %5d allocs/op %5d bytes/op",
	//			benchmark.name,
	//			function.name,
	//			int(results.T)/results.N,
	//			results.AllocsPerOp(),
	//			results.AllocedBytesPerOp(),
	//		)
	//		fmt.Println(out)
	//	}
	//}
}

func TestDBCil(t *testing.T) {
	benchmarkSubstrSize()
}
