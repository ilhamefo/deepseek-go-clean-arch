package service

import (
	"archive/tar"
	"compress/gzip"
	"event-registration/internal/core/domain"
	"event-registration/internal/request"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

const (
	BATCH_SIZE = 50000
)

var HeaderStyle excelize.Style = excelize.Style{
	Font: &excelize.Font{
		Size:  12,
		Color: "000000",
		Bold:  true,
	},
	Border: []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
	},
	Fill:      excelize.Fill{Type: "pattern", Color: []string{"e8f2a1"}, Pattern: 1},
	Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
}

func (s *ExporterService) setHeaders(sw *excelize.StreamWriter, f *excelize.File, headers []string) (err error) {
	var cells []interface{}

	headerStyle, err := f.NewStyle(&HeaderStyle)
	if err != nil {
		return err
	}

	for _, header := range headers {
		cells = append(cells, excelize.Cell{StyleID: headerStyle, Value: header})
	}

	return sw.SetRow("A1", cells, excelize.RowOpts{Height: 30, Hidden: false})
}

type ExporterService struct {
	repo   domain.ExporterRepository
	cache  domain.EventCache
	logger *zap.Logger
}

func NewExporterService(repo domain.ExporterRepository, cache domain.EventCache, logger *zap.Logger) *ExporterService {
	return &ExporterService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *ExporterService) ExportRekapTransaksi(req *request.TransaksiRequest) (err error) {

	var filename string
	var tanggal string = strings.ReplaceAll(req.DateStart+"_"+req.DateEnd, "/", "")

	s.logger.Info(
		"starting_query",
	)

	res, err := s.repo.FindTransaksi(req)
	if err != nil {
		return err
	}

	s.logger.Info(
		"done_get_data",
	)

	if len(req.Induk) > 0 {
		filename = "INDUK_" + req.Induk + "_" + tanggal
	} else if len(req.Area) > 0 {
		filename = "AREA_" + req.Area + "_" + tanggal
	} else if len(req.UnitCode) > 0 {
		filename = "UNIT_" + req.Area + "_" + tanggal
	} else {
		filename = "NASIONAL" + "_" + tanggal
	}

	files, err := s.generateXlsx(res, filename)
	if err != nil {
		s.logger.Error(
			"error_generate_excel",
			zap.Error(err),
		)
		return err
	}

	return s.compressFiles(files, "files/"+filename+".tar.gz")
}

func (s *ExporterService) generateXlsx(res []*domain.Transaksi, filename string) (files []string, err error) {
	sheetName := "Rekap Transaksi"
	batchSize := 25_000
	totalRows := len(res)

	s.logger.Info(
		"length_of_query_result",
		zap.Int("total_rows", totalRows),
	)

	numBatches := (totalRows + batchSize - 1) / batchSize

	for batch := 0; batch < numBatches; batch++ {
		start := batch * batchSize
		end := start + batchSize
		if end > totalRows {
			end = totalRows
		}

		// Create a new Excel file for each batch
		f := excelize.NewFile()

		defer func() {
			if err := f.Close(); err != nil {
				s.logger.Error(
					"error_create_new_file",
					zap.Error(err),
				)
			}
		}()

		err = f.SetSheetName("Sheet1", sheetName)
		if err != nil {
			s.logger.Error(
				"create_sheet",
				zap.Error(err),
			)
			return files, err
		}

		sw, err := f.NewStreamWriter(sheetName)
		if err != nil {
			s.logger.Error(
				"error_create_stream_writer",
				zap.Error(err),
			)
			return files, err
		}

		err = sw.SetColWidth(2, 15, 25)
		if err != nil {
			s.logger.Error(
				"error_set_col",
				zap.Error(err),
			)
			return files, err
		}

		err = s.setHeaders(sw, f, []string{"No.", "Nama Akun", "Nama Pelanggan", "Type", "Amount", "Status Code", "ID Pel", "Pembayaran", "Kanal Pembayaran", "Jenis Pembayaran", "Tanggal Transaksi", "Token", "Unit UPI", "Unit AP", "Unit UP"})
		if err != nil {
			s.logger.Error(
				"error_set_error",
				zap.Error(err),
			)
			return files, err
		}

		s.logger.Info(
			"length_of_batch",
			zap.Int("rows", end-start),
		)

		// Write rows for the current batch
		for i, data := range res[start:end] {
			rowIndex := i + 1

			cell, err := excelize.CoordinatesToCellName(1, rowIndex+1)
			if err != nil {
				s.logger.Error(
					"error_create_coordinate",
					zap.Error(err),
				)
				return files, err
			}

			var paymentType string = "Non Taglist"
			if data.Type != "nontaglis" {
				paymentType = "Taglist"
			}

			if err := sw.SetRow(cell, []interface{}{
				rowIndex,
				data.ConsumerName,
				data.Name,
				data.Type,
				data.Amount,
				data.StatusCode,
				data.MeterID,
				data.Title,
				data.PaymentGateway,
				paymentType,
				data.CreatedAt,
				data.Token,
				data.NameUnitUpi,
				data.NameUnitAP,
				data.NameUnitUP,
			}); err != nil {
				s.logger.Error(
					"error_set_row",
					zap.Error(err),
				)
				return files, fmt.Errorf("error set row : %s", err.Error())
			}
		}

		if err = sw.Flush(); err != nil {
			s.logger.Error(
				"error_flush_file",
				zap.Error(err),
			)
			return files, err
		}

		// Save the file with a batch-specific name
		batchFilename := fmt.Sprintf("files/REKAP_TRANSAKSI_EXPORT_%s_PART_%d.xlsx", filename, batch+1)
		if err := f.SaveAs(batchFilename); err != nil {
			s.logger.Error(
				"error_save_excel_file",
				zap.Error(err),
			)
			return files, err
		}

		s.logger.Info(
			"batch_saved",
			zap.Int("batch", batch+1),
			zap.String("to", batchFilename),
		)

		files = append(files, batchFilename)
	}

	return files, nil
}

func (s *ExporterService) compressFiles(files []string, outputFile string) error {
	s.logger.Info(
		"compress_started",
		zap.Int("file_count", len(files)),
		zap.Strings("files", files),
	)

	outFile, err := os.Create(outputFile)
	if err != nil {
		s.logger.Error(
			"error_create_output_file",
			zap.Error(err),
		)
		return err
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	gzipWriter, err = gzip.NewWriterLevel(outFile, gzip.BestCompression)
	if err != nil {
		s.logger.Error(
			"error_set_compression_level",
			zap.Error(err),
		)
		return err
	}

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, file := range files {
		// Open the input file
		inFile, err := os.Open(file)
		if err != nil {
			s.logger.Error(
				"error_open_file",
				zap.Error(err),
			)
			return err
		}
		defer inFile.Close()

		fileInfo, err := inFile.Stat()
		if err != nil {
			s.logger.Error(
				"error_get_file_info",
				zap.Error(err),
			)
			return err
		}

		header := &tar.Header{
			Name:    file,                   // File name
			Mode:    int64(fileInfo.Mode()), // File mode
			Size:    fileInfo.Size(),        // File size
			ModTime: fileInfo.ModTime(),     // Modification time
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			s.logger.Error(
				"error_writing_header",
				zap.Error(err),
			)
			return err
		}

		if _, err := io.Copy(tarWriter, inFile); err != nil {
			s.logger.Error(
				"error_copy_file",
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}
