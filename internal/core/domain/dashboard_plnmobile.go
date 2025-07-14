package domain

import "time"

type DashboardPLNMobileRepository interface {
	SummaryPengguna(filter Filter) (res *SummaryResponse, err error)

	// Pengguna
	// GetPengguna(filter Filter) (res TotalPengguna, err error)           // no export
	// GetPenggunaAktif(filter Filter) (res TotalPenggunaAktif, err error) // no export
	// GetDemografi(filter Filter) (res []*DemografiResponse, err error)
	// GetTotalIDPelTerdaftar(filter Filter) (res JumlahIDTerdaftar, err error)
	// GetTotalIDPelSwacam(filter Filter) (res JumlahIDTerdaftar, err error)

	GetLastUpdate(table string) (lastUpdate *time.Time, err error)

	// // Top ten page 1
	// GetTransaksiRupiah(filter Filter) (res []*BarResponse, err error)
	// GetTransaksiFrekuensi(filter Filter) (res []*BarResponse, err error)
	// GetTipeHp(filter Filter) (res []*BarResponse, err error)
	// GetOS(filter Filter) (res []*BarResponse, err error)

	// // Top ten page 2
	// GetChannelPaymentRupiah(filter Filter) (res []*BarResponse, err error)
	// GetChannelPaymentFrekuensi(filter Filter) (res []*BarResponse, err error)
	// GetTokenPaymentRupiah(filter Filter) (res []*BarResponse, err error)
	// GetTokenPaymentFrekuensi(filter Filter) (res []*BarResponse, err error)
	// GetTarifDayaRupiah(filter Filter) (res []*BarResponse, err error)
	// GetTarifDayaFrekuensi(filter Filter) (res []*BarResponse, err error)

	// // TREN
	// ConversionRate(filter Filter) (res *ConversionRate, err error)
	// ConversionRateExport(filter Filter) (res []*ConversionRateDetail, err error)
	// GetTrenPenambahanPengguna(filter Filter) (res []*LineResponse, err error)
	// GetTrenTransaksiRupiah(filter Filter) (res []*LineResponse, err error)
	// GetTrenTransaksiFrekuensi(filter Filter) (res []*LineResponse, err error)
	// GetKeluhanError(filter Filter) (res []*LineResponse, err error)
}

type Filter struct {
	Column    string
	Level     string
	Unit      []string
	StartDate *time.Time
	EndDate   *time.Time
	Type      string
	Limit     bool
	Export    bool

	// conversion rate only
	JenisLayanan string

	// sorting
	Rupiah             string
	Frekuensi          string
	TipeHpPelanggan    string
	ChannelRupiah      string
	ChannelFrekuensi   string
	TokenRupiah        string
	TokenFrekuensi     string
	TarifDayaRupiah    string
	TarifDayaFrekuensi string
}

type SummaryResponse struct {
	TotalPengguna     TotalPengguna        `json:"total_pengguna"`
	PenggunaAktif     TotalPenggunaAktif   `json:"total_pengguna_aktif"`
	Demografi         []*DemografiResponse `json:"demografi"`
	JumlahIDTerdaftar JumlahIDTerdaftar    `json:"jumlah_id_terdaftar"`
	JumlahIDPelSwacam JumlahIDTerdaftar    `json:"jumlah_id_swacam"`
}

type SummaryPageOneResponse struct {
	TotalPengguna TotalPengguna `json:"total_pengguna"`
}

type SummaryTopTenPageOneResponse struct {
	TransaksiRupiah    []*BarResponse `json:"transaksi_rupiah"`
	TransaksiFrekuensi []*BarResponse `json:"transaksi_frekuensi"`
	TipeHP             []*BarResponse `json:"tipe_hp"`
	OS                 []*BarResponse `json:"os"`
}

type SummaryTopTenPageTwoResponse struct {
	ChannelPaymentRupiah         []*BarResponse `json:"channel_payment_rupiah"`
	ChannelPaymentFrekuensi      []*BarResponse `json:"channel_payment_frekuensi"`
	TokenChannelPaymentRupiah    []*BarResponse `json:"token_channel_payment_rupiah"`
	TokenChannelPaymentFrekuensi []*BarResponse `json:"token_channel_payment_frekuensi"`
	TarifDayaRupiah              []*BarResponse `json:"tarif_daya_rupiah"`
	TarifDayaFrekuensi           []*BarResponse `json:"tarif_daya_frekuensi"`
}

type SummaryTrenResponse struct {
	PenambahanPengguna []*LineResponse `json:"penambahan_pengguna"`
	TransaksiRupiah    []*LineResponse `json:"transaksi_rupiah"`
	TransaksiFrekuensi []*LineResponse `json:"transaksi_frekuensi"`
	KeluhanError       []*LineResponse `json:"keluhan_error"`
}

type LineResponse struct {
	Jumlah  float64    `gorm:"column:jumlah" json:"jumlah"`
	Tanggal *time.Time `gorm:"column:tanggal" json:"tanggal"`
	Bulan   string     `gorm:"column:bulan" json:"bulan"`
}
type BarResponse struct {
	Title  string  `gorm:"column:title" json:"title"`
	Jumlah float64 `gorm:"column:jumlah" json:"jumlah"`

	// unit for export
	UnitUP  string     `gorm:"column:unit_up" json:"unit_up,omitempty"`
	UnitAP  string     `gorm:"column:unit_ap" json:"unit_ap,omitempty"`
	UnitUPI string     `gorm:"column:unit_upi" json:"unit_upi,omitempty"`
	Tanggal *time.Time `gorm:"column:tanggal" json:"tanggal,omitempty"`
}

type DemografiResponse struct {
	Total        int64  `gorm:"column:jumlah_user" json:"jumlah"`
	KategoriUmur string `gorm:"column:kategori_umur" json:"title"`

	// export fields
	Tanggal   *time.Time `gorm:"column:tanggal" json:"tanggal,omitempty"`
	Provinsi  string     `gorm:"column:provinsi" json:"provinsi,omitempty"`
	Kabupaten string     `gorm:"column:kabupaten" json:"kabupaten,omitempty"`
	Kecamatan string     `gorm:"column:kecamatan" json:"kecamatan,omitempty"`
}

type TotalPengguna struct {
	TotalPengguna int64   `gorm:"column:total" json:"total"`
	Prosentase    float32 `gorm:"column:prosentase" json:"prosentase"`
}
type TotalPenggunaAktif struct {
	Tahunan int64 `gorm:"column:tahunan" json:"tahunan"`
	Bulanan int64 `gorm:"column:bulanan" json:"bulanan"`
	Harian  int64 `gorm:"column:harian" json:"harian"`
}

type ConversionRate struct {
	ConversionRate float64 `gorm:"column:conversion_rate" json:"conversion_rate"`
	Rupiah         float64 `gorm:"column:rupiah" json:"rupiah"`
	FrekuensiAjuan int64   `gorm:"column:frekuensi_ajuan" json:"frekuensi_ajuan"`
	Rating         float32 `gorm:"column:rating_internal" json:"rating_internal"`
}

type ConversionRateDetail struct {
	Tanggal         *time.Time `gorm:"column:tanggal" json:"tanggal,omitempty"`
	UnitUPI         string     `gorm:"column:unit_upi" json:"unit_upi,omitempty"`
	UnitAP          string     `gorm:"column:unit_ap" json:"unit_ap,omitempty"`
	UnitUP          string     `gorm:"column:unit_up" json:"unit_up,omitempty"`
	Transaksi       string     `gorm:"column:transaksi" json:"transaksi,omitempty"`
	JumlahTransaksi int64      `gorm:"column:jumlah_transaksi" json:"jumlah_transaksi,omitempty"`
	JumlahRupiah    float64    `gorm:"column:jumlah_rupiah" json:"jumlah_rupiah,omitempty"`
	Status          string     `gorm:"column:status" json:"status,omitempty"`
}

type JumlahIDTerdaftar struct {
	Total      int64             `gorm:"column:total" json:"total"`
	TotalUnit  int64             `gorm:"column:total_unit" json:"total_unit"`
	Prosentase float64           `gorm:"column:prosentase" json:"prosentase"`
	Detail     []*PenggunaDetail `json:"detail,omitempty"`
}

type PenggunaDetail struct {
	UnitUP  string `gorm:"column:unit_up" json:"unit_up,omitempty"`
	UnitAP  string `gorm:"column:unit_ap" json:"unit_ap,omitempty"`
	UnitUPI string `gorm:"column:unit_upi" json:"unit_upi,omitempty"`
	Jumlah  string `gorm:"column:total" json:"total,omitempty"`
}
