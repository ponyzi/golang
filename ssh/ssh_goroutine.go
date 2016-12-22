package main

import (
        //"bufio"
        //"container/list"
        "golang.org/x/crypto/ssh"
        //"encoding/csv"
        "fmt"
        //"io"
        //"os"
        //"strings"
        "bytes"
        "log"
        "time"
)

type printOut struct {
        dest      string
        bufferOut string
}

type hostInfo struct {
        ip     string
        port   string
        user   string
        passwd string
}

func OneSsh(destInfo hostInfo, cmd string, resultChan chan printOut) {
        // hostInfo [ip[:port],user, passwd]
        config := &ssh.ClientConfig{
                User: destInfo.user,
                Auth: []ssh.AuthMethod{
                        ssh.Password(destInfo.passwd),
                },
        }
        var addr = ""
        if len(destInfo.port)>0 {
                addr = destInfo.ip + ":" + destInfo.port
        } else {
                addr = destInfo.ip
        }
        client, err := ssh.Dial("tcp", addr, config)
        checkErr(err)
        // Each ClientConn can support multiple interactive sessions,represented by a Session.
        session, err := client.NewSession()
        checkErr(err)
        defer session.Close()
        // Once Session created, you can execute a command on the remote side using the Run method.
        var b bytes.Buffer
        session.Stdout = &b
        if err := session.Run(cmd); err != nil {
                log.Fatal(destInfo.ip + "-- Failed to run: " + err.Error())
        }
        resultChan <- printOut{dest: destInfo.ip, bufferOut: b.String()}
        //fmt.Println(b.String())
}

func MultiSsh(hostInfos []hostInfo, cmd string) []printOut {
        fmt.Println(time.Now())
        resultChan := make(chan printOut, len(hostInfos))
        defer close(resultChan)
        var result []printOut
        for _, hostInfo := range hostInfos {
                go OneSsh(hostInfo, cmd, resultChan)
        }
        for i := 0; i < len(hostInfos); i++ {
                res := <-resultChan
                result = append(result, res)
        }
        fmt.Println(time.Now())
        return result
}

func main() {
        hosts := []hostInfo{
                hostInfo{"10.157.0.xxx","22","uuu","pppppp"},
                hostInfo{"10.157.xxx.xx","22","uuu","pppppp"},
                hostInfo{"10.157.xxx.xx","22","uuu","pppppp"},
                hostInfo{"10.157.xxx.xxx","22","uuu","pppppp"},
        }
        results := MultiSsh(hosts,"uptime")
        fmt.Println(results)
}

func checkErr(e error) {
        if e != nil {
                panic(e)
        }
}
