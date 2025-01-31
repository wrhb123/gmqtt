package thingspanel

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pborman/uuid"
)

type MqttClient struct {
	Client mqtt.Client
	IsFlag bool
}

var DefaultMqttClient *MqttClient = &MqttClient{}

func (c *MqttClient) MqttInit() error {
	// 掉线重连
	// var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	// 	fmt.Printf("Mqtt Connect lost: %v", err)
	// 	i := 0
	// 	for {

	// 		time.Sleep(5 * time.Second)
	// 		if !c.Client.IsConnectionOpen() {
	// 			i++
	// 			fmt.Println("Mqtt客户端掉线重连...", i)
	// 			if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
	// 				fmt.Println("Mqtt客户端连接失败...")
	// 			} else {
	// 				break
	// 			}
	// 		} else {
	// 			//subscribe(msgProc1, gatewayMsgProc)
	// 			break
	// 		}
	// 	}
	// }
	opts := mqtt.NewClientOptions()
	opts.SetUsername("root")
	opts.SetPassword("root")
	opts.AddBroker("localhost:1883")
	opts.SetAutoReconnect(true)
	opts.SetOrderMatters(false)
	//opts.OnConnectionLost = connectLostHandler
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		fmt.Println("Mqtt客户端已连接")
	})
	reconnec_number := 0
	uuid := uuid.New()
	opts.SetClientID(uuid)
	for { // 失败重连

		c.Client = mqtt.NewClient(opts)
		if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
			reconnec_number++
			fmt.Println("错误说明：", token.Error().Error())
			fmt.Println("Mqtt客户端连接失败...重试", reconnec_number)
		} else {
			fmt.Println("Mqtt客户端重连成功")
			c.IsFlag = true
			break
		}
		time.Sleep(5 * time.Second)
	}
	// Log.Error("连接MqttClIent...")
	// if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
	// 	Log.Error("MqttClIent连接失败...")
	// }
	return nil
}

func (c *MqttClient) SendData(topic string, data []byte) error {
	go func() {
		Log.Info("检查MqttClIent连接状态...")
		if !c.IsFlag {
			i := 1
			for {
				fmt.Println("等待...", i)
				if i == 10 || c.IsFlag {
					break
				}
				time.Sleep(1 * time.Second)
				i++
			}
		}
		Log.Info("发送设备状态...")
		token := c.Client.Publish(topic, 1, false, string(data))
		if !token.WaitTimeout(5 * time.Second) {
			Log.Warn("发送设备状态超时")
		} else if err := token.Error(); err != nil {
			Log.Warn("发送设备状态失败: " + err.Error())
		}
		Log.Info("发送设备状态完成")
	}()
	return nil
}
