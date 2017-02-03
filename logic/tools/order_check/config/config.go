package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/go-ini/ini"
)

type FxDBInfo struct {
	User   string
	Pass   string
	Host   string
	DBName string
}

type SettlementCommission struct {
	LevelPerString []string
	LevelPer       []int
}

type Score struct {
	EnlargeScale int
	FollowScore  int
}

type WorkerInfo struct {
	OrderCheckInterval int
	SWMaxWorker        int
	SWMaxChanLen       int
}

type Act struct {
	ActList []string
}

type DuobbSrv struct {
	Host string
}

type Config struct {
	Debug   bool
	Path    string
	Logging bool
	Version string

	FxDBInfo
	SettlementCommission
	WorkerInfo
	Score
	Act
	DuobbSrv
}

func NewConfig() *Config {
	c := new(Config)
	initFlag(c)

	if c.Path == "" {
		fmt.Println("order_check must run with config file, please check.")
		os.Exit(0)
	}

	cfg, err := ini.Load(c.Path)
	if err != nil {
		fmt.Printf("ini[%s] load error: %v\n", c.Path, err)
		os.Exit(1)
	}
	cfg.BlockMode = false
	err = cfg.MapTo(c)
	if err != nil {
		fmt.Printf("config MapTo error: %v\n", err)
		os.Exit(1)
	}

	for _, v := range c.SettlementCommission.LevelPerString {
		vi, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		c.SettlementCommission.LevelPer = append(c.SettlementCommission.LevelPer, vi)
	}

	fmt.Println(c)

	return c
}

func initFlag(c *Config) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	v := fs.Bool("v", false, "Print version and exit")
	fs.StringVar(&c.Path, "c", "", "api server config file.")

	fs.Parse(os.Args[1:])
	fs.Usage = func() {
		fmt.Println("Usage: " + os.Args[0] + " -c api.ini")
		fmt.Printf("\nglobal flags:\n")
		fs.PrintDefaults()
	}

	if *v {
		fmt.Println("version: 0.0.1")
		os.Exit(0)
	}
}
