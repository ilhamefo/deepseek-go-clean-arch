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

type Pelanggan struct {
	ID           string `json:"id" gorm:"column:id"`
	IDPel        string `json:"idpel" gorm:"column:idpel"`
	Name         string `json:"name" gorm:"column:name"`
	ConsumerName string `json:"consumer_name" gorm:"column:consumer_name"`
	EnergyType   string `json:"energy_type" gorm:"column:energy_type"`
	KWH          string `json:"kwh" gorm:"column:kwh"`
	Address      string `json:"address" gorm:"column:address"`
	MeterNo      string `json:"meter_no" gorm:"column:meter_no"`
	MeterType    string `json:"meter_type" gorm:"column:meter_type"`
	UnitUpi      string `json:"unit_upi" gorm:"column:unit_upi"`
	NamaUnitUpi  string `json:"nama_unit_upi" gorm:"column:nama_unit_upi"`
	UnitAp       string `json:"unit_ap" gorm:"column:unit_ap"`
	NamaUnitAp   string `json:"nama_unit_ap" gorm:"column:nama_unit_ap"`
	UnitUp       string `json:"unit_up" gorm:"column:unit_up"`
	NamaUnitUp   string `json:"nama_unit_up" gorm:"column:nama_unit_up"`
	CreatedAt    string `json:"created_at" gorm:"column:created_at"`
	LastUpdate   string `json:"last_update" gorm:"column:last_update"`
}

func (a *Pelanggan) TableName() string {
	return "public.mv_idpel_detail"
}

type Regional struct {
	ID           string  `json:"id" gorm:"column:id"`
	NamaRegional string  `json:"nama_regional" gorm:"column:nama_regional"`
	IsActive     bool    `json:"is_active" gorm:"column:is_active"`
	IDRegAPKT    string  `json:"id_reg_apkt" gorm:"column:id_reg_apkt;primaryKey"`
	Induk        []Induk `json:"induk" gorm:"foreignKey:IDRegAPKT"`
}

func (a *Regional) TableName() string {
	return "public.pln_unit_regional"
}

type Induk struct {
	IDUnitUPI   string `json:"id_unit_upi" gorm:"column:id_unit_upi;primaryKey"`
	Satuan      string `json:"satuan" gorm:"column:satuan"`
	NamaUnitUPI string `json:"nama_unit_upi" gorm:"column:nama_unit_upi"`
	IDRegAPKT   string `json:"id_reg_apkt" gorm:"column:id_reg_apkt"`
	Kontak      string `json:"kontak" gorm:"column:kontak"`
	Nama        string `json:"nama" gorm:"column:nama"`
	Area        []Area `json:"area" gorm:"foreignKey:IDUnitUPI"`
}

func (a *Induk) TableName() string {
	return "public.pln_unit_upi"
}

type Area struct {
	IDUnitAP   string `json:"id_unit_ap" gorm:"column:id_unit_ap;primaryKey"`
	Satuan     string `json:"satuan" gorm:"column:satuan"`
	NamaUnitAP string `json:"nama_unit_ap" gorm:"column:nama_unit_ap"`
	Alamat     string `json:"alamat" gorm:"column:alamat"`
	Telepon    string `json:"telepon" gorm:"column:telepon"`
	Fax        string `json:"fax" gorm:"column:fax"`
	IDUnitUPI  string `json:"id_unit_upi" gorm:"column:id_unit_upi"`
	MANAGER    string `json:"MANAGER" gorm:"column:MANAGER"`
	Kontak     string `json:"kontak" gorm:"column:kontak"`
	Nama       string `json:"nama" gorm:"column:nama"`
	Unit       []Unit `json:"unit" gorm:"foreignKey:IDUnitAP"`
}

func (a *Area) TableName() string {
	return "public.pln_unit_ap"
}

type Unit struct {
	IDUnitUP   string `json:"id_unit_up" gorm:"column:id_unit_up;primaryKey"`
	Satuan     string `json:"satuan" gorm:"column:satuan"`
	NamaUnitUP string `json:"nama_unit_up" gorm:"column:nama_unit_up"`
	Alamat     string `json:"alamat" gorm:"column:alamat"`
	Telepon    string `json:"telepon" gorm:"column:telepon"`
	Fax        string `json:"fax" gorm:"column:fax"`
	IDUnitAP   string `json:"id_unit_ap" gorm:"column:id_unit_ap"`
	MANAGER    string `json:"MANAGER" gorm:"column:MANAGER"`
	Kontak     string `json:"kontak" gorm:"column:kontak"`
	Nama       string `json:"nama" gorm:"column:nama"`
}

func (a *Unit) TableName() string {
	return "public.pln_unit_up"
}

type ExporterRepository interface {
	GetAllUnit() (result []*Regional, err error)
	FindTransaksi(req *request.RekapRequest) ([]*Transaksi, error)
	FindPelanggan(req *request.RekapRequest) ([]*Pelanggan, error)
}
