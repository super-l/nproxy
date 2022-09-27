package model

import (
	"github.com/super-l/nproxy/internal"
	"github.com/super-l/nproxy/services/rpc/bean"
	"reflect"
	"time"
)

type mProxy struct{}

var MProxy = mProxy{}

type Proxy struct {
	ID           uint64    `json:"id" gorm:"primary_key column:id"`
	ProtocolType string    `json:"protocol_type" gorm:"protocol_type"` // 协议类型
	LineType     int       `json:"line_type" gorm:"line_type"`         // 线路类型
	Value        string    `json:"value" gorm:"value"`                 // 值
	Country      string    `json:"country" gorm:"country"`             // 国家
	Source       string    `json:"source" gorm:"source"`               // 来源
	UsedTimes    int       `json:"used_times" gorm:"used_times"`       // 使用时间
	ExpiredAt    time.Time `json:"expired_at" gorm:"expired_at"`       // 到期时间
	CreatedAt    time.Time `json:"created_at" gorm:"created_at"`       // 添加时间
	UpdatedAt    time.Time `json:"updated_at" gorm:"updated_at"`       // 更新时间
}

func (mProxy) IsEmpty(mProxy Proxy) bool {
	return reflect.DeepEqual(mProxy, Proxy{})
}

func (mProxy) List() (datas []Proxy, err error) {
	var dataList []Proxy
	err = internal.GetDbInstance().Model(&Proxy{}).Find(&dataList).Error
	return dataList, err
}

func (m mProxy) Get(args bean.GetProxyArgs) (data []Proxy, err error) {
	db := internal.GetDbInstance().Model(&Proxy{})

	if args.ProtocolType != "" {
		db.Where("protocol_type", args.ProtocolType)
	}
	if args.LineType != 0 {
		db.Where("line_type", args.LineType)
	}
	if args.Country != "" {
		db.Where("country", args.Country)
	}
	db.Order("updated_at ASC")
	var limit = 1
	if args.Count != 0 {
		limit = args.Count
	}
	err = db.Limit(limit).Find(&data).Error

	// 更新数据
	for _, proxy := range data {
		proxy.UsedTimes = proxy.UsedTimes + 1
		proxy.UpdatedAt = time.Now()
		internal.GetDbInstance().Model(&Proxy{}).Where("id", proxy.ID).Updates(proxy)
	}
	return
}

func (mProxy) GetById(id int) (data Proxy, err error) {
	internal.GetDbInstance().Model(&Proxy{}).Where("id = ?", id).First(&data)
	return
}

func (mProxy) GetByValue(value string) (data Proxy, err error) {
	internal.GetDbInstance().Model(&Proxy{}).Where("value = ?", value).First(&data)
	return
}

func (mProxy) Add(data Proxy) (httpAgent *Proxy, err error) {
	err = internal.GetDbInstance().Model(&Proxy{}).Create(&data).Error
	if err != nil {
		return
	}
	return &data, nil
}

func (mProxy) Delete(id uint64) int64 {
	affect := internal.GetDbInstance().Model(&Proxy{}).Where("id = ?", id).Delete(&Proxy{}).RowsAffected
	return affect
}

func (mProxy) DeleteMore(ids []string) int64 {
	var affect int64
	for _, id := range ids {
		affect += internal.GetDbInstance().Model(&Proxy{}).Where("id = ?", id).Delete(&Proxy{}).RowsAffected
	}
	return affect
}

func (mProxy) Update(data Proxy) error {
	err := internal.GetDbInstance().Model(&Proxy{}).Updates(data).Error
	return err
}

func (mProxy) Count() (count int64, err error) {
	err = internal.GetDbInstance().Model(&Proxy{}).Count(&count).Error
	return
}

func (mProxy) DeleteAll() (err error) {
	err = internal.GetDbInstance().Model(&Proxy{}).Where("1 = 1").Delete(&Proxy{}).Error
	return
}
