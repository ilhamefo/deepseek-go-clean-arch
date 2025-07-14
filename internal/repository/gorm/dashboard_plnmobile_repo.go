package gorm

import (
	"errors"
	"event-registration/internal/core/domain"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardPLNMobileRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewDashboardPLNMobileRepo(
	db *gorm.DB, // `name:"authDB"`
	logger *zap.Logger,
) domain.DashboardPLNMobileRepository {
	return &DashboardPLNMobileRepo{db: db, logger: logger}
}

func (m DashboardPLNMobileRepo) SummaryPengguna(filter domain.Filter) (res *domain.SummaryResponse, err error) {

	res = &domain.SummaryResponse{}

	res.TotalPengguna, err = m.GetPengguna(filter)
	if err != nil {
		return res, err
	}

	res.PenggunaAktif, err = m.GetPenggunaAktif(filter)
	if err != nil {
		return res, err
	}

	res.Demografi, err = m.GetDemografi(filter)
	if err != nil {
		return res, err
	}

	res.JumlahIDTerdaftar, err = m.GetTotalIDPelTerdaftar(filter)
	if err != nil {
		return res, err
	}

	res.JumlahIDPelSwacam, err = m.GetTotalIDPelSwacam(filter)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (m DashboardPLNMobileRepo) whereUnit(level string, units []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch level {
		case "unit":
			db = db.Where("unit_up IN ?", units)
		case "area":
			db = db.Where("unit_ap IN ?", units)
		case "induk", "regional":
			db = db.Where("unit_upi IN ?", units)
		case "nasional":
			db = db.Where("unit_upi IS NOT NULL")
		}

		return db
	}
}

func (m DashboardPLNMobileRepo) whereUnitNoUnderscore(level string, units []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch level {
		case "unit":
			db = db.Where("unitup IN ?", units)
		case "area":
			db = db.Where("unitap IN ?", units)
		case "induk", "regional":
			db = db.Where("unitupi IN ?", units)
		case "nasional":
			db = db.Where("unitupi IS NOT NULL")
		}

		return db
	}
}

func (m DashboardPLNMobileRepo) GetPengguna(filter domain.Filter) (res domain.TotalPengguna, err error) {
	bindings := map[string]interface{}{
		"year":    filter.StartDate.Format("2006"),
		"subYear": filter.StartDate.AddDate(-1, 0, 0).Format("2006"),
	}

	// total pengguna (tak ada kolom unit)
	err = m.db.Raw(`SELECT
				current_year :: int as total,
				(current_year / NULLIF(current_year_minus_one, 0)) * 100 :: float as prosentase
				FROM (
				SELECT
					SUM(total_user) FILTER (WHERE DATE_PART('year', tanggal::date) = @year)::float as current_year,
					SUM(total_user) FILTER (WHERE DATE_PART('year', tanggal::date) = @subYear)::float as current_year_minus_one
				FROM "dashboard_plnmobile"."total_user_terdaftar"
				) as source_table
				`, bindings).
		Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, err
}

// total pengguna aktif per tahun
func (m DashboardPLNMobileRepo) GetPenggunaAktif(filter domain.Filter) (res domain.TotalPenggunaAktif, err error) {
	var addSelects string

	bindings := map[string]interface{}{
		"startYear":  time.Date(filter.StartDate.Year(), time.January, 1, 0, 0, 0, 0, filter.StartDate.Location()),
		"endYear":    time.Date(filter.StartDate.Year(), time.December, 31, 23, 59, 59, int(time.Nanosecond*time.Second-time.Nanosecond), filter.StartDate.Location()),
		"startMonth": time.Date(filter.StartDate.Year(), filter.StartDate.Month(), 1, 0, 0, 0, 0, filter.StartDate.Location()),
		"endMonth":   time.Date(filter.StartDate.Year(), filter.StartDate.Month(), 1, 0, 0, 0, 0, filter.StartDate.Location()).AddDate(0, 1, 0).Add(-time.Nanosecond),
		"startDate":  filter.StartDate,
		"endDate":    filter.EndDate,
	}

	switch filter.Type {
	case "harian":
		addSelects += `SUM ( total_user ) FILTER ( WHERE tanggal BETWEEN @startYear AND @endYear ) AS tahunan,
		SUM ( total_user ) FILTER ( WHERE tanggal BETWEEN @startMonth AND @endMonth ) :: INT AS bulanan,
		SUM ( total_user ) FILTER ( WHERE tanggal BETWEEN @startDate AND @endDate ) :: INT AS harian`
	case "bulanan":
		addSelects += `SUM ( total_user ) FILTER ( WHERE tanggal BETWEEN @startYear AND @endYear ) AS tahunan,
		SUM ( total_user ) FILTER ( WHERE tanggal BETWEEN @startMonth AND @endMonth ) :: INT AS bulanan,
		0 AS harian`
	case "tahunan":
		addSelects += `SUM ( total_user ) FILTER ( WHERE tanggal BETWEEN @startYear AND @endYear ) AS tahunan,
		0 AS bulanan,
		0 AS harian`
	default:
		return res, errors.New("invalid type, must be 'harian', 'bulanan', or 'tahunan'")
	}

	var selects string = `SELECT ` + addSelects + ` FROM dashboard_plnmobile.total_user_aktif`

	err = m.db.
		Raw(selects, bindings).
		Scan(&res).
		Error
	if err != nil {
		return res, err
	}

	return res, nil
}

// demografi
func (m DashboardPLNMobileRepo) GetDemografi(filter domain.Filter) (res []*domain.DemografiResponse, err error) {
	query := m.db.Debug().Table("dashboard_plnmobile.demografi_user").
		Where("tanggal_daftar_user BETWEEN ? AND ?", filter.StartDate, filter.EndDate)

	if filter.Export {
		query.Select(`kategori_umur, tanggal_daftar_user :: DATE AS tanggal, provinsi, kabupaten, kecamatan, SUM ( jumlah_user :: INT ) AS jumlah_user`).
			Group("kategori_umur, tanggal, provinsi, kabupaten, kecamatan")
	} else {
		query.Select("sum(jumlah_user :: int) as jumlah_user, kategori_umur").
			Group("kategori_umur")
	}

	err = query.Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

// total_idpel_terdaftar ( ada unit tapi tanggal kosong )
// idpel_ap2t_terdaftar
func (m DashboardPLNMobileRepo) GetTotalIDPelTerdaftar(filter domain.Filter) (res domain.JumlahIDTerdaftar, err error) {
	if filter.Export {
		return m.getTotalIDPelTerdaftarDetail(filter)
	} else {
		return m.getTotalIDPelTerdaftar(filter)
	}
}

func (m DashboardPLNMobileRepo) getTotalIDPelTerdaftarDetail(filter domain.Filter) (res domain.JumlahIDTerdaftar, err error) {
	err = m.db.
		Select("unit_upi, unit_ap, unit_up, SUM( jumlah_idpel_terdaftar ) AS total").
		Table("dashboard_plnmobile.total_idpel_terdaftar").
		Scopes(
			m.whereUnit(filter.Level, filter.Unit),
		).
		Order("total DESC").
		Group("unit_upi, unit_ap, unit_up").
		Scan(&res.Detail).Error

	if err != nil {
		return res, err
	}

	return res, nil
}

func (m DashboardPLNMobileRepo) getTotalIDPelTerdaftar(filter domain.Filter) (res domain.JumlahIDTerdaftar, err error) {
	var totalTerdaftar int64
	var totalIDpel int64

	// total terdaftar
	err = m.db.
		Select("SUM(jumlah_idpel_terdaftar) as total").
		Table("dashboard_plnmobile.total_idpel_terdaftar").
		Scopes(
			m.whereUnit(filter.Level, filter.Unit),
		).
		Scan(&totalTerdaftar).Error
	if err != nil {
		return res, err
	}

	// total idpel
	err = m.db.
		Select("SUM (jmlplg) as total").
		Table("dashboard_plnmobile.idpel_ap2t_terdaftar").
		Scopes(
			m.whereUnitNoUnderscore(filter.Level, filter.Unit),
		).
		Scan(&totalIDpel).Error
	if err != nil {
		return res, err
	}

	res.Total = totalIDpel
	res.TotalUnit = totalTerdaftar
	res.Prosentase = float64(totalTerdaftar) / float64(totalIDpel) * 100

	return res, nil
}

// total_idpel_swacam ( ada unit tapi tanggal kosong )
func (m DashboardPLNMobileRepo) GetTotalIDPelSwacam(filter domain.Filter) (res domain.JumlahIDTerdaftar, err error) {
	if filter.Export {
		return m.getTotalIDPelSwacamDetail(filter)
	} else {
		return m.getTotalIDPelSwacam(filter)
	}
}

func (m DashboardPLNMobileRepo) getTotalIDPelSwacamDetail(filter domain.Filter) (res domain.JumlahIDTerdaftar, err error) {
	err = m.db.Debug().
		Select("unit_upi, unit_ap, unit_up, SUM( total_idpel ) AS total").
		Table("dashboard_plnmobile.total_idpel_swacam").
		Scopes(
			m.whereUnit(filter.Level, filter.Unit),
		).
		Order("total DESC").
		Group("unit_upi, unit_ap, unit_up").
		Scan(&res.Detail).Error

	if err != nil {
		return res, err
	}

	return res, nil
}

func (m DashboardPLNMobileRepo) getTotalIDPelSwacam(filter domain.Filter) (res domain.JumlahIDTerdaftar, err error) {
	var totalTerdaftar int64
	var totalIDpel int64

	// total terdaftar
	err = m.db.
		Select("SUM(total_idpel) as total").
		Table("dashboard_plnmobile.total_idpel_swacam").
		Scopes(
			m.whereUnit(filter.Level, filter.Unit),
		).
		Scan(&totalTerdaftar).Error
	if err != nil {
		return res, err
	}

	// total idpel
	err = m.db.
		Select("SUM (jmlplg) as total").
		Table("dashboard_plnmobile.ap2t_swacam_terdaftar").
		Scopes(
			m.whereUnitNoUnderscore(filter.Level, filter.Unit),
		).
		Scan(&totalIDpel).Error
	if err != nil {
		return res, err
	}

	res.Total = totalIDpel
	res.TotalUnit = totalTerdaftar
	res.Prosentase = float64(totalTerdaftar) / float64(totalIDpel) * 100

	return res, nil
}
