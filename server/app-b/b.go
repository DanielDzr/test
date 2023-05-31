package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", handleRequest)
    log.Println("Service B is listening on port 80...")
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
    log.Println("Received message from Service A:", message)
    // 可以在这里编写处理消息的逻辑，执行相应的业务操作

    //// 响应给 Service C
    sendMessageToServiceC(message)
}

func sendMessageToServiceC(message string) string {
    // Service B 的 URL
    url := "http://app-c" // 替换为 Service C 的实际 URL

    // 构建请求体
    reqBody := []byte(message)

    // 发送 POST 请求给 Service C
    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(reqBody))
    if err != nil {
        return "failed to send message to Service C"
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Sprintf("Service C returned non-200 status: %v", resp.StatusCode)
    }

    return ""

}
