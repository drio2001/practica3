package taller

import (
	"testing"
	"time"
)

func baseConfig() Config {
	return Config{
		NumPlazas:    5,
		NumMecanicos: 3,
		MaxColaFase1: 10,
		MaxColaFase2: 10,
		MaxColaFase3: 10,
		MaxColaFase4: 10,
		VariacionMs:  500,
		BaseTiempoA:  5 * time.Second,
		BaseTiempoB:  3 * time.Second,
		BaseTiempoC:  1 * time.Second,
	}
}

func TestCaso1_Equitativo(t *testing.T) {
	cfg := baseConfig()
	cfg.CochesA = 10
	cfg.CochesB = 10
	cfg.CochesC = 10

	stats := RunSimulationRWMutex(cfg, true)
	t.Logf("Caso 1 -> %v", stats.Duracion)
}

func TestCaso2_Desigual_Mecanica(t *testing.T) {
	cfg := baseConfig()
	cfg.CochesA = 20
	cfg.CochesB = 5
	cfg.CochesC = 5

	stats := RunSimulationRWMutex(cfg, true)
	t.Logf("Caso 2 -> %v", stats.Duracion)
}

func TestCaso3_Desigual_Carroceria(t *testing.T) {
	cfg := baseConfig()
	cfg.CochesA = 5
	cfg.CochesB = 5
	cfg.CochesC = 20

	stats := RunSimulationRWMutex(cfg, true)
	t.Logf("Caso 3 -> %v", stats.Duracion)
}
