package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"math/rand"
	"time"
)

type Sudoku [9][9]int

// 生成完整数独
func generateFullSudoku(board *Sudoku) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] == 0 {
				nums := rand.Perm(9)
				for _, n := range nums {
					num := n + 1
					if isValid(board, i, j, num) {
						board[i][j] = num
						if generateFullSudoku(board) {
							return true
						}
						board[i][j] = 0
					}
				}
				return false
			}
		}
	}
	return true
}

// 校验数字有效性
func isValid(board *Sudoku, row, col, num int) bool {
	// 检查行和列
	for i := 0; i < 9; i++ {
		if board[row][i] == num || board[i][col] == num {
			return false
		}
	}

	// 检查3x3宫格
	startRow, startCol := row-row%3, col-col%3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[startRow+i][startCol+j] == num {
				return false
			}
		}
	}
	return true
}

// 挖空生成题目
func digHoles(board *Sudoku, holes int) {
	for count := 0; count < holes; {
		row, col := rand.Intn(9), rand.Intn(9)
		if board[row][col] != 0 {
			board[row][col] = 0
			count++
		}
	}
}

// 以下保持不变
func main() {
	// ... [保持原有main函数不变] ...
	rand.Seed(time.Now().UnixNano())

	// 创建GUI应用
	myApp := app.New()
	window := myApp.NewWindow("清新数独")
	window.Resize(fyne.NewSize(500, 500))

	// 生成初始数独
	var board Sudoku
	generateFullSudoku(&board)
	digHoles(&board, 40)

	// 创建数独网格
	grid := createSudokuGrid(board)

	// 控制面板
	// 控制面板
	toolbar := container.NewVBox(
		container.NewBorder(
			nil,
			widget.NewSeparator(), // 底部添加分割线
			nil,
			nil,
			container.NewHBox( // 水平排列按钮
				widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
					generateFullSudoku(&board)
					digHoles(&board, 40)
					refreshGrid(grid, board)
				}),
				widget.NewSeparator(),
				widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
					dialog.ShowInformation("规则", "每行、每列和每个3x3宫格需包含1-9不重复", window)
				}),
			),
		),
		widget.NewSeparator(), // 添加间隔线
	)
	// 创建游戏容器
	gameContainer := container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		container.NewScroll(grid), // 添加滚动容器
	)

	// 调整窗口尺寸（修改原有 500,500 的设置）
	window.Resize(fyne.NewSize(500, 600)) // 增加窗口高度
	gameContainer.Add(toolbar)
	gameContainer.Add(grid)
	window.SetContent(gameContainer)
	window.ShowAndRun()
}

// 创建数独网格界面
func createSudokuGrid(board Sudoku) *fyne.Container {
	cells := make([][]fyne.CanvasObject, 9)

	for i := 0; i < 9; i++ {
		cells[i] = make([]fyne.CanvasObject, 9)
		for j := 0; j < 9; j++ {
			entry := widget.NewEntry()
			entry.TextStyle = fyne.TextStyle{Bold: true}
			// 修复对齐方式
			entry.Text = " " // 添加空格占位符实现视觉居中
			entry.MultiLine = true
			entry.SetMinRowsVisible(3) // 增加可视区域

			if board[i][j] != 0 {
				entry.SetText(fmt.Sprintf(" %d ", board[i][j])) // 添加空格实现视觉居中
				entry.Disable()
				entry.TextStyle.Bold = false
			}

			// 设置验证器...
			entry.Validator = func(s string) error {
				if s == "" || (s >= "1" && s <= "9") {
					return nil
				}
				return fmt.Errorf("invalid value")
			}

			if board[i][j] != 0 {
				entry.SetText(fmt.Sprintf("%d", board[i][j]))
				entry.Disable()
				entry.TextStyle.Bold = false
				entry.Validator = nil
			}

			//entry.Alignment = fyne.TextAlignCenter
			entry.MultiLine = true
			cells[i][j] = entry
		}
	}
	rows := make([]fyne.CanvasObject, 9)
	for i := 0; i < 9; i++ {
		rows[i] = container.NewGridWithColumns(9, createRow(cells[i])...)
	}
	return container.NewGridWithRows(9, rows...)
}

// 创建带边框的行
func createRow(cells []fyne.CanvasObject) []fyne.CanvasObject {
	// ... [保持原有行创建逻辑不变] ...
	var row []fyne.CanvasObject
	for i := 0; i < 9; i++ {
		// 在 createSudokuGrid 的单元格创建处添加
		entry := widget.NewEntry()
		entry.MultiLine = false     // 禁用多行显示
		entry.SetMinRowsVisible(1)  // 减少最小行数
		entry.SetMinCharsVisible(1) // 减少最小字符宽度

		// 在 createRow 函数中调整单元格间距
		cell := container.NewPadded(cells[i])
		cell = container.NewPadded(cell) // 双重填充减少单元格尺寸

		if i%3 == 0 && i != 0 {
			cell = container.NewBorder(nil, nil,
				widget.NewSeparator(), nil, cell)
		}
		row = append(row, cell)
	}
	return row
}

// 刷新界面
func refreshGrid(grid *fyne.Container, board Sudoku) {
	// ... [保持原有刷新逻辑不变] ...
	newGrid := createSudokuGrid(board)
	grid.Objects = newGrid.Objects
	grid.Refresh()
}
