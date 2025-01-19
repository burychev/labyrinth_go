package main

import (
 "bufio"
 "fmt"
 "os"
 "strconv"
 "strings"
)

type Point struct {
 x, y int
}

var directions = []Point{
 {-1, 0}, 
 {1, 0}, 
 {0, -1}, 
 {0, 1},  
}

func readInputFromFile(inputFile string) (int, int, [][]int, Point, Point, error) {
 file, err := os.Open(inputFile)
 if err != nil {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("ошибка открытия файла: %v", err)
 }
 defer file.Close()

 scanner := bufio.NewScanner(file)


 if !scanner.Scan() {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("отсутствуют данные о размере лабиринта")
 }
 dimensions := strings.Fields(scanner.Text())
 if len(dimensions) != 2 {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректные размеры лабиринта")
 }
 rows, err := strconv.Atoi(dimensions[0])
 if err != nil {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректное значение строк")
 }
 cols, err := strconv.Atoi(dimensions[1])
 if err != nil {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректное значение столбцов")
 }

 for i := 0; i < rows; i++ {
  if !scanner.Scan() {
   return 0, 0, nil, Point{}, Point{}, fmt.Errorf("лабиринт содержит меньше строк, чем указано")
  }
  row := strings.Fields(scanner.Text())
  if len(row) != cols {
   return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректное количество столбцов в строке %d", i)
  }
  maze[i] = make([]int, cols)
  for j, value := range row {
   maze[i][j], err = strconv.Atoi(value)
   if err != nil {
    return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректное значение клетки: %v", value)
   }
  }
 }

 if !scanner.Scan() {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("отсутствуют стартовые и конечные координаты")
 }
 startEnd := strings.Fields(scanner.Text())
 if len(startEnd) != 4 {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректное количество координат старта и финиша")
 }
 startRow, err := strconv.Atoi(startEnd[0])
 if err != nil {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректная стартовая строка")
 }
 startCol, err := strconv.Atoi(startEnd[1])
 if err != nil {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректный стартовый столбец")
 }
 endRow, err := strconv.Atoi(startEnd[2])
 if err != nil {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректная конечная строка")
 }
 endCol, err := strconv.Atoi(startEnd[3])
 if err != nil {
  return 0, 0, nil, Point{}, Point{}, fmt.Errorf("некорректный конечный столбец")
 }

 return rows, cols, maze, Point{startRow, startCol}, Point{endRow, endCol}, nil
}

func bfs(maze [][]int, start, end Point) ([]Point, error) {
 rows := len(maze)
 cols := len(maze[0])

 visited := make([][]bool, rows)
 for i := range visited {
  visited[i] = make([]bool, cols)
 }

 parent := make(map[Point]Point)
 queue := []Point{start}
 visited[start.x][start.y] = true

 for len(queue) > 0 {
  current := queue[0]
  queue = queue[1:]

  if current == end {
   path := []Point{}
   for current != start {
    path = append([]Point{current}, path...)
    current = parent[current]
   }
   path = append([]Point{start}, path...)
   return path, nil
  }

  for _, dir := range directions {
   next := Point{current.x + dir.x, current.y + dir.y}
   if next.x >= 0 && next.x < rows && next.y >= 0 && next.y < cols &&
    maze[next.x][next.y] != 0 && !visited[next.x][next.y] {
    visited[next.x][next.y] = true
    parent[next] = current
    queue = append(queue, next)
   }
  }
 }

 return nil, fmt.Errorf("путь не найден")
}

func writeOutputToFile(outputFile string, path []Point) error {
 file, err := os.Create(outputFile)
 if err != nil {
  return fmt.Errorf("ошибка создания файла: %v", err)
 }
 defer file.Close()

 writer := bufio.NewWriter(file)
 for _, point := range path {
  fmt.Fprintf(writer, "%d %d\n", point.x, point.y)
 }
 fmt.Fprintln(writer, ".")
 return writer.Flush()
}

func main() {
 if len(os.Args) < 3 {
  fmt.Fprintln(os.Stderr, "использование: program <input_file> <output_file>")
  os.Exit(1)
 }

 inputFile := os.Args[1]
 outputFile := os.Args[2]

 rows, cols, maze, start, end, err := readInputFromFile(inputFile)
 if err != nil {
  fmt.Fprintln(os.Stderr, err)
  os.Exit(1)
 }

 path, err := bfs(maze, start, end)
 if err != nil {
  fmt.Fprintln(os.Stderr, err)
  os.Exit(1)
 }

 err = writeOutputToFile(outputFile, path)
 if err != nil {
  fmt.Fprintln(os.Stderr, err)
  os.Exit(1)
 }
}
