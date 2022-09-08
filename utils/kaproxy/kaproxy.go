package kaproxy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-demo/utils/common"
	"go-demo/utils/config"
	"go-demo/utils/http"
	"go-demo/utils/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Key   string
	Value string
}

type ProduceResp struct {
	Partition int   `json:"partition"`
	Offset    int64 `json:"offset"`
}

type KaProxyClient struct {
	context *gin.Context
	scheme  string
	host    string
	port    int
	token   string
}

type ConsumeResp struct {
	Encoding  string `json:"encoding"`
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Offset    int64  `json:"offset"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func NewClientDefault(c *gin.Context) *KaProxyClient {
	return &KaProxyClient{
		context: c,
		scheme:  config.Config.KaProxy.Scheme,
		host:    config.Config.KaProxy.Host,
		port:    config.Config.KaProxy.Port,
		token:   config.Config.KaProxy.Token,
	}
}

func NewClient(c *gin.Context, scheme, host string, port int, token string) *KaProxyClient {
	return &KaProxyClient{
		context: c,
		scheme:  scheme,
		host:    host,
		port:    port,
		token:   token,
	}
}

func (k *KaProxyClient) Produce(topic string, message Message) (*ProduceResp, error) {
	return k.produce(topic, -1, message, false, false)
}

func (k *KaProxyClient) ProduceWithoutReplicate(topic string, message Message) (*ProduceResp, error) {
	return k.produce(topic, -1, message, false, true)
}

func (k *KaProxyClient) ProduceWithHash(topic string, message Message) (*ProduceResp, error) {
	return k.produce(topic, -1, message, true, false)
}

func (k *KaProxyClient) ProduceWithPartition(topic string, partition int, message Message) (*ProduceResp, error) {
	return k.produce(topic, partition, message, false, false)
}

func (k *KaProxyClient) produce(topic string, partition int, message Message, hash, replicate bool) (*ProduceResp, error) {
	var requestUri string
	produceResp := &ProduceResp{}
	requestParams := make(map[string]string)

	requestParams["value"] = message.Value
	if message.Key != "" {
		requestParams["key"] = message.Key
	}
	if partition >= 0 {
		requestUri = fmt.Sprintf("/topic/%s/partition/%d", topic, partition)
	} else {
		requestUri = common.BufferJoin([]string{"/topic/", topic})
		if hash {
			requestParams["partitioner"] = "hash"
		}
	}

	requestUrl := common.BufferJoin([]string{fmt.Sprintf("%s://%s:%d", k.scheme, k.host, k.port), requestUri, "?token=", k.token})
	if replicate {
		common.BufferJoin([]string{requestUrl, "&replicate=no"})
	}
	err := http.HttpFormRequestReturnJsonByte(k.context, requestUrl, requestParams, produceResp)
	return produceResp, err
}

/*
* 消费者，timeout 可选，消费的阻塞时间，单位为ms，0或没有指定的话则不阻塞
* ttr 可选 消息的ttr，单位为ms，默认为300000ms（即5min）只有在atLeastOnce才能生效
 */

func (k *KaProxyClient) Consume(group, topic string, timeout, ttr uint32) (*ConsumeResp, error) {
	consumeResp := &ConsumeResp{}
	host := fmt.Sprintf("%s:%d", k.host, k.port)
	path := common.BufferJoin([]string{"/group/", group, "/topic/", topic})
	requesParams := make(map[string]string, 3)
	requesParams["token"] = k.token
	if timeout > 0 {
		requesParams["timeout"] = strconv.Itoa(int(timeout))
	}
	if ttr > 0 {
		requesParams["ttr"] = strconv.Itoa(int(ttr))
	}

	responseByte, err := http.Get(k.context, host, path, requesParams)
	if err != nil {
		if err.Error() == "" {
			logger.Infof(k.context, "kaProxy get nil")
			return nil, nil
		} else {
			logger.Warnf(k.context, "kaProxy get error:%+v, url:%s, data:%s", err.Error(), common.BufferJoin([]string{host, path}), string(responseByte))
			return nil, err
		}
	}

	err = json.Unmarshal(responseByte, &consumeResp)
	if err != nil {
		logger.Warnf(k.context, "kaProxy responseByte Unmarshal response error, data:%s err:%s ", string(responseByte), err.Error())
		return nil, err
	}

	if consumeResp.Encoding == "base64" {
		valueByte, err := base64.StdEncoding.DecodeString(consumeResp.Value)
		if err != nil {
			logger.Warnf(k.context, "kaProxy responseByte base64 decode error, data:%s err:%s ", consumeResp.Value, err.Error())
			return nil, err
		}
		consumeResp.Value = string(valueByte)
	}

	return consumeResp, nil
}
