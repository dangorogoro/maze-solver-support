package main
import (
  "fmt"
  "flag"
  "os"
  "image"
  _ "golang.org/x/image/bmp"
  "math"
)
func count_row_black(img image.Image, y int) (count int) {
  bounds := img.Bounds()
  for i := 0; i < bounds.Max.X; i++ {
    _,_,_,a := img.At(i,y).RGBA()
    if a != 0 {
      count++;
    }
  }
  return
}
func count_col_black(img image.Image, x int) (count int) {
  bounds := img.Bounds()
  for i := 0; i < bounds.Max.Y; i++ {
    _,_,_,a := img.At(x,i).RGBA()
    if a != 0 {
      count++;
    }
  }
  return
}
func calc_line_width(img image.Image) (width int) {
  min_x, _, _, _ := find_rectangle(img)
  last_count := count_col_black(img, min_x)
  width = 0
  for {
    width++
    count := count_col_black(img, min_x + width)
    if last_count / 2 > count {
      break
    }
  }
  return
}
func calc_maze_size(img image.Image) (x_size int, y_size int){
  const step = 20
  min_x, max_x, min_y, max_y := find_rectangle(img)
  width := calc_line_width(img)

  i := width
  start_length := width
  flag := false
  for {
    count := 0
    for j := min_y; j < min_y + step; j++ {
      _,_,_,a := img.At(min_x + i,j).RGBA()
      if a != 0 {
        count++;
      }
    }
    if count > step / 2 && flag == false {
      flag = true
      start_length = i
    }
    if flag == true && count < step / 2 {
      tmp := float64(max_x - min_x - width) / (float64(i + start_length) / 2 - float64(width) / 2)
      x_size = int(math.Floor(tmp + 0.5))
      break
    }
    i++
  }

  i = width
  start_length = width
  flag = false
  for {
    count := 0
    for j := min_x; j < min_x + step; j++ {
      _,_,_,a := img.At(j,min_y + i).RGBA()
      if a != 0 {
        count++;
      }
    }
    if count > step / 2 && flag == false {
      flag = true
      start_length = i
    }
    if flag == true && count < step / 2 {
      tmp := float64(max_y - min_y - width) / (float64(i + start_length) / 2 - float64(width) / 2)
      y_size = int(math.Floor(tmp + 0.5))
      break
    }
    i++
  }
  return
}
func find_rectangle(img image.Image) (min_x int, max_x int, min_y int, max_y int){
  bounds := img.Bounds()
  threshold := bounds.Max.Y / 2
  if bounds.Max.X < bounds.Max.Y {
    threshold = bounds.Max.X / 2
  }
  //fmt.Println(bounds.String())
  flag := false
  for i := 0; i < bounds.Max.Y; i++ {
    count := 0;
    for j := 0; j < bounds.Max.X; j++ {
      _,_,_,a := img.At(j,i).RGBA()
      if a != 0 {
        count++;
      }
    }
    if count > threshold {
      //fmt.Printf("Row %d: ", i)
      //fmt.Println(count)
      if flag == false {
        min_y = i
        flag = true
      } else {
        max_y = i
      }
    }
  }
  flag = false
  for i := 0; i < bounds.Max.X; i++ {
    count := 0;
    for j := 0; j < bounds.Max.Y; j++ {
      _,_,_,a := img.At(i,j).RGBA()
      if a != 0 {
        count++;
      }
    }
    if count > threshold {
      //fmt.Printf("Col %d: ", i)
      //fmt.Println(count)
      if flag == false {
        min_x = i
        flag = true
      } else {
        max_x = i
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
  fmt.Println(img.Bounds())
  fmt.Println(find_rectangle(img))
  fmt.Println(calc_maze_size(img))
  //var w io.writer
  //bmp.Encode(w, 
}
