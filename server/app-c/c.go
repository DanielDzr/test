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
    log.Println("Service C is listening on port 80...")
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
    log.Println("Received message from Service B:", message)
    // 可以在这里编写处理消息的逻辑，执行相应的业务操作

    //// 响应给 Service C
    sendMessageToServiceD(message)
}

func sendMessageToServiceD(message string) string {
    // Service B 的 URL
    url := "http://app-d" // 替换为 Service D 的实际 URL

    // 构建请求体
    reqBody := []byte(message)

    // 发送 POST 请求给 Service D
    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(reqBody))
    if err != nil {
        return "failed to send message to Service D"
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Sprintf("Service D returned non-200 status: %v", resp.StatusCode)
    }

    return ""

}
