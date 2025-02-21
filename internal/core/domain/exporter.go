package domain

import "event-registration/internal/request"

type Transaksi struct {
	ID             string `json:"id" gorm:"column:id"`
	Type           string `json:"type" gorm:"column:type"`
	Title          string `json:"title" gorm:"column:title"`
	Amount         string `json:"amount" gorm:"column:amount"`
	StatusCode     string `json:"status_code" gorm:"column:status_code"`
	Name           string `json:"name" gorm:"column:name"`
	ConsumerName   string `json:"consumer_name" gorm:"column:consumer_name"`
	MeterID        string `json:"meter_id" gorm:"column:meter_id"`
	PaymentGateway string `json:"payment_gateway" gorm:"column:payment_gateway"`
	UnitUP         string `json:"unit_up" gorm:"column:unit_up"`
	UnitAP         string `json:"unit_ap" gorm:"column:unit_ap"`
	UnitUPI        string `json:"unit_upi" gorm:"column:unit_upi"`
	CreatedAt      string `json:"created_at" gorm:"column:created_at"`
	UserID         string `json:"user_id" gorm:"column:user_id"`
	MeterAccountID string `json:"meter_account_id" gorm:"column:meter_account_id"`
	Token          string `json:"token" gorm:"column:token"`
	NameUnitUpi    string `json:"nama_unit_upi" gorm:"column:nama_unit_upi"`
	NameUnitAP     string `json:"nama_unit_ap" gorm:"column:nama_unit_ap"`
	NameUnitUP     string `json:"nama_unit_up" gorm:"column:nama_unit_up"`
}

func (a *Transaksi) TableName() string {
	return "public.transaksi"
}

type ExporterRepository interface {
	FindTransaksi(req *request.TransaksiRequest) ([]*Transaksi, error)
}
