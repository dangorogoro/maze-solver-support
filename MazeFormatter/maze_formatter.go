package main
import (
  "fmt"
  "flag"
  "os"
  "image"
  _ "golang.org/x/image/bmp"
  "math"
)
type Rectangle struct{
  min_x, max_x, min_y, max_y int
}
type Index struct{
  x_size, y_size int
}
type MazeParameter struct{
  rect Rectangle
  size Index
  width int
}
func judge_color(img image.Image, x int, y int) (flag bool) {
  r,g,b,a := img.At(x,y).RGBA()
  flag = false
  //if r == 0 && g == 0 && b == 0 {
  //  flag = true
  //}
  //fmt.Println(r,g,b,a)
  if a > 1000 {
    if a == 65535 {
      if r < 50000 && g < 50000 && b < 50000 {
        flag = true
      }
    } else {
      flag = true
    }
  }
  return
}
func culc_grid_length(param MazeParameter) (x_grid_length float64, y_grid_length float64) {
  width := param.width
  x_grid_length = float64(param.rect.max_x - param.rect.min_x - width) / float64(param.size.x_size)
  y_grid_length = float64(param.rect.max_y - param.rect.min_y - width) / float64(param.size.y_size)
  return
}
func check_word(img image.Image, param MazeParameter, index Index) (str string) {
  width := param.width
  _, y_grid_length := culc_grid_length(param)
  center_x, center_y := find_center_position(param, index)
  target_x := int(math.Floor(center_x + 0.5))
  target_y := int(math.Floor(center_y + 0.5))
  length := (int(y_grid_length) - width) / 3
  max_count := 0
  present_count := 0
  sum_count := 0
  for j := target_x - length; j <= target_x + length; j++ {
    if judge_color(img, j, target_y) == true {
      present_count++
      sum_count++
    } else{
      if max_count < present_count {
        max_count = present_count
      }
      present_count = 0
    }
  }
  str = "NONE"
  if sum_count > 0 {
    //fmt.Println(sum_count, max_count, length)
    if max_count == sum_count {
      str = "START"
    } else{
      str = "GOAL"
    }
  }
  return
}
func find_center_position(param MazeParameter, index Index) (center_x float64, center_y float64) {
  width := param.width
  x_grid_length ,y_grid_length := culc_grid_length(param)
  center_x = float64(param.rect.min_x) + x_grid_length * float64(index.x_size + 1) - (x_grid_length - float64(width)) / 2
  center_y = float64(param.rect.min_y) + y_grid_length * float64(index.y_size + 1) - (y_grid_length - float64(width)) / 2
  return
}
func generate_maze(img image.Image, param MazeParameter) {
  x_size := param.size.x_size
  y_size := param.size.y_size
  //First line
  for i := 0; i < x_size; i++ {
    fmt.Printf("o");
    fmt.Printf("---");
  }
  fmt.Printf("o\n");

  for j := 0; j < y_size; j++ {
    for i := 0; i < x_size; i++ {
      if find_wall(img, param, "WEST", Index{i, j}) {
        fmt.Printf("|")
      } else {
        fmt.Printf(" ")
      }
      fmt.Printf(" ")
      str := check_word(img, param, Index{i,j})
      switch str {
        case "NONE": fmt.Printf(" ")
        case "GOAL": fmt.Printf("G")
        case "START": fmt.Printf("S")
      }
      fmt.Printf(" ")
    }
    fmt.Printf("|\n")
    for i := 0; i < x_size; i++ {
      fmt.Printf("o")
      if find_wall(img, param, "SOUTH", Index{i, j}) {
        fmt.Printf("---")
      } else {
        fmt.Printf("   ")
      }
    }
    fmt.Printf("o\n")
  }
}
func set_maze_parameter(img image.Image) (param MazeParameter) {
  param.rect = find_rectangle(img)
  param.width = calc_line_width(img, param)
  param.size = calc_maze_size(img, param)
  return
}
func find_wall(img image.Image, param MazeParameter, str string, index Index) (flag bool){
  width := param.width
  search_band := width
  x_grid_length, y_grid_length := culc_grid_length(param)
  center_x, center_y := find_center_position(param, index)
  switch str {
  case "NORTH": center_y -= y_grid_length / 2
  case "SOUTH": center_y += y_grid_length / 2
  case "EAST": center_x += x_grid_length / 2
  case "WEST": center_x -= x_grid_length / 2
  }
  target_x := int(math.Floor(center_x + 0.5))
  target_y := int(math.Floor(center_y + 0.5))
  //fmt.Println(center_x, center_y)
  count := 0
  for i := target_x - search_band; i < target_x + search_band; i++ {
    for j := target_y - search_band; j < target_y + search_band; j++ {
      if judge_color(img,i,j) == true {
        count++;
      }
    }
  }
  //fmt.Println(count)
  flag = false
  if count > int(math.Pow(float64(width * 2), 2)) / 3 {
    flag = true
  }
  return
}
func count_line_black(img image.Image, str string, x int) (count int) {
  bounds := img.Bounds()
  if str == "col" {
    for i := 0; i < bounds.Max.Y; i++ {
      if judge_color(img, x, i) == true {
        count++;
      }
    }
  } else if str == "row" {
    for i := 0; i < bounds.Max.X; i++ {
      if judge_color(img, i, x) == true {
        count++;
      }
    }
  }
  return
}
func calc_line_width(img image.Image, param MazeParameter) (width int) {
  min_x := param.rect.min_x
  last_count := count_line_black(img, "col", min_x)
  width = 0
  for {
    width++
    count := count_line_black(img, "col", min_x + width)
    if count < last_count / 2 {
      break
    }
  }
  return
}
func calc_maze_size(img image.Image, param MazeParameter) (size Index) {
  min_x := param.rect.min_x
  max_x := param.rect.max_x
  min_y := param.rect.min_y
  max_y := param.rect.max_y
  width := param.width
  step := width * 5

  var x_size, y_size int
  i := width
  flag := false
  for {
    count := 0
    for j := min_y; j < min_y + step; j++ {
      if judge_color(img, min_x + i, j) == true {
        count++;
      }
    }
    if count > step * 2 / 3{
      if flag == false {
        x_size++
        //fmt.Println(min_x + i)
        flag = true
      }
    } else {
      flag = false
    }
    //if flag == true && count < step / 2 {
    //  fmt.Println(start_length, i)
    //  grid := (float64(i + start_length - 1) / 2 - float64(width) / 2)
    //  if i - start_length < width {
    //    grid = (float64(start_length) - float64(width) / 2)
    //  }
    //  tmp := float64(max_x - min_x - width) / grid
    //  fmt.Println("x_size", tmp)
    //  x_size = int(math.Floor(tmp + 0.5))
    //  break
    //}
    i++
    if max_x <= min_x + i {
      break
    }
  }

  i = width
  flag = false
  for {
    count := 0
    for j := min_x; j < min_x + step; j++ {
      if judge_color(img, j, min_y + i) == true {
        count++;
      }
    }
    if count > step * 2 / 3 {
      if flag == false {
        //fmt.Println(min_y + i)
        y_size++
        flag = true
      }
    } else {
      flag = false
    }
    //if flag == true && count < step / 2 {
    //  fmt.Println("y",start_length, i)
    //  grid := (float64(i + start_length - 1) / 2 - float64(width) / 2)
    //  if i - start_length < width {
    //    grid = (float64(i + start_length) / 2 - float64(width) / 2)
    //  }
    //  tmp := float64(max_y - min_y - width) / grid
    //  fmt.Println("y_size", tmp)
    //  y_size = int(math.Floor(tmp + 0.5))
    //  break
    //}
    i++
    if max_y <= min_y + i {
      break
    }
  }
  fmt.Println(x_size, y_size)
  size.x_size = x_size
  size.y_size = y_size
  return
}
func find_rectangle(img image.Image) (rect Rectangle) {
  bounds := img.Bounds()
  threshold := bounds.Max.Y / 2
  if bounds.Max.X < bounds.Max.Y {
    threshold = bounds.Max.X / 2
  }
  //fmt.Println(bounds.String())
  flag := false
  for i := 0; i < bounds.Max.Y; i++ {
    count := count_line_black(img, "row", i)
    if count > threshold {
      if flag == false {
        rect.min_y = i
        flag = true
      } else {
        rect.max_y = i
      }
    }
  }
  flag = false
  for i := 0; i < bounds.Max.X; i++ {
    count := count_line_black(img, "col", i)
    if count > threshold {
      if flag == false {
        rect.min_x = i
        flag = true
      } else {
        rect.max_x = i
      }
    }
  }
  return
}
//func check_maze_size(img image.Image) (x int, y, int){
//  
//}
func main() {
  fmt.Println("Parse start")
  flag.Parse()
  args := flag.Args()
  f, _ := os.Open(args[0])
  defer f.Close()
  img, _, err := image.Decode(f)
  if err != nil {
    fmt.Println(err)
  }
  param := set_maze_parameter(img)
  //fmt.Println(img.Bounds())
  //fmt.Println(find_rectangle(img))
  //fmt.Println(calc_maze_size(img, param))
  //find_wall(img, param, "NORTH", Index{0,0})
  //find_wall(img, param, "WEST", Index{0,1})
  //find_wall(img, param, "WEST", Index{0,2})
  generate_maze(img, param)
  //fmt.Println(check_word(img, param, Index{0, 15}))
  //fmt.Println(check_word(img, param, Index{6, 6}))
  //fmt.Println(check_word(img, param, Index{12, 15}))
  //fmt.Println(check_word(img, param, Index{2, 6}))
  //bmp.Encode(w, 
}
