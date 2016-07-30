package ssphttpapi

import "strings"
import "strconv"
import "log"
import "net/http"
import "net/url"

var accessKey string = ""

var RunSS func(string,string,bool)
var RunKCPTun func(string,string)

var isAuth bool = false

var portAlloc map[int]bool = make(map[int]bool,0)

func isLegalInt(i interface{}) bool {
	switch i.(type) {
		case int:
			return true
		default:
			return false
	}
	return false
}

func isLegalString(i interface{}) bool {
	switch i.(type) {
		case string:
			return true
		default:
			return false
	}
	return false
}

func doAddUser(args_str string) string {
	args,err := url.ParseQuery(args_str)
	if err!=nil {
		return err.Error()
	}
	port_arr,ok := args["port"]
	if !ok {
		return "Bad port"
	}
	port,_ := strconv.Atoi(port_arr[0])

	pw_arr,ok := args["pw"]
	if !ok {
		return "Bad pw"
	}
	pw := pw_arr[0]

	if port<1000 || port>19999 {
		return "Invalid port"
	}

	_,ok = portAlloc[port]
	if ok {
		return "This port is already allocated."
	}
	portAlloc[port]=true

	go RunSS(strconv.Itoa(port), pw, isAuth)
	go RunKCPTun("0.0.0.0:"+strconv.Itoa(port+10000),"127.0.0.1:"+strconv.Itoa(port))
	return "OK"
}

func onRequest(w http.ResponseWriter,r *http.Request) {
	parts := strings.Split(r.URL.Path,"/")
	if len(parts)!=4 {
		w.Write([]byte("Bad request"))
		return
	}
	if parts[1]!=accessKey {
		w.Write([]byte("Permission denied"))
		return
	}
	switch parts[2] {
		case "addUser":
			w.Write([]byte(doAddUser(parts[3])))
			break
		default:
			w.Write([]byte("Function not implemented"))
	}
}

func SetAccessKey(ak string) {
	accessKey=ak
}

func StartServer(listenAddr string) {
	log.Println("Listening on",listenAddr)
	http.HandleFunc("/",onRequest)
	http.ListenAndServe(listenAddr,nil)
}
