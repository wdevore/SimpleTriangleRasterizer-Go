# SimpleTriangleRasterizer-Go
A simple Bresenham triangle rasterizer that references on Ammeraal's graphics book and uses SDL for the window.

![Triangle Rasterizer](TriangleRasterizer.png)

This rasterizer renders a *single* triangle. It does not support shared edges based on Polygons. However, adding that feature is simply a matter of adding a Polygon class that tracks inside and outside edges. If you are considering using this for FPGAs then you will need to think about graphic pipelines and perhaps consider the Barycentric algorithm as well.

This rasterizer supports Translucency and Overdraw. It is based on the Top-Left algorithm. However, because this code only renders independent triangles it does draw the *right* side edges, but it doesn't draw the **Top**'s bottom edge which saves on overdraw at the X intercept.

# Code
Even in Go there is a large chunk of code needed just to provide an output for the rasterization. A majority of the code is simply just to construct and manage the window, **SDL** and render loop.

The main rasterizing code is contains within ...

# References
- [sunshine2k](http://www.sunshine2k.de/coding/java/TriangleRasterization/TriangleRasterization.html) specifics for scanline rasterization. It also covers the Barycentric Algorithm as well.
- [Ammeraal's book](https://smile.amazon.com/Computer-Graphics-Java-Programmers-Ammeraal/dp/0470031603/ref=sr_1_1?dchild=1&keywords=leen+ammeraal+graphics&qid=1625413592&sr=8-1) specifics for line drawing.
