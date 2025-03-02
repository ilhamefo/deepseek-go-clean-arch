package gorm

import (
	"event-registration/internal/core/domain"
	"event-registration/internal/helper"
	"event-registration/internal/request"

	"gorm.io/gorm"
)

type ExporterRepo struct {
	db *gorm.DB
}

func NewExporterRepo(db *gorm.DB) domain.ExporterRepository {
	return &ExporterRepo{db: db}
}

func (r *ExporterRepo) GetAllUnit() (result []*domain.Regional, err error) {
	err = r.db.Model(&domain.Regional{}).
		Preload("Induk.Area.Unit").
		Find(&result).Error

	return result, err
}

func (r *ExporterRepo) FindTransaksi(req *request.RekapRequest) (result []*domain.Transaksi, err error) {

	query := r.db.Select("name, consumer_name, type, amount, status_code, meter_id, title, payment_gateway, created_at, token, up.id_unit_up AS unit_up,	up.nama_unit_up, ap.nama_unit_ap, upi.nama_unit_upi, upi.id_unit_upi as unit_upi, ap.id_unit_ap as unit_ap").
		Joins("JOIN public.pln_unit_upi upi ON transaksi.unit_upi = upi.id_unit_upi :: text").
		Joins("JOIN public.pln_unit_ap ap ON transaksi.unit_ap = ap.id_unit_ap").
		Joins("JOIN public.pln_unit_up up ON transaksi.unit_up = up.id_unit_up")

	if len(req.Induk) > 0 {
		query.Where("unit_upi = ?", req.Induk)
	} else if len(req.Area) > 0 {
		query.Where("unit_ap = ?", req.Area)
	} else if len(req.UnitCode) > 0 {
		query.Where("unit_up = ?", req.UnitCode)
	}

	if len(req.DateStart) > 0 && len(req.DateEnd) > 0 {
		startDate, err := helper.StartDateParser(req.DateStart)
		if err != nil {
			return nil, err
		}

		endDate, err := helper.EndDateParser(req.DateEnd)
		if err != nil {
			return nil, err
		}

		query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err = query.Find(&result).Error

	return result, err
}

func (r *ExporterRepo) FindPelanggan(req *request.RekapRequest) (result []*domain.Pelanggan, err error) {

	query := r.db.Select("id, idpel, name, consumer_name, energy_type, kwh, address, meter_no, meter_type, unit_upi, nama_unit_upi, unit_ap, nama_unit_ap, unit_up, nama_unit_up, created_at")

	if len(req.Induk) > 0 {
		query.Where("unit_upi = ?", req.Induk)
	} else if len(req.Area) > 0 {
		query.Where("unit_ap = ?", req.Area)
	} else if len(req.UnitCode) > 0 {
		query.Where("unit_up = ?", req.UnitCode)
	}

	if len(req.DateStart) > 0 && len(req.DateEnd) > 0 {
		startDate, err := helper.StartDateParser(req.DateStart)
		if err != nil {
			return nil, err
		}

		endDate, err := helper.EndDateParser(req.DateEnd)
		if err != nil {
			return nil, err
		}

		query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err = query.Find(&result).Error

	return result, err
}
