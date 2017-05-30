package main

import (
    "log"
    "os"
    "bufio"
    "time"
    "strconv"
    "strings"
    "gopkg.in/errgo.v1"

    "github.com/posteo/go-agentx"
    "github.com/posteo/go-agentx/pdu"
    "github.com/posteo/go-agentx/value"
)

type Filename string

func readfile(f string) (id []int, magnitude []int, err error) {
	file, err := os.Open(f)
	if err != nil {
		return id,magnitude,err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
            s := strings.Split(scanner.Text(), ",")
            d := strings.Split(s[1], "e")
            x,_ :=strconv.Atoi(s[0])
            y,_ :=strconv.Atoi(d[0])
            id = append(id, x)
            magnitude = append(magnitude,y)
	}
	return id,magnitude,scanner.Err()
}

func main() {
    client := &agentx.Client{
        Net:               "tcp",
        Address:           "localhost:705",
        Timeout:           1 * time.Minute,
        ReconnectInterval: 1 * time.Second,
    }

    if err := client.Open(); err != nil {
        log.Fatalf(errgo.Details(err))
    }

    session, err := client.Session()
    if err != nil {
        log.Fatalf(errgo.Details(err))
    }

    listHandler := &agentx.ListHandler{}



    if err := session.Register(127, value.MustParseOID("1.3.6.1.4.1.4171.40")); err != nil {
        log.Fatalf(errgo.Details(err))
    }

    for {
	id,magnitude,_ := readfile("counters.conf")
        item := listHandler.Add("1.3.6.1.4.1.4171.40.1")
        item.Type = pdu.VariableTypeInteger
        item.Value = int32(time.Now().Unix())

        for i:=0;i<len(id);i++  {
            item = listHandler.Add("1.3.6.1.4.1.4171.40."+ strconv.Itoa(id[i]+1))
            item.Type = pdu.VariableTypeCounter32
		s := strconv.FormatInt(int64(magnitude[i])*int64(time.Now().Unix())*1000000, 2)
		last32  := s[len(s)-32:]
		 i,_ := strconv.ParseInt(last32, 2, 64)

            item.Value = uint32(i)
	}

        session.Handler = listHandler
        time.Sleep(500 * time.Millisecond)
    }


}
