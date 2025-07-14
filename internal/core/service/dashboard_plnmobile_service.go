package service

import (
	"event-registration/internal/common"
	"event-registration/internal/common/constant"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type DashboardPLNMobileService struct {
	repo   domain.DashboardPLNMobileRepository
	logger *zap.Logger
	config *common.Config
}

func NewDashboardPLNMobileService(repo domain.DashboardPLNMobileRepository, logger *zap.Logger, googleConfig *oauth2.Config, config *common.Config, sessionService *SessionService) *DashboardPLNMobileService {
	return &DashboardPLNMobileService{
		repo:   repo,
		logger: logger,
		config: config,
	}
}

func (s *DashboardPLNMobileService) getUnitFilter(req *request.GetFilterRequest) (filter models.Filter, err error) {
	switch req.Level {
	case "ulp":
		filter.Column = "unitup"
		filter.Unit = append(filter.Unit, strings.ToUpper(req.UnitID))
		filter.Level = "unit"
	case "up3":
		filter.Column = "unitap"
		filter.Unit = append(filter.Unit, strings.ToUpper(req.UnitID))
		filter.Level = "area"
	case "uid":
		filter.Column = "unitupi"
		filter.Unit = append(filter.Unit, strings.ToUpper(req.UnitID))
		filter.Level = "induk"
	case "reg":
		filter.Column = "unitupi"
		filter.Level = "regional"

		induks, err := s.IndukService.GetByReg(strings.ToUpper(req.UnitID))
		if err != nil {
			s.logger.Error("Error get induks by reg: ", zap.Error(err))
			return filter, err
		}

		for _, induk := range induks {
			filter.Unit = append(filter.Unit, strings.ToUpper(induk.ID))
		}
	default:
		filter.Column = "unitupi"
		filter.Level = "nasional"

		induks, _, _, err := s.IndukService.GetAll(0, 0, nil)
		if err != nil {
			s.logger.Error("Error get all induks: ", zap.Error(err))
			return filter, err
		}

		for _, induk := range induks {
			filter.Unit = append(filter.Unit, strings.ToUpper(induk.ID))
		}
	}

	tglAwal, tglAkhir, tipe, err := s.getPeriodeFilter(req)
	if err != nil {
		s.logger.Error("Error get periode filter: ", zap.Error(err))
		return filter, err
	}

	filter.StartDate = tglAwal
	filter.EndDate = tglAkhir
	filter.Type = tipe

	filter.Rupiah = req.Rupiah
	filter.Frekuensi = req.Frekuensi
	filter.TipeHpPelanggan = req.TipeHpPelanggan

	filter.ChannelRupiah = req.ChannelRupiah
	filter.ChannelFrekuensi = req.ChannelFrekuensi
	filter.TokenRupiah = req.TokenRupiah
	filter.TokenFrekuensi = req.TokenFrekuensi
	filter.TarifDayaRupiah = req.TarifDayaRupiah
	filter.TarifDayaFrekuensi = req.TarifDayaFrekuensi
	filter.Limit = true

	filter.JenisLayanan = s.getJenisLayanan(req.JenisLayanan)

	return filter, nil
}

func (s *DashboardPLNMobileService) getJenisLayanan(req string) string {
	switch req {
	case "iconnet":
		return "Iconnet"
	case "listriqu":
		return "ListriQu"
	case "pasang_baru":
		return "Pasang Baru"
	case "penerangan_sementara":
		return "Penerangan Sementara"
	case "perubahan_daya":
		return "Perubahan Daya"
	case "pulsa_dan_tagihan":
		return "Pulsa dan Tagihan"
	case "spbklu":
		return "SPBKLU"
	case "spklu":
		return "SPKLU"
	case "splu":
		return "SPLU"
	default:
		return ""
	}
}

func (s *DashboardPLNMobileService) getPeriodeFilter(req *request.GetFilterRequest) (tglAwal, tglAkhir *time.Time, tipe string, err error) {

	if len(req.TanggalAwal) > 0 && len(req.TanggalAkhir) > 0 {
		tanggalAwal, err := time.Parse("2006/01/02", req.TanggalAwal)
		if err != nil {
			s.logger.Error(constant.ErrFormatDate, zap.Error(err))
			return nil, nil, tipe, err
		}

		tanggalAkhir, err := time.Parse("2006/01/02", req.TanggalAkhir)
		if err != nil {
			s.logger.Error(constant.ErrFormatDate, zap.Error(err))
			return nil, nil, tipe, err
		}

		return &tanggalAwal, &tanggalAkhir, "harian", nil
	} else if len(req.Bulan) > 0 && len(req.Tahun) > 0 {
		tanggalAwal, err := time.Parse("200601", req.Tahun+req.Bulan)
		if err != nil {
			s.logger.Error(constant.ErrFormatDate, zap.Error(err))
			return nil, nil, tipe, err
		}

		tanggalAkhir := tanggalAwal.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)

		return &tanggalAwal, &tanggalAkhir, "bulanan", nil
	} else {
		tanggalAwal, err := time.Parse("2006", req.Tahun)
		if err != nil {
			s.logger.Error(constant.ErrFormatDate, zap.Error(err))
			return nil, nil, tipe, err
		}

		tanggalAkhir := tanggalAwal.AddDate(1, 0, 0).Add(-1 * time.Nanosecond)

		return &tanggalAwal, &tanggalAkhir, "tahunan", nil
	}

}

func (s *DashboardPLNMobileService) SummaryPengguna(req *request.DashboardPLNMobileRequest) (res *domain.SummaryResponse, lastUpdated *time.Time, err error) {
	// Get unit filter
	filter, err := s.getUnitFilter(&request.GetFilterRequest{
		Level:        req.Level,
		UnitID:       req.UnitID,
		TanggalAwal:  req.TanggalAwal,
		TanggalAkhir: req.TanggalAkhir,
		Bulan:        req.Bulan,
		Tahun:        req.Tahun,
	})

	if err != nil {
		s.logger.Error(constant.ErrUnitFilter, zap.Error(err))
		return res, nil, err
	}

	res, err = s.repo.SummaryPengguna(filter)
	if err != nil {
		s.logger.Error("Error get summary pengguna : ", zap.Error(err))
		return res, nil, err
	}

	lastUpdated, err = s.repo.GetLastUpdate(`dashboard_plnmobile.total_user_terdaftar`)
	if err != nil {
		s.logger.Error(constant.ErrGetLastUpdated, zap.Error(err))
		return res, nil, err
	}

	return res, lastUpdated, err
}
