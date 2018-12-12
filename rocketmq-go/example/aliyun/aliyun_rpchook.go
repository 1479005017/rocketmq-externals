package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/apache/incubator-rocketmq-externals/rocketmq-go/remoting"
	"sort"
)

type AliyunRPCHook struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	OnsChannel      string
}

func (h *AliyunRPCHook) DoBeforeRequest(remoteAddr string, request *remoting.RemotingCommand) {
	if request.ExtFields == nil {
		request.ExtFields = map[string]interface{}{}
	}

	request.ExtFields["AccessKey"] = h.AccessKeyId
	request.ExtFields["OnsChannel"] = h.OnsChannel
	request.ExtFields["SecurityToken"] = h.SecurityToken

	content := combineRequestContent(request)
	signature := hmacSHA1(content, h.AccessKeySecret)
	request.ExtFields["Signature"] = signature
}

func (h *AliyunRPCHook) DoAfterResponse(remoteAddr string, request *remoting.RemotingCommand, response *remoting.RemotingCommand) {
}

func combineRequestContent(request *remoting.RemotingCommand) []byte {
	extFields := request.ExtFields
	extFieldKeys := []string{}
	for key, _ := range extFields {
		extFieldKeys = append(extFieldKeys, key)
	}
	sort.Strings(extFieldKeys)

	buf := []byte{}
	for _, extFieldKey := range extFieldKeys {
		buf = combineBytes(buf, toBytes(extFields[extFieldKey]))
	}

	if request.Body != nil {
		buf = combineBytes(buf, request.Body)
	}

	return buf
}

func combineBytes(b1 []byte, b2 []byte) []byte {
	return append(b1, b2...)
}

func toBytes(data interface{}) []byte {
	switch data.(type) {
	case string:
		return []byte(data.(string))
	default:
		buf, err := json.Marshal(data)
		if err != nil {
			return []byte{}
		}
		return buf
	}
}

//hmac ,use sha1
func hmacSHA1(data []byte, key string) string {
	// hmac with sha1
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write(data)

	// base64 encode
	encoded := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return encoded
}
