package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-demo/utils/common"
	"go-demo/utils/config"
	"go-demo/utils/logger"

	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"crypto/tls"

	"github.com/gin-gonic/gin"
)

//http响应
type CommonData struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

// 并发请求
type ResultMapStruct struct {
	Url  string `json:"url"`
	Data []byte `json:"data"`
	Err  error  `json:"err"`
}

/**
 * 保存日志 请求参数和返回数据
 * @param  {[type]} c    *gin.Context              [description]
 * @param  {[type]} data map[string]interface{}) ([]byte,      error [description]
 * @return {[type]}      [description]
 */
func saveRespLog(c *gin.Context, data interface{}) error {
	requestParam, _ := c.Get("requestParam")
	requestUrl := c.Request.URL.String()
	reqStartTime, _ := c.Get("requestStartTime")
	costTime := common.End(reqStartTime.(time.Time))

	requestBody, err := json.Marshal(requestParam)
	if err != nil {
		logger.Warnf(c, "saveRespLog::errUrl$v requestParam:%v err%v", requestParam, err)
	}

	_, filePath, line, _ := runtime.Caller(3)
	errorFileSlice := strings.SplitAfterN(filePath, "/controller/", 2)
	if len(errorFileSlice) > 1 {
		filePath = "controller/" + errorFileSlice[1]
	}

	responseBody, err := json.Marshal(data)
	if err != nil {
		logger.Warnf(c, "saveRespLog::errUrl$v responseBody:%v err%v file:%s line:%d", responseBody, err, filePath, line)
	}

	if len(responseBody) > common.LogLenDefault {
		responseBody = []byte("log length too long")
	}

	hostName, _ := os.Hostname()
	logger.Infow("",
		"log_type", "request_log",
		"host_name", hostName,
		"request_time", reqStartTime.(time.Time).UnixNano()/1e6,
		"request_uri", requestUrl,
		"request_body", string(requestBody),
		"response_time", time.Now().UnixNano()/1e6,
		"response_body", string(responseBody),
		"cost_time", costTime,
		"file_path", filePath,
		"line_number", line,
		"req_id", c.GetString("req_id"),
		"req_source", c.GetString("req_source"),
	)
	return nil
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	ret := map[string]interface{}{
		"code": common.ERR_SUC.ErrorCode,
		"msg":  common.ERR_SUC.ErrorMsg,
		"data": data,
	}

	RenderJson(c, ret)
	return
}

// 返回空data Object
func ResponseErrorCodeAndMsg(c *gin.Context, err *common.Error, customMsg string) {
	ret := map[string]interface{}{
		"code": err.ErrorCode,
		"msg":  customMsg,
		"data": common.JsonEmptyObj,
	}

	RenderJson(c, ret)
	return
}

// 返回data数据及错误码
func ResponseErrorCodeAndData(c *gin.Context, err *common.Error, data interface{}) {
	ret := map[string]interface{}{
		"code": err.ErrorCode,
		"msg":  err.ErrorMsg,
		"data": data,
	}

	RenderJson(c, ret)
	return
}

func saveServiceLog(c *gin.Context, requestUrl string, requestTime time.Time, requestParam interface{}, responseData []byte) error {
	costTime := common.End(requestTime)

	if len(responseData) > common.LogLenDefault {
		responseData = []byte("log length too long")
	}
	requestBody, err := json.Marshal(requestParam)
	if err != nil {
		logger.Warnf(c, "saveServiceLog::errUrl$v requestParam:%v err%v", requestParam, err)
	}

	_, filePath, line, _ := runtime.Caller(3)
	errorFileSlice := strings.SplitAfterN(filePath, "", 2)
	if len(errorFileSlice) > 1 {
		filePath = "" + errorFileSlice[1]
	}

	hostName, _ := os.Hostname()
	logger.Infow("",
		"log_type", "service_log",
		"host_name", hostName,
		"request_time", requestTime.UnixNano()/1e6,
		"request_uri", requestUrl,
		"request_body", string(requestBody),
		"response_time", time.Now().UnixNano()/1e6,
		"response_body", string(responseData),
		"cost_time", costTime,
		"file_path", filePath,
		"line_number", line,
		"req_id", c.GetString("req_id"),
		"req_source", config.Config.AppName,
	)
	return nil
}

func RenderJson(c *gin.Context, data interface{}) {
	saveRespLog(c, data) //保存日志
	c.Header("Content-Type", "application/json;charset=UTF-8")
	c.Header("req_id", c.GetString("req_id"))
	c.Header("req_source", c.GetString("req_source"))
	c.JSON(200, data)
	c.Writer.Flush()
	c.Abort()
	return
}

func GetFullUrl(urlStr string, params map[string]string) string {
	v, _ := url.Parse(urlStr)

	paramsData := url.Values{}
	if params != nil {
		for key, value := range params {
			paramsData.Add(key, value)
		}
	}

	v.RawQuery = paramsData.Encode()
	urlPath := v.String()

	return urlPath
}

// HttpGet get请求
func Get(c *gin.Context, host, path string, params map[string]string) ([]byte, error) {
	query := url.Values{}
	if params != nil {
		for key, value := range params {
			query.Add(key, value)
		}
	}
	urlPath := url.URL{
		Scheme:   "http",
		Host:     host,
		Path:     path,
		RawQuery: query.Encode(),
	}
	urlPathStr := urlPath.String()
	req, err := http.NewRequest("GET", urlPathStr, nil)
	if err != nil {
		//logger.Warnf("http get::url:%s ,request:%+v  err:%v", urlStr, params, err)
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req.Header.Set("req_id", c.GetString("req_id"))
	req.Header.Set("req_source", config.Config.AppName)

	requestTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		//logger.Warnf("http get::url:%s ,request:%+v  err:%v", urlStr, params, err)
		return nil, err
	}

	defer resp.Body.Close()

	responseByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//logger.Warnf("http get::url:%s ,request:%+v  err:%v", urlStr, params, err)
		return nil, err
	}

	saveServiceLog(c, urlPathStr, requestTime, params, responseByte)

	if resp.StatusCode != http.StatusOK {
		//logger.Warnf("url:%s ,request:%+v, http statusCode:%d response_data:%+v", urlStr, params, resp.StatusCode, string(respData))
		return nil, errors.New(string(responseByte))
	}

	return responseByte, nil
}

//http get json 自定义header
func GetJsonCusHeader(c *gin.Context, host, path string, headerData, params map[string]string, responseData interface{}) (err error) {
	query := url.Values{}
	if params != nil {
		for key, value := range params {
			query.Add(key, value)
		}
	}
	urlPath := url.URL{
		Scheme:   "http",
		Host:     host,
		Path:     path,
		RawQuery: query.Encode(),
	}
	urlPathStr := urlPath.String()
	req, err := http.NewRequest("GET", urlPathStr, nil)
	if err != nil {
		return err
	}

	for index, val := range headerData {
		req.Header.Set(index, val)
	}
	req.Header.Set("req_id", c.GetString("req_id"))
	req.Header.Set("req_source", config.Config.AppName)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	requestTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	saveServiceLog(c, urlPathStr, requestTime, params, responseByte)

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(responseByte))
	}

	err = json.Unmarshal(responseByte, responseData)
	if err != nil {
		logger.Warnf(c, "Unmarshal response error, data:%s, err:%s", string(responseByte), err.Error())
		return err
	}
	return nil
}

func GetJsonRequest(c *gin.Context, host, path string, params map[string]string, responseData interface{}) (err error) {
	var responseByte []byte
	responseByte, err = Get(c, host, path, params)
	if err != nil {
		logger.Warnf(c, "GetJson error:%s, url:%s, data:%s", err.Error(), host+path, string(responseByte))
		return err
	}
	err = json.Unmarshal(responseByte, responseData)
	if err != nil {
		logger.Warnf(c, "Unmarshal response error, data:%s, err:%s", string(responseByte), err.Error())
		return err
	}
	return nil
}

// HttpPost post请求
func PostJson(c *gin.Context, url string, params []byte) ([]byte, error) {
	body := bytes.NewBuffer(params)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		//logger.Warnf("PostJson::url:%s ,request:%+v  err:%v", url, string(params), err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("req_id", c.GetString("req_id"))
	req.Header.Set("req_source", config.Config.AppName)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		//logger.Warnf("PostJson::url:%s ,request:%+v  err:%v", url, string(params), err)
		return nil, err
	}

	responseByte, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		//logger.Warnf("PostJson::url:%s ,request:%+v  err:%v", url, string(params), err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		//logger.Warnf("url:%s ,request:%+v, http statusCode:%d", url, string(params), resp.StatusCode)
		return nil, errors.New(string(responseByte))
	}

	return responseByte, nil
}

func PostJsonRequest(c *gin.Context, url string, requestData interface{}, responseData interface{}) (err error) {
	if !common.IsValidString(url) {
		return errors.New("PostJsonRequest url is invalid!")
	}

	var jsonData []byte
	jsonData, err = json.Marshal(requestData)
	if err != nil {
		logger.Warnf(c, "Marshal json error:%s", err.Error())
		return err
	}

	requestTime := time.Now()
	var responseByte []byte
	responseByte, err = PostJson(c, url, jsonData)
	saveServiceLog(c, url, requestTime, requestData, responseByte)
	if err != nil {
		logger.Warnf(c, "PostJson error:%s, url:%s, data:%s", err.Error(), url, string(responseByte))
		return err
	}
	//responseByte = []byte(strings.Replace(string(responseByte), "\"data\":[]", "\"data\":{}", 1))
	err = json.Unmarshal(responseByte, responseData)
	if err != nil {
		logger.Warnf(c, "Unmarshal response error, data:%s, err:%s", string(responseByte), err.Error())
		return err
	}
	return nil
}

// 和PostJsonRequest 区别 是 data不做任何处理 原样返回
func PostJsonRequestOrigin(c *gin.Context, url string, requestData interface{}, responseData interface{}) (err error) {
	if !common.IsValidString(url) {
		return errors.New("PostJsonRequest url is invalid!")
	}

	var jsonData []byte
	jsonData, err = json.Marshal(requestData)
	if err != nil {
		logger.Warnf(c, "Marshal json error:%s, req_id:%s", err.Error())
		return err
	}

	requestTime := time.Now()
	var responseByte []byte
	responseByte, err = PostJson(c, url, jsonData)
	saveServiceLog(c, url, requestTime, requestData, responseByte)
	if err != nil {
		logger.Warnf(c, "PostJson error:%s, url:%s, data:%s", err.Error(), url, string(responseByte))
		return err
	}

	err = json.Unmarshal(responseByte, responseData)
	if err != nil {
		logger.Warnf(c, "Unmarshal response error,data:%s, err:%s", string(responseByte), err.Error())
		return err
	}
	return nil
}

func PostJsonRequestReturnJsonByte(c *gin.Context, url string, requestData interface{}) ([]byte, error) {
	var jsonData []byte
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	requestTime := time.Now()
	var responseByte []byte
	responseByte, err = PostJson(c, url, jsonData)
	saveServiceLog(c, url, requestTime, requestData, responseByte)
	if err != nil {
		logger.Warnf(c, "PostJson error url:%s, data:%s, error:%s", url, string(responseByte), err.Error())
		return nil, err
	}

	return responseByte, nil
}

//GetBodyParam 获取json格式的请求数据，keyStruct需要传入指针，该方法会对keyStruct赋值
func GetBodyParam(c *gin.Context, keyStruct interface{}) (err error) {
	err = c.ShouldBind(keyStruct)
	c.Set("requestParam", keyStruct)
	return
}

// 并发请求
// key url
// value map[参数: 参数值]
func PostMapRequestAsync(c *gin.Context, params map[string]interface{}) map[string]*ResultMapStruct {
	finishNum := len(params)
	result := make(chan *ResultMapStruct, finishNum)
	resultMap := make(map[string]*ResultMapStruct, finishNum)
	for urlStr, param := range params {
		go func(urlStr string, param interface{}) {
			bytesArr, err := PostJsonRequestReturnJsonByte(c, urlStr, param)
			rms := ResultMapStruct{Data: bytesArr, Err: err, Url: urlStr}
			finishNum--
			result <- &rms
		}(urlStr, param)
	}

LOOP:
	for {
		select {
		case x := <-result:
			resultMap[x.Url] = x
			if finishNum <= 0 {
				break LOOP
			}
		}
	}

	return resultMap
}

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method //请求方法

		origin := c.Request.Header.Get("Origin") //请求头部
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,HEAD,OPTIONS,UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,X-Requested-With,Authorization,If-None-Match,sid,source,token,app-type,set-token,Cookie")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		c.Header("Access-Control-Max-Age", "1728000")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

func HttpsPostJsonRequestReturnJsonByte(c *gin.Context, url string, requestData interface{}, responseData interface{}) error {
	var jsonData []byte
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	var responseByte []byte

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v, err:%v", url, string(jsonData), err)
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("req_id", c.GetString("req_id"))
	req.Header.Set("req_source", config.Config.AppName)

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	requestTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v, err:%v", url, string(jsonData), err)
		return err
	}

	responseByte, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v, err:%v", url, string(jsonData), err)
		return err
	}

	saveServiceLog(c, url, requestTime, requestData, responseByte)

	if resp.StatusCode != http.StatusOK {
		logger.Warnf(c, "url:%s, request:%+v, http statusCode:%d, response:%+v", url, string(jsonData), resp.StatusCode, string(responseByte))
		return errors.New(string(responseByte))
	}

	err = json.Unmarshal(responseByte, responseData)
	if err != nil {
		logger.Warnf(c, "Unmarshal response error, data:%s, err:%s", string(responseByte), err.Error())
		return err
	}

	return nil
}

func HttpsPostHeaderRequestReturnJsonByte(c *gin.Context, url string, headerData map[string]string, responseData interface{}) error {
	var jsonData []byte
	jsonData, err := json.Marshal("")
	if err != nil {
		return err
	}

	var responseByte []byte

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v, err:%v", url, string(jsonData), err)
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("req_id", c.GetString("req_id"))
	req.Header.Set("req_source", config.Config.AppName)

	for index, val := range headerData {
		req.Header.Set(index, val)
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	requestTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v, err:%v", url, string(jsonData), err)
		return err
	}

	responseByte, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v, err:%v", url, string(jsonData), err)
		return err
	}

	saveServiceLog(c, url, requestTime, headerData, responseByte)

	if resp.StatusCode != http.StatusOK {
		logger.Warnf(c, "url:%s, request:%+v, http statusCode:%d, response:%+v", url, string(jsonData), resp.StatusCode, string(responseByte))
		return errors.New(string(responseByte))
	}

	err = json.Unmarshal(responseByte, responseData)
	if err != nil {
		logger.Warnf(c, "Unmarshal response error, data:%s, err:%s", string(responseByte), err.Error())
		return err
	}

	return nil
}

func HttpFormRequestReturnJsonByte(c *gin.Context, uri string, requestData map[string]string, responseData interface{}) error {
	formData := url.Values{}
	for k, v := range requestData {
		formData.Add(k, v)
	}
	data := formData.Encode()
	req, err := http.NewRequest("POST", uri, strings.NewReader(data))
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v, err:%v", uri, requestData, err)
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("req_id", c.GetString("req_id"))
	req.Header.Set("req_source", config.Config.AppName)

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	requestTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v, err:%v", uri, requestData, err)
		return err
	}

	responseByte, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logger.Warnf(c, "PostJson::url:%s, request:%+v,  err:%v", uri, string(responseByte), err)
		return err
	}

	saveServiceLog(c, uri, requestTime, requestData, responseByte)

	if resp.StatusCode != http.StatusOK {
		logger.Warnf(c, "url:%s, request:%+v, http statusCode:%d, response:%+v", uri, requestData, resp.StatusCode, string(responseByte))
		return errors.New(string(responseByte))
	}

	err = json.Unmarshal(responseByte, responseData)
	if err != nil {
		logger.Warnf(c, "Unmarshal response error, data:%s, err:%s", requestData, err.Error())
		return err
	}

	return nil
}
