package taller

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =======================
// Tipos y constantes
// =======================

type Categoria int

const (
	CatA Categoria = iota
	CatB
	CatC
)

func (c Categoria) String() string {
	switch c {
	case CatA:
		return "A"
	case CatB:
		return "B"
	case CatC:
		return "C"
	default:
		return "?"
	}
}

type Incidencia string

const (
	IncMecanica   Incidencia = "mecánica"
	IncElectrica  Incidencia = "eléctrica"
	IncCarroceria Incidencia = "carrocería"
)

type Fase int

const (
	FaseLlegada Fase = iota + 1
	FaseMecanico
	FaseLimpieza
	FaseEntrega
)

func (f Fase) String() string {
	switch f {
	case FaseLlegada:
		return "Llegada"
	case FaseMecanico:
		return "Mecánico"
	case FaseLimpieza:
		return "Limpieza"
	case FaseEntrega:
		return "Entrega"
	default:
		return "?"
	}
}

// =======================
// Estructuras
// =======================

type Coche struct {
	ID         int
	Categoria  Categoria
	Incidencia Incidencia
}

type Config struct {
	NumPlazas     int
	NumMecanicos  int
	MaxColaFase1  int
	MaxColaFase2  int
	MaxColaFase3  int
	MaxColaFase4  int
	CochesA       int
	CochesB       int
	CochesC       int
	VariacionMs   int
	SemillaRandom int64
	BaseTiempoA   time.Duration
	BaseTiempoB   time.Duration
	BaseTiempoC   time.Duration
}

type Stats struct {
	TotalCoches int
	Duracion    time.Duration
	Estrategia  string
}

// =======================
// Funciones auxiliares
// =======================

func incidenciaFromCat(cat Categoria) Incidencia {
	switch cat {
	case CatA:
		return IncMecanica
	case CatB:
		return IncElectrica
	case CatC:
		return IncCarroceria
	default:
		return IncMecanica
	}
}

func baseDurationForCat(cfg Config, cat Categoria) time.Duration {
	switch cat {
	case CatA:
		return cfg.BaseTiempoA
	case CatB:
		return cfg.BaseTiempoB
	case CatC:
		return cfg.BaseTiempoC
	default:
		return time.Second
	}
}

func variedDuration(base time.Duration, cfg Config, rnd *rand.Rand) time.Duration {
	if cfg.VariacionMs <= 0 {
		return base
	}
	delta := rnd.Intn(2*cfg.VariacionMs+1) - cfg.VariacionMs
	return base + time.Duration(delta)*time.Millisecond
}

// =======================
// Logger concurrente
// =======================

type Logger struct {
	mu      sync.RWMutex
	start   time.Time
	enabled bool
}

func NewLogger(enabled bool) *Logger {
	return &Logger{
		start:   time.Now(),
		enabled: enabled,
	}
}

func (l *Logger) Log(c Coche, fase Fase, estado string) {
	if !l.enabled {
		return
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	fmt.Printf("Tiempo %v | Coche %d | Incidencia %s | Fase %s | %s | Categoria %s\n",
		time.Since(l.start).Truncate(time.Millisecond),
		c.ID, c.Incidencia, fase.String(), estado, c.Categoria.String())
}
