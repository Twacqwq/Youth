package cmd

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	utils "github.com/Twacqwq/youth/pkg/utils"
	pkg "github.com/Twacqwq/youth/pkg/youth"
)

var (
	configDir = flag.String("c", ".", "load config.json")
)

func Youth() {
	flag.Parse()

	f, err := utils.LoadConfig(*configDir)
	if err != nil {
		log.Fatalf("failed load config.json. is exist? %v", err)
		os.Exit(1)
	}
	defer f.Close()

	var m []pkg.Member
	chYouth := make(chan pkg.Member, 10)
	results := make(chan pkg.Member, 10)

	dec := json.NewDecoder(f)
	err = dec.Decode(&m)
	if err != nil {
		log.Fatalf("decode failed %v", err)
		os.Exit(1)
	}

	for i := 0; i < cap(chYouth); i++ {
		go pkg.Worker(chYouth, results)
	}
	go pkg.Push(chYouth, m)
	pkg.CheckStatus(len(m), results)
}
