package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type terraform_config struct {
	Options []string `json:options`
	Command  string `json:command`
	Before []string `json:before`
	After  []string `json:after`
}

func loadConfig(c string) terraform_config {
	r:=terraform_config{}
	configFile,err:=os.Open(c)
	if err==nil {
		defer configFile.Close()
		byteValue, _ := ioutil.ReadAll(configFile)
		uerr:=json.Unmarshal(byteValue,&r)
		check(uerr,"Fatal error parsing options file")
	}
	return r
}

func main() {
	var args []string=os.Args
	config_file:="run_config.json"
	for i:=0;i<len(args);i++ {
		if args[i]=="---config-file" {
			config_file=args[i+1]
			args=append(args[0:i],args[i+2:]...)
		}
	}
	cfg:=loadConfig(config_file)
	for _,cmd:=range cfg.Before {
		RunCmd(cmd,nil)
	}
	args=append(args[1:],cfg.Options...)
	RunCmd(cfg.Command, args)
	for _,cmd:=range cfg.After {
		RunCmd(cmd,nil)
	}
}

func GetFile(fn string) ([]string, error) {
	temp, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	s := string(temp)
	return strings.Split(s, "\n"), err
}


func RunCmd(in string, args []string) error {
	cmd := exec.Command(in, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	// SecurityTokenScript -u ndb338  -p xxxxxx
	//out, err := cmd.CombinedOutput()
	err := cmd.Run()
	return err
}

func check(e error,m string) {
	if e!=nil {
		log.Fatal(m,"\n",e)
	}
}
