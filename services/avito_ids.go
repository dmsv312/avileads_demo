package services

import (
	"avileads-web/models"
	"avileads-web/utils"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego/orm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type avitoRow struct {
	AvitoID   string  `orm:"column(avito_id)"`
	StartDate *string `orm:"column(start_date)"`
}

const batchSize = 500

func ExtractAvitoIds(file multipart.File, hdr *multipart.FileHeader) ([]string, error) {
	tmp := filepath.Join(os.TempDir(), hdr.Filename)
	out, err := os.Create(tmp)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(out, file); err != nil {
		return nil, err
	}
	out.Close()

	xl, err := excelize.OpenFile(tmp)
	if err != nil {
		return nil, err
	}

	sheet := xl.GetSheetName(1)
	rows := xl.GetRows(sheet)

	var raw []string
	for _, r := range rows {
		if len(r) == 0 {
			break
		}
		id := strings.TrimSpace(r[0])
		if id == "" {
			break
		}
		if strings.ContainsAny(id, "Ee") {
			if f, err := strconv.ParseFloat(id, 64); err == nil {
				id = strconv.FormatInt(int64(f), 10)
			}
		}
		raw = append(raw, id)
	}

	return utils.CleanAvitoIds(strings.Join(raw, ",")), nil
}

func ProcessAvitoIds(vacancyId int, avitoIds []string, userId int, db orm.Ormer) error {
	if len(avitoIds) == 0 {
		return nil
	}

	if err := assertAvitoIdsFree(db, vacancyId, avitoIds); err != nil {
		return err
	}

	var existingRows []models.Avito_vacancy_ids
	if _, err := db.QueryTable("avito_vacancy_ids").
		Filter("vacancy_avito_id", vacancyId).
		All(&existingRows, "Avito_id"); err != nil {
		return err
	}

	existing := make(map[string]struct{}, len(existingRows))
	for _, r := range existingRows {
		existing[r.Avito_id] = struct{}{}
	}

	var toInsert []models.Avito_vacancy_ids
	for _, id := range avitoIds {
		clean := strings.ReplaceAll(id, " ", "")
		if _, ok := existing[clean]; ok {
			continue
		}
		toInsert = append(toInsert, models.Avito_vacancy_ids{
			Vacancy_avito_id: &models.Vacancy_avito{Id: vacancyId},
			Avito_id:         clean,
		})
		existing[clean] = struct{}{}
	}

	for start := 0; start < len(toInsert); start += batchSize {
		end := start + batchSize
		if end > len(toInsert) {
			end = len(toInsert)
		}
		batch := toInsert[start:end]
		if _, err := db.InsertMulti(len(batch), batch); err != nil {
			return err
		}
	}

	unionIds := make([]string, 0, len(existing))
	for id := range existing {
		unionIds = append(unionIds, id)
	}

	return models.UpdateAvitoVacancyIdsHistByAvitoIdsInVacancyUpdate(
		vacancyId, unionIds, userId, db)
}

func assertAvitoIdsFree(o orm.Ormer, vacancyId int, avitoIds []string) error {
	if len(avitoIds) == 0 {
		return nil
	}

	var (
		ph   []string
		args []interface{}
	)
	for _, id := range avitoIds {
		ph = append(ph, "?")
		args = append(args, id)
	}
	args = append(args, vacancyId)

	sql := fmt.Sprintf(`
		SELECT avito_id
		FROM avito_vacancy_ids
		WHERE avito_id IN (%s) AND vacancy_avito_id <> ?`,
		strings.Join(ph, ","))

	var duplicates []string
	if _, err := o.Raw(sql, args...).QueryRows(&duplicates); err != nil {
		return err
	}
	if len(duplicates) != 0 {
		return fmt.Errorf("avito_id %v уже привязан(ы) к другим РК", duplicates)
	}
	return nil
}

func ExportAvitoIdsToExcel(vacancyId int, o orm.Ormer) (string, error) {
	sql := `
		SELECT a.avito_id,
		       to_char(h.start_date, 'YYYY-MM-DD') AS start_date
		FROM avito_vacancy_ids a
		LEFT JOIN LATERAL (
		    SELECT start_date
		    FROM avito_vacancy_ids_hist h
		    WHERE h.vacancy_avito_id = a.vacancy_avito_id
		      AND h.avito_id = a.avito_id
		    ORDER BY h.start_date DESC NULLS LAST, h.id DESC
		    LIMIT 1
		) h ON TRUE
		WHERE a.vacancy_avito_id = ?
	`

	var rows []avitoRow
	if _, err := o.Raw(sql, vacancyId).QueryRows(&rows); err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", fmt.Errorf("у РК %d нет связанных Авито ID", vacancyId)
	}

	xl := excelize.NewFile()
	sheet := xl.GetSheetName(1)

	xl.SetCellValue(sheet, "A1", "Авито ID")
	xl.SetCellValue(sheet, "B1", "Дата привязки к РК")

	for i, r := range rows {
		rowIdx := i + 2
		dateCell := "—"
		if r.StartDate != nil && *r.StartDate != "" {
			dateCell = *r.StartDate
		}
		xl.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), r.AvitoID)
		xl.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), dateCell)
	}

	tmp, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("rk_%d_*.xlsx", vacancyId))
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	if err := xl.SaveAs(path); err != nil {
		return "", err
	}

	return path, nil
}

func PreviewNewAvitoIds(
	vacancyId int,
	inputType string,
	file multipart.File,
	hdr *multipart.FileHeader,
	manualRaw string,
	o orm.Ormer,
) (int, error) {
	var (
		avitoIds []string
		err      error
	)

	switch strings.ToLower(strings.TrimSpace(inputType)) {
	case "file":
		if file == nil {
			return 0, fmt.Errorf("файл не передан")
		}
		defer file.Close()
		avitoIds, err = ExtractAvitoIds(file, hdr)
	case "manual":
		avitoIds = utils.CleanAvitoIds(manualRaw)
	default:
		return 0, fmt.Errorf("inputType должен быть 'file' или 'manual'")
	}

	if err != nil {
		return 0, err
	}
	if len(avitoIds) == 0 {
		return 0, fmt.Errorf("не найдено ни одного Avito-ID")
	}

	if err := assertAvitoIdsFree(o, vacancyId, avitoIds); err != nil {
		return 0, err
	}

	if vacancyId == 0 {
		return 0, fmt.Errorf("РК ещё не создана")
	}

	var rows []models.Avito_vacancy_ids
	_, _ = o.QueryTable("avito_vacancy_ids").
		Filter("vacancy_avito_id", vacancyId).
		All(&rows, "Avito_id")

	existing := make(map[string]struct{}, len(rows))
	for _, r := range rows {
		existing[r.Avito_id] = struct{}{}
	}

	newCnt := 0
	for _, id := range avitoIds {
		if _, ok := existing[id]; !ok {
			newCnt++
		}
	}
	return newCnt, nil
}
