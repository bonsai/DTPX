# DTPX Ontology

## 1. Entity

```
Person
Agent
Project
Content
Asset
Document
Page
Element
Text
Image
Table
Graphic
Geometry
Color
ColorSpace
Colorant
Ink
Separation
Dot
Font
Profile
PDL
PDF
PostScript
TeX
HTML
CSS
DXF
RIP
Plate
Printer
Paper
Screen
Proof
Output
Experience
```

## 2. Relations

```
author -> creates -> content
editor -> structures -> content
content -> becomes -> document
document -> contains -> page
page -> contains -> element
element -> has_geometry -> geometry
element -> has_color -> color
color -> belongs_to -> color_space
color -> maps_to -> colorant
colorant -> realized_by -> ink
document -> described_by -> PDL
PostScript/PDF -> interpreted_by -> RIP
RIP -> produces -> raster
raster -> becomes -> dot
separation -> drives -> plate
plate -> drives -> printer
printer -> deposits -> ink
ink -> interacts_with -> paper
HTML/CSS -> rendered_by -> browser
TeX -> generates -> PDF
DXF -> describes -> geometry
proof -> observes -> output
output -> produces -> experience
experience -> feeds_back_to -> editor
```

## 3. Transformation

Every important edge is a transformation:

```
Meaning → Structure → Geometry → Color → Device data → Physical/Digital output → Perception
```

The ontology therefore distinguishes **entity** from **representation** and **representation** from **materialization**.
