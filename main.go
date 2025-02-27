package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	apiKey         string
	apiUrl         string
	searchApiUrl   string
	searchApiKey   string
	searchApiModel string
)

// 鉴权密钥
var Authorization string

// 请求deepseek API响应结构体
type DeepseekRJson struct {
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
			ToolCalls        []struct {
				Index    int    `json:"index"`
				Id       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"delta"`
	} `json:"choices"`
}

// 请求智普联网查询API响应结构体
type SearchOnlineStruct struct {
	Choices []struct {
		Message struct {
			Tool_calls []struct {
				Id            string `json:"id"`
				Search_result []struct {
					Content string `json:"content"`
					Icon    string `json:"icon"`
					Index   int    `json:"index"`
					Link    string `json:"link"`
					Media   string `json:"media"`
					Refer   string `json:"refer"`
					Title   string `json:"title"`
				} `json:"search_result"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

// 请求模型结构体
type Data struct {
	Model       string   `json:"model"`
	Messages    Messages `json:"messages"`
	Stream      bool     `json:"stream"`
	Temperature float64  `json:"temperature"`
	Tools       Tools    `json:"tools"`
}

type Messages []struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallId string `json:"tool_call_id"`
}

// 模型func call插件结构体
type Tools []struct {
	Function struct {
		Description string      `json:"description"`
		Name        string      `json:"name"`
		Parameters  interface{} `json:"parameters"`
		Required    []string    `json:"required"`
	} `json:"function"`
	Type string `json:"type"`
}

// 定义 Query 工具的参数结构
type QueryParameters struct {
	Query struct {
		Description string `json:"description"`
		Type        string `json:"type"`
	} `json:"query"`
}

// 定义 Prompt 工具的参数结构
type PromptParameters struct {
	Prompt struct {
		Description string `json:"description"`
		Type        string `json:"type"`
	} `json:"prompt"`
}

// 响应给客户端的数据格式
type RespData struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func main() {

	// 打开或创建日志文件
	file, err := os.OpenFile("searchOnline.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 设置日志输出到文件
	log.SetOutput(file)

	r := gin.Default()
	//加载env配置文件
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey = os.Getenv("APIKEY")
	apiUrl = os.Getenv("APIURL")
	searchApiKey = os.Getenv("SEARCHAPIKEY")
	searchApiUrl = os.Getenv("SEARCHAPIURL")
	searchApiModel = os.Getenv("SEARCHAPIMODEL")
	Authorization = os.Getenv("AUTHORIZATION")
	envMap := map[string]string{
		"APIKEY":           apiKey,
		"APIURL":           apiUrl,
		"SEARCHAPIMODEURL": searchApiUrl,
		"SEARCHAPIMODEL":   searchApiModel,
		"searchApiKEY":     searchApiKey,
	}
	for k, v := range envMap {
		if v == "" {
			log.Println("请检查根目录下env文件中" + k + "填写情况")
			return
		}
	}
	r.Use(AuthCheck())
	r.POST("/v1/chat/completions", Action)
	r.Run(":8000")

}

// 鉴权是否携带了token
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if Authorization == "" {
			c.Next()
		} else {
			AuthHeader := c.Request.Header.Get("Authorization")
			//从请求头读取到的是Bearer sk-xxxxxxxxxxx，所以需要删除前面的部分

			errAuth := func() {
				c.Header("Content-Type", "text/event-stream")
				c.Header("Transfer-Encoding", "chunked")
				data := RespData{
					Choices: []struct {
						Delta struct {
							Content string `json:"content"`
						} `json:"delta"`
					}{
						{
							Delta: struct {
								Content string `json:"content"`
							}{
								Content: "Authorization错误请检查！",
							},
						},
					},
				}
				marshal, err := json.Marshal(data)
				if err != nil {
					fmt.Println(err)
					log.Println(err)
				}
				rData := "data: " + string(marshal) + "\n\n"
				c.Writer.Write([]byte(rData))
			}

			if len(AuthHeader) > 7 && Authorization == AuthHeader[7:] {
				c.Next()
			} else {
				errAuth()
				c.Abort()
			}

		}

	}
}

func Action(c *gin.Context) {
	tools := Tools{
		{
			Type: "function",
			Function: struct {
				Description string      `json:"description"`
				Name        string      `json:"name"`
				Parameters  interface{} `json:"parameters"`
				Required    []string    `json:"required"`
			}{
				Description: "The function sends a query to the browser and returns relevant results based on the search terms provided. The model should avoid using this function if it already possesses the required information or can provide a confident answer without external data",
				Name:        "searchOnline",
				Parameters: QueryParameters{
					Query: struct {
						Description string `json:"description"`
						Type        string `json:"type"`
					}{
						Description: "What to search for",
						Type:        "string",
					},
				},
				Required: []string{"query"},
			},
		},
		{
			Type: "function",
			Function: struct {
				Description string      `json:"description"`
				Name        string      `json:"name"`
				Parameters  interface{} `json:"parameters"`
				Required    []string    `json:"required"`
			}{
				Description: "Generate an image based on a given prompt",
				Name:        "generateImage",
				Parameters: PromptParameters{
					Prompt: struct {
						Description string `json:"description"`
						Type        string `json:"type"`
					}{
						Description: "A text prompt describing the image to be generated",
						Type:        "string",
					},
				},
				Required: []string{"prompt"},
			},
		},
	}

	data := Data{}
	c.BindJSON(&data)

	data.Tools = tools

	byteData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	response := PostApi(apiUrl, byteData) //转发用户请求给硅基流动

	defer response.Body.Close()

	reader := bufio.NewReader(response.Body)
	funcName := ""  //func call的函数名称
	arguments := "" //func call的参数
	c.Header("Content-Type", "text/event-stream")
	c.Header("Transfer-Encoding", "chunked")
	for {
		readBytes, err := reader.ReadBytes('\n')
		rjson := DeepseekRJson{}
		if err != nil {
			if err == io.EOF {
				if funcName == "" {
					c.Writer.Write([]byte(" "))
				}
				break
			}
		}

		strBuffer := strings.TrimSpace(string(readBytes))

		strBuffer = strings.TrimPrefix(strBuffer, "data: ")
		strBuffer = strings.TrimPrefix(strBuffer, "data:")

		if strBuffer == "[DONE]" {
			//解析完毕了
			if funcName == "" {
				readBytes = append(readBytes, '\n')
				c.Writer.Write(byteData)
			}
			break
		}

		if strBuffer == "" {
			continue
		}

		if err := json.Unmarshal([]byte(strBuffer), &rjson); err == nil {
			if rjson.Choices != nil {
				//防止索引越界异常
				if len(rjson.Choices) > 0 {
					//如果有func call响应
					if rjson.Choices[0].Delta.ToolCalls != nil {
						//检查方法名称是否存在
						if rjson.Choices[0].Delta.ToolCalls[0].Function.Name != "" {
							funcName += rjson.Choices[0].Delta.ToolCalls[0].Function.Name
						}
						//拼接参数
						arguments += rjson.Choices[0].Delta.ToolCalls[0].Function.Arguments
						continue
					} else {
						if rjson.Choices[0].Delta.Content != "" {
							readBytes = append(readBytes, '\n')
							c.Writer.Write(readBytes)
						}
						continue
					}
				}

				continue
			} else {
				c.Writer.Write([]byte(strBuffer))
				continue
			}

		} else {
			//c.Writer.Write([]byte("系统解析响应时遇到异常：" + err.Error()))
			errMap := map[string]interface{}{
				"choices": []map[string]interface{}{
					{
						"delta": map[string]interface{}{
							"content": "系统解析响应时遇到异常：" + err.Error(),
						},
					},
				},
			}
			marshal, err := json.Marshal(errMap)
			if err != nil {
				fmt.Println(err)
				log.Println(err)
			}
			rData := "data: " + string(marshal) + "\n\n"
			c.Writer.Write([]byte(rData))

			return
		}
	}
	switch funcName {
	case "searchOnline":
		DoSearchOnline(arguments, data, c)
		break
	case "generateImage":
		DoGenerateImage(arguments, c)
		break
	}

}

// 请求智普搜索api
func SearchOnline(query string) (SearchOnlineStruct, error) {
	searchJson := map[string]interface{}{
		"tool": searchApiModel,
		"messages": []map[string]interface{}{
			{"role": "user", "content": query},
		},
		"stream": false,
	}
	marshal, err := json.Marshal(searchJson)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return SearchOnlineStruct{}, err
	}
	request, err := http.NewRequest("POST", searchApiUrl, bytes.NewReader(marshal))
	if err != nil {
		return SearchOnlineStruct{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+searchApiKey)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return SearchOnlineStruct{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return SearchOnlineStruct{}, err
	}
	searchResp := SearchOnlineStruct{}

	json.Unmarshal(body, &searchResp)
	return searchResp, nil
}

// 调用联网func call
func DoSearchOnline(arguments string, data Data, c *gin.Context) {
	argumentsJson := struct {
		Query string `json:"query"`
	}{}
	err := json.Unmarshal([]byte(arguments), &argumentsJson)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	searchOnline, err := SearchOnline(argumentsJson.Query)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	searchResult := searchOnline.Choices[0].Message.Tool_calls[1].Search_result
	for _, result := range searchResult {

		datas := RespData{
			Choices: []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			}{
				{
					Delta: struct {
						Content string `json:"content"`
					}{
						Content: "标题：" + result.Title + "\n\n内容：" + result.Content + "\n\n来源：" + "[" + result.Media + "](" + result.Link + ")" + "\n\n",
					},
				},
			},
		}

		data.Messages = append(data.Messages, struct {
			Role       string `json:"role"`
			Content    string `json:"content"`
			ToolCallId string `json:"tool_call_id"`
		}{Role: "tool", Content: result.Content, ToolCallId: searchOnline.Choices[0].Message.Tool_calls[1].Id})

		marshal, err2 := json.Marshal(datas)
		if err2 != nil {
			fmt.Println(err2)
		}
		rData := "data: " + string(marshal) + "\n\n"

		c.Writer.Write([]byte(rData))
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	resp := PostApi(apiUrl, marshal) //携带联网查询后的结果重新请求硅基流动模型
	defer resp.Body.Close()
	newReader := bufio.NewReader(resp.Body)
	for {
		readString, err := newReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		c.Writer.Write([]byte(readString))
	}
	c.Writer.Write([]byte(" "))

}

// 调用画图func call
func DoGenerateImage(arguments string, c *gin.Context) {
	argumentsJson := struct {
		Prompt string `json:"prompt"`
	}{}
	err := json.Unmarshal([]byte(arguments), &argumentsJson)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	url := fmt.Sprintf("https://image.pollinations.ai/prompt/%s?width=1024&height=1024&seed=100&model=flux&nologo=true", argumentsJson.Prompt)

	data := RespData{
		Choices: []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		}{
			{
				Delta: struct {
					Content string `json:"content"`
				}{
					Content: fmt.Sprintf("![%s](%s)", argumentsJson.Prompt, url),
				},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}

	rData := "data: " + string(jsonData) + "\n\n"
	c.Writer.Write([]byte(rData))

}

func PostApi(url string, data []byte) *http.Response {

	request, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	return response
}
