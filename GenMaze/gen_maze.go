package main

import (
  "bufio"
  "fmt"
  "os"
  "flag"
  "strings"
  "regexp"
)
var EAST  int8 = 0x01
var NORTH int8 = 0x02
var WEST  int8 = 0x04
var SOUTH int8 = 0x08
var FILE_PATH string = "cppCode/"
var START = make([]int, 0)
var GOAL = make([][]int, 0)

func maze_parse(mazeSize int, cnt int, line string, mazeData [][]int8) {
  for position, word := range line {
    x := position / 4
    y := mazeSize - cnt / 2 - 1
    if cnt % 2 == 0 {
      //NORTH
      if position % 4 == 2 {
        if word == '-' && y >= 0 {
          mazeData[x][y] |= NORTH
        }
        if y < mazeSize - 1 && word == '-' {
          mazeData[x][y + 1] |= SOUTH
        }
      }
    } else if cnt % 2 == 1 {
      //WEST
      if position % 4 == 0 {
          if word == '|' && y >= 0 && x < mazeSize {
            mazeData[x][y] |= WEST
          }
          if x > 0 && word == '|' {
            mazeData[x - 1][y] |= EAST
        }
      } else if position % 4 == 2 && word != ' ' {
        switch word {
          case 'G':
            GOAL = append(GOAL, []int{x,y})
          case 'S':
            START = append(START, x, y)
        }
      }
    }
  }
}
func file_write(file_name string, mazeData [][]int8) {
  fileID := FILE_PATH + file_name + ".cpp"
  fp, err := os.Create(fileID)
  if err != nil {
    panic(err)
  }
  defer fp.Close()
  start_format := "#include <MazeData.h>\nuint8_t " + file_name + "[] = {"
  data_start_format := "\n"
  data_format := "0x%x, "
  end_format := "};\n"
  fp.WriteString(start_format)
  for y :=  len(mazeData[0]) - 1; y >= 0; y-- {
    fp.WriteString(data_start_format)
    for x := 0; x <= len(mazeData[0]) - 1; x++ {
      if x == len(mazeData[0]) - 1 && y == 0 {
        fp.WriteString(fmt.Sprintf("0x%x", mazeData[x][y]))
      }else {
        fp.WriteString(fmt.Sprintf(data_format, mazeData[x][y]))
      }
    }
  }
  fp.WriteString(end_format)
  if len(START) > 0 {
    fp.WriteString("IndexVec "+ file_name + "_start" + " = ")
    fp.WriteString(fmt.Sprintf("IndexVec(%d, %d);\n", START[0], START[1]))
    fp.WriteString("std::set<IndexVec> "+ file_name + "_goal" + " = {")
  }
  for i := 0; i < len(GOAL); i++ {
    x := GOAL[i][0]
    y := GOAL[i][1]
    fp.WriteString(fmt.Sprintf("IndexVec(%d, %d)", x, y))
    if i != len(GOAL) - 1 {
      fp.WriteString(", ")
    }
  }
  fp.WriteString("};")
}
func file_check(file_name string) {
  fileID := FILE_PATH + file_name + ".cpp"
  fp, err := os.Open(fileID)
  if err != nil {
    panic(err)
  }
  defer fp.Close()
  scanner := bufio.NewScanner(fp)
  for scanner.Scan(){
    fmt.Println(scanner.Text())
  }
}
func file_read(filename string) (mazeData [][]int8) {
  fp, err := os.Open(filename)
  if err != nil {
    panic(err)
  }
  defer fp.Close()
  cnt := 0
  scanner := bufio.NewScanner(fp)
  scanner.Scan()
  line := scanner.Text()
  mazeSize := strings.Count(line, "o") - 1
  mazeData = make([][]int8, mazeSize)
  for i := range mazeData {
    mazeData[i] = make([]int8, mazeSize)
  }
  maze_parse(mazeSize, cnt, line, mazeData)
  cnt++
  for scanner.Scan(){
    line = scanner.Text()
    maze_parse(mazeSize, cnt, line, mazeData)
    cnt++
  }
  return
}

func main(){
  flag.Parse()
  args := flag.Args()
  mazeData := file_read(args[0])
  str := args[0]
  rep := regexp.MustCompile(`(.*/)*([^/]+?)?\..*$`)
  str = rep.ReplaceAllString(str, "$2")
  fmt.Println(str)
  file_write(str, mazeData)
  file_check(str)
}
