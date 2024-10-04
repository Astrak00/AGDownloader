package courses

import (
	"testing"
)

// Test the function getCoursesNamesLanguages
func TestGetCoursesNamesLanguages(t *testing.T) {
	var tests = []struct {
		fullName string
		nameEs   string
		nameEN   string
	}{
		{"Curso de Prueba 24/25-1CTestCourse 24/25-1S", "Curso de Prueba 24/25-1C", "TestCourse 24/25-1S"},
		{"MAG Cálculo Diferencial Aplicado 24/25-2CApplied Differential Calculus 24/25-2S", "MAG Cálculo Diferencial Aplicado 24/25-2C", "Applied Differential Calculus 24/25-2S"},
		{"Sala de Estudiantes Grado en Ingeniería InformáticaBachelor in Computer Science and Engineering 24/25", "Sala de Estudiantes Grado en Ingeniería Informática", "Bachelor in Computer Science and Engineering 24/25"},
		{"Sala Convenio-Bilateral de Paquito el Chocolatero 24/25Convenio-Bilateral students room - Paquito el Chocolatero 24/25", "Sala Convenio-Bilateral de Paquito el Chocolatero 24/25", "Convenio-Bilateral students room - Paquito el Chocolatero 24/25"},
		{"MAG. Sistemas Distribuidos 24/25-2CMAG. Distributed Systems 24/25-S2", "MAG. Sistemas Distribuidos 24/25-2C", "MAG. Distributed Systems 24/25-S2"},
		{"Fundamentos de internet de las cosas 24/25-2CFoundations of internet of things 24/25-S2", "Fundamentos de internet de las cosas 24/25-2C", "Foundations of internet of things 24/25-S2"},
	}

	for _, tt := range tests {
		t.Run(tt.fullName, func(t *testing.T) {
			ans1, ans2 := getCoursesNamesLanguages(tt.fullName)
			if ans1 != tt.nameEs {
				t.Errorf("got %s, want %s", ans1, tt.nameEs)
			}
			if ans2 != tt.nameEN {
				t.Errorf("got %s, want %s", ans2, tt.nameEN)
			}
		})
	}
}
