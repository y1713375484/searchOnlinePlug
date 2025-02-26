# DeepSeek 联网插件

这是一个为硅基流动大模型（DeepSeek）设计的联网插件，允许模型通过调用外部API进行实时信息查询。该插件通过HTTP请求与外部API交互，并将查询结果返回给模型，以便模型能够提供最新的信息和数据。

## 功能

- **实时信息查询**：通过调用外部API，模型可以实时获取最新的信息。
- **插件集成**：插件与硅基流动大模型无缝集成，支持通过函数调用（Function Calling）机制触发查询。
- **日志记录**：所有操作和错误信息都会被记录到日志文件中，便于调试和问题排查。

### 兼容性
目前测试过chatbox是完全兼容的，其他遵循openAI格式的客户端原则上是都兼容的

**chatbox配置方法**

![WechatIMG1289](https://github.com/user-attachments/assets/49b69db0-f119-477e-9a9a-b5c8041681d6)



### 效果
![12141739862175_ pic_hd](https://github.com/user-attachments/assets/85458b54-25f0-41db-81f0-09fa78e079b6)<br><br>
这是大模型为联网之前，无法查询到最新的联网信息<br><br>
![12121739862135_ pic](https://github.com/user-attachments/assets/45b57d25-fedd-4bcb-8840-6cf0afe631a7)
![12131739862154_ pic](https://github.com/user-attachments/assets/577e288b-bc05-4687-814a-c3b55d8e94b3)
这是大模型联网后的效果，可以查询到最新的联网信息，并且可以列出相关的文章信息以及链接

### 支持的模型有
#### Deepseek 系列：<br>
deepseek-ai/DeepSeek-V2.5

#### 书生系列：<br>
internlm/internlm2_5-20b-chat

internlm/internlm2_5-7b-chat

Pro/internlm/internlm2_5-7b-chat

#### Qwen系列：<br>
Qwen/Qwen2.5-72B-Instruct

Qwen/Qwen2.5-32B-Instruct

Qwen/Qwen2.5-14B-Instruct

Qwen/Qwen2.5-7B-Instruct

Pro/Qwen/Qwen2.5-7B-Instruct

#### GLM 系列：<br>
THUDM/glm-4-9b-chat
Pro/THUDM/glm-4-9b-chat



## 环境变量配置

在使用该插件之前，您需要配置以下环境变量：

- `APIKEY`：硅基流动大模型的API密钥。
- `APIURL`：硅基流动大模型的API地址。
- `SEARCHAPIKEY`：智普联网搜索API的密钥。
- `SEARCHAPIURL`：智普联网搜索API的地址。
- `SEARCHAPIMODEL`：智普联网搜索API使用的模型名称。

您可以在项目根目录下的 `.env` 文件中配置这些环境变量。

## 安装与运行

1. **在Releases里下载对应的系统可执行文件**
2. **确保根目录.env文件填写完整**
3. **windows系统直接点击exe文件执行，如果运行失败可以在根目录下log文件查看报错信息**
4. **配置环境变量**

   在项目根目录下创建 `.env` 文件，并填写相应的环境变量：

   ```env
   APIKEY=your_deepseek_api_key
   APIURL=https://api.siliconflow.cn/v1/chat/completions
   SEARCHAPIKEY=your_search_api_key
   SEARCHAPIURL=https://open.bigmodel.cn/api/paas/v4/tools
   SEARCHAPIMODEL=web-search-pro
   ```
智普联网搜索apikey获取方法:[内含aff](https://www.bigmodel.cn/invite?icode=yT8eVZEpgS7b5z7C%2B87nKbC%2Fk7jQAKmT1mpEiZXXnFw%3D)


项目启动后，将在 `localhost:8000` 上运行，并监听 `/v1/chat/completions` 路径的POST请求。





## API 接口

### POST `/v1/chat/completions`

该接口接收来自硅基流动大模型的请求，并根据请求内容决定是否触发联网查询。

**请求体示例**：

```bash
curl -X POST http://127.0.0.1:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
  "model": "deepseek-ai/DeepSeek-V2.5",
  "messages": [
    {
      "role": "user",
      "content": "查询今天的天气"
    }
  ],
  "stream": true,
  "temperature": 0.7

  }' 
```

**响应**：

响应为流式输出，包含模型生成的文本或联网查询的结果。

## 日志

所有操作和错误信息都会被记录到 `searchOnline.log` 文件中，您可以通过查看该文件来调试和排查问题。

## 依赖

- [Gin](https://github.com/gin-gonic/gin)：用于构建HTTP服务器。
- [godotenv](https://github.com/joho/godotenv)：用于加载环境变量。

## 贡献

欢迎提交Issue和Pull Request来改进该项目。

## 许可证

该项目采用MIT许可证，详情请参阅 [LICENSE](LICENSE) 文件。

---

通过该插件，您可以让硅基流动大模型具备实时联网查询的能力，从而提供更加准确和及时的信息。
