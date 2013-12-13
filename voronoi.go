package main

import (
   "image"
   "image/color"
   "image/jpeg"
   "math"
   "math/rand"
   "time"
   "os"
)

var t = time.Now().UnixNano()
var r = rand.New(rand.NewSource(t))

func randPoint(bounds image.Point) image.Point {
   x := r.Intn(bounds.X)
   y := r.Intn(bounds.Y)
   return image.Point{x, y}
}

func dist(p1, p2 image.Point) float64 {
   a1, a2 := float64(p1.X), float64(p1.Y)
   b1, b2 := float64(p2.X), float64(p2.Y)

   d := math.Pow(a1 - b1, 2) + math.Pow(a2 - b2, 2)
   d = math.Sqrt(d)
   return d
}

func dist2(p1, p2 image.Point) float64 {
   a1, a2 := float64(p1.X), float64(p1.Y)
   b1, b2 := float64(p2.X), float64(p2.Y)

   d := math.Abs(a1-b1) + math.Abs(a2-b2)
   return d
}

func drawSquare(target image.Point, m *image.RGBA) {
   size := 7
   c := color.RGBA{255, 255, 255, 0}

   for i := 0; i < size; i++ {
      for j := 0; j < size; j++ {
         m.SetRGBA(target.X + j, target.Y+i, c)
      }
   }
}

func writeImage(wname string, m *image.RGBA) {
   outFile, _ := os.Create(wname)
   jpeg.Encode(outFile, m, &jpeg.Options{jpeg.DefaultQuality})
   defer outFile.Close()
}

func getRandPoints(num int, bounds image.Point) []image.Point {
   ret := make([]image.Point, num)
   for i := 0; i < num; i++ {
      ret[i] = randPoint(bounds)
   }
   return ret
}

func getMin(list []float64) (float64, int) {
   index := 0
   min := list[index]

   for i, f := range list {
      if min > f {
         min = f
         index = i
      }
   }
   return min, index
}

func computeRegion(regions []vRegion, m *image.RGBA) {
   x := m.Bounds().Max.X
   y := m.Bounds().Max.Y
   for i := 0; i < x; i++ {
      for j := 0; j < y; j++ {
         p1 := image.Point{i, j}
         distLst := make([]float64, len(regions))
         for i, r := range regions {
            distLst[i] = dist(p1, r.site)
         }
         _, minIndex := getMin(distLst)
         regions[minIndex].region = append(regions[minIndex].region, p1)
      }
   }
}

func randRGBA() color.RGBA {
   red := uint8(r.Intn(255))
   green := uint8(r.Intn(255))
   blue := uint8(r.Intn(255))
   alpha := uint8(r.Intn(10))
   return color.RGBA{red, green, blue, alpha}
}

func drawRegion(v vRegion, target *image.RGBA) {
   regionValues := v.region
   c := randRGBA()
   for _, r := range regionValues {
      target.SetRGBA(r.X, r.Y, c)
   }
}

type vRegion struct {
   site image.Point
   region []image.Point
}

func main() {
   width, height := 700, 900
   bounds := image.Point{700, 900}

   m := image.NewRGBA(image.Rect(0, 0, width, height))
   regionCount := 512
   regions := make([]vRegion, regionCount)
   sites := getRandPoints(regionCount, bounds)

   for i, s := range sites {
      regions[i].site = s
   }

   computeRegion(regions, m)

   for _, r := range regions {
      drawRegion(r, m)
      //drawSquare(r.site, m)
   }

   writeImage("out.jpg", m)
}
