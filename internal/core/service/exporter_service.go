package service

import (
	"archive/tar"
	"compress/gzip"
	"event-registration/internal/common/helper"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

const (
	BATCH_SIZE = 50000
	filesDir   = "files/"
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

func (s *ExporterService) ExportRekapTransaksi(req *request.RekapRequest) (err error) {

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

	return s.compressFiles(files, filesDir+filename+".tar.gz")
}

func (s *ExporterService) ExportAllRekapTransaksi(req *request.RekapRequest) (err error) {
	var payload []Payload
	units, err := s.repo.GetAllUnit()
	if err != nil {
		s.logger.Error(
			"error_get_all_units",
			zap.Error(err),
		)
		return err
	}

	for _, unit := range units {
		if len(unit.Induk) > 0 {
			for _, induk := range unit.Induk {

				filename := helper.NormalizeString(induk.Satuan + " " + induk.NamaUnitUPI)

				payload = append(payload, Payload{
					filename: filename,
					req: &request.RekapRequest{
						UnitCode:  "",
						Area:      "",
						Induk:     induk.IDUnitUPI,
						Pusat:     "",
						DateStart: req.DateStart,
						DateEnd:   req.DateEnd,
					},
				})

				if len(induk.Area) > 0 {
					for _, area := range induk.Area {
						filename := helper.NormalizeString(area.Satuan + " " + area.NamaUnitAP)

						payload = append(payload, Payload{
							filename: filename,

							req: &request.RekapRequest{
								UnitCode:  "",
								Area:      area.IDUnitAP,
								Induk:     "",
								Pusat:     "",
								DateStart: req.DateStart,
								DateEnd:   req.DateEnd,
							},
						})

						if len(area.Unit) > 0 {
							for _, unit := range area.Unit {
								filename := helper.NormalizeString(unit.Satuan + " " + unit.NamaUnitUP)

								payload = append(payload, Payload{
									filename: filename,

									req: &request.RekapRequest{
										UnitCode:  unit.IDUnitUP,
										Area:      "",
										Induk:     "",
										Pusat:     "",
										DateStart: req.DateStart,
										DateEnd:   req.DateEnd,
									},
								})
							}
						}
					}
				}
			}
		}
	}

	s.logger.Info(
		"payload_info",
		zap.Any("payload", payload),
	)

	s.ProcessIndukDataWithWorkerPool(payload)

	s.logger.Info(
		"done_export",
	)

	return nil
}

type Payload struct {
	filename string
	req      *request.RekapRequest
}

func (s *ExporterService) ProcessIndukDataWithWorkerPool(data []Payload) {
	workerCount := 10
	jobs := make(chan Payload, len(data))
	errorsChan := make(chan error, len(data)) // Channel for collecting errors

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for d := range jobs {
				s.logger.Info(
					"worker_processing",
					zap.Int("worker_id", workerID),
					zap.String("filename", d.filename),
				)
				if err := s.process(d); err != nil {
					errorsChan <- fmt.Errorf("worker %d: %w", workerID, err) // Add worker ID to error
				}
			}
		}(i + 1)
	}

	for _, d := range data {
		jobs <- d
	}
	close(jobs)

	wg.Wait()
	close(errorsChan)

	var allErrors []error
	for err := range errorsChan {
		allErrors = append(allErrors, err)
	}

	if len(allErrors) > 0 {
		fmt.Println("\nErrors encountered:")
		for _, err := range allErrors {
			fmt.Println(err)
			s.logger.Error(
				"error_from_channel",
				zap.Error(err),
			)
		}
	}
}

func (s *ExporterService) process(data Payload) error {
	res, err := s.repo.FindTransaksi(data.req)
	if err != nil {
		s.logger.Error(
			"error_find_transaksi",
			zap.Error(err),
			zap.String("filename", data.filename),
		)

		return err
	}

	_, err = s.generateXlsx(res, data.filename)
	if err != nil {
		s.logger.Error(
			"error_find_transaksi",
			zap.Error(err),
			zap.String("filename", data.filename),
		)

		return err
	}

	return nil
}

func (s *ExporterService) ExportRekapPelanggan(req *request.RekapRequest) (err error) {

	var filename string
	var tanggal string = strings.ReplaceAll(req.DateStart+"_"+req.DateEnd, "/", "")

	s.logger.Info(
		"starting_query",
	)

	res, err := s.repo.FindPelanggan(req)
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
		filename = "UNIT_" + req.UnitCode + "_" + tanggal
	} else {
		filename = "NASIONAL" + "_" + tanggal
	}

	files, err := s.generateXlsxPelanggan(res, filename)
	if err != nil {
		s.logger.Error(
			"error_generate_excel",
			zap.Error(err),
		)
		return err
	}

	return s.compressFiles(files, filesDir+filename+".tar.gz")
}

func (s *ExporterService) generateXlsx(res []*domain.Transaksi, filename string) (files []string, err error) {
	sheetName := "Rekap Transaksi"
	batchSize := 25_000 * 6
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

			var paymentType string = ""
			if data.Type == "" {
				paymentType = "-"
			} else {
				paymentType = data.Type
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

		path := fmt.Sprintf("%s%s", filesDir, filename)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			errMakeDir := os.Mkdir(path, 0o750) // Secure permissions: owner=rwx, group=rx
			if errMakeDir != nil {
				s.logger.Error(
					"error_make_dir",
					zap.Error(errMakeDir),
				)
				return nil, errMakeDir
			}
		}

		// Save the file with a batch-specific name
		batchFilename := fmt.Sprintf("files/%s/REKAP_TRANSAKSI_EXPORT_%s_PART_%d.xlsx", filename, filename, batch+1)
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

var pelangganHeaders = []string{
	"No.", "ID PELANGGAN", "NAMA", "CONSUMER NAME", "TIPE ENERGI", "KWH", "ALAMAT", "METER NO", "TIPE METER", "UNIT UPI", "NAMA UNIT UPI", "UNIT AP", "NAMA UNIT AP", "UNIT UP", "NAMA UNIT UP", "CREATED AT",
}

func (s *ExporterService) generateXlsxPelanggan(res []*domain.Pelanggan, filename string) (files []string, err error) {
	sheetName := "Rekap Pelanggan"
	batchSize := 25_000 * 6
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

		batchFilename, err := s.createPelangganBatchFile(res[start:end], filename, batch, sheetName)
		if err != nil {
			return files, err
		}
		files = append(files, batchFilename)
	}

	return files, nil
}

func (s *ExporterService) createPelangganBatchFile(batch []*domain.Pelanggan, filename string, batchNum int, sheetName string) (string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			s.logger.Error("error_create_new_file", zap.Error(err))
		}
	}()

	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		s.logger.Error("create_sheet", zap.Error(err))
		return "", err
	}

	sw, err := f.NewStreamWriter(sheetName)
	if err != nil {
		s.logger.Error("error_create_stream_writer", zap.Error(err))
		return "", err
	}

	if err := sw.SetColWidth(2, 15, 26); err != nil {
		s.logger.Error("error_set_col", zap.Error(err))
		return "", err
	}

	if err := s.setHeaders(sw, f, pelangganHeaders); err != nil {
		s.logger.Error("error_set_error", zap.Error(err))
		return "", err
	}

	s.logger.Info("length_of_batch", zap.Int("rows", len(batch)))

	if err := s.writePelangganRows(sw, batch); err != nil {
		return "", err
	}

	if err := sw.Flush(); err != nil {
		s.logger.Error("error_flush_file", zap.Error(err))
		return "", err
	}

	batchFilename := fmt.Sprintf("files/REKAP_PELANGGAN_EXPORT_%s_PART_%d.xlsx", filename, batchNum+1)
	if err := f.SaveAs(batchFilename); err != nil {
		s.logger.Error("error_save_excel_file", zap.Error(err))
		return "", err
	}

	s.logger.Info("batch_saved", zap.Int("batch", batchNum+1), zap.String("to", batchFilename))
	return batchFilename, nil
}

func (s *ExporterService) writePelangganRows(sw *excelize.StreamWriter, batch []*domain.Pelanggan) error {
	for i, data := range batch {
		rowIndex := i + 1
		cell, err := excelize.CoordinatesToCellName(1, rowIndex+1)
		if err != nil {
			s.logger.Error("error_create_coordinate", zap.Error(err))
			return err
		}
		if err := sw.SetRow(cell, []interface{}{
			rowIndex,
			data.IDPel,
			data.Name,
			data.ConsumerName,
			data.EnergyType,
			data.KWH,
			data.Address,
			data.MeterNo,
			data.MeterType,
			data.UnitUpi,
			data.NamaUnitUpi,
			data.UnitAp,
			data.NamaUnitAp,
			data.UnitUp,
			data.NamaUnitUp,
			data.CreatedAt,
		}); err != nil {
			s.logger.Error("error_set_row", zap.Error(err))
			return fmt.Errorf("error set row : %s", err.Error())
		}
	}
	return nil
}

func (s *ExporterService) compressFiles(files []string, outputFile string) error {
	s.logger.Info(
		"compress_started",
		zap.Int("file_count", len(files)),
		zap.Strings("files", files),
	)

	// Security: Prevent path traversal and absolute paths
	if strings.Contains(outputFile, "..") || strings.HasPrefix(outputFile, "/") || strings.HasPrefix(outputFile, "\\") {
		s.logger.Error("invalid_output_file_path", zap.String("outputFile", outputFile))
		return fmt.Errorf("invalid output file path: %s", outputFile)
	}

	// Only allow output in the 'files/' directory
	if !strings.HasPrefix(outputFile, "files/") && !strings.HasPrefix(outputFile, "files\\") {
		s.logger.Error("output_file_must_be_in_files_dir", zap.String("outputFile", outputFile))
		return fmt.Errorf("output file must be in 'files/' directory: %s", outputFile)
	}

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
