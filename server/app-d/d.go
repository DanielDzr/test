package main

import (
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", handleRequest)
    log.Println("Service D is listening on port 80...")
    log.Fatal(http.ListenAndServe(":80", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 检查请求方法是否为 POST
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // 读取请求体
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Println("Failed to read request body:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // 处理接收到的消息
    message := string(body)
    log.Println("Received message from Service C:", message)

}
