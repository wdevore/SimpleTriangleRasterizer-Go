package api

// ITriangle defines triangles
type ITriangle interface {
	Set(x1, y1, x2, y2, x3, y3 int)
	Draw(raster IRasterBuffer)
	Fill(raster IRasterBuffer)
}
