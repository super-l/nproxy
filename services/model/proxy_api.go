package model

import (
	"github.com/super-l/nproxy/internal"
	"reflect"
	"time"
)

type mProxyApi struct{}

var MProxyApi = mProxyApi{}

type ProxyApi struct {
	ID             uint64    `json:"id" gorm:"primary_key column:id"`
	ProtocolType   string    `json:"protocol_type" gorm:"protocol_type"`       // 协议类型
	LineType       int       `json:"line_type" gorm:"line_type"`               // 线路类型
	Value          string    `json:"value" gorm:"value"`                       // 值
	Source         string    `json:"source" gorm:"source"`                     // 来源
	GetTimes       int       `json:"used_times" gorm:"get_times"`              // 使用时间
	PeriodValidity int64     `json:"period_validity"  gorm:"gperiod_validity"` // 周期有效性
	CreatedAt      time.Time `json:"created_at" gorm:"created_at"`             // 添加时间
	UpdatedAt      time.Time `json:"updated_at" gorm:"updated_at"`             // 更新时间
}

func (mProxyApi) IsEmpty(data ProxyApi) bool {
	return reflect.DeepEqual(data, ProxyApi{})
}

func (mProxyApi) List() (datas []ProxyApi, err error) {
	var dataList []ProxyApi
	err = internal.GetDbInstance().Model(&ProxyApi{}).Find(&dataList).Error
	return dataList, err
}

func (mProxyApi) GetById(id int) (data ProxyApi, err error) {
	internal.GetDbInstance().Model(&ProxyApi{}).Where("id = ?", id).First(&data)
	return
}

func (mProxyApi) GetByValue(value string) (data ProxyApi, err error) {
	internal.GetDbInstance().Model(&ProxyApi{}).Where("value = ?", value).First(&data)
	return
}

func (mProxyApi) Add(data ProxyApi) (httpAgent *ProxyApi, err error) {
	err = internal.GetDbInstance().Model(&ProxyApi{}).Create(&data).Error
	if err != nil {
		return
	}
	return &data, nil
}

func (mProxyApi) Delete(id uint64) int64 {
	affect := internal.GetDbInstance().Model(&ProxyApi{}).Where("id = ?", id).Delete(&ProxyApi{}).RowsAffected
	return affect
}

func (mProxyApi) Update(data ProxyApi) error {
	err := internal.GetDbInstance().Model(&ProxyApi{}).Save(data).Error
	return err
}

func (mProxyApi) Count() (count int64, err error) {
	err = internal.GetDbInstance().Model(&ProxyApi{}).Count(&count).Error
	return
}

func (mProxyApi) DeleteAll() (err error) {
	err = internal.GetDbInstance().Model(&ProxyApi{}).Where("1 = 1").Delete(&ProxyApi{}).Error
	return
}

func (mProxyApi) DeleteMore(ids []string) int64 {
	var affect int64
	for _, id := range ids {
		affect += internal.GetDbInstance().Model(&ProxyApi{}).Where("id = ?", id).Delete(&Proxy{}).RowsAffected
	}
	return affect
}
