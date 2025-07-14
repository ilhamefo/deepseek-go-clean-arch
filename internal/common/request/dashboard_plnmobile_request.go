package request

type DashboardPLNMobileRequest struct {
	Level        string `json:"level" form:"level" binding:"omitempty,oneof=nas reg uid up3 ulp,max=100" example:"up3"`
	UnitID       string `json:"unit_id" form:"unit_id" binding:"max=100" example:"52000"`
	TanggalAwal  string `form:"tanggal_awal" json:"tanggal_awal" binding:"omitempty,datetime=2006/01/02,max=100"`
	TanggalAkhir string `form:"tanggal_akhir" json:"tanggal_akhir" binding:"omitempty,datetime=2006/01/02,max=100"`
	Bulan        string `form:"bulan" json:"bulan" binding:"omitempty,datetime=01,max=100"`
	Tahun        string `form:"tahun" json:"tahun" binding:"required,datetime=2006,max=100"`
	Export       string `form:"export" json:"export" binding:"omitempty,max=100,oneof=transaksi_rupiah transaksi_frekuensi tipe_hp tipe_os channel_rupiah channel_frekuensi token_rupiah token_frekuensi tarif_daya_rupiah tarif_daya_frekuensi tren_penambahan_pengguna tren_transaksi_rupiah tren_transaksi_frekuensi tren_keluhan_error demografi_user jumlah_id_terdaftar jumlah_idpel_swacam convertion_rate" example:"transaksi_rupiah"`
	JenisLayanan string `json:"jenis_layanan" form:"jenis_layanan" binding:"omitempty,oneof=iconnet listriqu pasang_baru penerangan_sementara perubahan_daya pulsa_dan_tagihan spbklu spklu splu,max=100" example:"up3"`
}

type DashboardPLNMobileToptenRequest struct {
	Level              string `json:"level" form:"level" binding:"omitempty,oneof=nas reg uid up3 ulp,max=100" example:"up3"`
	UnitID             string `json:"unit_id" form:"unit_id" binding:"max=100" example:"52000"`
	TanggalAwal        string `form:"tanggal_awal" json:"tanggal_awal" binding:"omitempty,datetime=2006/01/02,max=100"`
	TanggalAkhir       string `form:"tanggal_akhir" json:"tanggal_akhir" binding:"omitempty,datetime=2006/01/02,max=100"`
	Bulan              string `form:"bulan" json:"bulan" binding:"omitempty,datetime=01,max=100"`
	Tahun              string `form:"tahun" json:"tahun" binding:"required,datetime=2006,max=100"`
	Rupiah             string `json:"rupiah" form:"rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	Frekuensi          string `json:"frekuensi" form:"frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TipeHpPelanggan    string `json:"tipe_hp" form:"tipe_hp" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	ChannelRupiah      string `json:"channel_rupiah" form:"channel_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	ChannelFrekuensi   string `json:"channel_frekuensi" form:"channel_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TokenRupiah        string `json:"token_rupiah" form:"token_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TokenFrekuensi     string `json:"token_frekuensi" form:"token_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TarifDayaRupiah    string `json:"tarif_rupiah" form:"tarif_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TarifDayaFrekuensi string `json:"tarif_frekuensi" form:"tarif_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
}

type DashboardPLNMobileToptenPage2Request struct {
	Level              string `json:"level" form:"level" binding:"omitempty,oneof=nas reg uid up3 ulp,max=100" example:"up3"`
	UnitID             string `json:"unit_id" form:"unit_id" binding:"max=100" example:"52000"`
	TanggalAwal        string `form:"tanggal_awal" json:"tanggal_awal" binding:"omitempty,datetime=2006/01/02,max=100"`
	TanggalAkhir       string `form:"tanggal_akhir" json:"tanggal_akhir" binding:"omitempty,datetime=2006/01/02,max=100"`
	Bulan              string `form:"bulan" json:"bulan" binding:"omitempty,datetime=01,max=100"`
	Tahun              string `form:"tahun" json:"tahun" binding:"required,datetime=2006,max=100"`
	ChannelRupiah      string `json:"channel_rupiah" form:"channel_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	ChannelFrekuensi   string `json:"channel_frekuensi" form:"channel_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TokenRupiah        string `json:"token_rupiah" form:"token_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TokenFrekuensi     string `json:"token_frekuensi" form:"token_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TarifDayaRupiah    string `json:"tarif_rupiah" form:"tarif_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TarifDayaFrekuensi string `json:"tarif_frekuensi" form:"tarif_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
}

type GetFilterRequest struct {
	Level        string `json:"level" form:"level" binding:"max=100"`
	UnitID       string `json:"unit_id" form:"unit_id" binding:"max=100"`
	TanggalAwal  string `form:"tanggal_awal" json:"tanggal_awal" binding:"datetime=2006/01/02,max=100"`
	TanggalAkhir string `form:"tanggal_akhir" json:"tanggal_akhir" binding:"datetime=2006/01/02,max=100"`
	Bulan        string `form:"bulan" json:"bulan" binding:"max=100"`
	Tahun        string `form:"tahun" json:"tahun" binding:"max=100"`
	JenisLayanan string `json:"jenis_layanan" form:"jenis_layanan" binding:"omitempty,oneof=iconnet listriqu pasang_baru penerangan_sementara perubahan_daya pulsa_dan_tagihan spbklu spklu splu,max=100" example:"up3"`

	Rupiah             string `json:"rupiah" form:"rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	Frekuensi          string `json:"frekuensi" form:"frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TipeHpPelanggan    string `json:"tipe_hp" form:"tipe_hp" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	ChannelRupiah      string `json:"channel_rupiah" form:"channel_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	ChannelFrekuensi   string `json:"channel_frekuensi" form:"channel_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TokenRupiah        string `json:"token_rupiah" form:"token_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TokenFrekuensi     string `json:"token_frekuensi" form:"token_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TarifDayaRupiah    string `json:"tarif_rupiah" form:"tarif_rupiah" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	TarifDayaFrekuensi string `json:"tarif_frekuensi" form:"tarif_frekuensi" binding:"omitempty,oneof=asc desc,max=100" example:"asc"`
	Export             bool   `json:"export" form:"export" binding:"omitempty"`
}
