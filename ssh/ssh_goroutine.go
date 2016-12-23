package main


import (
        //"bufio"
        //"container/list"
        "golang.org/x/crypto/ssh"
        "regexp"
        "fmt"
        "io/ioutil"
        //"os"
        "strings"
        "bytes"
        "log"
        "time"
)

/*
host 文件格式： ip  [port]  user passwd
*/

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

func GetHosts(lines []string) []hostInfo {
        var hosts []hostInfo
        for _, line := range lines {
                infos := MatchOneHostInfo(line)
                if len(infos) < 9 {
                        continue
                }
                host := hostInfo{ip: infos[1], port: infos[6], user: infos[7], passwd: infos[8]}
                hosts = append(hosts, host)
        }
        return hosts
}

func MatchOneHostInfo(line string) []string {

        r, err := regexp.Compile(`((([12][0-9][0-9]|[1-9][0-9]|[0-9])\.){3,3}([12][0-9][0-9]|[1-9][0-9]|[0-9]))(\s+(\d+))?\s+(.*)\s+(.*)\b`)
        checkErr(err)
        // 匹配内容列表，包含了子串，正常的话 ar[1]:ip,ar[6]:port,ar[7]:user,ar[8]:passwd
        ar := r.FindStringSubmatch(line)
        fmt.Println(ar)
        for i, v := range ar {
                fmt.Println(i, ":", v)
        }
        return ar
}

func ReadHostFile(filename string) []string {
        dat, err := ioutil.ReadFile(filename) 
        checkErr(err)
        lines := strings.Split(string(dat), "\n")
        return lines
}
func main() {
        destLines:=ReadHostFile("host")
        hosts:=GetHosts(destLines)
        results := MultiSsh(hosts,"uptime")
        fmt.Println(results)
}

func checkErr(e error) {
        if e != nil {
                panic(e)
        }
}
