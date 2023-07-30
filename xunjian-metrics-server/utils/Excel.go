package utils

import (
	"fmt"
	"log"
	"os"
	"reflect"

	excelize "github.com/xuri/excelize/v2"
)

func checkFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func removeFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func openExcel(filename string) *excelize.File {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	return f
}

func CreateExcel(filename string) {
	if checkFileExist(filename) {
		removeFile(filename)
	}
	excelf := excelize.NewFile()
	excelf.SetSheetName("Sheet1", "POD巡检")

	err := excelf.SaveAs(filename)
	if err != nil {
		fmt.Println("excel create error")
		fmt.Println(err.Error())
	}
}

func deleteSheet(f *excelize.File, sheetName string) {
	f.DeleteSheet(sheetName)
}

func createSheet(f *excelize.File, sheetName string) int {
	index := f.NewSheet(sheetName)
	return index
}

func closeWorkBook(f *excelize.File) error {
	err := f.Close()
	if err != nil {
		fmt.Println("close workbook fail")
		return fmt.Errorf(err.Error())
	}
	return nil
}

func titileStyle(excelF *excelize.File) int {
	// 设置title样式
	headerStyle, err := excelF.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Bold: false, Size: 11},
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return headerStyle
}

func contentStyle(excelF *excelize.File, content interface{}, flag interface{}) int {
	// 判断二个interface类型是否一致
	if reflect.TypeOf(&content).Kind() != reflect.TypeOf(&flag).Kind() {
		log.Fatal("type error")
	}
	switch content.(type) {
	case int:
		// 如果给出的状态与flag模式匹配上了，则赋予单元格红色
		value, _ := content.(int)
		pattern, _ := flag.(int)
		if value > pattern {
			return contentRedStyle(excelF)
		} else {
			return contentNormalStyle(excelF)
		}
	case float64:
		value, _ := content.(float64)
		pattern, _ := flag.(float64)
		if value > pattern {
			return contentRedStyle(excelF)
		} else {
			return contentNormalStyle(excelF)
		}
	case string:
		// 如果给出的状态与flag模式匹配上了，则赋予单元格红色
		value, _ := content.(string)
		pattern, _ := flag.(string)
		re := matchString(pattern, value)
		if re {
			return contentRedStyle(excelF)
		} else {
			return contentNormalStyle(excelF)
		}
	case nil:
		return contentNormalStyle(excelF)
	}
	return 0
}

func contentRedStyle(excelF *excelize.File) int {
	leftStyle, err := excelF.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FF0000"},
		},
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return leftStyle
}

func contentNormalStyle(excelF *excelize.File) int {
	leftStyle, err := excelF.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return leftStyle
}
