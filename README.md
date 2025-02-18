该项目是基于硅基流动+智普搜索的一个ai联网插件，目前支持的联网模型有如下：<br>

![image](https://github.com/user-attachments/assets/1764828b-0685-4e42-8755-f6cdb29bcd2b)

## 原理
1.用户请求插件api<br>
2.插件给用户的提问请求添加一个func call格式如下，并转发请求给硅基流动<br>
```
{
	"model": "Qwen/Qwen2.5-72B-Instruct",
	"messages": [{
			"role": "user",
			"content": "英伟达最新的显卡型号"
		}

	],
	"stream": true,
	"temperature": 0.6,
	"tools": [{
		"function": {
			"description": "The function sends a query to the browser and returns relevant results based on the search terms provided. The model should avoid using this function if it already possesses the required information or can provide a confident answer without external data",
			"name": "searchOnline",
			"parameters": {
				"query": {
					"description": "What to search for",
					"type": "string"
				}
			},
			"required": ["searchContent"]
		},
		"type": "function"
	}]

}
```
<br>
3.根据硅基流动响应判断是否调用了func call，如果没调用那么直接将硅基流动的响应发送给客户端<br>
4.如果func call调用了，那么将参数提取，请求智普搜索api，然后将搜索内容响应给客户端，并携带搜索内容重新请求硅基流动，最后将硅基流动响应转发给客户端<br>

## 兼容性
测试了chatbox客户端实测没问题，别的客户端理论上应该也是兼容的。

## 联网效果
轨迹流动Qwen模型联网前提问：<br>

![12141739862175_ pic_hd](https://github.com/user-attachments/assets/ca1d589c-bf11-4797-88e6-92eb4334788d)

chatbox使用了本插件后提问效果：<br>

![12121739862135_ pic](https://github.com/user-attachments/assets/452e76b9-797f-4e42-9590-272df1350ea8)

![12131739862154_ pic](https://github.com/user-attachments/assets/31974686-636a-4324-9fa3-a2e5844a851f)



## 智扑联网搜索api key获取方法
[获取key（内含aff）](https://www.bigmodel.cn/invite?icode=yT8eVZEpgS7b5z7C%2B87nKbC%2Fk7jQAKmT1mpEiZXXnFw%3D)<br>

![image](https://github.com/user-attachments/assets/9d5fb3b9-79ca-4d53-940c-1581298047a7)

## 硅基流动api key获取方法
[获取key（内含aff）](https://cloud.siliconflow.cn/i/zE8h2FaP)<br>
![image](https://github.com/user-attachments/assets/bac3ceac-6c73-48af-a41e-1cb4c24f6906)

## 运行方法
运行前确保根目录下有.env配置文件，并完整填写好配置信息<br>

如果下载代码运行确保电脑上安装了go环境，然后用命令行切换到项目根目录下执行 go run main.go<br>

或者下载releases中打包好的二进制文件，直接运行就行了。

## 使用方法
与gpt使用方法一致，无需传入Authorization
```bash
curl -X POST http://127.0.0.1:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "Qwen/Qwen2.5-32B-Instruct",
    "messages": [
        {
            "role": "user",
            "content": "帮我在浏览器查询一下东北的美食有哪些"
        }
    ],
    "stream":true,
   "temperature": 0.6
}' 
```



