package main

import (
    "bytes"
    "fmt"
    "log"
    "net/http"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
    // 创建 AWS 会话
    sess, err := session.NewSession(&aws.Config{
        // 添加 AWS 区域
        Region: aws.String("YOUR_REGION"),
        // 添加 AWS 访问凭证信息
        Credentials: credentials.NewStaticCredentials("<YOUR_ACCESS_KEY_ID>", "<YOUR_SECRET_ACCESS_KEY>", ""),
    })
    if err != nil {
        log.Fatal("Failed to create AWS session:", err)
    }

    // 创建 SQS 服务客户端
    svc := sqs.New(sess)

    // 启动服务 A
    go startServiceA(svc)

    // 防止程序退出
    select {}
}

func startServiceA(svc *sqs.SQS) {
    // 循环接收来自 SQS 的消息
    for {
        result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
            QueueUrl:            aws.String("YOUR_QUEUE_URL"), // 替换为您的 SQS 队列 URL
            MaxNumberOfMessages: aws.Int64(1),
            WaitTimeSeconds:     aws.Int64(20),
        })
        if err != nil {
            log.Println("Error receiving message from SQS:", err)
            continue
        }

        // 检查是否接收到消息
        if len(result.Messages) > 0 {
            // 获取第一条消息
            message := result.Messages[0]
            println("Received message from SQS:", *message.Body)
            // 将消息传递给服务 BC
            sendMessageToServiceB(*message.Body)

            // 删除已处理的消息
            _, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
                QueueUrl:      aws.String("YOUR_QUEUE_URL"), // 替换为您的 SQS 队列 URL
                ReceiptHandle: message.ReceiptHandle,
            })
            if err != nil {
                log.Println("Error deleting message from SQS:", err)
            }
        }
    }
}

func sendMessageToServiceB(message string) string {
    // Service B 的 URL
    url := "http://app-b:8081" // 替换为 Service B 的实际 URL

    // 构建请求体
    reqBody := []byte(message)

    // 发送 POST 请求给 Service B
    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(reqBody))
    if err != nil {
        return "failed to send message to Service B"
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Sprintf("Service B returned non-200 status: %v", resp.StatusCode)
    }

    return ""

}
