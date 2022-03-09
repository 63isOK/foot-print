package main

import "errors"

// ErrInstanceNotFound 实例未找到,需同步
var ErrInstanceNotFound = errors.New("根据实例ID未找到实例")

// CommonResult 通用的响应参数
type CommonResult struct {
	Success bool
	Message string
}

// PostData 上报数据
type PostData struct {
	InstanceID       uint             `json:"instance_id"`        // 云云对接实例id
	Type             string           `json:"type"`               // 上报的数据类型，取值为"property"，"event"，"status"，"alarm"
	DeviceID         uint             `json:"device_id"`          // 设备id (ps:对接服务需要进行转换)
	ProductModelName string           `json:"product_model_name"` // 设备型号名，对应第三方云平台的产品
	DeviceIdentifier DeviceIdentifier `json:"device_identifier"`  // 设备唯一标识，供第三方云搜索用,具体字段详见下文
	Identifier       string           `json:"identifier"`         // 上报数据的标识， 物模型标识或告警名称
	Data             interface{}      `json:"data"`               // 上报的数据
}

// DeviceIdentifier 设备唯一标识
type DeviceIdentifier struct {
	Value                   string `json:"value"`                      // 标识值,可以是设备编号/设备mac地址/ip/串口号
	ProductModelName        string `json:"product_model_name"`         // 型号名
	ParentValue             string `json:"parent_value"`               // 父设备标识值
	ParentProductModelName  string `json:"parent_product_model_name"`  // 父设备型号名
	GatewayValue            string `json:"gateway_value"`              // 网关标识值
	GatewayProductModelName string `json:"gateway_product_model_name"` // 网关型号名
}

// InstanceConfigs  实例配置列表
type InstanceConfigs []string
