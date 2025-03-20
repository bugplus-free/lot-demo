package broker

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)
var Broker_URL = "tcp://192.168.57.231:1883"
var Firefox_URL = "http://www.baidu.com"
var Client_ID string
func init(){
	Client_ID=getIPWithoutDots()
}
func getIPWithoutDots() string {
	// 获取本机所有网络接口
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	// 遍历所有地址，找到非 loopback 的 IPv4 地址
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			// 获取 IP 地址字符串
			ip := ipNet.IP.String()
			// 分割 IP 地址段
			parts := strings.Split(ip, ".")
			// 补齐每个段到三位
			for i := range parts {
				parts[i] = fmt.Sprintf("%03s", parts[i])
			}
			// 将所有部分合并成一个字符串，并返回
			return strings.Join(parts, "")
		}
	}

	// 如果没有找到有效的IPv4地址，返回空字符串
	return ""
}
func openFirefox(url string) error {
	// 调用 Firefox 浏览器并传递 URL
	cmd := exec.Command("firefox", url)
	return cmd.Start() // 执行并返回结果
}
func createClientOptions(broker string, clientID string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(500 * time.Millisecond) // 低功耗设备心跳间隔
	opts.SetAutoReconnect(true)         // 自动重连
	return opts
}
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// 获取消息的主题
	topic := msg.Topic()

	// 根据主题进行不同的处理
	switch topic {
	case "topic/beginfirefox":
		if err := openFirefox(Firefox_URL); err != nil {
			fmt.Println("Error opening Firefox:", err)
			handleTopic(client,msg,1)
		} else {
			fmt.Println("Firefox opened successfully")
			handleTopic(client,msg,0)
		}
		// fmt.Printf("Received message from topic/one: %s\n", msg.Payload())
		// 对 topic/one 的消息进行处理
	case "topic/two":
	default:

	}
}
func handleTopic(client mqtt.Client, msg mqtt.Message,op int) {
	// 获取消息的主题
	topic := msg.Topic()

	// 在原主题后加上 "/response" 作为响应的主题
	responseTopic := topic + "/response"

	// 处理主题的特定逻辑（这里只是打印）
	fmt.Println("Handling message for topic:", topic)
	var responseTopicmessage string
	if op==0{
		responseTopicmessage = fmt.Sprintf("%s/%s/ack", topic, Client_ID)
	}else{
		responseTopicmessage = fmt.Sprintf("%s/%s/fail", topic, Client_ID)
	}
	
	// 发送响应消息，告知消息已处理完
	sendResponse(client, responseTopic, responseTopicmessage)
}
// 发送响应消息
func sendResponse(client mqtt.Client, topic string, message string) {
	var maxRetries = 5 // 最大重试次数
	var retryInterval = 500 * time.Millisecond // 重试间隔
    
	for retries := 0; retries < maxRetries; retries++ {
	    token := client.Publish(topic, 2, false, message) // QoS 2 (消息确认)，retained=false
	    token.Wait()
	    
	    if token.Error() != nil {
		fmt.Printf("Error publishing response: %v. Retrying...\n", token.Error())
		time.Sleep(retryInterval) // 等待重试
	    } else {
		fmt.Printf("Response sent to topic: %s with message: %s\n", topic, message)
		return // 成功发送后退出
	    }
	}
    
	// 如果达到最大重试次数仍未成功，输出错误
	fmt.Printf("Failed to send response after %d retries\n", maxRetries)
}
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	// fmt.Println("Connected")
	// 连接成功后订阅多个主题
	if token := client.Subscribe("topic/beginfirefox", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("Failed to subscribe to topic one: %v\n", token.Error())
	}
	// if token := client.Subscribe("topic/two", 0, nil); token.Wait() && token.Error() != nil {
	// 	fmt.Printf("Failed to subscribe to topic two: %v\n", token.Error())
	// }
	// 你可以继续订阅更多主题
}
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

func ConnectDevice() mqtt.Client {
	opts := createClientOptions(Broker_URL, Client_ID)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
	    panic(token.Error())
	}
	return client
}
