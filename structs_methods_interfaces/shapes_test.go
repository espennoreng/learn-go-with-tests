package structsmethodsinterfaces

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

type Shape interface {
	Area() float64
}

func TestArea(t *testing.T) {

	areaTests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{name: "rectangle", shape: Rectangle{Width: 20.0, Height: 10.0}, want: 200.0},
		{name: "circle", shape: Circle{Radius: 10.0}, want: 314.1592653589793},
		{name: "triangle", shape: Triangle{10.0, 6}, want: 30.0},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			hastArea := tt.shape.Area()
			if hastArea != tt.want {
				t.Errorf("%#v hastArea %g want %g", tt.shape, hastArea, tt.want)
			}
		})

	}
}
