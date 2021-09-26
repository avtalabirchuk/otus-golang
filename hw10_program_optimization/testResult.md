## without optimization test

- `go test -v -count=1 -tags bench . `
```
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:44: time used: 759.975029ms
    stats_optimization_test.go:45: memory used: 291Mb
    assertion_compare.go:332: 
                Error Trace:    stats_optimization_test.go:47
                Error:          "759975029" is not less than "300000000"
                Test:           TestGetDomainStat_Time_And_Memory
                Messages:       [the program is too slow]
--- FAIL: TestGetDomainStat_Time_And_Memory (35.70s)
FAIL
FAIL    github.com/avtalabirchuk/otus-golang/hw10_program_optimization  35.962s
FAIL
```
- `go test -v -count=1 -timeout=30s -tags bench .`
```
=== RUN   TestGetDomainStat_Time_And_Memory
panic: test timed out after 30s

goroutine 38 [running]:
testing.(*M).startAlarm.func1()
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1700 +0xe5
created by time.goFunc
        /usr/local/Cellar/go/1.16.4/libexec/src/time/sleep.go:180 +0x45

goroutine 1 [chan receive]:
testing.(*T).Run(0xc00019a300, 0x123dc54, 0x21, 0x1248700, 0x108bee6)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1239 +0x2da
testing.runTests.func1(0xc00019a180)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1511 +0x78
testing.tRunner(0xc00019a180, 0xc000207de0)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1193 +0xef
testing.runTests(0xc0001a81c8, 0x170a380, 0x1, 0x1, 0xc04c34bd1f542828, 0x6fc348dc4, 0x1714800, 0x1237411)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1509 +0x2fe
testing.(*M).Run(0xc0001fc080, 0x0)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1417 +0x1eb
main.main()
        _testmain.go:43 +0x138

goroutine 19 [chan receive]:
testing.(*B).doBench(0xc0001cc900, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/benchmark.go:281 +0x73
testing.(*B).run(0xc0001cc900)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/benchmark.go:275 +0x78
testing.Benchmark(0xc0001810b0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/benchmark.go:815 +0x125
github.com/avtalabirchuk/otus-golang/hw10_program_optimization.TestGetDomainStat_Time_And_Memory(0xc00019a300)
        /Users/andrey.talabirchuk/golang/otus-golang/hw10_program_optimization/stats_optimization_test.go:42 +0x70
testing.tRunner(0xc00019a300, 0x1248700)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1193 +0xef
created by testing.(*T).Run
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/testing.go:1238 +0x2b3

goroutine 37 [runnable]:
strings.(*Builder).WriteRune(0xc0099b97b0, 0x62, 0x1, 0x0, 0x0)
        /usr/local/Cellar/go/1.16.4/libexec/src/strings/builder.go:104 +0x2ac
regexp/syntax.(*Prog).Prefix(0xc0003c2330, 0xc0003bfea0, 0x0, 0x0)
        /usr/local/Cellar/go/1.16.4/libexec/src/regexp/syntax/prog.go:158 +0x115
regexp.compile(0xc00090b298, 0x5, 0xd4, 0xc00090b298, 0xc00090b29a, 0xc0099b9920)
        /usr/local/Cellar/go/1.16.4/libexec/src/regexp/regexp.go:198 +0x353
regexp.Compile(...)
        /usr/local/Cellar/go/1.16.4/libexec/src/regexp/regexp.go:133
regexp.Match(0xc00090b298, 0x5, 0xc00090b320, 0x10, 0x10, 0xc00090b298, 0x5, 0x0)
        /usr/local/Cellar/go/1.16.4/libexec/src/regexp/regexp.go:560 +0x59
github.com/avtalabirchuk/otus-golang/hw10_program_optimization.countDomains(0x1, 0xc0000201c0, 0xe, 0xc000020078, 0x7, 0xc00001a160, 0x1e, 0xc0000201d0, 0xf, 0xc0000201e0, ...)
        /Users/andrey.talabirchuk/golang/otus-golang/hw10_program_optimization/stats.go:55 +0x175
github.com/avtalabirchuk/otus-golang/hw10_program_optimization.GetDomainStat(0x9a95d98, 0xc00019c050, 0x1234056, 0x3, 0xc00019c050, 0x0, 0x0)
        /Users/andrey.talabirchuk/golang/otus-golang/hw10_program_optimization/stats.go:29 +0x193
github.com/avtalabirchuk/otus-golang/hw10_program_optimization.TestGetDomainStat_Time_And_Memory.func1(0xc0001cc900)
        /Users/andrey.talabirchuk/golang/otus-golang/hw10_program_optimization/stats_optimization_test.go:35 +0x216
testing.(*B).runN(0xc0001cc900, 0x11e3b4)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/benchmark.go:192 +0xeb
testing.(*B).launch(0xc0001cc900)
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/benchmark.go:325 +0xea
created by testing.(*B).doBench
        /usr/local/Cellar/go/1.16.4/libexec/src/testing/benchmark.go:280 +0x55
FAIL    github.com/avtalabirchuk/otus-golang/hw10_program_optimization  30.458s
FAIL
```