package main

import (
	"flag"

	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/log"
	"github.com/tidusant/c3m-common/mycrypto"
	rpsex "github.com/tidusant/chadmin-repo/session"
	//"io" ret

	"net"
	"net/http"
	"net/rpc"
	//	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {

}

func main() {
	var port int
	var debug bool
	//fmt.Println(mycrypto.Encode("abc,efc", 5))
	flag.IntVar(&port, "port", 8081, "help message for flagname")
	flag.BoolVar(&debug, "debug", false, "Indicates if debug messages should be printed in log files")

	flag.Parse()

	//logLevel := log.DebugLevel
	if !debug {
		//logLevel = log.InfoLevel
		gin.SetMode(gin.ReleaseMode)
	}

	// log.SetOutputFile(fmt.Sprintf("portal-"+strconv.Itoa(port)), logLevel)
	// defer log.CloseOutputFile()
	// log.RedirectStdOut()

	log.Infof("running with port:" + strconv.Itoa(port))

	//init config

	router := gin.Default()

	router.POST("/*name", func(c *gin.Context) {
		requestDomain := c.Request.Header.Get("Origin")
		allowDomain := c3mcommon.CheckDomain(requestDomain)
		strrt := ""
		c.Header("Access-Control-Allow-Origin", "*")
		if allowDomain != "" {
			c.Header("Access-Control-Allow-Origin", allowDomain)
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers,access-control-allow-credentials")
			c.Header("Access-Control-Allow-Credentials", "true")
			log.Debugf("check request:%s", c.Request.URL.Path)
			if rpsex.CheckRequest(c.Request.URL.Path, c.Request.UserAgent(), c.Request.Referer(), c.Request.RemoteAddr, "POST") {
				strrt = myRoute(c, "")
			} else {
				log.Debugf("check request error")
			}
		} else {
			log.Debugf("Not allow " + requestDomain)
		}
		if strrt == "" {
			strrt = c3mcommon.Fake64()
		}
		c.String(http.StatusOK, strrt)

	})

	router.Run(":" + strconv.Itoa(port))

}

func myRoute(c *gin.Context, rpcname string) string {
	strrt := ""
	name := c.Param("name")
	name = name[1:] //remove slash
	data := c.PostForm("data")
	log.Debugf("header:%v", c.Request.Header)
	log.Debugf("Request:%v", c.Request)
	userIP, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	log.Debugf("decode name:%s", mycrypto.Decode(name))
	args := strings.Split(mycrypto.Decode(name), "|")
	RPCname := args[0]
	if rpcname != "" {
		RPCname = rpcname
	}

	session := ""
	if len(args) > 1 {
		session = args[1]
	}
	log.Debugf("decode name:%s", mycrypto.Decode(name))
	if RPCname == "CreateSex" {
		data = rpsex.CreateSession()
		log.Debugf("create sex:%s", data)
		strrt = data
	} else {
		log.Debugf("session:%s", session)
		log.Debugf("ip:%s", userIP)
		reply := c3mcommon.ReturnJsonMessage("-5", "unknown error", "", "")
		data = mycrypto.Decode(data)
		if rpsex.CheckSession(session) {

			if RPCname != "aut" {
				//check login
				userid := ""

				client, err := rpc.Dial("tcp", viper.GetString("RPCname.aut"))
				if c3mcommon.CheckError("dial RPCAuth", err) {
					autCall := client.Go("Arith.Run", session+"|"+userIP+"|"+"aut", &userid, nil)
					autreplyCall := <-autCall.Done
					c3mcommon.CheckError("RPCAuth aut ", autreplyCall.Error)
					client.Close()
				} else {
					reply = c3mcommon.ReturnJsonMessage("-1", "service not run", "", "")
				}

				//RPC call
				if userid != "" {

					// Synchronous call
					//								log.Debugf("data:%s", mycrypto.Decode(data))
					//								client, err := rpc.Dial("tcp", ":"+viper.GetString("RPCname.Auth"))
					//								c3mcommon.CheckError("dial RPCAuth", err)
					//								err = client.Call("Arith.Run", mycrypto.Decode(data), &reply)
					//								client.Close()
					//								c3mcommon.CheckError("RPCAuth.Call", err)

					//Asyn call only for http
					client, err := rpc.Dial("tcp", viper.GetString("RPCname."+RPCname))
					if c3mcommon.CheckError("dial RPC"+RPCname+"."+data, err) {
						log.Debugf("Call RPC " + RPCname + " data:" + data)
						log.Debugf("Call RPC " + RPCname + " userid:" + userid)
						rpcCall := client.Go("Arith.Run", session+"|"+userid+"|"+data, &reply, nil)
						rpcreplyCall := <-rpcCall.Done
						c3mcommon.CheckError("RPC"+RPCname+"."+data, rpcreplyCall.Error)
						client.Close()
					} else {
						reply = c3mcommon.ReturnJsonMessage("-1", "service not run", "", "")
					}
				} else {
					//not authorize
					reply = c3mcommon.ReturnJsonMessage("-3", "not authorize", "", "")
				}
			} else {
				client, err := rpc.Dial("tcp", viper.GetString("RPCname.aut"))
				if c3mcommon.CheckError("dial RPCAuth"+"."+data, err) {
					autCall := client.Go("Arith.Run", session+"|"+userIP+"|"+data, &reply, nil)
					autreplyCall := <-autCall.Done
					c3mcommon.CheckError("RPCAuth."+data, autreplyCall.Error)
					client.Close()
				} else {
					reply = c3mcommon.ReturnJsonMessage("-1", "service not run", "", "")
				}
			}

		} else {
			reply = c3mcommon.ReturnJsonMessage("-2", "session not found", "", "")

		}
		if reply != "" {
			// args = strings.Split(data, "|")
			// if RPCname != "news" || args[0] != "l" {

			// 	strrt = mycrypto.Encode(reply, RPCname+"|"+session+data)
			// } else {
			// 	strrt = lzjs.CompressToBase64(reply)
			// }
			log.Debugf("reply", reply)
			strrt = mycrypto.Encode(reply, 8)
		}

	}
	return strrt
}
